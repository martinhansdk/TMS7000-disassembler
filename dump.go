package main

import (
        "intelhex"
        "log"
        "os"
        "tms7000"
)

func main() {
        file, err := os.Open("OP6E000.IHX")
        if err != nil {
                log.Fatal(err)
        }

        hexreader := intelhex.NewReader(file)
        var disassembler tms7000.Disassembler

        err = hexreader.Iterate(func(val intelhex.LocatedByte) {
                disassembler.NextByte(val)
        })

        if err != nil {
                log.Fatal(err)
        }
}
