package util

import (
	"errors"
)

// GetCommonPrefixSize gets the size of the longest common prefix of bits between the two slices
// If the two slices are of different lengths, an error occurs
func GetCommonPrefixSize(first []byte, second []byte) (uint16, error) {

	if len(first) != len(second) {
		return 0, errors.New("Slices are of different lengths")
	}

	var idx uint16

	for {
		if GetBit(first[:], uint(idx)) != GetBit(second[:], uint(idx)) {
			break
		}

		idx++
	}

	return idx, nil
}

// GetBit gets the value of the bit from the byte array data at the given index.
func GetBit(data []byte, index uint) bool {
	byteIdx := index / 8
	return GetFromByte(data[byteIdx], uint8(index%8))
}

// GetFromByte gets the value of the bit from within a given byte.
func GetFromByte(data byte, index uint8) bool {
	var startMask byte = 1
	var idx uint8 = 7

	for {
		if idx == index {
			break
		}

		idx--
		startMask = startMask << 1
	}

	bit := data & startMask
	return !(bit == 0)
}
