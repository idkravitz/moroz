package common

import "os"

func PanicIf(error os.Error) {
    if error != nil {
        panic(error)
    }
}

func GetFileSize(fobj *os.File) int64 {
    fileinfo, error := fobj.Stat()
    PanicIf(error)
    return fileinfo.Size
}

func GetSeek(fd *os.File) int64 {
    pos, err := fd.Seek(0, 1)
    PanicIf(err)
    return pos
}
