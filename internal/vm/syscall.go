package vm

import (
	"fishy/pkg/utils"
	"fmt"
	"os"
	"syscall"
)

type SyscallIndex int

type SyscallFunction func(m *Machine)

const (
	SYS_EXIT SyscallIndex = iota + 1
	SYS_OPEN
	SYS_READ
	SYS_WRITE
	SYS_CLOSE
)

var Syscalls = map[SyscallIndex]SyscallFunction{
	SYS_EXIT: func(m *Machine) {
		status := m.getRegister(utils.RegisterToIndex("x0"))
		os.Exit(int(status))
	},
	SYS_OPEN: func(m *Machine) {
		addr := m.getRegister(utils.RegisterToIndex("x0"))
		length := m.getRegister(utils.RegisterToIndex("x1"))
		mode := m.getRegister(utils.RegisterToIndex("x2"))
		perm := m.getRegister(utils.RegisterToIndex("x3"))

		path := m.memory[addr : addr+length]

		fd, err := syscall.Open(string(path), int(mode), perm)
		if err != nil {
			panic(err)
			// fd = -1
		}

		m.setRegister(utils.RegisterToIndex("x0"), uint32(fd))
	},
	SYS_READ: func(m *Machine) {
		fd := m.getRegister(utils.RegisterToIndex("x0"))
		addr := m.getRegister(utils.RegisterToIndex("x1"))
		length := m.getRegister(utils.RegisterToIndex("x2"))

		buffer := make([]byte, length)
		n, err := syscall.Read(int(fd), buffer)
		if err != nil {
			panic(err)
			// n = -1
		}

		for i := 0; i < int(length); i++ {
			m.memory[int(addr)+i] = buffer[i]
		}

		m.setRegister(utils.RegisterToIndex("x0"), uint32(n))
	},
	SYS_WRITE: func(m *Machine) {
		fd := m.getRegister(utils.RegisterToIndex("x0"))
		addr := m.getRegister(utils.RegisterToIndex("x1"))
		length := m.getRegister(utils.RegisterToIndex("x2"))

		start := int(addr)
		end := start + int(length)
		buffer := m.memory[start:end]

		n, err := syscall.Write(int(fd), buffer)
		if err != nil {
			panic(err)
			// n = -1
		}

		m.setRegister(utils.RegisterToIndex("x0"), uint32(n))
	},
	SYS_CLOSE: func(m *Machine) {
		fd := m.getRegister(utils.RegisterToIndex("x0"))
		syscall.Close(int(fd))
	},
}

func (m *Machine) handleSyscall() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	index := m.getRegister(utils.RegisterToIndex("x15"))
	sc := SyscallIndex(index)

	if call, ok := Syscalls[sc]; ok {
		call(m)
	} else {
		panic(fmt.Sprintf("unknown syscall: %d", sc))
	}
}
