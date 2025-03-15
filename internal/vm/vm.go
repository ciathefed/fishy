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
	"sync"
)

type Thread struct {
	registers []uint64
	isRunning bool
	done      chan bool
}

type Machine struct {
	threads     map[int]*Thread
	mainThread  *Thread
	memory      []byte
	symbolTable map[uint64]datatype.DataType
	wg          *sync.WaitGroup
	debug       bool
}

func New(bytecode []byte, memorySize int, debug bool) *Machine {
	m := &Machine{
		threads:     make(map[int]*Thread),
		memory:      bytecode,
		symbolTable: make(map[uint64]datatype.DataType),
		wg:          &sync.WaitGroup{},
		debug:       debug,
	}

	thread := m.CreateThread()
	m.mainThread = m.threads[0]

	m.parserHeaderStart()
	m.parseHeaderSymbolTable()

	// fmt.Println(memorySize)
	// fmt.Println(len(m.memory))
	// fmt.Println(memorySize - len(m.memory))

	m.memory = append(m.memory, make([]byte, memorySize-len(m.memory))...)

	m.setRegister(thread, utils.RegisterToIndex("sp"), uint64(len(m.memory)))
	m.setRegister(thread, utils.RegisterToIndex("fp"), uint64(len(m.memory)))

	return m
}

func (m *Machine) CreateThread() *Thread {
	thread := &Thread{
		registers: make([]uint64, 21),
		isRunning: true,
		done:      make(chan bool),
	}
	m.threads[len(m.threads)] = thread
	return thread
}

func (m *Machine) GetThread(index int) (*Thread, bool) {
	thread, ok := m.threads[index]
	return thread, ok
}

func (m *Machine) GetThreadIndex(thread *Thread) (int, bool) {
	for key, value := range m.threads {
		if thread == value {
			return key, true
		}
	}
	return -1, false
}

func (m *Machine) RunThread(thread *Thread) {
	defer func() {
		thread.isRunning = false
		close(thread.done)
	}()

	for thread.isRunning {
		pos := m.position(thread)
		instruction := m.decodeNumber("word", pos)
		op := opcode.Opcode(instruction)

		switch op {
		case opcode.NOP:
			m.incRegister(thread, utils.RegisterToIndex("ip"), 1)
		case opcode.HLT:
			thread.isRunning = false
			return
		case opcode.BRK:
			if m.debug {
				panic("todo")
			} else {
				m.incRegister(thread, utils.RegisterToIndex("ip"), 1)
			}
		case opcode.SYSCALL:
			m.handleSyscall(thread)
		case opcode.MOV_REG_REG:
			m.handleMovRegReg(thread)
		case opcode.MOV_REG_LIT:
			m.handleMovRegLit(thread)
		case opcode.MOV_REG_ADR:
			m.handleMovRegAdr(thread)
		case opcode.MOV_REG_AOF:
			m.handleMovRegAof(thread)
		case opcode.MOV_AOF_REG:
			m.handleMovAofReg(thread)
		case opcode.MOV_AOF_LIT:
			m.handleMovAofLit(thread)
		case opcode.ADD_REG_LIT, opcode.ADD_REG_REG, opcode.ADD_REG_AOF,
			opcode.SUB_REG_LIT, opcode.SUB_REG_REG, opcode.SUB_REG_AOF,
			opcode.MUL_REG_LIT, opcode.MUL_REG_REG, opcode.MUL_REG_AOF,
			opcode.DIV_REG_LIT, opcode.DIV_REG_REG, opcode.DIV_REG_AOF:
			m.handleArithmetic(thread, op)
		case opcode.AND_REG_LIT, opcode.AND_REG_REG,
			opcode.OR_REG_LIT, opcode.OR_REG_REG,
			opcode.XOR_REG_LIT, opcode.XOR_REG_REG,
			opcode.SHL_REG_LIT, opcode.SHL_REG_REG,
			opcode.SHR_REG_LIT, opcode.SHR_REG_REG:
			m.handleBitwise(thread, op)
		case opcode.CMP_REG_LIT, opcode.CMP_REG_REG:
			m.handleCompare(thread, op)
		case opcode.JMP_LIT, opcode.JMP_REG,
			opcode.JEQ_LIT, opcode.JEQ_REG,
			opcode.JNE_LIT, opcode.JNE_REG,
			opcode.JLT_LIT, opcode.JLT_REG,
			opcode.JGT_LIT, opcode.JGT_REG,
			opcode.JLE_LIT, opcode.JLE_REG,
			opcode.JGE_LIT, opcode.JGE_REG:
			m.handleJump(thread, op)
		case opcode.PUSH_LIT:
			m.handlePushLit(thread)
		case opcode.PUSH_REG:
			m.handlePushReg(thread)
		case opcode.PUSH_AOF:
			m.handlePushAof(thread)
		case opcode.POP_REG:
			m.handlePopReg(thread)
		case opcode.POP_AOF:
			m.handlePopAof(thread)
		case opcode.CALL_LIT:
			m.handleCallLit(thread)
		case opcode.RET:
			m.handleRet(thread)
		default:
			log.Fatal("unhandled instruction", "op", op.String())
		}
	}
}

func (m *Machine) Run() {
	m.RunThread(m.mainThread)

}

func (m *Machine) decodeNumber(dataType string, index int) int {
	dataType = strings.ToLower(dataType)
	switch dataType {
	case "byte":
		return int(m.memory[index])
	case "word":
		bytes := m.memory[index : index+2]
		return int(binary.BigEndian.Uint16(bytes))
	case "dword":
		bytes := m.memory[index : index+4]
		return int(binary.BigEndian.Uint32(bytes))
	case "qword", "unset":
		bytes := m.memory[index : index+8]
		return int(binary.BigEndian.Uint64(bytes))
	default:
		log.Fatal("unknown data type", "type", dataType)
	}
	return -1
}

func (m *Machine) decodeNumberBytes(dataType string, index int) []byte {
	dataType = strings.ToLower(dataType)
	switch dataType {
	case "byte":
		return []byte{m.memory[index]}
	case "word":
		return m.memory[index : index+2]
	case "dword":
		return m.memory[index : index+4]
	case "qword", "unset":
		return m.memory[index : index+8]
	default:
		log.Fatal("unknown data type", "type", dataType)
	}
	return nil
}

func (m *Machine) decodeRegister(index int) int {
	v := m.memory[index]
	return int(v)
}

func (m *Machine) decodeValue(thread *Thread, dataType datatype.DataType) ast.Value {
	pos := m.position(thread)
	indexValue := m.memory[pos]
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

	switch indexValue {
	case 2:
		pos := m.position(thread)
		reg := m.decodeRegister(pos)
		m.incRegister(thread, utils.RegisterToIndex("ip"), 1)
		return &ast.Register{Value: reg}
	case 0, 1, 3, 4:
		pos := m.position(thread)
		num := m.decodeNumber(dataType.String(), pos)
		m.incRegister(thread, utils.RegisterToIndex("ip"), uint64(dataType.Size()))
		return &ast.NumberLiteral{Value: strconv.Itoa(num)}
	case 5:
		pos := m.position(thread)
		reg := m.decodeRegister(pos)
		m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

		pos = m.position(thread)
		op := ast.Operator(m.decodeNumber("byte", pos))
		m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

		pos = m.position(thread)
		offset := m.decodeNumber(dataType.String(), pos)
		m.incRegister(thread, utils.RegisterToIndex("ip"), uint64(dataType.Size()))

		return &ast.RegisterOffsetNumber{
			Left:     ast.Register{Value: reg},
			Operator: op,
			Right:    ast.NumberLiteral{Value: strconv.Itoa(offset)},
		}
	case 6:
		pos := m.position(thread)
		reg0 := m.decodeRegister(pos)
		m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

		pos = m.position(thread)
		op := ast.Operator(m.decodeNumber("byte", pos))
		m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

		pos = m.position(thread)
		reg1 := m.decodeRegister(pos)
		m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

		return &ast.RegisterOffsetRegister{
			Left:     ast.Register{Value: reg0},
			Operator: op,
			Right:    ast.Register{Value: reg1},
		}
	case 7:
		pos := m.position(thread)
		addr := m.decodeNumber(dataType.String(), pos)
		m.incRegister(thread, utils.RegisterToIndex("ip"), uint64(dataType.Size()))

		pos = m.position(thread)
		op := ast.Operator(m.decodeNumber("byte", pos))
		m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

		pos = m.position(thread)
		offset := m.decodeNumber(dataType.String(), pos)
		m.incRegister(thread, utils.RegisterToIndex("ip"), uint64(dataType.Size()))

		return &ast.LabelOffsetNumber{
			Left:     &ast.NumberLiteral{Value: strconv.Itoa(addr)},
			Operator: op,
			Right:    ast.NumberLiteral{Value: strconv.Itoa(offset)},
		}
	case 8:
		pos := m.position(thread)
		addr := m.decodeNumber(dataType.String(), pos)
		m.incRegister(thread, utils.RegisterToIndex("ip"), uint64(dataType.Size()))

		pos = m.position(thread)
		op := ast.Operator(m.decodeNumber("byte", pos))
		m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

		pos = m.position(thread)
		reg := m.decodeRegister(pos)
		m.incRegister(thread, utils.RegisterToIndex("ip"), 1)

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
	start := m.readLiteral(m.mainThread, datatype.QWORD)
	m.memory = m.memory[8:]
	m.setRegister(m.mainThread, utils.RegisterToIndex("ip"), uint64(start))
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

func (m *Machine) position(thread *Thread) int {
	index := utils.RegisterToIndex("ip")
	return int(m.getRegister(thread, index))
}

func (m *Machine) setRegister(thread *Thread, index int, value uint64) {
	thread.registers[index] = value
}

func (m *Machine) getRegister(thread *Thread, index int) uint64 {
	return thread.registers[index]
}

func (m *Machine) incRegister(thread *Thread, index int, amount uint64) {
	thread.registers[index] += amount
}

func (m *Machine) readRegister(thread *Thread) int {
	pos := m.position(thread)
	reg := m.decodeRegister(pos)
	m.incRegister(thread, utils.RegisterToIndex("ip"), 1)
	return reg
}

func (m *Machine) readLiteral(thread *Thread, dataType datatype.DataType) uint64 {
	pos := m.position(thread)
	lit := m.decodeNumber(dataType.String(), pos)
	m.incRegister(thread, utils.RegisterToIndex("ip"), uint64(dataType.Size()))
	return uint64(lit)
}

func (m *Machine) stackPush(thread *Thread, v []byte) {
	spIndex := utils.RegisterToIndex("sp")
	spValue := m.getRegister(thread, spIndex)

	// byteArray := utils.Bytes8(v)

	memIndex := int(spValue) - len(v)

	copy(m.memory[memIndex:memIndex+len(v)], v)

	m.setRegister(thread, spIndex, spValue-uint64(len(v)))
}

func (m *Machine) stackPopBytes(thread *Thread, dataType datatype.DataType) []byte {
	spIndex := utils.RegisterToIndex("sp")
	spValue := m.getRegister(thread, spIndex)

	memIndex := int(spValue)

	var value []byte
	switch dataType {
	case datatype.BYTE:
		value = []byte{m.memory[memIndex]}
	case datatype.WORD:
		value = m.memory[memIndex : memIndex+dataType.Size()]
	case datatype.DWORD:
		value = m.memory[memIndex : memIndex+dataType.Size()]
	case datatype.QWORD, datatype.UNSET:
		value = m.memory[memIndex : memIndex+dataType.Size()]
	default:
		log.Fatal("unknown data type", "type", dataType)
	}

	m.setRegister(thread, spIndex, spValue+uint64(dataType.Size()))

	return value
}

func (m *Machine) stackPop(thread *Thread, dataType datatype.DataType) uint64 {
	spIndex := utils.RegisterToIndex("sp")
	spValue := m.getRegister(thread, spIndex)

	memIndex := int(spValue)

	var value uint64
	switch dataType {
	case datatype.BYTE:
		value = uint64(m.memory[memIndex])
	case datatype.WORD:
		value = uint64(binary.BigEndian.Uint16(m.memory[memIndex : memIndex+dataType.Size()]))
	case datatype.DWORD:
		value = uint64(binary.BigEndian.Uint32(m.memory[memIndex : memIndex+dataType.Size()]))
	case datatype.QWORD, datatype.UNSET:
		value = binary.BigEndian.Uint64(m.memory[memIndex : memIndex+dataType.Size()])
	default:
		log.Fatal("unknown data type", "type", dataType)
	}

	m.setRegister(thread, spIndex, spValue+uint64(dataType.Size()))

	return value
}

func (m *Machine) DumpRegisters(index int) {
	for i, register := range m.threads[index].registers {
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
