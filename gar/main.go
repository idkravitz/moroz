package main

import (
    "os"
    "bufio"
    "fmt"
    "flag"
)

var (
    useRLE = flag.Bool("r", false, "use RLE compression")
    useHuffman = flag.Bool("h", false, "use Huffman compression")
    useCompress = flag.Bool("c", false, "compress")
    useExtract = flag.Bool("x", false, "extract")
    outputName = flag.String("f", "-", "set output file (default is \"-\" -- stdin)")
)

const (
    _ = iota
    HuffmanCompress
    RLECompress
)

type CompressionMethod int

func printUsage() {
    flag.Usage()
    os.Exit(0)
}

func rleCompressedWrite(b byte, c byte, out *bufio.Writer) {
    out.WriteByte(b)
    out.WriteByte(c)
}

func compressRLE(in *bufio.Reader, out *bufio.Writer) {
    var buff [512]byte
    prev,err := in.ReadByte()
    if err == nil {
        out.WriteByte(prev)
    }
    var count byte = 0
    for {
        n, ok := in.Read(buff[:])
        if ok == os.EOF {
            break
        }
        for _,v := range buff[:n] {
            switch {
            case v == prev:
                count++
                if count != 255 {
                    break
                }
                rleCompressedWrite(prev, count, out)
                count = 0
                prev = (v + 1) % 255
            case count != 0:
                rleCompressedWrite(prev, count, out)
                count = 0
                fallthrough
            default:
                out.WriteByte(v)
                prev = v
            }
        }
    }
}

func extractRLE(in *bufio.Reader, out *bufio.Writer) (error os.Error) {
    error = nil
    var buff [512]byte
    var match bool = false
    prev,err := in.ReadByte()
    if err == nil {
        out.WriteByte(prev)
    }
    for {
        n, ok := in.Read(buff[:])
        if ok == os.EOF {
            if match {
                error = os.NewError("Corrupted archive")
            }
            break
        }
        for _,v := range buff[:n] {
            if match {
                for ; v > 1; v-- {
                    out.WriteByte(prev)
                }
                match = false
            } else {
                out.WriteByte(v)
                match = prev == v
                prev = v
            }
        }
    }
    return
}

func extract(in *buffio.Reader, method CompressionMethod) (err os.Error) {
    file, err := os.Open(name, os.O_RDONLY, 0777)
    if err != nil {
        return
    }
    defer file.Close()
    rd := bufio.NewReader(file)
    if method == RLECompress {
        compressRLE(rd, out)
    }
    return
}

func compress (in *bufio.Reader, out *bufio.Writer, method CompressionMethod) {
    if method == RLECompress {
        compressRLE(in, out)
    }
}

func handleError(err os.Error) {
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func main() {
    flag.Parse()
    if flag.NArg() == 0 ||
        (!*useCompress && !*useExtract || *useCompress && *useExtract) ||
        (!*useHuffman  && !*useRLE     || *useHuffman  && *useRLE) {
        printUsage()
    }
    var out *bufio.Writer
    if *outputName == "-" {
        out = bufio.NewWriter(os.Stdout)
        defer out.Flush()
    } else {
        file, err := os.Open(*outputName, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0666)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        out = bufio.NewWriter(file)
        defer func() {
            out.Flush()
            file.Close()
        }()
    }
    for _, name := range flag.Args() {
        file, err := os.Open(name, os.O_RDONLY, 0777)
        handleError(err)
        defer file.Close()
        in := bufio.NewReader(file)
        if useCompress {
            err := compress(name, RLECompress, out)
        } else {
            err := extract(name, RLECompress, out)
        }
        handleError(err)
    }
}
