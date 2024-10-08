package vm

import (
	"bytes"
	"encoding/binary"
	"fishy/pkg/log"
	"fishy/pkg/utils"
	"fmt"
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
	SYS_NET_LISTEN_TCP
	SYS_NET_CONNECT_TCP
	SYS_NET_ACCEPT
	SYS_NET_GETPEERNAME
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

		fd, err := syscall.Open(string(path), int(mode), uint32(perm))
		if err != nil {
			log.Fatal(err, "syscall", SYS_OPEN)
			// fd = -1
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
			log.Fatal(err, "syscall", SYS_READ)
			// n = -1
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

		start := int(addr)
		end := start + int(length)
		buffer := m.memory[start:end]

		n, err := syscall.Write(int(fd), buffer)
		if err != nil {
			log.Fatal(err, "syscall", SYS_WRITE)
			// n = -1
		}

		m.setRegister(utils.RegisterToIndex("x0"), uint64(n))
	},
	SYS_CLOSE: func(m *Machine) {
		fd := m.getRegister(utils.RegisterToIndex("x0"))
		err := syscall.Close(int(fd))
		if err != nil {
			log.Fatal(err, "syscall", SYS_CLOSE)
		}
	},
	SYS_NET_LISTEN_TCP: func(m *Machine) {
		listenOptsAddr := m.getRegister(utils.RegisterToIndex("x0"))
		listenOptsData := m.memory[listenOptsAddr : listenOptsAddr+7]

		var listenOpts SocketListenTcpOpts
		buffer := bytes.NewReader(listenOptsData)
		err := binary.Read(buffer, binary.BigEndian, &listenOpts)
		if err != nil {
			log.Fatal(err, "syscall", SYS_NET_LISTEN_TCP)
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
			log.Fatal(err, "syscall", SYS_NET_LISTEN_TCP)
		}
		file, err := ln.File()
		if err != nil {
			log.Fatal(err, "syscall", SYS_NET_LISTEN_TCP)
		}

		m.setRegister(utils.RegisterToIndex("x0"), uint64(file.Fd()))
	},
	SYS_NET_CONNECT_TCP: func(m *Machine) {
		connectOptsAddr := m.getRegister(utils.RegisterToIndex("x0"))
		connectOptsData := m.memory[connectOptsAddr : connectOptsAddr+7]

		var connectOpts SocketConnectTcpOpts
		buffer := bytes.NewReader(connectOptsData)
		err := binary.Read(buffer, binary.BigEndian, &connectOpts)
		if err != nil {
			log.Fatal(err, "syscall", SYS_NET_CONNECT_TCP)
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
			log.Fatal(err, "syscall", SYS_NET_CONNECT_TCP)
		}
		file, err := dial.File()
		if err != nil {
			log.Fatal(err, "syscall", SYS_NET_CONNECT_TCP)
		}

		m.setRegister(utils.RegisterToIndex("x0"), uint64(file.Fd()))
	},
	SYS_NET_ACCEPT: func(m *Machine) {
		fd := m.getRegister(utils.RegisterToIndex("x0"))

		conn, _, err := syscall.Accept(int(fd))
		if err != nil {
			log.Fatal(err, "syscall", SYS_NET_ACCEPT)
		}

		m.setRegister(utils.RegisterToIndex("x0"), uint64(conn))
	},
	SYS_NET_GETPEERNAME: func(m *Machine) {
		fd := m.getRegister(utils.RegisterToIndex("x0"))
		returnAddr := m.getRegister(utils.RegisterToIndex("x1"))

		file := os.NewFile(uintptr(fd), "")
		if file == nil {
			log.Fatal(fmt.Errorf("failed to create file from fd"), "syscall", SYS_NET_GETPEERNAME)
		}
		conn, err := net.FileConn(file)
		if err != nil {
			log.Fatal(err, "syscall", SYS_NET_GETPEERNAME)
		}

		remoteAddr := conn.RemoteAddr()
		ipAddr, portStr, err := net.SplitHostPort(remoteAddr.String())
		if err != nil {
			log.Fatal(err, "syscall", SYS_NET_GETPEERNAME)
		}
		port, err := strconv.ParseInt(portStr, 10, 32)
		if err != nil {
			log.Fatal(err, "syscall", SYS_NET_GETPEERNAME)
		}

		finalAddr := net.ParseIP(ipAddr)
		if finalAddr == nil {
			log.Fatal(fmt.Errorf("failed to parse IP address"), "syscall", SYS_NET_GETPEERNAME)
		}

		result := make([]byte, 6)
		copy(result[0:4], finalAddr.To4())
		result[4] = byte(port >> 8)
		result[5] = byte(port & 0xff)

		copy(m.memory[int(returnAddr):], result)
	},
}

func (m *Machine) handleSyscall() {
	m.incRegister(utils.RegisterToIndex("ip"), 2)

	index := m.getRegister(utils.RegisterToIndex("x15"))
	sc := SyscallIndex(index)

	if call, ok := Syscalls[sc]; ok {
		call(m)
	} else {
		log.Fatal("unknown syscall", "index", sc)
	}
}
