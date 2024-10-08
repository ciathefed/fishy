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
	ADD_REG_AOF
	SUB_REG_LIT
	SUB_REG_REG
	SUB_REG_AOF
	MUL_REG_LIT
	MUL_REG_REG
	MUL_REG_AOF
	DIV_REG_LIT
	DIV_REG_REG
	DIV_REG_AOF

	AND_REG_LIT
	AND_REG_REG
	OR_REG_LIT
	OR_REG_REG
	XOR_REG_LIT
	XOR_REG_REG
	SHL_REG_LIT
	SHL_REG_REG
	SHR_REG_LIT
	SHR_REG_REG

	CMP_REG_LIT
	CMP_REG_REG

	JMP_LIT
	JMP_REG
	JEQ_LIT
	JEQ_REG
	JNE_LIT
	JNE_REG
	JLT_LIT
	JLT_REG
	JGT_LIT
	JGT_REG
	JLE_LIT
	JLE_REG
	JGE_LIT
	JGE_REG

	PUSH_LIT
	PUSH_REG
	PUSH_AOF
	POP_REG
	POP_AOF

	CALL_LIT
	RET
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
	case AND_REG_LIT:
		return "AND_REG_LIT"
	case AND_REG_REG:
		return "AND_REG_REG"
	case OR_REG_LIT:
		return "OR_REG_LIT"
	case OR_REG_REG:
		return "OR_REG_REG"
	case XOR_REG_LIT:
		return "XOR_REG_LIT"
	case XOR_REG_REG:
		return "XOR_REG_REG"
	case SHL_REG_LIT:
		return "SHL_REG_LIT"
	case SHL_REG_REG:
		return "SHL_REG_REG"
	case SHR_REG_LIT:
		return "SHR_REG_LIT"
	case SHR_REG_REG:
		return "SHR_REG_REG"
	case CMP_REG_LIT:
		return "CMP_REG_LIT"
	case CMP_REG_REG:
		return "CMP_REG_REG"
	case JMP_LIT:
		return "JMP_LIT"
	case JMP_REG:
		return "JMP_REG"
	case JEQ_LIT:
		return "JEQ_LIT"
	case JEQ_REG:
		return "JEQ_REG"
	case JNE_LIT:
		return "JNE_LIT"
	case JNE_REG:
		return "JNE_REG"
	case JLT_LIT:
		return "JLT_LIT"
	case JLT_REG:
		return "JLT_REG"
	case JGT_LIT:
		return "JGT_LIT"
	case JGT_REG:
		return "JGT_REG"
	case JLE_LIT:
		return "JLE_LIT"
	case JLE_REG:
		return "JLE_REG"
	case JGE_LIT:
		return "JGE_LIT"
	case JGE_REG:
		return "JGE_REG"
	case PUSH_LIT:
		return "PUSH_LIT"
	case PUSH_REG:
		return "PUSH_REG"
	case PUSH_AOF:
		return "PUSH_AOF"
	case POP_REG:
		return "POP_REG"
	case POP_AOF:
		return "POP_AOF"
	case CALL_LIT:
		return "CALL_LIT"
	case RET:
		return "RET"
	default:
		return fmt.Sprintf("0x%04X", int(o))
	}
}
