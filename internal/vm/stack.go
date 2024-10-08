package vm

import (
	"fishy/pkg/ast"
	"fishy/pkg/datatype"
	"fishy/pkg/log"
	"fishy/pkg/utils"
	"strconv"
)

func (m *Machine) handlePushLit() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	pos := m.position()
	rdt := datatype.DataType(m.decodeNumber("u8", pos))
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	pos = m.position()
	lit := m.decodeNumberBytes(rdt.String(), pos)
	m.incRegister(utils.RegisterToIndex("ip"), uint64(rdt.Size()))
	m.stackPush(lit)
}

func (m *Machine) handlePushReg() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	pos := m.position()
	reg := m.decodeRegister(pos)
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	pos = m.position()
	rdt := datatype.DataType(m.decodeNumber("u8", pos))
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	m.stackPush(rdt.MakeBytes(m.getRegister(reg)))
}

func (m *Machine) handlePushAof() {
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
	case *ast.RegisterOffsetNumber:
		num := m.getRegister(v.Left.Value)
		addr += int(num)
	case *ast.RegisterOffsetRegister:
		num := m.getRegister(v.Left.Value)
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
		addr = applyOffset(addr, v.Operator, strconv.Itoa(int(m.getRegister(v.Right.Value))))
	case *ast.LabelOffsetNumber:
		addr = applyOffset(addr, v.Operator, v.Right.Value)
	case *ast.LabelOffsetRegister:
		addr = applyOffset(addr, v.Operator, strconv.Itoa(int(m.getRegister(v.Right.Value))))
	}

	switch dt {
	case datatype.U8:
		m.stackPush([]byte{m.memory[addr]})
	case datatype.U16:
		m.stackPush(m.memory[addr : addr+2])
	case datatype.U32:
		m.stackPush(m.memory[addr : addr+4])
	case datatype.U64, datatype.UNSET:
		m.stackPush(m.memory[addr : addr+8])
	}
}

func (m *Machine) handlePopReg() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	pos := m.position()
	rdt := datatype.DataType(m.decodeNumber("u8", pos))
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	pos = m.position()
	reg := m.decodeRegister(pos)
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	m.setRegister(reg, m.stackPop(rdt))
}

func (m *Machine) handlePopAof() {
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
	case *ast.RegisterOffsetNumber:
		num := m.getRegister(v.Left.Value)
		addr += int(num)
	case *ast.RegisterOffsetRegister:
		num := m.getRegister(v.Left.Value)
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
		addr = applyOffset(addr, v.Operator, strconv.Itoa(int(m.getRegister(v.Right.Value))))
	case *ast.LabelOffsetNumber:
		addr = applyOffset(addr, v.Operator, v.Right.Value)
	case *ast.LabelOffsetRegister:
		addr = applyOffset(addr, v.Operator, strconv.Itoa(int(m.getRegister(v.Right.Value))))
	}

	bytes := m.stackPopBytes(rdt)

	switch dt {
	case datatype.U8:
		m.memory[addr] = bytes[0]
	case datatype.U16:
		copy(m.memory[addr:addr+dt.Size()], bytes[:])
	case datatype.U32:
		copy(m.memory[addr:addr+dt.Size()], bytes[:])
	case datatype.U64, datatype.UNSET:
		copy(m.memory[addr:addr+dt.Size()], bytes[:])
	}
}
