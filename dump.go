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

	hexreader := intelhex.NewReader(file, 0x0000)
	locatedBytes := intelhex.NewQueue(64353)
	labels := tms7000.NewLabelTable()
	labels.Add(0xf559, 0, "START")
	labels.Add(0xe059, 0, "MAIN")

	labels.Add(0xecf0, 0, "INT5_HANDLER")
	labels.Add(0xefa6, 0, "INT4_HANDLER")
	labels.Add(0xf005, 0, "INT3_HANDLER")
	labels.Add(0xe8bf, 0, "INT2_HANDLER")
	labels.Add(0xe2f9, 0, "INT1_HANDLER")

	labels.Add(0xffd8, 2, "UNKNOWN_VECTOR")
	labels.Add(0xffda, 2, "UNKNOWN_VECTOR")
	labels.Add(0xffdc, 2, "UNKNOWN_VECTOR")
	labels.Add(0xffde, 2, "UNKNOWN_VECTOR")
	labels.Add(0xffe0, 2, "UNKNOWN_VECTOR")
	labels.Add(0xffe2, 2, "UNKNOWN_VECTOR")
	labels.Add(0xffe4, 2, "UNKNOWN_VECTOR")
	labels.Add(0xffe6, 2, "UNKNOWN_VECTOR")
	labels.Add(0xffe8, 2, "UNKNOWN_VECTOR")
	labels.Add(0xffea, 2, "UNKNOWN_VECTOR")
	labels.Add(0xffec, 2, "UNKNOWN_VECTOR")
	labels.Add(0xffee, 2, "UNKNOWN_VECTOR")
	labels.Add(0xfff0, 2, "UNKNOWN_VECTOR")
	labels.Add(0xfff2, 2, "UNKNOWN_VECTOR")
	labels.Add(0xfff4, 2, "INT5_VECTOR")
	labels.Add(0xfff6, 2, "INT4_VECTOR")
	labels.Add(0xfff8, 2, "INT3_VECTOR")
	labels.Add(0xfffa, 2, "INT2_VECTOR")
	labels.Add(0xfffc, 2, "INT1_VECTOR")
	labels.Add(0xfffe, 2, "RESET_VECTOR")

	disassembler := tms7000.NewDisassembler(tms7000.TMS7000InstructionSet, locatedBytes, labels)

	err = hexreader.Iterate(func(val intelhex.LocatedByte) {
		locatedBytes.Push(&val)
	})

	disassembler.Do()

	if err != nil {
		log.Fatal(err)
	}
}
