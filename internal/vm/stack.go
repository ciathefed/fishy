package vm

import (
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
)

func (m *Machine) handlePush(op opcode.Opcode) {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	switch op {
	case opcode.PUSH_REG:
		pos := m.position()
		reg := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		m.stackPush(m.getRegister(reg))
	case opcode.PUSH_LIT:
		pos := m.position()
		lit := m.decodeNumber("u32", pos)
		m.incRegister(utils.RegisterToIndex("ip"), 4)
		m.stackPush(uint32(lit))
	}
}

func (m *Machine) handlePop(op opcode.Opcode) {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	switch op {
	case opcode.POP_REG:
		pos := m.position()
		reg := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		m.setRegister(reg, m.stackPop())
	}
}
