package utils

func Bytes2(value uint16) []byte {
	return []byte{
		byte(value >> 8),   // High byte
		byte(value & 0xFF), // Low byte
	}
}

func Bytes4(value uint32) []byte {
	return []byte{
		byte(value >> 24),  // Byte 3
		byte(value >> 16),  // Byte 2
		byte(value >> 8),   // Byte 1
		byte(value & 0xFF), // Byte 0
	}
}

func Bytes8(value uint64) []byte {
	return []byte{
		byte(value >> 56),  // Byte 7
		byte(value >> 48),  // Byte 6
		byte(value >> 40),  // Byte 5
		byte(value >> 32),  // Byte 4
		byte(value >> 24),  // Byte 3
		byte(value >> 16),  // Byte 2
		byte(value >> 8),   // Byte 1
		byte(value & 0xFF), // Byte 0
	}
}
