package intelhex

import (
        "bufio"
        "encoding/hex"
        "io"
        "strconv"
)

type LocatedByte struct {
        Value   byte
        Address uint
}

type IntelHexReader struct {
        rd      *bufio.Reader
        offset  uint
        address uint
        err     error
}

// NewReader returns a new Reader
func NewReader(rd io.Reader, offset uint) *IntelHexReader {
        return &IntelHexReader{
                rd:     bufio.NewReader(rd),
                offset: offset,
        }
}

const (
        dataRecord = iota
        endOfFileRecord
        extendedSegmentAddressRecord
        StartSegmentAddressRecord
        ExtendedLinearAddressRecord
        StartLinearAdddressRecord
)

type IteratorFunc func(val LocatedByte)

func (self *IntelHexReader) Iterate(fun IteratorFunc) error {

        for {
                line, err := self.rd.ReadString('\n')
                if err == io.EOF {
                        break
                }

                err = self.iterateLine(line, fun)
                if err != nil {
                        return err
                }
        }

        return nil
}

func (self *IntelHexReader) iterateLine(line string, fun IteratorFunc) error {
        if len(line) > 11 && line[0] == ':' {
                var byte_count, address, record_type int64
                var err error
                byte_count, err = strconv.ParseInt(line[1:3], 16, 8)
                if err != nil {
                        return err
                }

                address, err = strconv.ParseInt(line[3:7], 16, 16)
                if err != nil {
                        return err
                }

                record_type, err = strconv.ParseInt(line[7:9], 16, 8)
                if err != nil {
                        return err
                }

                hexpayload := line[9 : 9+byte_count*2]

                switch record_type {
                case dataRecord:
                        self.address = uint(address) + self.offset
                        payload, err := hex.DecodeString(hexpayload)
                        if err != nil {
                                return err
                        }

                        for i := 0; i < len(payload); i++ {
                                element := LocatedByte{
                                        Value:   payload[i],
                                        Address: self.address,
                                }
                                fun(element)
                                self.address++
                        }
                case endOfFileRecord:
                        break
                default:

                }
        }

        return nil
}
