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

	m.incRegister(utils.RegisterToIndex("ip"), 1)

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
	dt := datatype.DataType(m.decodeNumber("u8", pos))
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	pos = m.position()
	reg := m.decodeRegister(pos)
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	pos = m.position()
	lit := m.decodeNumber(dt.String(), pos)
	m.incRegister(utils.RegisterToIndex("ip"), uint64(dt.Size()))

	m.setRegister(reg, uint64(lit))
}

func (m *Machine) handleMovRegAdr() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	pos := m.position()
	dt := datatype.DataType(m.decodeNumber("u8", pos))
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	pos = m.position()
	reg := m.decodeRegister(pos)
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	pos = m.position()
	addr := m.decodeNumber(dt.String(), pos)
	m.incRegister(utils.RegisterToIndex("ip"), uint64(dt.Size()))

	m.setRegister(reg, uint64(addr))
}

func (m *Machine) handleMovRegAof() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	pos := m.position()
	rdt := datatype.DataType(m.decodeNumber("u8", pos))
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	pos = m.position()
	reg := m.decodeRegister(pos)
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	value := m.decodeValue(rdt)
	addr := 0
	switch v := value.(type) {
	case *ast.NumberLiteral:
		num, _ := strconv.ParseUint(v.Value, 10, 64)
		addr += int(num)
	case *ast.Register:
		num := m.getRegister(v.Value)
		addr += int(num)
	case *ast.RegisterOffset:
		num := m.getRegister(v.Left.Value)
		addr += int(num)
	case *ast.LabelOffset:
		num, _ := strconv.ParseUint(v.Left.(*ast.NumberLiteral).Value, 10, 64)
		addr += int(num)
	default:
		log.Fatal("unknown value to get address of", "value", value)
	}

	dt := datatype.DataType(datatype.UNSET)
	if adt, ok := m.symbolTable[uint64(addr)]; ok {
		dt = adt
	}
	if rdt != datatype.UNSET {
		dt = rdt
	}

	switch v := value.(type) {
	case *ast.RegisterOffset:
		addr = applyOffset(addr, v.Operator, v.Right.Value)
	case *ast.LabelOffset:
		addr = applyOffset(addr, v.Operator, v.Right.Value)
	}

	switch dt {
	case datatype.U8:
		m.setRegister(reg, uint64(m.memory[addr]))
	case datatype.U16:
		num := binary.BigEndian.Uint16(m.memory[addr : addr+2])
		m.setRegister(reg, uint64(num))
	case datatype.U32:
		num := binary.BigEndian.Uint32(m.memory[addr : addr+4])
		m.setRegister(reg, uint64(num))
	case datatype.U64, datatype.UNSET:
		num := binary.BigEndian.Uint64(m.memory[addr : addr+8])
		m.setRegister(reg, num)
	}
}

func (m *Machine) handleMovAofReg() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	pos := m.position()
	rdt := datatype.DataType(m.decodeNumber("u8", pos))
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	value := m.decodeValue(rdt)
	addr := 0
	switch v := value.(type) {
	case *ast.NumberLiteral:
		num, _ := strconv.ParseUint(v.Value, 10, 64)
		addr += int(num)
	case *ast.Register:
		num := m.getRegister(v.Value)
		addr += int(num)
	case *ast.RegisterOffset:
		num := m.getRegister(v.Left.Value)
		addr += int(num)
	case *ast.LabelOffset:
		num, _ := strconv.ParseUint(v.Left.(*ast.NumberLiteral).Value, 10, 64)
		addr += int(num)
	default:
		log.Fatal("unknown value to get address of", "value", value)
	}

	pos = m.position()
	reg := m.decodeRegister(pos)
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	dt := datatype.DataType(datatype.UNSET)
	if adt, ok := m.symbolTable[uint64(addr)]; ok {
		dt = adt
	}
	if rdt != datatype.UNSET {
		dt = rdt
	}

	switch v := value.(type) {
	case *ast.RegisterOffset:
		addr = applyOffset(addr, v.Operator, v.Right.Value)
	case *ast.LabelOffset:
		addr = applyOffset(addr, v.Operator, v.Right.Value)
	}

	switch dt {
	case datatype.U8:
		m.memory[addr] = byte(m.getRegister(reg))
	case datatype.U16:
		bytes := utils.Bytes2(uint16(m.getRegister(reg)))
		copy(m.memory[addr:addr+2], bytes[:])
	case datatype.U32:
		bytes := utils.Bytes4(uint32(m.getRegister(reg)))
		copy(m.memory[addr:addr+4], bytes[:])
	case datatype.U64, datatype.UNSET:
		bytes := utils.Bytes8(m.getRegister(reg))
		copy(m.memory[addr:addr+8], bytes[:])
	}
}

func (m *Machine) handleMovAofLit() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	pos := m.position()
	rdt := datatype.DataType(m.decodeNumber("u8", pos))
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	value := m.decodeValue(rdt)
	addr := 0
	switch v := value.(type) {
	case *ast.NumberLiteral:
		num, _ := strconv.ParseUint(v.Value, 10, 64)
		addr += int(num)
	case *ast.Register:
		num := m.getRegister(v.Value)
		addr += int(num)
	case *ast.RegisterOffset:
		num := m.getRegister(v.Left.Value)
		addr += int(num)
	case *ast.LabelOffset:
		num, _ := strconv.ParseUint(v.Left.(*ast.NumberLiteral).Value, 10, 64)
		addr += int(num)
	default:
		log.Fatal("unknown value to get address of", "value", value)
	}

	pos = m.position()
	lit := m.decodeNumber(rdt.String(), pos)
	m.incRegister(utils.RegisterToIndex("ip"), uint64(rdt.Size()))

	dt := datatype.DataType(datatype.UNSET)
	if adt, ok := m.symbolTable[uint64(addr)]; ok {
		dt = adt
	}
	if rdt != datatype.UNSET {
		dt = rdt
	}

	switch v := value.(type) {
	case *ast.RegisterOffset:
		addr = applyOffset(addr, v.Operator, v.Right.Value)
	case *ast.LabelOffset:
		addr = applyOffset(addr, v.Operator, v.Right.Value)
	}

	switch dt {
	case datatype.U8:
		m.memory[addr] = byte(lit)
	case datatype.U16:
		bytes := utils.Bytes2(uint16(lit))
		copy(m.memory[addr:addr+2], bytes[:])
	case datatype.U32:
		bytes := utils.Bytes4(uint32(lit))
		copy(m.memory[addr:addr+4], bytes[:])
	case datatype.U64, datatype.UNSET:
		bytes := utils.Bytes8(uint64(lit))
		copy(m.memory[addr:addr+8], bytes[:])
	}
}

func applyOffset(addr int, operator ast.Operator, value string) int {
	num, _ := strconv.ParseUint(value, 10, 64)
	switch operator {
	case ast.ADD:
		return addr + int(num)
	case ast.SUBTRACT:
		return addr - int(num)
	case ast.MULTIPLY:
		return addr * int(num)
	case ast.DIVIDE:
		return addr / int(num)
	default:
		return addr
	}
}
