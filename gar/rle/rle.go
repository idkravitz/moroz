package rle

import (
    "os"
    "bufio"
)

func Compress(in *bufio.Reader, out *bufio.Writer) {
    var curr, prev, count byte = 0, 0, 0
    found := false
    var error os.Error = nil
    for {
        curr, error = in.ReadByte()
        if error != nil {
            break
        }
        if found {
            if curr == prev && count < 255 {
                count++
            } else {
                out.WriteByte(count)
                out.WriteByte(curr)
                count = 0
                found = false
            }
        } else {
            out.WriteByte(curr)
            found = curr == prev
        }
        prev = curr
    }
    if count > 0 {
        out.WriteByte(count)
    }
}

func Extract(in *bufio.Reader, out *bufio.Writer) {
    var (
        curr, prev byte = 0, 0
        found, valid_prev bool = false, true
        error os.Error = nil
    )
    for error != os.EOF {
        curr, error = in.ReadByte()
        if error != nil {
            if found {
                panic("Archive corrupted")
            }
            break
        }
        if found {
            for ; curr > 0; curr-- {
                out.WriteByte(prev)
            }
            found = false
            valid_prev = false
        } else {
            out.WriteByte(curr)
            found = curr == prev && valid_prev
            prev = curr
            valid_prev = true
        }
    }
}
