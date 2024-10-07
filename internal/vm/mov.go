package vm

import (
	"encoding/binary"
	"fishy/pkg/ast"
	"fishy/pkg/datatype"
	"fishy/pkg/log"
	"fishy/pkg/utils"
	"strconv"
)

// TODO: change registers to uint64?
// TODO: floats
// TODO: handle memory corruption (moving 4-byte number to 2-byte number)

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
		num, _ := strconv.ParseInt(v.Value, 10, 64)
		addr += int(num)
	case *ast.Register:
		num := m.getRegister(v.Value)
		addr += int(num)
	default:
		log.Fatal("unknown value to get address of", "value", value)
	}

	if dt, ok := m.symbolTable[uint32(addr)]; ok {
		switch dt {
		case datatype.U8, datatype.UNKNOWN:
			m.setRegister(reg, uint32(m.memory[addr]))
		case datatype.U16:
			num := binary.BigEndian.Uint16(m.memory[addr : addr+2])
			m.setRegister(reg, uint32(num))
		case datatype.U32:
			num := binary.BigEndian.Uint32(m.memory[addr : addr+4])
			m.setRegister(reg, num)
		}
	} else {
		m.setRegister(reg, uint32(m.memory[addr]))
	}
}

func (m *Machine) handleMovAofReg() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	value := m.decodeValue()
	addr := 0
	switch v := value.(type) {
	case *ast.NumberLiteral:
		num, _ := strconv.ParseInt(v.Value, 10, 64)
		addr += int(num)
	case *ast.Register:
		num := m.getRegister(v.Value)
		addr += int(num)
	default:
		log.Fatal("unknown value to get address of", "value", value)
	}

	pos := m.position()
	reg := m.decodeRegister(pos)
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	bytes := utils.Bytes4(m.getRegister(reg))
	for i := 0; i < 4; i++ {
		m.memory[addr+i] = bytes[i]
	}
}

func (m *Machine) handleMovAofLit() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	value := m.decodeValue()
	addr := 0
	switch v := value.(type) {
	case *ast.NumberLiteral:
		num, _ := strconv.ParseInt(v.Value, 10, 64)
		addr += int(num)
	case *ast.Register:
		num := m.getRegister(v.Value)
		addr += int(num)
	default:
		log.Fatal("unknown value to get address of", "value", value)
	}

	pos := m.position()
	lit := m.decodeNumber("u32", pos)
	m.incRegister(utils.RegisterToIndex("ip"), 4)

	bytes := utils.Bytes4(uint32(lit))
	for i := 0; i < 4; i++ {
		m.memory[addr+i] = bytes[i]
	}
}
