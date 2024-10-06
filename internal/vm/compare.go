package vm

import (
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
		m.applyRegLitCompare(func(reg, lit uint32) Flag {
			if m.getRegister(int(reg)) == lit {
				return FLAG_EQ
			}
			if m.getRegister(int(reg)) < lit {
				return FLAG_LT
			}
			return FLAG_GT
		})
	case opcode.CMP_REG_REG:
		m.applyRegRegCompare(func(reg0, reg1 uint32) Flag {
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

func (m *Machine) applyRegLitCompare(operation func(uint32, uint32) Flag) {
	reg := m.readRegister()
	lit := m.readLiteral()
	result := operation(uint32(reg), lit)
	m.setRegister(utils.RegisterToIndex("cp"), uint32(result))
}

func (m *Machine) applyRegRegCompare(operation func(uint32, uint32) Flag) {
	reg0 := m.readRegister()
	reg1 := m.readRegister()
	result := operation(m.getRegister(reg0), m.getRegister(reg1))
	m.setRegister(utils.RegisterToIndex("cp"), uint32(result))
}
