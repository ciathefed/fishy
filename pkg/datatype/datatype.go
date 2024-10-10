package datatype

import (
	"encoding/binary"
	"fishy/pkg/utils"
	"log"
)

type DataType int

const (
	BYTE DataType = iota
	WORD
	DWORD
	QWORD
	UNSET = 0xff
)

func (d DataType) String() string {
	switch d {
	case BYTE:
		return "BYTE"
	case WORD:
		return "WORD"
	case DWORD:
		return "DWORD"
	case QWORD:
		return "QWORD"
	case UNSET:
		return "UNSET"
	default:
		log.Fatal("invalid data type", "type", int(d))
		return ""
	}
}

func (d DataType) Size() int {
	switch d {
	case BYTE:
		return 1
	case WORD:
		return 2
	case DWORD:
		return 4
	case QWORD:
		return 8
	case UNSET:
		return 8
	default:
		log.Fatal("unknown data type", "type", int(d))
		return -1
	}
}

func (d DataType) MakeBytes(value uint64) []byte {
	switch d {
	case BYTE:
		return []byte{byte(value)}
	case WORD:
		return utils.Bytes2(uint16(value))
	case DWORD:
		return utils.Bytes4(uint32(value))
	case QWORD:
		return utils.Bytes8(value)
	case UNSET:
		return utils.Bytes8(value)
	default:
		log.Fatal("invalid data type", "type", int(d))
		return nil
	}
}

func (d DataType) ReadBytes(byteArray []byte, index int) uint64 {
	switch d {
	case BYTE:
		return uint64(byteArray[index])
	case WORD:
		return uint64(binary.BigEndian.Uint16(byteArray[index : index+d.Size()]))
	case DWORD:
		return uint64(binary.BigEndian.Uint32(byteArray[index : index+d.Size()]))
	case QWORD:
		return uint64(binary.BigEndian.Uint32(byteArray[index : index+d.Size()]))
	case UNSET:
		return uint64(binary.BigEndian.Uint32(byteArray[index : index+d.Size()]))
	default:
		log.Fatal("invalid data type", "type", int(d))
		return 0
	}
}

func FromString(ident string) DataType {
	switch ident {
	case "byte":
		return BYTE
	case "word":
		return WORD
	case "dword":
		return DWORD
	case "qword":
		return QWORD
	default:
		log.Fatal("unknown data type", "type", ident)
		return UNSET
	}
}
