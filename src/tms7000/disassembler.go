package tms7000

import (
	"fmt"
	"intelhex"
)

type Disassembler struct {
	input          *intelhex.Queue
	instructionSet InstructionSet
	labels         *LabelTable
}

func NewDisassembler(instructionSet InstructionSet, input *intelhex.Queue, labels *LabelTable) *Disassembler {
	return &Disassembler{
		input:          input,
		instructionSet: instructionSet,
		labels:         labels,
	}
}

func (self *Disassembler) Do() {

	for !self.input.Empty() {
		b := self.input.Pop()

		label := self.labels.Get(b.Address)

		if label != nil {
			fmt.Printf("%s:\n", label.name)
		}

		if label != nil && label.length > 0 {
			// data entry
			str := fmt.Sprintf("%02x", b.Value)
			for i := uint(0); i < label.length-1; i++ {
				str += fmt.Sprintf("%02x", self.input.Pop().Value)
			}
			fmt.Printf("%04x        %-20s\n",
				b.Address, str)
		} else {
			// decode instruction
			instruction, ok := self.instructionSet[b.Value]

			if ok {
				args := make([]byte, instruction.Args())
				var argstring string
				for i := range args {
					args[i] = self.input.Pop().Value
					argstring += fmt.Sprintf(" %02x", args[i])
				}

				fmt.Printf("%04x  %02x %-9s  %-20s -- %s\n",
					b.Address, b.Value, argstring, instruction.String(args, b.Address, self.labels), instruction.comment)
			} else {
				fmt.Printf("%04x  %02x                                 -- unknown instruction\n", b.Address, b.Value)
			}
		}
	}
}
