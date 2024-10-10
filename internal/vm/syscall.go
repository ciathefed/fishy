package vm

import (
	"bytes"
	"encoding/binary"
	"fishy/pkg/utils"
	"net"
	"os"
	"strconv"
	"syscall"
)

type SocketListenTcpOpts struct {
	Type    uint8
	Address [4]byte
	Port    uint16
}

type SocketConnectTcpOpts struct {
	Type    uint8
	Address [4]byte
	Port    uint16
}

type SyscallIndex int

type SyscallFunction func(m *Machine)

const (
	SYS_EXIT SyscallIndex = iota + 1
	SYS_OPEN
	SYS_READ
	SYS_WRITE
	SYS_CLOSE
	SYS_STRERR
	SYS_INT_TO_STR
	SYS_NET_LISTEN_TCP
	SYS_NET_CONNECT_TCP
	SYS_NET_ACCEPT
	SYS_NET_GETPEERNAME
	SYS_NET_IP_TO_STR
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

		n := -1
		if addr >= uint64(len(m.memory)) || addr+length > uint64(len(m.memory)) {
			m.SetErrorCodeRegister(EADDROUTOFBOUNDS)
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		if length == 0 {
			m.SetErrorCodeRegister(EINVALIDLENGTH)
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		path := m.memory[addr : addr+length]

		fd, err := syscall.Open(string(path), int(mode), uint32(perm))
		if err != nil {
			m.SetErrorCodeRegister(MatchString(err.Error()))
		}

		m.setRegister(utils.RegisterToIndex("x0"), uint64(fd))
	},
	SYS_READ: func(m *Machine) {
		fd := m.getRegister(utils.RegisterToIndex("x0"))
		addr := m.getRegister(utils.RegisterToIndex("x1"))
		length := m.getRegister(utils.RegisterToIndex("x2"))

		buffer := make([]byte, length)
		n, err := syscall.Read(int(fd), buffer)
		if err != nil {
			m.SetErrorCodeRegister(MatchString(err.Error()))
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		for i := 0; i < int(length); i++ {
			m.memory[int(addr)+i] = buffer[i]
		}

		m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
	},
	SYS_WRITE: func(m *Machine) {
		fd := m.getRegister(utils.RegisterToIndex("x0"))
		addr := m.getRegister(utils.RegisterToIndex("x1"))
		length := m.getRegister(utils.RegisterToIndex("x2"))

		n := -1
		start := int(addr)
		end := start + int(length)

		if start >= len(m.memory) || end > len(m.memory) {
			m.SetErrorCodeRegister(EADDROUTOFBOUNDS)
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		if start < 0 || start >= end {
			m.SetErrorCodeRegister(EINVALIDLENGTH)
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		buffer := m.memory[start:end]

		n, err := syscall.Write(int(fd), buffer)
		if err != nil {
			m.SetErrorCodeRegister(MatchString(err.Error()))
		}

		m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
	},
	SYS_CLOSE: func(m *Machine) {
		fd := m.getRegister(utils.RegisterToIndex("x0"))
		err := syscall.Close(int(fd))
		n := 0
		if err != nil {
			n = -1
			m.SetErrorCodeRegister(MatchString(err.Error()))
		}
		m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
	},
	SYS_STRERR: func(m *Machine) {
		er := m.getRegister(utils.RegisterToIndex("er"))
		addr := m.getRegister(utils.RegisterToIndex("x0"))
		length := m.getRegister(utils.RegisterToIndex("x1"))

		code := ErrorCode(er)
		message := []byte(code.String())
		n := -1

		if addr >= uint64(len(m.memory)) {
			m.SetErrorCodeRegister(EADDROUTOFBOUNDS)
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		if length == 0 || addr+length > uint64(len(m.memory)) {
			m.SetErrorCodeRegister(EINVALIDLENGTH)
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		copy(m.memory[addr:addr+length], message[:])

		m.setRegister(utils.RegisterToIndex("x0"), uint64(len(message)))
	},
	SYS_INT_TO_STR: func(m *Machine) {
		number := m.getRegister(utils.RegisterToIndex("x0"))
		addr := m.getRegister(utils.RegisterToIndex("x1"))
		length := m.getRegister(utils.RegisterToIndex("x2"))

		str := strconv.Itoa(int(number))
		n := -1

		if addr >= uint64(len(m.memory)) {
			m.SetErrorCodeRegister(EADDROUTOFBOUNDS)
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		if length == 0 || addr+length > uint64(len(m.memory)) {
			m.SetErrorCodeRegister(EINVALIDLENGTH)
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		copy(m.memory[addr:addr+length], str[:])

		m.setRegister(utils.RegisterToIndex("x0"), uint64(len(str)))
	},
	SYS_NET_LISTEN_TCP: func(m *Machine) {
		listenOptsAddr := m.getRegister(utils.RegisterToIndex("x0"))
		listenOptsData := m.memory[listenOptsAddr : listenOptsAddr+7]

		n := -1
		var listenOpts SocketListenTcpOpts
		buffer := bytes.NewReader(listenOptsData)
		err := binary.Read(buffer, binary.BigEndian, &listenOpts)
		if err != nil {
			m.SetErrorCodeRegister(MatchString(err.Error()))
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		address := net.IPv4(listenOpts.Address[0], listenOpts.Address[1], listenOpts.Address[2], listenOpts.Address[3])
		network := "tcp"
		switch listenOpts.Type {
		case 0:
			network = "tcp"
		case 1:
			network = "tcp4"
		case 2:
			network = "tcp6"
		}

		ln, err := net.ListenTCP(network, &net.TCPAddr{
			IP:   address,
			Port: int(listenOpts.Port),
		})
		if err != nil {
			m.SetErrorCodeRegister(MatchString(err.Error()))
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}
		file, err := ln.File()
		if err != nil {
			m.SetErrorCodeRegister(MatchString(err.Error()))
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		m.setRegister(utils.RegisterToIndex("x0"), uint64(file.Fd()))
	},
	SYS_NET_CONNECT_TCP: func(m *Machine) {
		connectOptsAddr := m.getRegister(utils.RegisterToIndex("x0"))
		connectOptsData := m.memory[connectOptsAddr : connectOptsAddr+7]

		n := -1
		var connectOpts SocketConnectTcpOpts
		buffer := bytes.NewReader(connectOptsData)
		err := binary.Read(buffer, binary.BigEndian, &connectOpts)
		if err != nil {
			m.SetErrorCodeRegister(MatchString(err.Error()))
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		address := net.IPv4(connectOpts.Address[0], connectOpts.Address[1], connectOpts.Address[2], connectOpts.Address[3])
		network := "tcp"
		switch connectOpts.Type {
		case 0:
			network = "tcp"
		case 1:
			network = "tcp4"
		case 2:
			network = "tcp6"
		}

		dial, err := net.DialTCP(network, nil, &net.TCPAddr{
			IP:   address,
			Port: int(connectOpts.Port),
		})
		if err != nil {
			m.SetErrorCodeRegister(MatchString(err.Error()))
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}
		file, err := dial.File()
		if err != nil {
			m.SetErrorCodeRegister(MatchString(err.Error()))
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		m.setRegister(utils.RegisterToIndex("x0"), uint64(file.Fd()))
	},
	SYS_NET_ACCEPT: func(m *Machine) {
		fd := m.getRegister(utils.RegisterToIndex("x0"))

		conn, _, err := syscall.Accept(int(fd))
		if err != nil {
			m.SetErrorCodeRegister(MatchString(err.Error()))
		}

		m.setRegister(utils.RegisterToIndex("x0"), uint64(conn))
	},
	SYS_NET_GETPEERNAME: func(m *Machine) {
		fd := m.getRegister(utils.RegisterToIndex("x0"))
		returnAddr := m.getRegister(utils.RegisterToIndex("x1"))

		n := -1
		sa, err := syscall.Getpeername(int(fd))
		if err != nil {
			m.SetErrorCodeRegister(MatchString(err.Error()))
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		sockAddr := sa.(*syscall.SockaddrInet4)

		result := make([]byte, 6)
		copy(result[0:4], sockAddr.Addr[:])
		result[4] = byte(sockAddr.Port >> 8)
		result[5] = byte(sockAddr.Port & 0xff)

		copy(m.memory[int(returnAddr):], result)
	},
	SYS_NET_IP_TO_STR: func(m *Machine) {
		ipAddr := m.getRegister(utils.RegisterToIndex("x0"))
		returnAddr := m.getRegister(utils.RegisterToIndex("x1"))
		length := m.getRegister(utils.RegisterToIndex("x2"))

		n := -1
		if ipAddr >= uint64(len(m.memory)) ||
			ipAddr+4 > uint64(len(m.memory)) ||
			returnAddr >= uint64(len(m.memory)) ||
			returnAddr+length > uint64(len(m.memory)) {
			m.SetErrorCodeRegister(EADDROUTOFBOUNDS)
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		if length == 0 {
			m.SetErrorCodeRegister(EINVALIDLENGTH)
			m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
			return
		}

		ipBuffer := m.memory[ipAddr : ipAddr+4]

		ip := net.IP(ipBuffer).String()

		copy(m.memory[returnAddr:returnAddr+length], ip[:])

		m.setRegister(utils.RegisterToIndex("x0"), uint64(len(ip)))
	},
}

func (m *Machine) handleSyscall() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	index := m.getRegister(utils.RegisterToIndex("x15"))
	sc := SyscallIndex(index)

	if call, ok := Syscalls[sc]; ok {
		call(m)
	} else {
		m.SetErrorCodeRegister(UNKNOWN_SYSCAlL)
	}
}

func (m *Machine) SetErrorCodeRegister(code ErrorCode) {
	m.setRegister(utils.RegisterToIndex("er"), uint64(code))
}
