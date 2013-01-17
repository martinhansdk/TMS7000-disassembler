package tms7000

type Instruction struct {
        name string
}

var InstructionSet = map[byte]Instruction{
        0x05: {"EINT"},
}
