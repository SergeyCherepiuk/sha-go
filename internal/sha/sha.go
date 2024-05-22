package sha

import (
	"bytes"
	"fmt"
	"math"
	"math/bits"

	"github.com/SergeyCherepiuk/sha-go/internal/ascii"
)

const (
	ChunkSize         = 512
	LengthSegmentSize = 64
)

var (
	H0 uint32 = 0b01101010000010011110011001100111
	H1 uint32 = 0b10111011011001111010111010000101
	H2 uint32 = 0b00111100011011101111001101110010
	H3 uint32 = 0b10100101010011111111010100111010
	H4 uint32 = 0b01010001000011100101001001111111
	H5 uint32 = 0b10011011000001010110100010001100
	H6 uint32 = 0b00011111100000111101100110101011
	H7 uint32 = 0b01011011111000001100110100011001
)

var K = []uint32{
	0b01000010100010100010111110011000,
	0b01110001001101110100010010010001,
	0b10110101110000001111101111001111,
	0b11101001101101011101101110100101,
	0b00111001010101101100001001011011,
	0b01011001111100010001000111110001,
	0b10010010001111111000001010100100,
	0b10101011000111000101111011010101,
	0b11011000000001111010101010011000,
	0b00010010100000110101101100000001,
	0b00100100001100011000010110111110,
	0b01010101000011000111110111000011,
	0b01110010101111100101110101110100,
	0b10000000110111101011000111111110,
	0b10011011110111000000011010100111,
	0b11000001100110111111000101110100,
	0b11100100100110110110100111000001,
	0b11101111101111100100011110000110,
	0b00001111110000011001110111000110,
	0b00100100000011001010000111001100,
	0b00101101111010010010110001101111,
	0b01001010011101001000010010101010,
	0b01011100101100001010100111011100,
	0b01110110111110011000100011011010,
	0b10011000001111100101000101010010,
	0b10101000001100011100011001101101,
	0b10110000000000110010011111001000,
	0b10111111010110010111111111000111,
	0b11000110111000000000101111110011,
	0b11010101101001111001000101000111,
	0b00000110110010100110001101010001,
	0b00010100001010010010100101100111,
	0b00100111101101110000101010000101,
	0b00101110000110110010000100111000,
	0b01001101001011000110110111111100,
	0b01010011001110000000110100010011,
	0b01100101000010100111001101010100,
	0b01110110011010100000101010111011,
	0b10000001110000101100100100101110,
	0b10010010011100100010110010000101,
	0b10100010101111111110100010100001,
	0b10101000000110100110011001001011,
	0b11000010010010111000101101110000,
	0b11000111011011000101000110100011,
	0b11010001100100101110100000011001,
	0b11010110100110010000011000100100,
	0b11110100000011100011010110000101,
	0b00010000011010101010000001110000,
	0b00011001101001001100000100010110,
	0b00011110001101110110110000001000,
	0b00100111010010000111011101001100,
	0b00110100101100001011110010110101,
	0b00111001000111000000110010110011,
	0b01001110110110001010101001001010,
	0b01011011100111001100101001001111,
	0b01101000001011100110111111110011,
	0b01110100100011111000001011101110,
	0b01111000101001010110001101101111,
	0b10000100110010000111100000010100,
	0b10001100110001110000001000001000,
	0b10010000101111101111111111111010,
	0b10100100010100000110110011101011,
	0b10111110111110011010001111110111,
	0b11000110011100010111100011110010,
}

type Hash [32]byte

func (h Hash) String() string {
	bytes := make([]byte, 0, 64)
	for _, b := range [32]byte(h) {
		hex := fmt.Sprintf("%.2x", b)
		bytes = append(bytes, hex...)
	}
	return string(bytes)
}

func (h Hash) Bits() string {
	bits := make([]byte, 0, 256)
	for _, b := range [32]byte(h) {
		hex := fmt.Sprintf("%.8b", b)
		bits = append(bits, hex...)
	}
	return string(bits)
}

func Sum(message []byte) Hash {
	block := messageBlock(message)
	integers := integerBlock(block)

	for i := 0; i < len(integers)/16; i++ {
		var (
			start = i * 16
			end   = (i + 1) * 16
		)
		processChunk(integers[start:end])
	}

	return toHash(H0, H1, H2, H3, H4, H5, H6, H7)
}

func messageBlock(message []byte) []byte {
	block := ascii.ToBits(string(message))
	block = append(block, '1')

	var (
		sizeWithoutPadding = len(block) + LengthSegmentSize
		chunksCount        = sizeWithoutPadding/ChunkSize + 1
		messageBlockSize   = ChunkSize * chunksCount
		paddingSize        = messageBlockSize - sizeWithoutPadding
	)

	block = append(block, bytes.Repeat([]byte{'0'}, paddingSize)...)

	format := fmt.Sprintf("%%.%db", LengthSegmentSize)
	blockLength := fmt.Sprintf(format, len(message)*8)
	block = append(block, blockLength...)

	return block
}

func integerBlock(block []byte) []uint32 {
	integers := make([]uint32, len(block)/32)

	for i := range integers {
		var (
			start   = i * 32
			end     = (i + 1) * 32
			bits    = block[start:end]
			integer = uint32(0)
		)

		for j, b := range bits {
			weight := uint32(math.Pow(2, float64(31-j)))
			integer += weight * uint32(b-48)
		}

		integers[i] = integer
	}

	return integers
}

func processChunk(chunk []uint32) {
	messageSchedule := make([]uint32, 64)
	copy(messageSchedule[:16], chunk)

	for i := range 48 {
		nextRowFirstStage(messageSchedule, i)
	}

	state := SecondStageState{H0, H1, H2, H3, H4, H5, H6, H7}
	for i := range 64 {
		nextRowSecondStage(messageSchedule, &state, i)
	}

	H0 += state.A
	H1 += state.B
	H2 += state.C
	H3 += state.D
	H4 += state.E
	H5 += state.F
	H6 += state.G
	H7 += state.H
}

func nextRowFirstStage(messageSchedule []uint32, i int) {
	var (
		w0  = messageSchedule[i+0]
		w1  = messageSchedule[i+1]
		w9  = messageSchedule[i+9]
		w14 = messageSchedule[i+14]
	)

	var (
		w1rr7  = bits.RotateLeft32(w1, -7)
		w1rr18 = bits.RotateLeft32(w1, -18)
		w1rs3  = w1 >> 3
		s0     = w1rr7 ^ w1rr18 ^ w1rs3
	)

	var (
		w14rr17 = bits.RotateLeft32(w14, -17)
		w14rr19 = bits.RotateLeft32(w14, -19)
		w14rs10 = w14 >> 10
		s1      = w14rr17 ^ w14rr19 ^ w14rs10
	)

	messageSchedule[i+16] = w0 + s0 + w9 + s1
}

type SecondStageState struct {
	A, B, C, D, E, F, G, H uint32
}

func nextRowSecondStage(messageSchedule []uint32, state *SecondStageState, i int) {
	var (
		err6  = bits.RotateLeft32(state.E, -6)
		err11 = bits.RotateLeft32(state.E, -11)
		err25 = bits.RotateLeft32(state.E, -25)
		s1    = err6 ^ err11 ^ err25
	)

	var (
		c1 = choice(state.E, state.F, state.G)
		t1 = state.H + s1 + c1 + K[i] + messageSchedule[i]
	)

	var (
		arr2  = bits.RotateLeft32(state.A, -2)
		arr13 = bits.RotateLeft32(state.A, -13)
		arr22 = bits.RotateLeft32(state.A, -22)
		s0    = arr2 ^ arr13 ^ arr22
	)

	var (
		m1 = majority(state.A, state.B, state.C)
		t2 = m1 + s0
	)

	var (
		a = t1 + t2
		e = state.D + t1
	)

	state.H = state.G
	state.G = state.F
	state.F = state.E
	state.E = e
	state.D = state.C
	state.C = state.B
	state.B = state.A
	state.A = a
}

func choice(a, b, c uint32) uint32 {
	return (a & b) ^ (^a & c)
}

func majority(a, b, c uint32) uint32 {
	return (a & b) ^ (a & c) ^ (b & c)
}

func toHash(parts ...uint32) Hash {
	hash := Hash{}
	for i, part := range parts {
		hash[i*4+0] = byte(part >> 24)
		hash[i*4+1] = byte(part >> 16)
		hash[i*4+2] = byte(part >> 8)
		hash[i*4+3] = byte(part >> 0)
	}
	return hash
}
