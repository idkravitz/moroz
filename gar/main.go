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

func compress (name string, method CompressionMethod) (err os.Error) {
    file, err := os.Open(name, os.O_RDONLY, 0777)
    if err != nil {
        return
    }
    defer file.Close()
    return
}

func main() {
    flag.Parse()
    if flag.NArg() == 0 ||
        (!*useCompress && !*useExtract || *useCompress && *useExtract) ||
        (!*useHuffman  && !*useRLE     || *useHuffman  && *useRLE) {
        printUsage()
    }
    for _,name := range flag.Args() {
        err := compress(name, RLECompress)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        //fmt.Println(name)
    }
}
