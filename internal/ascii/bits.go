package ascii

import (
	"fmt"
	"strconv"
)

func ToBits(s string) []uint8 {
	bits := make([]uint8, len(s)*8)

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

func byteToBits(ch byte) [8]uint8 {
	var bits [8]uint8

	bitsString := fmt.Sprintf("%.8b", ch)
	for i, ch := range bitsString {
		bit, _ := strconv.ParseUint(string(ch), 2, 8)
		bits[i] = uint8(bit)
	}

	return bits
}

func bitsToByte(bits [8]uint8) byte {
	degrees := [8]uint8{128, 64, 32, 16, 8, 4, 2, 1}

	var sum byte
	for i, b := range bits {
		sum += degrees[i] * b
	}

	return sum
}
