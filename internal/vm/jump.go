package vm

import (
	"fishy/pkg/datatype"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
)

func (m *Machine) handleJump(thread *Thread, op opcode.Opcode) {
	m.incRegister(thread, utils.RegisterToIndex("ip"), 2)

	switch op {
	case opcode.JMP_LIT:
		target := m.readLiteral(thread, datatype.QWORD)
		m.setRegister(thread, utils.RegisterToIndex("ip"), target)
	case opcode.JMP_REG:
		targetReg := m.readRegister(thread)
		target := m.getRegister(thread, targetReg)
		m.setRegister(thread, utils.RegisterToIndex("ip"), target)
	case opcode.JEQ_LIT:
		target := m.readLiteral(thread, datatype.QWORD)
		if m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_EQ) {
			m.setRegister(thread, utils.RegisterToIndex("ip"), target)
		}
	case opcode.JEQ_REG:
		targetReg := m.readRegister(thread)
		target := m.getRegister(thread, targetReg)
		if m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_EQ) {
			m.setRegister(thread, utils.RegisterToIndex("ip"), target)
		}
	case opcode.JNE_LIT:
		target := m.readLiteral(thread, datatype.QWORD)
		if m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_LT) || m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_GT) {
			m.setRegister(thread, utils.RegisterToIndex("ip"), target)
		}
	case opcode.JNE_REG:
		targetReg := m.readRegister(thread)
		target := m.getRegister(thread, targetReg)
		if m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_LT) || m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_GT) {
			m.setRegister(thread, utils.RegisterToIndex("ip"), target)
		}
	case opcode.JLT_LIT:
		target := m.readLiteral(thread, datatype.QWORD)
		if m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_LT) {
			m.setRegister(thread, utils.RegisterToIndex("ip"), target)
		}
	case opcode.JLT_REG:
		targetReg := m.readRegister(thread)
		target := m.getRegister(thread, targetReg)
		if m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_LT) {
			m.setRegister(thread, utils.RegisterToIndex("ip"), target)
		}
	case opcode.JGT_LIT:
		target := m.readLiteral(thread, datatype.QWORD)
		if m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_GT) {
			m.setRegister(thread, utils.RegisterToIndex("ip"), target)
		}
	case opcode.JGT_REG:
		targetReg := m.readRegister(thread)
		target := m.getRegister(thread, targetReg)
		if m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_GT) {
			m.setRegister(thread, utils.RegisterToIndex("ip"), target)
		}
	case opcode.JLE_LIT:
		target := m.readLiteral(thread, datatype.QWORD)
		if m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_EQ) || m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_LT) {
			m.setRegister(thread, utils.RegisterToIndex("ip"), target)
		}
	case opcode.JLE_REG:
		targetReg := m.readRegister(thread)
		target := m.getRegister(thread, targetReg)
		if m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_EQ) || m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_LT) {
			m.setRegister(thread, utils.RegisterToIndex("ip"), target)
		}
	case opcode.JGE_LIT:
		target := m.readLiteral(thread, datatype.QWORD)
		if m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_EQ) || m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_GT) {
			m.setRegister(thread, utils.RegisterToIndex("ip"), target)
		}
	case opcode.JGE_REG:
		targetReg := m.readRegister(thread)
		target := m.getRegister(thread, targetReg)
		if m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_EQ) || m.getRegister(thread, utils.RegisterToIndex("cp")) == uint64(FLAG_GT) {
			m.setRegister(thread, utils.RegisterToIndex("ip"), target)
		}
	}
}
