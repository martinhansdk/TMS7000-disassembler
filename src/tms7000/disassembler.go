package tms7000

import (
	"fmt"
	"intelhex"
)

type Disassembler struct {
	input          *intelhex.Queue
	instructionSet InstructionSet
}

func NewDisassembler(instructionSet InstructionSet, input *intelhex.Queue) *Disassembler {
	return &Disassembler{
		input:          input,
		instructionSet: instructionSet,
	}
}

func (self *Disassembler) Do() {

	for !self.input.Empty() {

		b := self.input.Pop()
		instruction, ok := self.instructionSet[b.Value]

		if ok {
			args := make([]byte, instruction.Args())
			var argstring string
			for i := range args {
				args[i] = self.input.Pop().Value
				argstring += fmt.Sprintf(" %02x", args[i])
			}

			fmt.Printf("%04x  %02x %-9s  %-20s -- %s\n",
				b.Address, b.Value, argstring, instruction.String(args, b.Address), instruction.comment)
		} else {
			fmt.Printf("%04x  %02x                                 -- unknown instruction\n", b.Address, b.Value)
		}
	}
}
