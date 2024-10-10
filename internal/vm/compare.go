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

func (m *Machine) handleCompare(thread *Thread, op opcode.Opcode) {
	m.incRegister(thread, utils.RegisterToIndex("ip"), 2)

	switch op {
	case opcode.CMP_REG_LIT:
		m.applyRegLitCompare(thread, func(reg, lit uint64) Flag {
			if reg == lit {
				return FLAG_EQ
			}
			if reg < lit {
				return FLAG_LT
			}
			return FLAG_GT
		})
	case opcode.CMP_REG_REG:
		m.applyRegRegCompare(thread, func(reg0, reg1 uint64) Flag {
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

func (m *Machine) applyRegLitCompare(thread *Thread, operation func(uint64, uint64) Flag) {
	reg := m.readRegister(thread)
	lit := m.readLiteral(thread, datatype.QWORD)
	result := operation(m.getRegister(thread, reg), lit)
	m.setRegister(thread, utils.RegisterToIndex("cp"), uint64(result))
}

func (m *Machine) applyRegRegCompare(thread *Thread, operation func(uint64, uint64) Flag) {
	reg0 := m.readRegister(thread)
	reg1 := m.readRegister(thread)
	result := operation(m.getRegister(thread, reg0), m.getRegister(thread, reg1))
	m.setRegister(thread, utils.RegisterToIndex("cp"), uint64(result))
}
