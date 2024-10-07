package datatype

import (
	"fishy/pkg/utils"
	"log"
)

type DataType int

const (
	U8 DataType = iota
	U16
	U32
	U64
	// I8
	// I16
	// I32
	// I64
	// F32
	// F64
	UNSET = 0xff
)

func (d DataType) String() string {
	switch d {
	case U8:
		return "U8"
	case U16:
		return "U16"
	case U32:
		return "U32"
	case U64:
		return "U64"
	case UNSET:
		return "UNSET"
	default:
		log.Fatal("invalid data type", "type", int(d))
		return ""
	}
}

func (d DataType) Size() int {
	switch d {
	case U8:
		return 1
	case U16:
		return 2
	case U32:
		return 4
	case U64:
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
	case U8:
		return []byte{byte(value)}
	case U16:
		return utils.Bytes2(uint16(value))
	case U32:
		return utils.Bytes4(uint32(value))
	case U64:
		return utils.Bytes8(value)
	case UNSET:
		return utils.Bytes8(value)
	default:
		log.Fatal("invalid data type", "type", int(d))
		return nil
	}
}

func FromString(ident string) DataType {
	switch ident {
	case "u8":
		return U8
	case "u16":
		return U16
	case "u32":
		return U32
	case "u64":
		return U64
	default:
		log.Fatal("unknown data type", "type", ident)
		return UNSET
	}
}
