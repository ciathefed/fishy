package vm

import (
	"fishy/pkg/utils"
	"os"
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
	SYS_WRITE: func(m *Machine) {
		fd := m.getRegister(utils.RegisterToIndex("x0"))
		addr := m.getRegister(utils.RegisterToIndex("x1"))
		length := m.getRegister(utils.RegisterToIndex("x2"))

		start := int(addr)
		end := start + int(length)
		buffer := m.memory[start:end]

		file := os.NewFile(uintptr(fd), "pipe")
		n, err := file.Write(buffer)
		if err != nil {
			panic(err)
		}

		m.setRegister(utils.RegisterToIndex("x0"), uint32(n))
	},
}
