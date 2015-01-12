package main

func index(count uint16) []uint16 {
	index := make([]uint16, count)

	for i := uint16(0); i < count; i++ {
		index[i] = i
	}

	return index
}
