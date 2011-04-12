package huffman

import (
    "os"
    "fmt"
    "gob"
    "bufio"
    "container/heap"
)
import . "common"

type (
    node struct {
        count uint64
        val byte
        left *node
        right *node
    }
    nHeap []*node
    cbRecord struct {
        Len, Code uint
    }
    dbRecord struct {
        char byte
        len uint
    }
    codeBook [256]cbRecord
    decodeBook map[uint] *dbRecord
    nMeta struct {
        Cb []cbRecord
        Fsize int64
    }
)

func (heap *nHeap) Len() int {
    return len(*heap)
}

func (heap *nHeap) Less(i, j int) bool {
    return (*heap)[i].count < (*heap)[j].count
}

func (heap *nHeap) Swap(i, j int) {
    (*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i]
}

func (heap *nHeap) Push(n interface{}) {
    (*heap) = append((*heap), n.(*node))
}

func (heap *nHeap) Pop() interface{} {
    n := heap.Len() - 1
    elm := (*heap)[n]
    (*heap) = (*heap)[0:n]
    return elm
}

func countFreqs(in *os.File) [256]uint64 {
    var (
        freqs [256]uint64
        curr byte
        error os.Error
    )
    for i := 0; i < len(freqs); i++ {
        freqs[i] = 0
    }
    for bin := bufio.NewReader(in); true; {
        curr, error = bin.ReadByte()
        if error != nil {
            break
        }
        freqs[curr]++
    }
    in.Seek(0, 0)
    return freqs
}

func makeCodeBook(node *node, len, code uint, cb []cbRecord) {
    if node.left == nil && node.right == nil {
        cb[node.val] = cbRecord{len, code}
    } else {
        makeCodeBook(node.left,  len + 1, code, cb[:])
        makeCodeBook(node.right, len + 1, code | 1 << len, cb[:])
    }
}

func serializeMetaInfo(fin, fout *os.File, cb []cbRecord) {
    var meta nMeta
    meta.Cb = make([]cbRecord, 256)
    for k, v := range cb {
        meta.Cb[k] = v
    }
    meta.Fsize = GetFileSize(fin)
    PanicIf(gob.NewEncoder(fout).Encode(meta))
}

func Compress(fin, fout *os.File) {
    nodes := new(nHeap)
    freqs := countFreqs(fin)
    for b, f := range freqs {
        if f != 0 {
            nodes.Push(&node{count: f, val: byte(b)})
        }
    }
    heap.Init(nodes)
    var cb codeBook
    if nodes.Len() != 1 {
        for nodes.Len() > 1 {
            l := heap.Pop(nodes).(*node)
            r := heap.Pop(nodes).(*node)
            parent := &node{count: l.count + r.count, left: l, right: r}
            heap.Push(nodes, parent)
        }
        tree := (*nodes)[0]
        makeCodeBook(tree, 0, 0, cb[:])
    } else {
       cb[(*nodes)[0].val] = cbRecord{Len: 1, Code: 0}
    }

    serializeMetaInfo(fin, fout, cb[:])

    // encode
    var (
        outbyte, outlen byte = 0, 0
        i uint = 0
    )
    in := bufio.NewReader(fin)
    out := bufio.NewWriter(fout)
    defer out.Flush()

    for {
        inbyte, error := in.ReadByte()
        if error != nil {
            break
        }
        record := cb[inbyte]
        for i = 0; i < record.Len; {
            for ; i < record.Len && outlen < 8; i++ {
                if (record.Code & (1 << i)) != 0 {
                    outbyte |= 1 << outlen
                }
                outlen++
            }
            if outlen == 8 {
                out.WriteByte(outbyte)
                outbyte, outlen = 0, 0
            }
        }
    }
    fmt.Println(outlen)
    if outlen != 0 {
        out.WriteByte(outbyte)
    }
}

func deserializeMetaInfo(fin, fout *os.File) (int64, decodeBook) {
    var meta nMeta
    PanicIf(gob.NewDecoder(fin).Decode(&meta))
    db := make(decodeBook)
    for i, record := range meta.Cb {
        if record.Len != 0 {
            db[record.Code] = &dbRecord{char: byte(i), len: record.Len}
        }
    }
    return meta.Fsize, db
}

func Extract(fin, fout *os.File) int64 {
    var (
        code, code_len uint = 0, 0
        outptr *dbRecord= nil
        cursize, readBytes int64 = 0, 0
    )
    filesize, db := deserializeMetaInfo(fin, fout)
    pos := GetSeek(fin)

    in := bufio.NewReader(fin)
    out := bufio.NewWriter(fout)
    defer out.Flush()

    for cursize < filesize {
        curr, error := in.ReadByte()
        if error != nil {
            break
        }
        readBytes++
        for i := uint(0); i < 8; i++ {
            if (curr & (1 << i)) != 0 {
                code |= 1 << code_len
            }
            outptr = db[code]
            code_len++
            if outptr != nil && outptr.len == code_len && cursize < filesize {
                out.WriteByte(outptr.char)
                cursize++
                code, code_len = 0, 0
            }
        }
    }
    return pos + readBytes
}
