package vm

import (
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
)

func (m *Machine) handleJump(op opcode.Opcode) {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	switch op {
	case opcode.JMP_LIT:
		target := m.readLiteral()
		m.setRegister(utils.RegisterToIndex("ip"), target)
	case opcode.JMP_REG:
		targetReg := m.readRegister()
		target := m.getRegister(targetReg)
		m.setRegister(utils.RegisterToIndex("ip"), target)
	case opcode.JEQ_LIT:
		target := m.readLiteral()
		if m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_EQ) {
			m.setRegister(utils.RegisterToIndex("ip"), target)
		}
	case opcode.JEQ_REG:
		targetReg := m.readRegister()
		target := m.getRegister(targetReg)
		if m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_EQ) {
			m.setRegister(utils.RegisterToIndex("ip"), target)
		}
	case opcode.JNE_LIT:
		target := m.readLiteral()
		if m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_LT) || m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_GT) {
			m.setRegister(utils.RegisterToIndex("ip"), target)
		}
	case opcode.JNE_REG:
		targetReg := m.readRegister()
		target := m.getRegister(targetReg)
		if m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_LT) || m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_GT) {
			m.setRegister(utils.RegisterToIndex("ip"), target)
		}
	case opcode.JLT_LIT:
		target := m.readLiteral()
		if m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_LT) {
			m.setRegister(utils.RegisterToIndex("ip"), target)
		}
	case opcode.JLT_REG:
		targetReg := m.readRegister()
		target := m.getRegister(targetReg)
		if m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_LT) {
			m.setRegister(utils.RegisterToIndex("ip"), target)
		}
	case opcode.JGT_LIT:
		target := m.readLiteral()
		if m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_GT) {
			m.setRegister(utils.RegisterToIndex("ip"), target)
		}
	case opcode.JGT_REG:
		targetReg := m.readRegister()
		target := m.getRegister(targetReg)
		if m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_GT) {
			m.setRegister(utils.RegisterToIndex("ip"), target)
		}
	case opcode.JLE_LIT:
		target := m.readLiteral()
		if m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_EQ) || m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_LT) {
			m.setRegister(utils.RegisterToIndex("ip"), target)
		}
	case opcode.JLE_REG:
		targetReg := m.readRegister()
		target := m.getRegister(targetReg)
		if m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_EQ) || m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_LT) {
			m.setRegister(utils.RegisterToIndex("ip"), target)
		}
	case opcode.JGE_LIT:
		target := m.readLiteral()
		if m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_EQ) || m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_GT) {
			m.setRegister(utils.RegisterToIndex("ip"), target)
		}
	case opcode.JGE_REG:
		targetReg := m.readRegister()
		target := m.getRegister(targetReg)
		if m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_EQ) || m.getRegister(utils.RegisterToIndex("cp")) == uint32(FLAG_GT) {
			m.setRegister(utils.RegisterToIndex("ip"), target)
		}
	}
}
