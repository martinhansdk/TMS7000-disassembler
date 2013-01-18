package tms7000

import (
        "fmt"
)

type AddressingMode struct {
        args      uint
        argformat func([]byte) string
}

func noArgs(output string) func([]byte) string {
        return func([]byte) string {
                return output
        }
}

func oneArg(format string) func([]byte) string {
        return func(args []byte) string {
                return fmt.Sprintf(format, args[0])
        }
}

func twoArgs(format string) func([]byte) string {
        return func(args []byte) string {
                return fmt.Sprintf(format, args[0], args[1])
        }
}

func label16bit(format string) func([]byte) string {
        return func(args []byte) string {
                return fmt.Sprintf(format, uint16(args[0])<<8+uint16(args[1]))
        }
}

func iop16bitOneReg(format string) func([]byte) string {
        return func(args []byte) string {
                bigval := uint16(args[0])<<8 + uint16(args[1])
                return fmt.Sprintf(format, bigval, args[2])
        }
}

var AMOther = AddressingMode{0, noArgs("")}

var AMSingleRegisterA = AddressingMode{0, noArgs("A")}
var AMSingleRegisterB = AddressingMode{0, noArgs("B")}
var AMSingleRegisterRn = AddressingMode{1, oneArg("R%d")}
var AMASingleRegisterRn = AddressingMode{1, oneArg("R%d,A")}
var AMBSingleRegisterRn = AddressingMode{1, oneArg("R%d,B")}
var AMRnSingleRegisterA = AddressingMode{1, oneArg("A,R%d")}
var AMRnSingleRegisterB = AddressingMode{1, oneArg("B,R%d")}

var AMImmediateOpA = AddressingMode{1, oneArg("%%%d,A")}
var AMImmediateOpB = AddressingMode{1, oneArg("%%%d,B")}
var AMImmediateOpRn = AddressingMode{2, twoArgs("%%%d,R%d")}
var AM16bitImmediateOpRegister = AddressingMode{3, iop16bitOneReg("%%>%04x,R%d")}

var AMABRegister = AddressingMode{1, noArgs("A,B")}
var AMBARegister = AddressingMode{1, noArgs("B,A")}

var AMDualRegister = AddressingMode{2, twoArgs("R%d,R%d")}
var AMAPeripheral = AddressingMode{1, oneArg("A,P%d")}
var AMBPeripheral = AddressingMode{1, oneArg("B,P%d")}
var AMPeripheralA = AddressingMode{1, oneArg("P%d,A")}
var AMPeripheralB = AddressingMode{1, oneArg("P%d,B")}
var AMPeripheralImmediate = AddressingMode{2, twoArgs("%%%d,P%d")}

var AMDirectLabel = AddressingMode{2, label16bit("@>%04x")}
var AMDirectLabelIndexedB = AddressingMode{2, label16bit("@>%04x(B)")}
var AMDirectLabelIndexedBReg = AddressingMode{3, iop16bitOneReg("%>%04x(B),R%d")}

var AMRegisterFileIndirect AddressingMode = AddressingMode{1, oneArg("*R%d")}

/*
var AMProgramCounterRelative AddressingMode = AddressingMode{"PC relative", "%%>%s,%s,%s"}
var AMDirectMemory AddressingMode = AddressingMode{"direct memory", "@%s"}
var AMIndexed AddressingMode = AddressingMode{"indexed", "@%s(%s)"}
*/

type Instruction struct {
        name    string
        comment string
        amode   AddressingMode
}

func (self *Instruction) String(args []byte) string {
        return fmt.Sprintf("%s %s", self.name, self.amode.argformat(args))
}

func (self *Instruction) Args() uint {
        return self.amode.args
}

type InstructionSet map[byte]Instruction

var TMS7000InstructionSet = InstructionSet{

        0x69:   {"ADC", "(A) + (B) + (C) -> (B)", AMABRegister},
        0x19:   {"ADC", "(Rs) + (A) + (C) -> (A)", AMASingleRegisterRn},
        0x39:   {"ADC", "(Rs) + (B) + (C) -> (B)", AMBSingleRegisterRn},
        0x49:   {"ADC", "(Rs) + (Rd) + (C) -> (Rd)", AMDualRegister},
        0x29:   {"ADC", "%iop + (A) + (C) -> (A)", AMImmediateOpA},
        0x59:   {"ADC", "%iop + (B) + (C) -> (B)", AMImmediateOpB},
        0x79:   {"ADC", "%iop + (Rd) + (C) -> (Rd)", AMImmediateOpRn},

        0x68:   {"ADD", "(A) + (B) -> (B)", AMABRegister},
        0x18:   {"ADD", "(Rs) + (A) -> (A)", AMASingleRegisterRn},
        0x38:   {"ADD", "(Rs) + (B) -> (B)", AMBSingleRegisterRn},
        0x48:   {"ADD", "(Rs) + (Rd) -> (Rd)", AMDualRegister},
        0x28:   {"ADD", "%iop + (A) -> (A)", AMImmediateOpA},
        0x58:   {"ADD", "%iop + (B) -> (B)", AMImmediateOpB},
        0x78:   {"ADD", "%iop + (Rd) -> (Rd)", AMImmediateOpRn},

        0x63:   {"AND", "(A) AND (B) -> (B)", AMABRegister},
        0x13:   {"AND", "(Rs) AND (A) -> (A)", AMASingleRegisterRn},
        0x33:   {"AND", "(Rs) AND (B) -> (B)", AMBSingleRegisterRn},
        0x43:   {"AND", "(Rs) AND (Rd) -> (Rd)", AMDualRegister},
        0x23:   {"AND", "%iop AND (A) -> (A)", AMImmediateOpA},
        0x53:   {"AND", "%iop AND (B) -> (B)", AMImmediateOpB},
        0x73:   {"AND", "%iop AND (Rd) -> (Rd)", AMImmediateOpRn},

        0x83:   {"ANDP", "(A) AND (Pn) -> (Pn)", AMAPeripheral},
        0x93:   {"ANDP", "(B) AND (Pn) -> (Pn)", AMBPeripheral},
        0xA3:   {"ANDP", "%iop AND (Pn) -> (Pn)", AMPeripheralImmediate},

        0x8c:   {"BR", "@LABEL -> PC", AMDirectLabel},
        0xac:   {"BR", "@LABEL(B) -> PC", AMDirectLabelIndexedB},
        0x9c:   {"BR", "*Rn -> PC", AMRegisterFileIndirect},

        0x8e:   {"CALL", "PUSH PC, @LABEL -> PC", AMDirectLabel},
        0xae:   {"CALL", "PUSH PC, @LABEL(B) -> PC", AMDirectLabelIndexedB},
        0x9e:   {"CALL", "PUSH PC, *Rn -> PC", AMRegisterFileIndirect},

        0xb5:   {"CLR", "0 -> (A)", AMSingleRegisterA},
        0xc6:   {"CLR", "0 -> (B)", AMSingleRegisterB},
        0xd5:   {"CLR", "0 -> (Rd)", AMSingleRegisterRn},

        0xb2:   {"DEC", "", AMSingleRegisterA},
        0xc2:   {"DEC", "", AMSingleRegisterB},
        0xd2:   {"DEC", "", AMSingleRegisterRn},

        0x06:   {"DINT", "Clear global interrupt enable bit", AMOther},

        0x05:   {"EINT", "Set global interrupt enable bit", AMOther},

        0xb3:   {"INC", "", AMSingleRegisterA},
        0xc3:   {"INC", "", AMSingleRegisterB},
        0xd3:   {"INC", "", AMSingleRegisterRn},

        0x0d:   {"LDSP", "(B) -> (SP) Load SP with register B's contents", AMOther},

        0x00:   {"NOP", "", AMOther},

        0xc0:   {"MOV", "", AMABRegister},
        0xd0:   {"MOV", "", AMASingleRegisterRn},
        0x62:   {"MOV", "", AMBARegister},
        0xd1:   {"MOV", "", AMBSingleRegisterRn},
        0x12:   {"MOV", "", AMRnSingleRegisterA},
        0x32:   {"MOV", "", AMRnSingleRegisterB},
        0x42:   {"MOV", "", AMDualRegister},
        0x22:   {"MOV", "", AMImmediateOpA},
        0x52:   {"MOV", "", AMImmediateOpB},
        0x72:   {"MOV", "", AMImmediateOpRn},

        0x88:   {"MOVD", "", AM16bitImmediateOpRegister},
        0xa8:   {"MOVD", "", AMDirectLabelIndexedBReg},
        0x98:   {"MOVD", "", AMDualRegister},

        0x82:   {"MOVP", "", AMAPeripheral},
        0x92:   {"MOVP", "", AMBPeripheral},
        0xA2:   {"MOVP", "", AMPeripheralImmediate},
        0x80:   {"MOVP", "", AMPeripheralA},
        0x91:   {"MOVP", "", AMPeripheralB},

        0x0b:   {"RETI", "return from interrupt", AMOther},

        0x0a:   {"RETS", "return from subroutine", AMOther},

        0x8b:   {"STA", "(A) -> @LABEL", AMDirectLabel},
        0xab:   {"STA", "(A) -> @LABEL(B)", AMDirectLabelIndexedB},
        0x9b:   {"STA", "(A) -> *Rn", AMRegisterFileIndirect},

        0x09:   {"STSP", "(SP) -> (B) Copy the SP into Register B", AMOther},
}
