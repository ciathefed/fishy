package vm

import (
	"encoding/binary"
	"fishy/pkg/ast"
	"fishy/pkg/opcode"
	"fishy/pkg/utils"
	"fmt"
	"strconv"
)

type Machine struct {
	registers []uint32
	memory    []byte
	debug     bool
}

func New(bytecode []byte, memorySize int, debug bool) *Machine {
	memory := make([]byte, memorySize-len(bytecode)+4)
	memory = append(bytecode, memory...)

	m := &Machine{registers: make([]uint32, 23), memory: memory, debug: debug}
	m.parseHeader()
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
		case opcode.ADD_REG_LIT, opcode.ADD_REG_REG,
			opcode.SUB_REG_LIT, opcode.SUB_REG_REG,
			opcode.MUL_REG_LIT, opcode.MUL_REG_REG,
			opcode.DIV_REG_LIT, opcode.DIV_REG_REG:
			m.handleArithmetic(op)

		default:
			panic(fmt.Sprintf("unknown instruction: %s", op.String()))
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

func (m *Machine) parseHeader() {
	pos := m.position()
	start := m.decodeNumber("u32", pos)
	m.memory = m.memory[4:]
	m.setRegister(utils.RegisterToIndex("ip"), uint32(start))
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

func (m *Machine) DumpRegisters() {
	for i, register := range m.registers {
		name := utils.IndexToRegister(i)
		fmt.Printf("%-3s: 0x%08X\n", name, register)
	}
}

func (m *Machine) DumpMemory(start int, end int) {
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
