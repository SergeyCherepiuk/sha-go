package splitter

func Split(bytes []byte, blockSize int) [][]byte {
	blocks := make([][]byte, 0)

	var block []byte
	for i, b := range bytes {
		if i%blockSize == 0 {
			block = make([]byte, 0)
		}

		block = append(block, b)

		if len(block) == blockSize {
			blocks = append(blocks, block)
		}
	}

	if len(block) != 0 && len(block) != blockSize {
		blocks = append(blocks, block)
	}

	return blocks
}
