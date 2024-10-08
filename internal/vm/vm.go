package vm

import (
	"encoding/binary"
	"fishy/pkg/ast"
	"fishy/pkg/datatype"
	"fishy/pkg/log"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"fmt"
	"strconv"
	"strings"
)

type Machine struct {
	registers   []uint64
	memory      []byte
	symbolTable map[uint64]datatype.DataType
	debug       bool
}

func New(bytecode []byte, memorySize int, debug bool) *Machine {
	m := &Machine{
		registers:   make([]uint64, 21),
		memory:      bytecode,
		symbolTable: make(map[uint64]datatype.DataType),
		debug:       debug,
	}

	m.parserHeaderStart()
	m.parseHeaderSymbolTable()

	m.memory = append(m.memory, make([]byte, memorySize-len(m.memory))...)

	m.setRegister(utils.RegisterToIndex("sp"), uint64(len(m.memory)))
	m.setRegister(utils.RegisterToIndex("fp"), uint64(len(m.memory)))

	return m
}

func (m *Machine) Run() {
	for {
		pos := m.position()
		instruction := m.decodeNumber("u16", pos)
		op := opcode.Opcode(instruction)

		switch op {
		case opcode.NOP:
			m.incRegister(utils.RegisterToIndex("ip"), 1)
		case opcode.HLT:
			return
		case opcode.BRK:
			if m.debug {
				panic("todo")
			} else {
				m.incRegister(utils.RegisterToIndex("ip"), 1)
			}
		case opcode.SYSCALL:
			m.handleSyscall()
		case opcode.MOV_REG_REG:
			m.handleMovRegReg()
		case opcode.MOV_REG_LIT:
			m.handleMovRegLit()
		case opcode.MOV_REG_ADR:
			m.handleMovRegAdr()
		case opcode.MOV_REG_AOF:
			m.handleMovRegAof()
		case opcode.MOV_AOF_REG:
			m.handleMovAofReg()
		case opcode.MOV_AOF_LIT:
			m.handleMovAofLit()
		case opcode.ADD_REG_LIT, opcode.ADD_REG_REG,
			opcode.SUB_REG_LIT, opcode.SUB_REG_REG,
			opcode.MUL_REG_LIT, opcode.MUL_REG_REG,
			opcode.DIV_REG_LIT, opcode.DIV_REG_REG:
			m.handleArithmetic(op)
		case opcode.AND_REG_LIT,
			opcode.AND_REG_REG,
			opcode.OR_REG_LIT,
			opcode.OR_REG_REG,
			opcode.XOR_REG_LIT,
			opcode.XOR_REG_REG,
			opcode.SHL_REG_LIT,
			opcode.SHL_REG_REG,
			opcode.SHR_REG_LIT,
			opcode.SHR_REG_REG:
			m.handleBitwise(op)
		case opcode.CMP_REG_LIT, opcode.CMP_REG_REG:
			m.handleCompare(op)
		case opcode.JMP_LIT, opcode.JMP_REG,
			opcode.JEQ_LIT, opcode.JEQ_REG,
			opcode.JNE_LIT, opcode.JNE_REG,
			opcode.JLT_LIT, opcode.JLT_REG,
			opcode.JGT_LIT, opcode.JGT_REG,
			opcode.JLE_LIT, opcode.JLE_REG,
			opcode.JGE_LIT, opcode.JGE_REG:
			m.handleJump(op)
		case opcode.PUSH_LIT, opcode.PUSH_REG:
			m.handlePush(op)
		case opcode.POP_REG:
			m.handlePop(op)
		case opcode.CALL_LIT:
			m.handleCallLit()
		case opcode.RET:
			m.handleRet()
		default:
			log.Fatal("unhandled instruction", "op", op.String())
		}
	}
}

func (m *Machine) decodeNumber(dataType string, index int) int {
	dataType = strings.ToLower(dataType)
	switch dataType {
	case "u8":
		return int(m.memory[index])
	case "u16":
		bytes := m.memory[index : index+2]
		return int(binary.BigEndian.Uint16(bytes))
	case "u32":
		bytes := m.memory[index : index+4]
		return int(binary.BigEndian.Uint32(bytes))
	case "u64", "unset":
		bytes := m.memory[index : index+8]
		return int(binary.BigEndian.Uint64(bytes))
	default:
		log.Fatal("unknown data type", "type", dataType)
	}
	return -1
}

func (m *Machine) decodeRegister(index int) int {
	v := m.memory[index]
	return int(v)
}

func (m *Machine) decodeValue(dataType datatype.DataType) ast.Value {
	pos := m.position()
	indexValue := m.memory[pos]
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	switch indexValue {
	case 2:
		pos := m.position()
		reg := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		return &ast.Register{Value: reg}
	case 0, 1, 3, 4:
		pos := m.position()
		num := m.decodeNumber(dataType.String(), pos)
		m.incRegister(utils.RegisterToIndex("ip"), uint64(dataType.Size()))
		return &ast.NumberLiteral{Value: strconv.Itoa(num)}
	case 5:
		pos := m.position()
		reg := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)

		pos = m.position()
		op := ast.Operator(m.decodeNumber("u8", pos))
		m.incRegister(utils.RegisterToIndex("ip"), 1)

		pos = m.position()
		offset := m.decodeNumber(dataType.String(), pos)
		m.incRegister(utils.RegisterToIndex("ip"), uint64(dataType.Size()))

		return &ast.RegisterOffsetNumber{
			Left:     ast.Register{Value: reg},
			Operator: op,
			Right:    ast.NumberLiteral{Value: strconv.Itoa(offset)},
		}
	case 6:
		pos := m.position()
		reg0 := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)

		pos = m.position()
		op := ast.Operator(m.decodeNumber("u8", pos))
		m.incRegister(utils.RegisterToIndex("ip"), 1)

		pos = m.position()
		reg1 := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)

		return &ast.RegisterOffsetRegister{
			Left:     ast.Register{Value: reg0},
			Operator: op,
			Right:    ast.Register{Value: reg1},
		}
	case 7:
		pos := m.position()
		addr := m.decodeNumber(dataType.String(), pos)
		m.incRegister(utils.RegisterToIndex("ip"), uint64(dataType.Size()))

		pos = m.position()
		op := ast.Operator(m.decodeNumber("u8", pos))
		m.incRegister(utils.RegisterToIndex("ip"), 1)

		pos = m.position()
		offset := m.decodeNumber(dataType.String(), pos)
		m.incRegister(utils.RegisterToIndex("ip"), uint64(dataType.Size()))

		return &ast.LabelOffsetNumber{
			Left:     &ast.NumberLiteral{Value: strconv.Itoa(addr)},
			Operator: op,
			Right:    ast.NumberLiteral{Value: strconv.Itoa(offset)},
		}
	case 8:
		pos := m.position()
		addr := m.decodeNumber(dataType.String(), pos)
		m.incRegister(utils.RegisterToIndex("ip"), uint64(dataType.Size()))

		pos = m.position()
		op := ast.Operator(m.decodeNumber("u8", pos))
		m.incRegister(utils.RegisterToIndex("ip"), 1)

		pos = m.position()
		reg := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)

		return &ast.LabelOffsetRegister{
			Left:     &ast.NumberLiteral{Value: strconv.Itoa(addr)},
			Operator: op,
			Right:    ast.Register{Value: reg},
		}
	default:
		log.Fatal("unknown value index", "index", indexValue)
	}
	return nil
}

func (m *Machine) parserHeaderStart() {
	start := m.readLiteral(datatype.U64)
	m.memory = m.memory[8:]
	m.setRegister(utils.RegisterToIndex("ip"), uint64(start))
}

func (m *Machine) parseHeaderSymbolTable() {
	size := int(binary.BigEndian.Uint64(m.memory[:8]))
	start := binary.BigEndian.Uint64(m.memory[8:16]) - 8
	end := binary.BigEndian.Uint64(m.memory[16:24]) - 8

	keyValues := m.memory[start:end]

	for i := 0; i < len(keyValues); i += size + 1 {
		if i+size >= len(keyValues) {
			log.Fatal(fmt.Sprintf("symbol table is not multiple of %d", size+1), "length", len(keyValues))
		}

		var key uint64
		switch size {
		case 1:
			key = uint64(keyValues[i])
		case 2:
			key = uint64(binary.BigEndian.Uint16(keyValues[i : i+size]))
		case 4:
			key = uint64(binary.BigEndian.Uint32(keyValues[i : i+size]))
		default:
			key = binary.BigEndian.Uint64(keyValues[i : i+size])
		}
		value := keyValues[i+size]

		m.symbolTable[key] = datatype.DataType(value)
	}

	m.memory = m.memory[end:]
}

func (m *Machine) position() int {
	index := utils.RegisterToIndex("ip")
	return int(m.getRegister(index))
}

func (m *Machine) setRegister(index int, value uint64) {
	m.registers[index] = value
}

func (m *Machine) getRegister(index int) uint64 {
	return m.registers[index]
}

func (m *Machine) incRegister(index int, amount uint64) {
	m.registers[index] += amount
}

func (m *Machine) readRegister() int {
	pos := m.position()
	reg := m.decodeRegister(pos)
	m.incRegister(utils.RegisterToIndex("ip"), 1)
	return reg
}

func (m *Machine) readLiteral(dataType datatype.DataType) uint64 {
	pos := m.position()
	lit := m.decodeNumber(dataType.String(), pos)
	m.incRegister(utils.RegisterToIndex("ip"), uint64(dataType.Size()))
	return uint64(lit)
}

func (m *Machine) stackPush(v uint64) {
	spIndex := utils.RegisterToIndex("sp")
	spValue := m.getRegister(spIndex)

	byteArray := utils.Bytes8(v)

	memIndex := int(spValue) - 8

	copy(m.memory[memIndex:memIndex+8], byteArray)

	m.setRegister(spIndex, spValue-8)
}

func (m *Machine) stackPop() uint64 {
	spIndex := utils.RegisterToIndex("sp")
	spValue := m.getRegister(spIndex)

	memIndex := int(spValue)

	value := binary.BigEndian.Uint64(m.memory[memIndex : memIndex+8])

	m.setRegister(spIndex, spValue+8)

	return value
}

func (m *Machine) DumpRegisters() {
	for i, register := range m.registers {
		name := utils.IndexToRegister(i)
		fmt.Printf("%-3s: 0x%016X\n", name, register)
	}
}

func (m *Machine) DumpMemory(start int, end int) {
	if end >= len(m.memory) {
		end = len(m.memory)
	}
	bytecode := m.memory[start:end]
	numLines := (len(bytecode) + 15) / 16

	for line := 0; line < numLines; line++ {
		startIndex := line * 16
		endIndex := startIndex + 16
		if endIndex > len(bytecode) {
			endIndex = len(bytecode)
		}
		lineBytes := bytecode[startIndex:endIndex]

		fmt.Printf("0x%04X: ", startIndex)

		for _, b := range lineBytes {
			fmt.Printf("%02X ", b)
		}
		fmt.Println()
	}
}
