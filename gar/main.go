package main

import (
    "os"
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

func compress (name string, method CompressionMethod) os.Error {
    file, ok := os.Open(name, os.O_RDONLY, 0777)
    if ok != nil {
        return ok
    }
    return nil // no error
}

func main() {
    flag.Parse()
    if flag.NArg() == 0 ||
        (!*useCompress && !*useExtract || *useCompress && *useExtract) ||
        (!*useHuffman  && !*useRLE     || *useHuffman  && *useRLE) {
        printUsage()
    }
    for _,name := range flag.Args() {
        //fmt.Println(name)
    }
}
