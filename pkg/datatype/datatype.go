package datatype

import "fmt"

type DataType int

const (
	U8 DataType = iota
	U16
	U32
	// U64
	// I8
	// I16
	// I32
	// I64
	// F32
	// F64
	UNKNOWN = 0xff
)

func (d DataType) String() string {
	switch d {
	case U8:
		return "U8"
	case U16:
		return "U16"
	case U32:
		return "U32"
	// case U64:
	// 	return "U64"
	// case I8:
	// 	return "I8"
	// case I16:
	// 	return "I16"
	// case I32:
	// 	return "I32"
	// case I64:
	// 	return "I64"
	// case F32:
	// 	return "F32"
	// case F64:
	// 	return "F64"
	case UNKNOWN:
		return "UNKNOWN"
	default:
		return fmt.Sprintf("INVALID: %d", int(d))
	}
}
