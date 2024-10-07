package vm

import (
	"fishy/pkg/datatype"
	"fishy/pkg/utils"
)

func (m *Machine) handleCallLit() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)
	functionAddress := m.readLiteral(datatype.U64)

	m.stackPush(m.getRegister(utils.RegisterToIndex("ip")))
	m.setRegister(utils.RegisterToIndex("ip"), functionAddress)
}

func (m *Machine) handleRet() {
	returnAddress := m.stackPop()
	m.setRegister(utils.RegisterToIndex("ip"), returnAddress)
}
