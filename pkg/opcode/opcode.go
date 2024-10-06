package opcode

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
	default:
		return "UNKNOWN"
	}
}
