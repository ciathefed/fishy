package vm

import (
	"fishy/pkg/datatype"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
)

type Flag int

const (
	FLAG_EQ Flag = iota
	FLAG_LT
	FLAG_GT
)

func (m *Machine) handleCompare(op opcode.Opcode) {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	switch op {
	case opcode.CMP_REG_LIT:
		m.applyRegLitCompare(func(reg, lit uint64) Flag {
			if reg == lit {
				return FLAG_EQ
			}
			if reg < lit {
				return FLAG_LT
			}
			return FLAG_GT
		})
	case opcode.CMP_REG_REG:
		m.applyRegRegCompare(func(reg0, reg1 uint64) Flag {
			if reg0 == reg1 {
				return FLAG_EQ
			}
			if reg0 < reg1 {
				return FLAG_LT
			}
			return FLAG_GT
		})
	}
}

func (m *Machine) applyRegLitCompare(operation func(uint64, uint64) Flag) {
	reg := m.readRegister()
	lit := m.readLiteral(datatype.QWORD)
	result := operation(m.getRegister(reg), lit)
	m.setRegister(utils.RegisterToIndex("cp"), uint64(result))
}

func (m *Machine) applyRegRegCompare(operation func(uint64, uint64) Flag) {
	reg0 := m.readRegister()
	reg1 := m.readRegister()
	result := operation(m.getRegister(reg0), m.getRegister(reg1))
	m.setRegister(utils.RegisterToIndex("cp"), uint64(result))
}
