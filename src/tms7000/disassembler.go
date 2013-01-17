package tms7000

import (
        "fmt"
        "intelhex"
)

type Disassembler struct {
}

func (self *Disassembler) NextByte(b intelhex.LocatedByte) {
        instruction, ok := InstructionSet[b.Value]

        if ok {
                fmt.Printf("%04x    %s\n", b.Address, instruction.name)
        }
}
