package vm

import (
	"fishy/pkg/datatype"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
)

func (m *Machine) handleBitwise(op opcode.Opcode) {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	pos := m.position()
	rdt := datatype.DataType(m.decodeNumber("byte", pos))
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	switch op {
	case opcode.AND_REG_LIT:
		m.applyRegLitBitwise(rdt, func(reg, lit uint64) uint64 { return reg & lit })
	case opcode.AND_REG_REG:
		m.applyRegRegBitwise(func(reg0, reg1 uint64) uint64 { return reg0 & reg1 })
	case opcode.OR_REG_LIT:
		m.applyRegLitBitwise(rdt, func(reg, lit uint64) uint64 { return reg | lit })
	case opcode.OR_REG_REG:
		m.applyRegRegBitwise(func(reg0, reg1 uint64) uint64 { return reg0 | reg1 })
	case opcode.XOR_REG_LIT:
		m.applyRegLitBitwise(rdt, func(reg, lit uint64) uint64 { return reg ^ lit })
	case opcode.XOR_REG_REG:
		m.applyRegRegBitwise(func(reg0, reg1 uint64) uint64 { return reg0 ^ reg1 })
	case opcode.SHL_REG_LIT:
		m.applyRegLitBitwise(rdt, func(reg, lit uint64) uint64 { return reg << lit })
	case opcode.SHL_REG_REG:
		m.applyRegRegBitwise(func(reg0, reg1 uint64) uint64 { return reg0 << reg1 })
	case opcode.SHR_REG_LIT:
		m.applyRegLitBitwise(rdt, func(reg, lit uint64) uint64 { return reg >> lit })
	case opcode.SHR_REG_REG:
		m.applyRegRegBitwise(func(reg0, reg1 uint64) uint64 { return reg0 >> reg1 })
	}
}

func (m *Machine) applyRegLitBitwise(dataType datatype.DataType, operation func(uint64, uint64) uint64) {
	reg := m.readRegister()
	lit := m.readLiteral(dataType)
	temp := m.getRegister(reg)
	m.setRegister(reg, operation(temp, lit))
}

func (m *Machine) applyRegRegBitwise(operation func(uint64, uint64) uint64) {
	reg0 := m.readRegister()
	reg1 := m.readRegister()
	temp0 := m.getRegister(reg0)
	temp1 := m.getRegister(reg1)
	m.setRegister(reg0, operation(temp0, temp1))
}
