package vm

import (
	"bytes"
	"encoding/binary"
	"fishy/pkg/utils"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
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

type SyscallFunction func(m *Machine, thread *Thread)

const (
	SYS_EXIT SyscallIndex = iota + 1
	SYS_OPEN
	SYS_READ
	SYS_WRITE
	SYS_CLOSE
	SYS_STRERR
	SYS_INT_TO_STR
	SYS_STR_TO_INT
	SYS_CLOCK

	SYS_NET_LISTEN_TCP
	SYS_NET_CONNECT_TCP
	SYS_NET_ACCEPT
	SYS_NET_GETPEERNAME
	SYS_NET_IP_TO_STR

	SYS_THREAD_SPAWN
	SYS_THREAD_START
	SYS_THREAD_STOP
	SYS_THREAD_JOIN
)

func (m *Machine) handleSyscall(thread *Thread) {
	var Syscalls = map[SyscallIndex]SyscallFunction{
		SYS_EXIT: func(m *Machine, thread *Thread) {
			status := m.getRegister(thread, utils.RegisterToIndex("x0"))
			os.Exit(int(status))
		},
		SYS_OPEN: func(m *Machine, thread *Thread) {
			addr := m.getRegister(thread, utils.RegisterToIndex("x0"))
			length := m.getRegister(thread, utils.RegisterToIndex("x1"))
			mode := m.getRegister(thread, utils.RegisterToIndex("x2"))
			perm := m.getRegister(thread, utils.RegisterToIndex("x3"))

			n := -1
			if addr >= uint64(len(m.memory)) || addr+length > uint64(len(m.memory)) {
				m.SetErrorCodeRegister(thread, EADDROUTOFBOUNDS)
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			if length == 0 {
				m.SetErrorCodeRegister(thread, EINVALIDLENGTH)
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			path := m.memory[addr : addr+length]

			fd, err := syscall.Open(string(path), int(mode), uint32(perm))
			if err != nil {
				m.SetErrorCodeRegister(thread, MatchString(err.Error()))
			}

			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(fd))
		},
		SYS_READ: func(m *Machine, thread *Thread) {
			fd := m.getRegister(thread, utils.RegisterToIndex("x0"))
			addr := m.getRegister(thread, utils.RegisterToIndex("x1"))
			length := m.getRegister(thread, utils.RegisterToIndex("x2"))

			buffer := make([]byte, length)
			n, err := syscall.Read(int(fd), buffer)
			if err != nil {
				m.SetErrorCodeRegister(thread, MatchString(err.Error()))
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			for i := 0; i < int(length); i++ {
				m.memory[int(addr)+i] = buffer[i]
			}

			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
		},
		SYS_WRITE: func(m *Machine, thread *Thread) {
			fd := m.getRegister(thread, utils.RegisterToIndex("x0"))
			addr := m.getRegister(thread, utils.RegisterToIndex("x1"))
			length := m.getRegister(thread, utils.RegisterToIndex("x2"))

			n := -1
			start := int(addr)
			end := start + int(length)

			if start >= len(m.memory) || end > len(m.memory) {
				m.SetErrorCodeRegister(thread, EADDROUTOFBOUNDS)
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			if start < 0 || start >= end {
				m.SetErrorCodeRegister(thread, EINVALIDLENGTH)
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			buffer := m.memory[start:end]

			n, err := syscall.Write(int(fd), buffer)
			if err != nil {
				m.SetErrorCodeRegister(thread, MatchString(err.Error()))
			}

			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
		},
		SYS_CLOSE: func(m *Machine, thread *Thread) {
			fd := m.getRegister(thread, utils.RegisterToIndex("x0"))
			err := syscall.Close(int(fd))
			n := 0
			if err != nil {
				n = -1
				m.SetErrorCodeRegister(thread, MatchString(err.Error()))
			}
			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
		},
		SYS_STRERR: func(m *Machine, thread *Thread) {
			er := m.getRegister(thread, utils.RegisterToIndex("er"))
			addr := m.getRegister(thread, utils.RegisterToIndex("x0"))
			length := m.getRegister(thread, utils.RegisterToIndex("x1"))

			code := ErrorCode(er)
			message := []byte(code.String())
			n := -1

			if addr >= uint64(len(m.memory)) {
				m.SetErrorCodeRegister(thread, EADDROUTOFBOUNDS)
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			if length == 0 || addr+length > uint64(len(m.memory)) {
				m.SetErrorCodeRegister(thread, EINVALIDLENGTH)
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			copy(m.memory[addr:addr+length], message[:])

			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(len(message)))
		},
		SYS_THREAD_SPAWN: func(m *Machine, thread *Thread) {
			startAddr := m.getRegister(thread, utils.RegisterToIndex("x0"))
			stackOffset := m.getRegister(thread, utils.RegisterToIndex("x1"))

			n := -1
			currentPos := m.getRegister(m.mainThread, utils.RegisterToIndex("sp"))

			newThread := m.CreateThread()
			threadIndex, ok := m.GetThreadIndex(newThread)
			if !ok {
				m.SetErrorCodeRegister(thread, MatchString("failed to spawn thread"))
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			m.setRegister(newThread, utils.RegisterToIndex("sp"), currentPos-stackOffset)
			m.setRegister(newThread, utils.RegisterToIndex("fp"), currentPos-stackOffset)
			m.setRegister(newThread, utils.RegisterToIndex("ip"), startAddr)

			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(threadIndex))
		},
		SYS_THREAD_START: func(m *Machine, thread *Thread) {
			threadIndex := m.getRegister(thread, utils.RegisterToIndex("x0"))

			n := -1
			if workingThread, ok := m.threads[int(threadIndex)]; ok {
				m.wg.Add(1)
				go func() {
					m.RunThread(workingThread)
					defer func() {
						m.wg.Done()
					}()
				}()
			} else {
				m.SetErrorCodeRegister(thread, MatchString("failed to get thread"))
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
			}

			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(0))
		},
		SYS_THREAD_STOP: func(m *Machine, thread *Thread) {
			threadIndex := m.getRegister(thread, utils.RegisterToIndex("x0"))

			n := -1
			if workingThread, ok := m.threads[int(threadIndex)]; ok {
				workingThread.isRunning = false
			} else {
				m.SetErrorCodeRegister(thread, MatchString("failed to get thread"))
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
			}

			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(0))
		},
		SYS_THREAD_JOIN: func(m *Machine, thread *Thread) {
			threadIndex := m.getRegister(thread, utils.RegisterToIndex("x0"))

			n := -1
			if workingThread, ok := m.threads[int(threadIndex)]; ok {
				<-workingThread.done
			} else {
				m.SetErrorCodeRegister(thread, MatchString("failed to get thread"))
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
			}

			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(0))
		},
		SYS_INT_TO_STR: func(m *Machine, thread *Thread) {
			number := m.getRegister(thread, utils.RegisterToIndex("x0"))
			addr := m.getRegister(thread, utils.RegisterToIndex("x1"))
			length := m.getRegister(thread, utils.RegisterToIndex("x2"))

			n := -1

			if addr >= uint64(len(m.memory)) || addr+length > uint64(len(m.memory)) {
				m.SetErrorCodeRegister(thread, EADDROUTOFBOUNDS)
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			if length == 0 {
				m.SetErrorCodeRegister(thread, EINVALIDLENGTH)
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			str := strconv.Itoa(int(number))
			copy(m.memory[addr:addr+length], str[:])

			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(len(str)))
		},
		SYS_STR_TO_INT: func(m *Machine, thread *Thread) {
			numberAddr := m.getRegister(thread, utils.RegisterToIndex("x0"))
			numberLength := m.getRegister(thread, utils.RegisterToIndex("x1"))
			returnAddr := m.getRegister(thread, utils.RegisterToIndex("x2"))

			n := -1

			if numberAddr >= uint64(len(m.memory)) || numberAddr+numberLength > uint64(len(m.memory)) ||
				returnAddr >= uint64(len(m.memory)) || returnAddr+8 > uint64(len(m.memory)) {
				m.SetErrorCodeRegister(thread, EADDROUTOFBOUNDS)
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			if numberLength == 0 || numberLength > numberAddr {
				m.SetErrorCodeRegister(thread, EINVALIDLENGTH)
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			numberBuffer := string(m.memory[numberAddr : numberAddr+numberLength])
			number, err := strconv.ParseUint(numberBuffer, 10, 64)
			if err != nil {
				m.SetErrorCodeRegister(thread, MatchString(err.Error()))
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			numberBytes := utils.Bytes8(number)

			copy(m.memory[returnAddr:returnAddr+8], numberBytes[:])

			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(0))
		},
		SYS_CLOCK: func(m *Machine, thread *Thread) {
			millis := time.Now().UnixMilli()
			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(millis))
		},
		SYS_NET_LISTEN_TCP: func(m *Machine, thread *Thread) {
			listenOptsAddr := m.getRegister(thread, utils.RegisterToIndex("x0"))
			listenOptsData := m.memory[listenOptsAddr : listenOptsAddr+7]

			n := -1
			var listenOpts SocketListenTcpOpts
			buffer := bytes.NewReader(listenOptsData)
			err := binary.Read(buffer, binary.BigEndian, &listenOpts)
			if err != nil {
				m.SetErrorCodeRegister(thread, MatchString(err.Error()))
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
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
				m.SetErrorCodeRegister(thread, MatchString(err.Error()))
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}
			file, err := ln.File()
			if err != nil {
				m.SetErrorCodeRegister(thread, MatchString(err.Error()))
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(file.Fd()))
		},
		SYS_NET_CONNECT_TCP: func(m *Machine, thread *Thread) {
			connectOptsAddr := m.getRegister(thread, utils.RegisterToIndex("x0"))
			connectOptsData := m.memory[connectOptsAddr : connectOptsAddr+7]

			n := -1
			var connectOpts SocketConnectTcpOpts
			buffer := bytes.NewReader(connectOptsData)
			err := binary.Read(buffer, binary.BigEndian, &connectOpts)
			if err != nil {
				m.SetErrorCodeRegister(thread, MatchString(err.Error()))
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
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
				m.SetErrorCodeRegister(thread, MatchString(err.Error()))
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}
			file, err := dial.File()
			if err != nil {
				m.SetErrorCodeRegister(thread, MatchString(err.Error()))
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(file.Fd()))
		},
		SYS_NET_ACCEPT: func(m *Machine, thread *Thread) {
			fd := m.getRegister(thread, utils.RegisterToIndex("x0"))

			conn, _, err := syscall.Accept(int(fd))
			if err != nil {
				m.SetErrorCodeRegister(thread, MatchString(err.Error()))
			}

			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(conn))
		},
		SYS_NET_GETPEERNAME: func(m *Machine, thread *Thread) {
			fd := m.getRegister(thread, utils.RegisterToIndex("x0"))
			returnAddr := m.getRegister(thread, utils.RegisterToIndex("x1"))

			n := -1
			sa, err := syscall.Getpeername(int(fd))
			if err != nil {
				m.SetErrorCodeRegister(thread, MatchString(err.Error()))
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			sockAddr := sa.(*syscall.SockaddrInet4)

			result := make([]byte, 6)
			copy(result[0:4], sockAddr.Addr[:])
			result[4] = byte(sockAddr.Port >> 8)
			result[5] = byte(sockAddr.Port & 0xff)

			copy(m.memory[int(returnAddr):], result)

			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(0))
		},
		SYS_NET_IP_TO_STR: func(m *Machine, thread *Thread) {
			ipAddr := m.getRegister(thread, utils.RegisterToIndex("x0"))
			returnAddr := m.getRegister(thread, utils.RegisterToIndex("x1"))
			length := m.getRegister(thread, utils.RegisterToIndex("x2"))

			n := -1
			if ipAddr >= uint64(len(m.memory)) ||
				ipAddr+4 > uint64(len(m.memory)) ||
				returnAddr >= uint64(len(m.memory)) ||
				returnAddr+length > uint64(len(m.memory)) {
				m.SetErrorCodeRegister(thread, EADDROUTOFBOUNDS)
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			if length == 0 {
				m.SetErrorCodeRegister(thread, EINVALIDLENGTH)
				m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(n))
				return
			}

			ipBuffer := m.memory[ipAddr : ipAddr+4]

			ip := net.IP(ipBuffer).String()

			copy(m.memory[returnAddr:returnAddr+length], ip[:])

			m.setRegister(thread, utils.RegisterToIndex("x0"), uint64(len(ip)))
		},
	}

	m.incRegister(thread, utils.RegisterToIndex("ip"), 2)

	index := m.getRegister(thread, utils.RegisterToIndex("x15"))
	sc := SyscallIndex(index)

	if call, ok := Syscalls[sc]; ok {
		call(m, thread)
	} else {
		m.SetErrorCodeRegister(thread, UNKNOWN_SYSCAlL)
	}
}

func (m *Machine) SetErrorCodeRegister(thread *Thread, code ErrorCode) {
	m.setRegister(thread, utils.RegisterToIndex("er"), uint64(code))
}
