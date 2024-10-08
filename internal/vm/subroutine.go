package vm

import (
	"fishy/pkg/datatype"
	"fishy/pkg/utils"
)

func (m *Machine) handleCallLit() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	functionAddress := m.readLiteral(datatype.U64)
	bytes := utils.Bytes8(m.getRegister(utils.RegisterToIndex("ip")))
	m.stackPush(bytes)
	m.setRegister(utils.RegisterToIndex("ip"), functionAddress)
}

func (m *Machine) handleRet() {
	returnAddress := m.stackPop(datatype.U64)
	m.setRegister(utils.RegisterToIndex("ip"), returnAddress)
}
