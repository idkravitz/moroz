package main

import (
    "os"
    "bufio"
    "fmt"
    "flag"
    "rle"
)

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

type CompressionMethod int

func printUsage() {
    flag.Usage()
    os.Exit(0)
}

/*func extract(in *buffio.Reader, method CompressionMethod) (err os.Error) {
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
}*/

func compress (in *bufio.Reader, out *bufio.Writer, method CompressionMethod) {
    if method == RLECompress {
       rle.Compress(in, out)
    }
}

func handleError(err os.Error) {
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func main() {
    defer func() {
        if error := recover(); error != nil {
            fmt.Printf("Error: %s", error)
        }
    }()
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
            panic(err)
        }
        out = bufio.NewWriter(file)
        defer func() {
            out.Flush()
            file.Close()
        }()
    }
    for _, name := range flag.Args() {
        file, err := os.Open(name, os.O_RDONLY, 0777)
        if err != nil {
            panic(err)
        }
        defer file.Close()
        in := bufio.NewReader(file)
        if *useCompress {
            compress(in, out, RLECompress)
        } else {
            rle.Extract(in, out)
        }
    }
}
