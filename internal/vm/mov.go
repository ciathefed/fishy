package vm

import (
	"fishy/pkg/ast"
	"fishy/pkg/utils"
	"fmt"
	"strconv"
)

func (m *Machine) handleMovRegReg() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	pos := m.position()
	reg0 := m.decodeRegister(pos)
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	pos = m.position()
	reg1 := m.decodeRegister(pos)
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	temp := m.getRegister(reg1)
	m.setRegister(reg0, temp)
}

func (m *Machine) handleMovRegLit() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	pos := m.position()
	reg := m.decodeRegister(pos)
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	pos = m.position()
	lit := m.decodeNumber("u32", pos)
	m.incRegister(utils.RegisterToIndex("ip"), 4)

	m.setRegister(reg, uint32(lit))
}

func (m *Machine) handleMovRegAdr() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	pos := m.position()
	reg := m.decodeRegister(pos)
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	pos = m.position()
	addr := m.decodeNumber("u32", pos)
	m.incRegister(utils.RegisterToIndex("ip"), 4)

	m.setRegister(reg, uint32(addr))
}

func (m *Machine) handleMovRegAof() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	pos := m.position()
	reg := m.decodeRegister(pos)
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	value := m.decodeValue()
	addr := 0
	switch v := value.(type) {
	case *ast.NumberLiteral:
		num, _ := strconv.ParseInt(v.Value, 10, 32)
		addr += int(num)
	case *ast.Register:
		num := m.getRegister(v.Value)
		addr += int(num)
	default:
		panic(fmt.Sprintf("unknown value to get address of: %#v", value))
	}

	m.setRegister(reg, uint32(m.memory[addr]))
}

func (m *Machine) handleMovAofReg() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	value := m.decodeValue()
	addr := 0
	switch v := value.(type) {
	case *ast.NumberLiteral:
		num, _ := strconv.ParseInt(v.Value, 10, 32)
		addr += int(num)
	case *ast.Register:
		num := m.getRegister(v.Value)
		addr += int(num)
	default:
		panic(fmt.Sprintf("unknown value to get address of: %#v", value))
	}

	pos := m.position()
	reg := m.decodeRegister(pos)
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	m.memory[addr] = byte(m.getRegister(reg))
}

func (m *Machine) handleMovAofLit() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	value := m.decodeValue()
	addr := 0
	switch v := value.(type) {
	case *ast.NumberLiteral:
		num, _ := strconv.ParseInt(v.Value, 10, 32)
		addr += int(num)
	case *ast.Register:
		num := m.getRegister(v.Value)
		addr += int(num)
	default:
		panic(fmt.Sprintf("unknown value to get address of: %#v", value))
	}

	pos := m.position()
	lit := m.decodeNumber("u32", pos)
	m.incRegister(utils.RegisterToIndex("ip"), 4)

	m.memory[addr] = byte(lit)
}
