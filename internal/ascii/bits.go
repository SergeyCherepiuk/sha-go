package ascii

import (
	"fmt"
)

func ToBits(s string) []byte {
	bits := make([]byte, len(s)*8)

	for i, ch := range s {
		chBits := byteToBits(byte(ch))
		start, end := i*8, (i+1)*8
		copy(bits[start:end], chBits[:])
	}

	return bits
}

func FromBits(bits []uint8) string {
	bytes := make([]byte, len(bits)/8)

	for i := 0; i < len(bits)/8; i++ {
		var chBits [8]uint8
		start, end := i*8, (i+1)*8
		copy(chBits[:], bits[start:end])
		bytes[i] = bitsToByte(chBits)
	}

	return string(bytes)
}

func byteToBits(ch byte) [8]byte {
	bits := [8]byte{}
	for i, r := range fmt.Sprintf("%.8b", ch) {
		bits[i] = byte(r)
	}
	return bits
}

func bitsToByte(bits [8]byte) byte {
	degrees := [8]uint8{128, 64, 32, 16, 8, 4, 2, 1}

	var sum byte
	for i, b := range bits {
		sum += degrees[i] * (b - 48)
	}

	return sum
}
