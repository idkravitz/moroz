package main

import (
    "os"
    "gob"
    "fmt"
    "flag"
    "path"
)

import (
    "./rle"
    "./huffman"
)

import . "common"

var (
    useRLE = flag.Bool("r", false, "use RLE compression")
    useHuffman = flag.Bool("h", false, "use Huffman compression")
    useCompress = flag.Bool("c", false, "compress")
    useExtract = flag.Bool("x", false, "extract")
    outputName = flag.String("f", "-", "set output file (default is \"-\" -- stdout)")
)

const (
    _ = iota
    HuffmanCompress
    RLECompress
)

type (
    CompressionMethod int
    headerMeta struct {
        Method CompressionMethod
    }
    commonMeta struct {
        Name string
    }
)

func printUsage() {
    flag.Usage()
    os.Exit(0)
}

func compress (in *os.File, out *os.File, method CompressionMethod) {
    var act func(*os.File, *os.File);
    if method == RLECompress {
        act = rle.Compress
    } else {
        act = huffman.Compress
    }
    act(in, out)
}

func extract (in *os.File, out *os.File, method CompressionMethod) {
    var act func(*os.File, *os.File) int64;
    if method == RLECompress {
        act = rle.Extract
    } else {
        act = huffman.Extract
    }
    in.Seek(act(in, out), 0)
}

func panicHandler() {
    if error := recover(); error != nil {
        fmt.Printf("Error: %s", error)
    }
}

func getCompressionMethod() CompressionMethod {
    if *useRLE {
        return RLECompress
    }
    return HuffmanCompress
}

func openForWrite(name string) *os.File {
    out, err := os.Open(name, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0666)
    PanicIf(err)
    return out
}

func openForRead(name string) *os.File {
    in, err := os.Open(name, os.O_RDONLY, 0777)
    PanicIf(err)
    return in
}

func isEOF(fin *os.File) bool {
    pos, err := fin.Seek(0, 1)
    PanicIf(err)
    return pos == GetFileSize(fin)
}

func main() {
    defer panicHandler();

    var in, out *os.File
    flag.Parse()
    if flag.NArg() == 0 || (!*useCompress && !*useExtract || *useCompress && *useExtract) ||
        *useCompress && (!*useHuffman  && !*useRLE || *useHuffman  && *useRLE) {
        printUsage()
    }
    if *useCompress {
        if *outputName == "-" {
            out = os.Stdout
        } else {
            out = openForWrite(*outputName)
            defer out.Close()
        }
        PanicIf(gob.NewEncoder(out).Encode(headerMeta{getCompressionMethod()}))
    }
    for _, name := range flag.Args() {
        in = openForRead(name)
        defer in.Close()
        if *useCompress {
            PanicIf(gob.NewEncoder(out).Encode(commonMeta{path.Base(name)}))
            compress(in, out, getCompressionMethod())
        } else {
            var hmeta headerMeta
            PanicIf(gob.NewDecoder(in).Decode(&hmeta))
            for !isEOF(in) {
                var cmeta commonMeta
                PanicIf(gob.NewDecoder(in).Decode(&cmeta))
                out = openForWrite(cmeta.Name)
                defer out.Close()
                extract(in, out, hmeta.Method)
            }
        }
    }
}
