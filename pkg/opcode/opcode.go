package opcode

import "fmt"

type Opcode int

const (
	NOP Opcode = iota
	HLT
	BRK
	SYSCALL

	MOV_REG_REG
	MOV_REG_LIT
	MOV_REG_ADR
	MOV_REG_AOF
	MOV_AOF_REG
	MOV_AOF_LIT

	ADD_REG_LIT
	ADD_REG_REG
	SUB_REG_LIT
	SUB_REG_REG
	MUL_REG_LIT
	MUL_REG_REG
	DIV_REG_LIT
	DIV_REG_REG
)

func (o Opcode) String() string {
	switch o {
	case NOP:
		return "NOP"
	case HLT:
		return "HLT"
	case BRK:
		return "BRK"
	case SYSCALL:
		return "SYSCALL"
	case MOV_REG_REG:
		return "MOV_REG_REG"
	case MOV_REG_LIT:
		return "MOV_REG_LIT"
	case MOV_REG_ADR:
		return "MOV_REG_ADR"
	case MOV_REG_AOF:
		return "MOV_REG_AOF"
	case MOV_AOF_REG:
		return "MOV_AOF_REG"
	case MOV_AOF_LIT:
		return "MOV_AOF_LIT"
	case ADD_REG_LIT:
		return "ADD_REG_LIT"
	case ADD_REG_REG:
		return "ADD_REG_REG"
	case SUB_REG_LIT:
		return "SUB_REG_LIT"
	case SUB_REG_REG:
		return "SUB_REG_REG"
	case MUL_REG_LIT:
		return "MUL_REG_LIT"
	case MUL_REG_REG:
		return "MUL_REG_REG"
	case DIV_REG_LIT:
		return "DIV_REG_LIT"
	case DIV_REG_REG:
		return "DIV_REG_REG"
	default:
		return fmt.Sprintf("0x%04X", int(o))
	}
}
