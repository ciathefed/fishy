package vm

import (
	"encoding/binary"
	"fishy/pkg/ast"
	"fishy/pkg/datatype"
	"fishy/pkg/log"
	"fishy/pkg/utils"
	"strconv"
)

func (m *Machine) handleMovRegReg(thread *Thread) {
	m.incRegister(thread, utils.RegisterToIndex("ip"), 2)

	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	pos := m.position(thread)
	reg0 := m.decodeRegister(pos)
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	pos = m.position(thread)
	reg1 := m.decodeRegister(pos)
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	temp := m.getRegister(thread, reg1)
	m.setRegister(thread, reg0, temp)
}

func (m *Machine) handleMovRegLit(thread *Thread) {
	m.incRegister(thread, utils.RegisterToIndex("ip"), 2)

	pos := m.position(thread)
	dt := datatype.DataType(m.decodeNumber("byte", pos))
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	pos = m.position(thread)
	reg := m.decodeRegister(pos)
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	pos = m.position(thread)
	lit := m.decodeNumber(dt.String(), pos)
	m.incRegister(thread, utils.RegisterToIndex("ip"), uint64(dt.Size()))

	m.setRegister(thread, reg, uint64(lit))
}

func (m *Machine) handleMovRegAdr(thread *Thread) {
	m.incRegister(thread, utils.RegisterToIndex("ip"), 2)

	pos := m.position(thread)
	dt := datatype.DataType(m.decodeNumber("byte", pos))
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	pos = m.position(thread)
	reg := m.decodeRegister(pos)
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	pos = m.position(thread)
	addr := m.decodeNumber(dt.String(), pos)
	m.incRegister(thread, utils.RegisterToIndex("ip"), uint64(dt.Size()))

	m.setRegister(thread, reg, uint64(addr))
}

func (m *Machine) handleMovRegAof(thread *Thread) {
	m.incRegister(thread, utils.RegisterToIndex("ip"), 2)

	pos := m.position(thread)
	rdt := datatype.DataType(m.decodeNumber("byte", pos))
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	pos = m.position(thread)
	reg := m.decodeRegister(pos)
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	value := m.decodeValue(thread, rdt)
	addr := 0
	switch v := value.(type) {
	case *ast.NumberLiteral:
		num, _ := strconv.ParseUint(v.Value, 10, 64)
		addr += int(num)
	case *ast.Register:
		num := m.getRegister(thread, v.Value)
		addr += int(num)
	case *ast.RegisterOffsetNumber:
		num := m.getRegister(thread, v.Left.Value)
		addr += int(num)
	case *ast.RegisterOffsetRegister:
		num := m.getRegister(thread, v.Left.Value)
		addr += int(num)
	case *ast.LabelOffsetNumber:
		num, _ := strconv.ParseUint(v.Left.(*ast.NumberLiteral).Value, 10, 64)
		addr += int(num)
	case *ast.LabelOffsetRegister:
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
	case *ast.RegisterOffsetNumber:
		addr = applyOffset(addr, v.Operator, v.Right.Value)
	case *ast.RegisterOffsetRegister:
		addr = applyOffset(addr, v.Operator, strconv.Itoa(int(m.getRegister(thread, v.Right.Value))))
	case *ast.LabelOffsetNumber:
		addr = applyOffset(addr, v.Operator, v.Right.Value)
	case *ast.LabelOffsetRegister:
		addr = applyOffset(addr, v.Operator, strconv.Itoa(int(m.getRegister(thread, v.Right.Value))))
	}

	switch dt {
	case datatype.BYTE:
		m.setRegister(thread, reg, uint64(m.memory[addr]))
	case datatype.WORD:
		num := binary.BigEndian.Uint16(m.memory[addr : addr+2])
		m.setRegister(thread, reg, uint64(num))
	case datatype.DWORD:
		num := binary.BigEndian.Uint32(m.memory[addr : addr+4])
		m.setRegister(thread, reg, uint64(num))
	case datatype.QWORD, datatype.UNSET:
		num := binary.BigEndian.Uint64(m.memory[addr : addr+8])
		m.setRegister(thread, reg, num)
	}
}

func (m *Machine) handleMovAofReg(thread *Thread) {
	m.incRegister(thread, utils.RegisterToIndex("ip"), 2)

	pos := m.position(thread)
	rdt := datatype.DataType(m.decodeNumber("byte", pos))
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	value := m.decodeValue(thread, rdt)
	addr := 0
	switch v := value.(type) {
	case *ast.NumberLiteral:
		num, _ := strconv.ParseUint(v.Value, 10, 64)
		addr += int(num)
	case *ast.Register:
		num := m.getRegister(thread, v.Value)
		addr += int(num)
	case *ast.RegisterOffsetNumber:
		num := m.getRegister(thread, v.Left.Value)
		addr += int(num)
	case *ast.RegisterOffsetRegister:
		num := m.getRegister(thread, v.Left.Value)
		addr += int(num)
	case *ast.LabelOffsetNumber:
		num, _ := strconv.ParseUint(v.Left.(*ast.NumberLiteral).Value, 10, 64)
		addr += int(num)
	case *ast.LabelOffsetRegister:
		num, _ := strconv.ParseUint(v.Left.(*ast.NumberLiteral).Value, 10, 64)
		addr += int(num)
	default:
		log.Fatal("unknown value to get address of", "value", value)
	}

	pos = m.position(thread)
	reg := m.decodeRegister(pos)
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	dt := datatype.DataType(datatype.UNSET)
	if adt, ok := m.symbolTable[uint64(addr)]; ok {
		dt = adt
	}
	if rdt != datatype.UNSET {
		dt = rdt
	}

	switch v := value.(type) {
	case *ast.RegisterOffsetNumber:
		addr = applyOffset(addr, v.Operator, v.Right.Value)
	case *ast.RegisterOffsetRegister:
		addr = applyOffset(addr, v.Operator, strconv.Itoa(int(m.getRegister(thread, v.Right.Value))))
	case *ast.LabelOffsetNumber:
		addr = applyOffset(addr, v.Operator, v.Right.Value)
	case *ast.LabelOffsetRegister:
		addr = applyOffset(addr, v.Operator, strconv.Itoa(int(m.getRegister(thread, v.Right.Value))))
	}

	switch dt {
	case datatype.BYTE:
		m.memory[addr] = byte(m.getRegister(thread, reg))
	case datatype.WORD:
		bytes := utils.Bytes2(uint16(m.getRegister(thread, reg)))
		copy(m.memory[addr:addr+dt.Size()], bytes[:])
	case datatype.DWORD:
		bytes := utils.Bytes4(uint32(m.getRegister(thread, reg)))
		copy(m.memory[addr:addr+dt.Size()], bytes[:])
	case datatype.QWORD, datatype.UNSET:
		bytes := utils.Bytes8(m.getRegister(thread, reg))
		copy(m.memory[addr:addr+dt.Size()], bytes[:])
	}
}

func (m *Machine) handleMovAofLit(thread *Thread) {
	m.incRegister(thread, utils.RegisterToIndex("ip"), 2)

	pos := m.position(thread)
	rdt := datatype.DataType(m.decodeNumber("byte", pos))
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	value := m.decodeValue(thread, rdt)
	addr := 0
	switch v := value.(type) {
	case *ast.NumberLiteral:
		num, _ := strconv.ParseUint(v.Value, 10, 64)
		addr += int(num)
	case *ast.Register:
		num := m.getRegister(thread, v.Value)
		addr += int(num)
	case *ast.RegisterOffsetNumber:
		num := m.getRegister(thread, v.Left.Value)
		addr += int(num)
	case *ast.RegisterOffsetRegister:
		num := m.getRegister(thread, v.Left.Value)
		addr += int(num)
	case *ast.LabelOffsetNumber:
		num, _ := strconv.ParseUint(v.Left.(*ast.NumberLiteral).Value, 10, 64)
		addr += int(num)
	case *ast.LabelOffsetRegister:
		num, _ := strconv.ParseUint(v.Left.(*ast.NumberLiteral).Value, 10, 64)
		addr += int(num)
	default:
		log.Fatal("unknown value to get address of", "value", value)
	}

	pos = m.position(thread)
	lit := m.decodeNumber(rdt.String(), pos)
	m.incRegister(thread, utils.RegisterToIndex("ip"), uint64(rdt.Size()))

	dt := datatype.DataType(datatype.UNSET)
	if adt, ok := m.symbolTable[uint64(addr)]; ok {
		dt = adt
	}
	if rdt != datatype.UNSET {
		dt = rdt
	}

	switch v := value.(type) {
	case *ast.RegisterOffsetNumber:
		addr = applyOffset(addr, v.Operator, v.Right.Value)
	case *ast.RegisterOffsetRegister:
		addr = applyOffset(addr, v.Operator, strconv.Itoa(int(m.getRegister(thread, v.Right.Value))))
	case *ast.LabelOffsetNumber:
		addr = applyOffset(addr, v.Operator, v.Right.Value)
	case *ast.LabelOffsetRegister:
		addr = applyOffset(addr, v.Operator, strconv.Itoa(int(m.getRegister(thread, v.Right.Value))))
	}

	switch dt {
	case datatype.BYTE:
		m.memory[addr] = byte(lit)
	case datatype.WORD:
		bytes := utils.Bytes2(uint16(lit))
		copy(m.memory[addr:addr+dt.Size()], bytes[:])
	case datatype.DWORD:
		bytes := utils.Bytes4(uint32(lit))
		copy(m.memory[addr:addr+dt.Size()], bytes[:])
	case datatype.QWORD, datatype.UNSET:
		bytes := utils.Bytes8(uint64(lit))
		copy(m.memory[addr:addr+dt.Size()], bytes[:])
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
