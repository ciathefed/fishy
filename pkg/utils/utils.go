package utils

var Instructions = []string{
	".section",

	"nop",
	"hlt",
	"brk",
	"syscall",

	"mov",
	"add", "sub", "mul", "div",
	"cmp",
	"jmp", "jeq", "jne", "jlt", "jgt", "jle", "jge", "jz",
	"push", "pop",
}

var Sequences = []string{
	"db",
}

var Registers = []string{
	"x0", "x1", "x2", "x3", "x4", "x5", "x6", "x7", "x8", "x9", "x10", "x11", "x12", "x13", "x14", "x15",
	"ip", "fp", "sp", "cp", "er",
}

func RegisterToIndex(name string) int {
	for i, register := range Registers {
		if name == register {
			return i
		}
	}
	return -1
}

func IndexToRegister(index int) string {
	for i, register := range Registers {
		if i == index {
			return register
		}
	}
	return ""
}

func Bytes2(value uint16) []byte {
	return []byte{
		byte(value >> 8),
		byte(value & 0xFF),
	}
}

func Bytes4(value uint32) []byte {
	return []byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value & 0xFF),
	}
}

func Bytes8(value uint64) []byte {
	return []byte{
		byte(value >> 56),
		byte(value >> 48),
		byte(value >> 40),
		byte(value >> 32),
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value & 0xFF),
	}
}
