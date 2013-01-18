package tms7000

import (
	"fmt"
)

type Format struct {
	args      uint
	argformat func([]byte, uint) string
}

func noArgs(output string) func([]byte, uint) string {
	return func(args []byte, pc uint) string {
		return output
	}
}

func oneArg(format string) func([]byte, uint) string {
	return func(args []byte, pc uint) string {
		return fmt.Sprintf(format, args[0])
	}
}

func twoArgs(format string) func([]byte, uint) string {
	return func(args []byte, pc uint) string {
		return fmt.Sprintf(format, args[0], args[1])
	}
}

func label16bit(format string) func([]byte, uint) string {
	return func(args []byte, pc uint) string {
		return fmt.Sprintf(format, uint16(args[0])<<8+uint16(args[1]))
	}
}

func iop16bitOneReg(format string) func([]byte, uint) string {
	return func(args []byte, pc uint) string {
		bigval := uint16(args[0])<<8 + uint16(args[1])
		return fmt.Sprintf(format, bigval, args[2])
	}
}

func pcPlusOffset(format string) func([]byte, uint) string {
	return func(args []byte, pc uint) string {
		address := pc + uint(args[0])
		return fmt.Sprintf(format, address)
	}
}

func argPcPlusOffset(format string) func([]byte, uint) string {
	return func(args []byte, pc uint) string {
		address := pc + uint(args[1])
		return fmt.Sprintf(format, args[0], address)
	}
}

func argArgPcPlusOffset(format string) func([]byte, uint) string {
	return func(args []byte, pc uint) string {
		address := pc + uint(args[1])
		return fmt.Sprintf(format, args[0], args[1], address)
	}
}

var F_None = Format{0, noArgs("")}

var F_A = Format{0, noArgs("A")}
var F_B = Format{0, noArgs("B")}
var F_Rn = Format{1, oneArg("R%d")}
var F_Rn_A = Format{1, oneArg("R%d,A")}
var F_Rn_B = Format{1, oneArg("R%d,B")}
var F_A_Rn = Format{1, oneArg("A,R%d")}
var F_B_Rn = Format{1, oneArg("B,R%d")}
var F_ST = Format{0, noArgs("ST")}

var F_iop_A = Format{1, oneArg("%%%d,A")}
var F_iop_B = Format{1, oneArg("%%%d,B")}
var F_iop_Rn = Format{2, twoArgs("%%%d,R%d")}
var F_iop16_Rn = Format{3, iop16bitOneReg("%%>%04x,R%d")}

var F_A_B = Format{1, noArgs("A,B")}
var F_B_A = Format{1, noArgs("B,A")}

var F_Rn_Rn = Format{2, twoArgs("R%d,R%d")}
var F_A_Pn = Format{1, oneArg("A,P%d")}
var F_B_Pn = Format{1, oneArg("B,P%d")}
var F_Pn_A = Format{1, oneArg("P%d,A")}
var F_Pn_B = Format{1, oneArg("P%d,B")}
var F_iop_Pn = Format{2, twoArgs("%%%d,P%d")}

var F_iop16 = Format{2, label16bit("@>%04x")}
var F_iop16idxB = Format{2, label16bit("@>%04x(B)")}
var F_iop16idxB_Rn = Format{3, iop16bitOneReg("%>%04x(B),R%d")}

var F_starRn = Format{1, oneArg("*R%d")}

var F_offst = Format{1, pcPlusOffset("@>%04x")}
var F_A_offst = Format{1, pcPlusOffset("A,@>%04x")}
var F_B_offst = Format{1, pcPlusOffset("B,@>%04x")}
var F_B_A_offst = Format{1, pcPlusOffset("B,A,@>%04x")}
var F_Rn_offst = Format{2, argPcPlusOffset("R%d,@>%04x")}
var F_Rn_A_offst = Format{2, argPcPlusOffset("R%d,A,@>%04x")}
var F_Rn_B_offst = Format{2, argPcPlusOffset("R%d,B,@>%04x")}
var F_Rn_Rn_offst = Format{3, argArgPcPlusOffset("R%d,R%d,@>%04x")}
var F_iop_A_offst = Format{2, argPcPlusOffset("%%>%d,A,@>%04x")}
var F_iop_B_offst = Format{2, argPcPlusOffset("%%>%d,B,@>%04x")}
var F_iop_Rn_offst = Format{3, argArgPcPlusOffset("%%>%d,R%d,@>%04x")}
var F_A_Pn_offst = Format{2, argPcPlusOffset("A,P%d,@>%04x")}
var F_B_Pn_offst = Format{2, argPcPlusOffset("B,P%d,@>%04x")}
var F_iop_Pn_offst = Format{3, argArgPcPlusOffset("%%>%d,P%d,@>%04x")}

type Instruction struct {
	name    string
	comment string
	amode   Format
}

func (self *Instruction) String(args []byte, pc uint) string {
	return fmt.Sprintf("%s %s", self.name, self.amode.argformat(args, pc))
}

func (self *Instruction) Args() uint {
	return self.amode.args
}

type InstructionSet map[byte]Instruction

var TMS7000InstructionSet = InstructionSet{

	0x69: {"ADC", "", F_B_A},
	0x19: {"ADC", "", F_Rn_A},
	0x39: {"ADC", "", F_Rn_B},
	0x49: {"ADC", "", F_Rn_Rn},
	0x29: {"ADC", "", F_iop_A},
	0x59: {"ADC", "", F_iop_B},
	0x79: {"ADC", "", F_iop_Rn},

	0x68: {"ADD", "", F_B_A},
	0x18: {"ADD", "", F_Rn_A},
	0x38: {"ADD", "", F_Rn_B},
	0x48: {"ADD", "", F_Rn_Rn},
	0x28: {"ADD", "", F_iop_A},
	0x58: {"ADD", "", F_iop_B},
	0x78: {"ADD", "", F_iop_Rn},

	0x63: {"AND", "", F_B_A},
	0x13: {"AND", "", F_Rn_A},
	0x33: {"AND", "", F_Rn_B},
	0x43: {"AND", "", F_Rn_Rn},
	0x23: {"AND", "", F_iop_A},
	0x53: {"AND", "", F_iop_B},
	0x73: {"AND", "", F_iop_Rn},

	0x83: {"ANDP", "", F_A_Pn},
	0x93: {"ANDP", "", F_B_Pn},
	0xA3: {"ANDP", "", F_iop_Pn},

	0x8c: {"BR", "", F_iop16},
	0xac: {"BR", "", F_iop16idxB},
	0x9c: {"BR", "", F_starRn},

	0x66: {"BTJO", "", F_B_A_offst},
	0x16: {"BTJO", "", F_Rn_A_offst},
	0x36: {"BTJO", "", F_Rn_B_offst},
	0x46: {"BTJO", "", F_Rn_Rn_offst},
	0x26: {"BTJO", "", F_iop_A_offst},
	0x56: {"BTJO", "", F_iop_B_offst},
	0x76: {"BTJO", "", F_iop_Rn_offst},

	0x86: {"BTJOP", "", F_A_Pn_offst},
	0x96: {"BTJOP", "", F_B_Pn_offst},
	0xA6: {"BTJOP", "", F_iop_Pn_offst},

	0x67: {"BTJZ", "", F_B_A_offst},
	0x17: {"BTJZ", "", F_Rn_A_offst},
	0x37: {"BTJZ", "", F_Rn_B_offst},
	0x47: {"BTJZ", "", F_Rn_Rn_offst},
	0x27: {"BTJZ", "", F_iop_A_offst},
	0x57: {"BTJZ", "", F_iop_B_offst},
	0x77: {"BTJZ", "", F_iop_Rn_offst},

	0x87: {"BTJZP", "", F_A_Pn_offst},
	0x97: {"BTJZP", "", F_B_Pn_offst},
	0xA7: {"BTJZP", "", F_iop_Pn_offst},

	0x8e: {"CALL", "", F_iop16},
	0xae: {"CALL", "", F_iop16idxB},
	0x9e: {"CALL", "", F_starRn},

	0xb5: {"CLR", "", F_A},
	0xc6: {"CLR", "", F_B},
	0xd5: {"CLR", "", F_Rn},

	0x6d: {"CMP", "", F_B_A},
	0x1d: {"CMP", "", F_Rn_A},
	0x3d: {"CMP", "", F_Rn_B},
	0x4d: {"CMP", "", F_Rn_Rn},
	0x2d: {"CMP", "", F_iop_A},
	0x5d: {"CMP", "", F_iop_B},
	0x7d: {"CMP", "", F_iop_Rn},

	0x8d: {"CMPA", "", F_iop16},
	0xad: {"CMPA", "", F_iop16idxB},
	0x9d: {"CMPA", "", F_starRn},

	0x6e: {"DAC", "", F_A_B},
	0x1e: {"DAC", "", F_Rn_A},
	0x3e: {"DAC", "", F_Rn_B},
	0x4e: {"DAC", "", F_Rn_Rn},
	0x2e: {"DAC", "", F_iop_A},
	0x5e: {"DAC", "", F_iop_B},
	0x7e: {"DAC", "", F_iop_Rn},

	0xb2: {"DEC", "", F_A},
	0xc2: {"DEC", "", F_B},
	0xd2: {"DEC", "", F_Rn},

	0xba: {"DECD", "", F_A_offst},
	0xca: {"DECD", "", F_B_offst},
	0xda: {"DECD", "", F_Rn_offst},

	0x06: {"DINT", "Clear global interrupt enable bit", F_None},

	0xbb: {"DJNZ", "", F_A},
	0xcb: {"DJNZ", "", F_B},
	0xdb: {"DJNZ", "", F_Rn},

	0x6f: {"DSB", "", F_B_A},
	0x1f: {"DSB", "", F_Rn_A},
	0x3f: {"DSB", "", F_Rn_B},
	0x4f: {"DSB", "", F_Rn_Rn},
	0x2f: {"DSB", "", F_iop_A},
	0x5f: {"DSB", "", F_iop_B},
	0x7f: {"DSB", "", F_iop_Rn},

	0x05: {"EINT", "Set global interrupt enable bit", F_None},

	0x01: {"IDLE", "Sleep until interrupt", F_None},

	0xb3: {"INC", "", F_A},
	0xc3: {"INC", "", F_B},
	0xd3: {"INC", "", F_Rn},

	0xb4: {"INV", "", F_A},
	0xc4: {"INV", "", F_B},
	0xd4: {"INV", "", F_Rn},

	0xe0: {"JMP", "", F_offst},
	0xe2: {"JEQ", "", F_offst},
	0xe5: {"JGE", "", F_offst},
	0xe4: {"JGT", "", F_offst},
	0xe3: {"JHS", "", F_offst},
	0xe7: {"JL", "", F_offst},
	0xe6: {"JNE", "", F_offst},

	0x8a: {"LDA", "", F_iop16},
	0xaa: {"LDA", "", F_iop16idxB},
	0x9a: {"LDA", "", F_starRn},

	0x0d: {"LDSP", "(B) -> (SP)", F_None},

	0xc0: {"MOV", "", F_A_B},
	0xd0: {"MOV", "", F_A_Rn},
	0x62: {"MOV", "", F_B_A},
	0xd1: {"MOV", "", F_B_Rn},
	0x12: {"MOV", "", F_Rn_A},
	0x32: {"MOV", "", F_Rn_B},
	0x42: {"MOV", "", F_Rn_Rn},
	0x22: {"MOV", "", F_iop_A},
	0x52: {"MOV", "", F_iop_B},
	0x72: {"MOV", "", F_iop_Rn},

	0x88: {"MOVD", "", F_iop16_Rn},
	0xa8: {"MOVD", "", F_iop16idxB_Rn},
	0x98: {"MOVD", "", F_Rn_Rn},

	0x82: {"MOVP", "", F_A_Pn},
	0x92: {"MOVP", "", F_B_Pn},
	0xA2: {"MOVP", "", F_iop_Pn},
	0x80: {"MOVP", "", F_Pn_A},
	0x91: {"MOVP", "", F_Pn_B},

	0x6c: {"MOV", "", F_B_A},
	0x1c: {"MOV", "", F_Rn_A},
	0x3c: {"MOV", "", F_Rn_B},
	0x4c: {"MOV", "", F_Rn_Rn},
	0x2c: {"MOV", "", F_iop_A},
	0x5c: {"MOV", "", F_iop_B},
	0x7c: {"MOV", "", F_iop_Rn},

	0x00: {"NOP", "", F_None},

	0x64: {"OR", "", F_A_B},
	0x14: {"OR", "", F_Rn_A},
	0x34: {"OR", "", F_Rn_B},
	0x44: {"OR", "", F_Rn_Rn},
	0x24: {"OR", "", F_iop_A},
	0x54: {"OR", "", F_iop_B},
	0x74: {"OR", "", F_iop_Rn},

	0x84: {"ORP", "", F_A_Pn},
	0x94: {"ORP", "", F_B_Pn},
	0xA4: {"ORP", "", F_iop_Pn},

	0xb9: {"POP", "", F_A},
	0xc9: {"POP", "", F_B},
	0xd9: {"POP", "", F_Rn},
	0x08: {"POP", "", F_ST},

	0xb8: {"PUSH", "", F_A},
	0xc8: {"PUSH", "", F_B},
	0xd8: {"PUSH", "", F_Rn},
	0x0e: {"PUSH", "", F_ST},

	0x0b: {"RETI", "return from interrupt", F_None},

	0x0a: {"RETS", "return from subroutine", F_None},

	0xbe: {"RL", "", F_A},
	0xce: {"RL", "", F_B},
	0xde: {"RL", "", F_Rn},

	0xbf: {"RLC", "", F_A},
	0xcf: {"RLC", "", F_B},
	0xdf: {"RLC", "", F_Rn},

	0xbc: {"RR", "", F_A},
	0xcc: {"RR", "", F_B},
	0xdc: {"RR", "", F_Rn},

	0xbd: {"RRC", "", F_A},
	0xcd: {"RRC", "", F_B},
	0xdd: {"RRC", "", F_Rn},

	0x6b: {"SBB", "", F_B_A},
	0x1b: {"SBB", "", F_Rn_A},
	0x3b: {"SBB", "", F_Rn_B},
	0x4b: {"SBB", "", F_Rn_Rn},
	0x2b: {"SBB", "", F_iop_A},
	0x5b: {"SBB", "", F_iop_B},
	0x7b: {"SBB", "", F_iop_Rn},

	0x07: {"SETC", "", F_None},

	0x8b: {"STA", "", F_iop16},
	0xab: {"STA", "", F_iop16idxB},
	0x9b: {"STA", "", F_starRn},

	0x09: {"STSP", "(SP) -> (B)", F_None},

	0x6a: {"SUB", "", F_B_A},
	0x1a: {"SUB", "", F_Rn_A},
	0x3a: {"SUB", "", F_Rn_B},
	0x4a: {"SUB", "", F_Rn_Rn},
	0x2a: {"SUB", "", F_iop_A},
	0x5a: {"SUB", "", F_iop_B},
	0x7a: {"SUB", "", F_iop_Rn},

	0xb7: {"SWAP", "", F_A},
	0xc7: {"SWAP", "", F_B},
	0xd7: {"SWAP", "", F_Rn},

	0xe8: {"TRAP-0", "", F_None},
	0xe9: {"TRAP-1", "", F_None},
	0xea: {"TRAP-2", "", F_None},
	0xeb: {"TRAP-3", "", F_None},
	0xec: {"TRAP-4", "", F_None},
	0xed: {"TRAP-5", "", F_None},
	0xee: {"TRAP-6", "", F_None},
	0xef: {"TRAP-7", "", F_None},
	0xf0: {"TRAP-8", "", F_None},
	0xf1: {"TRAP-9", "", F_None},
	0xf2: {"TRAP-10", "", F_None},
	0xf3: {"TRAP-11", "", F_None},
	0xf4: {"TRAP-12", "", F_None},
	0xf5: {"TRAP-13", "", F_None},
	0xf6: {"TRAP-14", "", F_None},
	0xf7: {"TRAP-15", "", F_None},
	0xf8: {"TRAP-16", "", F_None},
	0xf9: {"TRAP-17", "", F_None},
	0xfa: {"TRAP-18", "", F_None},
	0xfb: {"TRAP-19", "", F_None},
	0xfc: {"TRAP-20", "", F_None},
	0xfd: {"TRAP-21", "", F_None},
	0xfe: {"TRAP-22", "", F_None},
	0xff: {"TRAP-23", "", F_None},

	0xb0: {"TSTA", "", F_None},

	0xc1: {"TSTB", "", F_None},

	0xb6: {"XCHB", "", F_A},
	0xd6: {"XCHB", "", F_Rn},

	0x65: {"XOR", "", F_B_A},
	0x15: {"XOR", "", F_Rn_A},
	0x35: {"XOR", "", F_Rn_B},
	0x45: {"XOR", "", F_Rn_Rn},
	0x25: {"XOR", "", F_iop_A},
	0x55: {"XOR", "", F_iop_B},
	0x75: {"XOR", "", F_iop_Rn},

	0x85: {"XORP", "", F_A_Pn},
	0x95: {"XORP", "", F_B_Pn},
	0xA5: {"XORP", "", F_iop_Pn},
}
