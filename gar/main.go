package main

import (
    "os"
    "fmt"
    "flag"
)

var (
    useRLE = flag.Bool("r", false, "use RLE compression")
    useHuffman = flag.Bool("h", false, "use Huffman compression")
    outputName = flag.String("f", "-", "set output file (default is \"-\" -- stdin)")
)

func main() {
    flag.Parse()
    if flag.NArg() == 0 {
        flag.Usage()
        os.Exit(0)
    }
    for _,name := range flag.Args() {
        fmt.Printf(name)
    }
}
