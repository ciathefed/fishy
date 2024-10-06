package vm

import (
	"fishy/pkg/utils"
)

func (m *Machine) handleCallLit() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)
	functionAddress := m.readLiteral()

	m.stackPush(m.getRegister(utils.RegisterToIndex("ip")))
	m.setRegister(utils.RegisterToIndex("ip"), functionAddress)
}

func (m *Machine) handleRet() {
	returnAddress := m.stackPop()
	m.setRegister(utils.RegisterToIndex("ip"), returnAddress)
}
