package vm

import (
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
)

func (m *Machine) handleArithmetic(op opcode.Opcode) {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	switch op {
	case opcode.ADD_REG_LIT:
		pos := m.position()
		reg := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		pos = m.position()
		lit := m.decodeNumber("u32", pos)
		m.incRegister(utils.RegisterToIndex("ip"), 4)
		temp := m.getRegister(reg)
		m.setRegister(reg, temp+uint32(lit))
	case opcode.ADD_REG_REG:
		pos := m.position()
		reg0 := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		pos = m.position()
		reg1 := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		temp0 := m.getRegister(reg0)
		temp1 := m.getRegister(reg1)
		m.setRegister(reg0, temp0+temp1)
	case opcode.SUB_REG_LIT:
		pos := m.position()
		reg := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		pos = m.position()
		lit := m.decodeNumber("u32", pos)
		m.incRegister(utils.RegisterToIndex("ip"), 4)
		temp := m.getRegister(reg)
		m.setRegister(reg, temp-uint32(lit))
	case opcode.SUB_REG_REG:
		pos := m.position()
		reg0 := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		pos = m.position()
		reg1 := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		temp0 := m.getRegister(reg0)
		temp1 := m.getRegister(reg1)
		m.setRegister(reg0, temp0-temp1)
	case opcode.MUL_REG_LIT:
		pos := m.position()
		reg := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		pos = m.position()
		lit := m.decodeNumber("u32", pos)
		m.incRegister(utils.RegisterToIndex("ip"), 4)
		temp := m.getRegister(reg)
		m.setRegister(reg, temp*uint32(lit))
	case opcode.MUL_REG_REG:
		pos := m.position()
		reg0 := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		pos = m.position()
		reg1 := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		temp0 := m.getRegister(reg0)
		temp1 := m.getRegister(reg1)
		m.setRegister(reg0, temp0*temp1)
	case opcode.DIV_REG_LIT:
		pos := m.position()
		reg := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		pos = m.position()
		lit := m.decodeNumber("u32", pos)
		m.incRegister(utils.RegisterToIndex("ip"), 4)
		temp := m.getRegister(reg)
		m.setRegister(reg, temp/uint32(lit))
	case opcode.DIV_REG_REG:
		pos := m.position()
		reg0 := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		pos = m.position()
		reg1 := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		temp0 := m.getRegister(reg0)
		temp1 := m.getRegister(reg1)
		m.setRegister(reg0, temp0/temp1)

	}
}
