package vm

import (
	"fishy/pkg/ast"
	"fishy/pkg/datatype"
	"fishy/pkg/log"
	"fishy/pkg/utils"
	"strconv"
)

func (m *Machine) handlePushLit(thread *Thread) {
	m.incRegister(thread, utils.RegisterToIndex("ip"), 2)

	pos := m.position(thread)
	rdt := datatype.DataType(m.decodeNumber("byte", pos))
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	pos = m.position(thread)
	lit := m.decodeNumberBytes(rdt.String(), pos)
	m.incRegister(thread, utils.RegisterToIndex("ip"), uint64(rdt.Size()))
	m.stackPush(thread, lit)
}

func (m *Machine) handlePushReg(thread *Thread) {
	m.incRegister(thread, utils.RegisterToIndex("ip"), 2)

	pos := m.position(thread)
	rdt := datatype.DataType(m.decodeNumber("byte", pos))
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	pos = m.position(thread)
	reg := m.decodeRegister(pos)
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	m.stackPush(thread, rdt.MakeBytes(m.getRegister(thread, reg)))
}

func (m *Machine) handlePushAof(thread *Thread) {
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
		m.stackPush(thread, []byte{m.memory[addr]})
	case datatype.WORD:
		m.stackPush(thread, m.memory[addr:addr+2])
	case datatype.DWORD:
		m.stackPush(thread, m.memory[addr:addr+4])
	case datatype.QWORD, datatype.UNSET:
		m.stackPush(thread, m.memory[addr:addr+8])
	}
}

func (m *Machine) handlePopReg(thread *Thread) {
	m.incRegister(thread, utils.RegisterToIndex("ip"), 2)

	pos := m.position(thread)
	rdt := datatype.DataType(m.decodeNumber("byte", pos))
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	pos = m.position(thread)
	reg := m.decodeRegister(pos)
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	m.setRegister(thread, reg, m.stackPop(thread, rdt))
}

func (m *Machine) handlePopAof(thread *Thread) {
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

	bytes := m.stackPopBytes(thread, rdt)

	switch dt {
	case datatype.BYTE:
		m.memory[addr] = bytes[0]
	case datatype.WORD:
		copy(m.memory[addr:addr+dt.Size()], bytes[:])
	case datatype.DWORD:
		copy(m.memory[addr:addr+dt.Size()], bytes[:])
	case datatype.QWORD, datatype.UNSET:
		copy(m.memory[addr:addr+dt.Size()], bytes[:])
	}
}
