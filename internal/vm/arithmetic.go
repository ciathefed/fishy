package vm

import (
	"encoding/binary"
	"fishy/pkg/ast"
	"fishy/pkg/datatype"
	"fishy/pkg/log"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"strconv"
)

func (m *Machine) handleArithmetic(thread *Thread, op opcode.Opcode) {
	m.incRegister(thread, utils.RegisterToIndex("ip"), 2)

	pos := m.position(thread)
	rdt := datatype.DataType(m.decodeNumber("byte", pos))
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	switch op {
	case opcode.ADD_REG_LIT:
		m.applyRegLitArithmetic(thread, rdt, func(reg, lit uint64) uint64 { return reg + lit })
	case opcode.ADD_REG_REG:
		m.applyRegRegArithmetic(thread, func(reg0, reg1 uint64) uint64 { return reg0 + reg1 })
	case opcode.ADD_REG_AOF:
		m.applyRegAof(thread, rdt, func(reg0, value uint64) uint64 { return reg0 + value })
	case opcode.SUB_REG_LIT:
		m.applyRegLitArithmetic(thread, rdt, func(reg, lit uint64) uint64 { return reg - lit })
	case opcode.SUB_REG_REG:
		m.applyRegRegArithmetic(thread, func(reg0, reg1 uint64) uint64 { return reg0 - reg1 })
	case opcode.SUB_REG_AOF:
		m.applyRegAof(thread, rdt, func(reg0, value uint64) uint64 { return reg0 - value })
	case opcode.MUL_REG_LIT:
		m.applyRegLitArithmetic(thread, rdt, func(reg, lit uint64) uint64 { return reg * lit })
	case opcode.MUL_REG_REG:
		m.applyRegRegArithmetic(thread, func(reg0, reg1 uint64) uint64 { return reg0 * reg1 })
	case opcode.MUL_REG_AOF:
		m.applyRegAof(thread, rdt, func(reg0, value uint64) uint64 { return reg0 * value })
	case opcode.DIV_REG_LIT:
		m.applyRegLitArithmetic(thread, rdt, func(reg, lit uint64) uint64 { return reg / lit })
	case opcode.DIV_REG_REG:
		m.applyRegRegArithmetic(thread, func(reg0, reg1 uint64) uint64 { return reg0 / reg1 })
	case opcode.DIV_REG_AOF:
		m.applyRegAof(thread, rdt, func(reg0, value uint64) uint64 { return reg0 / value })
	}
}

func (m *Machine) applyRegLitArithmetic(thread *Thread, dataType datatype.DataType, operation func(uint64, uint64) uint64) {
	reg := m.readRegister(thread)
	lit := m.readLiteral(thread, dataType)
	temp := m.getRegister(thread, reg)
	m.setRegister(thread, reg, operation(temp, lit))
}

func (m *Machine) applyRegRegArithmetic(thread *Thread, operation func(uint64, uint64) uint64) {
	reg0 := m.readRegister(thread)
	reg1 := m.readRegister(thread)
	temp0 := m.getRegister(thread, reg0)
	temp1 := m.getRegister(thread, reg1)
	m.setRegister(thread, reg0, operation(temp0, temp1))
}

func (m *Machine) applyRegAof(thread *Thread, dataType datatype.DataType, operation func(uint64, uint64) uint64) {
	reg0 := m.readRegister(thread)

	value := m.decodeValue(thread, dataType)
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
	if dataType != datatype.UNSET {
		dt = dataType
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
		temp0 := m.getRegister(thread, reg0)
		temp1 := uint64(m.memory[addr])
		m.setRegister(thread, reg0, operation(temp0, temp1))
	case datatype.WORD:
		temp0 := m.getRegister(thread, reg0)
		temp1 := binary.BigEndian.Uint16(m.memory[addr : addr+dt.Size()])
		m.setRegister(thread, reg0, operation(temp0, uint64(temp1)))
	case datatype.DWORD:
		temp0 := m.getRegister(thread, reg0)
		temp1 := binary.BigEndian.Uint32(m.memory[addr : addr+dt.Size()])
		m.setRegister(thread, reg0, operation(temp0, uint64(temp1)))
	case datatype.QWORD, datatype.UNSET:
		temp0 := m.getRegister(thread, reg0)
		temp1 := binary.BigEndian.Uint64(m.memory[addr : addr+dt.Size()])
		m.setRegister(thread, reg0, operation(temp0, temp1))
	}
}
