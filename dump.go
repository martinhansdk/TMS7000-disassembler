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

        hexreader := intelhex.NewReader(file, 0xe000)
        locatedBytes := intelhex.NewQueue(64353)
        disassembler := tms7000.NewDisassembler(tms7000.TMS7000InstructionSet, locatedBytes)

        err = hexreader.Iterate(func(val intelhex.LocatedByte) {
                locatedBytes.Push(&val)
        })

        disassembler.Do()

        if err != nil {
                log.Fatal(err)
        }
}
