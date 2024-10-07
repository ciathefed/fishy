package vm

import (
	"encoding/binary"
	"fishy/pkg/ast"
	"fishy/pkg/datatype"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"fmt"
	"log"
	"strconv"
)

type Machine struct {
	registers   []uint32
	memory      []byte
	symbolTable map[uint32]datatype.DataType
	debug       bool
}

func New(bytecode []byte, memorySize int, debug bool) *Machine {
	m := &Machine{
		registers:   make([]uint32, 21),
		memory:      bytecode,
		symbolTable: make(map[uint32]datatype.DataType),
		debug:       debug,
	}

	m.parserHeaderStart()
	m.parseHeaderSymbolTable()

	m.memory = append(m.memory, make([]byte, memorySize-len(m.memory))...)

	m.setRegister(utils.RegisterToIndex("sp"), uint32(len(m.memory)))
	m.setRegister(utils.RegisterToIndex("fp"), uint32(len(m.memory)))

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
			panic(fmt.Sprintf("unhandled instruction: %s", op.String()))
		}
	}
}

func (m *Machine) decodeNumber(dataType string, index int) int {
	switch dataType {
	case "u8":
		return int(m.memory[index])
	case "u16":
		bytes := m.memory[index : index+2]
		return int(binary.BigEndian.Uint16(bytes))
	case "u32":
		bytes := m.memory[index : index+4]
		return int(binary.BigEndian.Uint32(bytes))
	case "u64":
		bytes := m.memory[index : index+8]
		return int(binary.BigEndian.Uint64(bytes))
	default:
		panic(fmt.Sprintf("unknown data type: %s", dataType))
	}
}

func (m *Machine) decodeRegister(index int) int {
	v := m.memory[index]
	return int(v)
}

func (m *Machine) decodeValue() ast.Value {
	pos := m.position()
	indexValue := m.memory[pos]
	m.incRegister(utils.RegisterToIndex("ip"), 1)

	switch indexValue {
	case 2:
		pos := m.position()
		reg := m.decodeRegister(pos)
		m.incRegister(utils.RegisterToIndex("ip"), 1)
		return &ast.Register{Value: reg}
	case 0, 1, 3, 4, 5:
		pos := m.position()
		num := m.decodeNumber("u32", pos)
		m.incRegister(utils.RegisterToIndex("ip"), 4)
		return &ast.NumberLiteral{Value: strconv.Itoa(num)}
	default:
		panic(fmt.Sprintf("unknown value index: %d", indexValue))
	}
}

func (m *Machine) parserHeaderStart() {
	start := m.readLiteral()
	m.memory = m.memory[4:]
	m.setRegister(utils.RegisterToIndex("ip"), uint32(start))
}

func (m *Machine) parseHeaderSymbolTable() {
	start := binary.BigEndian.Uint32(m.memory[:4]) - 4
	end := binary.BigEndian.Uint32(m.memory[4:8]) - 4

	keyValues := m.memory[start:end]

	for i := 0; i < len(keyValues); i += 8 {
		if i+7 >= len(keyValues) {
			log.Fatal("symbol table is not multiple of 8", "length", len(keyValues))
		}

		key := binary.BigEndian.Uint32(keyValues[i : i+4])
		value := binary.BigEndian.Uint32(keyValues[i+4 : i+8])

		m.symbolTable[key] = datatype.DataType(value)
	}

	m.memory = m.memory[end:]
}

func (m *Machine) position() int {
	index := utils.RegisterToIndex("ip")
	return int(m.getRegister(index))
}

func (m *Machine) setRegister(index int, value uint32) {
	m.registers[index] = value
}

func (m *Machine) getRegister(index int) uint32 {
	return m.registers[index]
}

func (m *Machine) incRegister(index int, amount int) {
	m.registers[index] += uint32(amount)
}

func (m *Machine) readRegister() int {
	pos := m.position()
	reg := m.decodeRegister(pos)
	m.incRegister(utils.RegisterToIndex("ip"), 1)
	return reg
}

func (m *Machine) readLiteral() uint32 {
	pos := m.position()
	lit := m.decodeNumber("u32", pos)
	m.incRegister(utils.RegisterToIndex("ip"), 4)
	return uint32(lit)
}

func (m *Machine) stackPush(v uint32) {
	spIndex := utils.RegisterToIndex("sp")
	spValue := m.getRegister(spIndex)

	byteArray := utils.Bytes4(v)

	memIndex := int(spValue) - 4

	copy(m.memory[memIndex:memIndex+4], byteArray)

	m.setRegister(spIndex, spValue-4)
}

func (m *Machine) stackPop() uint32 {
	spIndex := utils.RegisterToIndex("sp")
	spValue := m.getRegister(spIndex)

	memIndex := int(spValue)

	value := binary.BigEndian.Uint32(m.memory[memIndex : memIndex+4])

	m.setRegister(spIndex, spValue+4)

	return value
}

func (m *Machine) DumpRegisters() {
	for i, register := range m.registers {
		name := utils.IndexToRegister(i)
		fmt.Printf("%-3s: 0x%08X\n", name, register)
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
