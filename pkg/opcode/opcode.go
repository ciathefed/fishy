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
