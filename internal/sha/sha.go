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

func Hash(message []byte) []byte {
	block := messageBlock(message)
	integers := integerBlock(block)

	messageSchedule := make([]uint32, 64)
	copy(messageSchedule[0:16], integers[0:16])

	for i := range 48 {
		nextRow(messageSchedule, i)
	}

	for i, integer := range messageSchedule {
		fmt.Printf("w%d\t%.32b\n", i, integer)
	}

	return nil
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

func nextRow(messageSchedule []uint32, i int) {
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
