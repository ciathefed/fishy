package vm

import (
	"fishy/pkg/datatype"
	"fishy/pkg/utils"
)

func (m *Machine) handleCallLit(thread *Thread) {
	m.incRegister(thread, utils.RegisterToIndex("ip"), 2)

	functionAddress := m.readLiteral(thread, datatype.QWORD)
	bytes := utils.Bytes8(m.getRegister(thread, utils.RegisterToIndex("ip")))
	m.stackPush(thread, bytes)
	m.setRegister(thread, utils.RegisterToIndex("ip"), functionAddress)
}

func (m *Machine) handleRet(thread *Thread) {
	returnAddress := m.stackPop(thread, datatype.QWORD)
	m.setRegister(thread, utils.RegisterToIndex("ip"), returnAddress)
}
