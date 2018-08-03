package main

import (
	"net"

	"github.com/golang-collections/go-datastructures/bitarray"
)

func main() {
}

// TODO: DO NOT FORGET TO CHANGE IT BACK TO 20 AFTER TESTING YA DINGUS
const dhtAddrByteCount uint = 1

type KTree struct {
	left  *innerKTree
	right *innerKTree
	k     uint
}

type innerKTree struct {
	left  *innerKTree
	right *innerKTree
	data  []NodeData
}

type NodeData struct {
	dhtAddress [dhtAddrByteCount]byte
	netAddress net.Addr
}

func Make(k uint) KTree {
	return KTree{
		left:  makeInnerKTree(),
		right: makeInnerKTree(),
		k:     k,
	}
}

func makeInnerKTree() *innerKTree {
	return &innerKTree{
		left:  nil,
		right: nil,
		data:  nil,
	}
}

func (tree *KTree) Insert(data NodeData) {
	var runningTree *innerKTree
	var runningIdx uint = 1
	prefix := bitarray.NewBitArray(uint64(dhtAddrByteCount * 8))
	dhtAddr := data.dhtAddress

	if !getBit(dhtAddr[:], 0) {
		runningTree = tree.left
	} else {
		runningTree = tree.right
	}

	for {
		runningBit := getBit(dhtAddr[:], runningIdx)

		if runningBit {
			prefix.SetBit(uint64(runningIdx))
		}

		if !runningBit && runningTree.left == nil ||
			runningBit && runningTree.right == nil {

			if uint(len(runningTree.data)) == tree.k {
				// TODO: Add recursive branching if the K-Bucket is full

				// The address space is fully explored. Cannot recursively split
				if runningIdx == uint((dhtAddrByteCount*8)-1) {

				} else {

					leftData, rightData := splitByBit(runningIdx+1, runningTree.data)

					runningTree.data = nil

					runningTree.left = makeInnerKTree()
					runningTree.right = makeInnerKTree()

					runningTree.left.data = leftData
					runningTree.right.data = rightData
				}

			} else {
				runningTree.data = append(runningTree.data, data)
			}

			break
		} else {

			if !runningBit {
				runningTree = runningTree.left
			} else {
				runningTree = runningTree.right
			}
		}

		runningIdx++
	}
}

// splitByBit takes the node data and splits it into two slices, depending on the value of the bit at the bitIdx location in the pastry address.
// Nodes which have the bit set are part of the rightData slice and the rest are part of the leftData slice.
func splitByBit(bitIdx uint, data []NodeData) (leftData []NodeData, rightData []NodeData) {
	leftData = nil
	rightData = nil

	for _, nodeData := range data {
		if !getBit(nodeData.dhtAddress[:], bitIdx) {
			leftData = append(leftData, nodeData)
		} else {
			rightData = append(rightData, nodeData)
		}
	}

	return
}

func (tree *KTree) LocateClosest(data [dhtAddrByteCount]byte) []NodeData {
	var runningTree *innerKTree
	var runningIdx uint = 1

	if !getBit(data[:], 0) {
		runningTree = tree.left
	} else {
		runningTree = tree.right
	}

	for {
		runningBit := getBit(data[:], runningIdx)

		if !runningBit {

			if runningTree.left == nil {
				return runningTree.data
			}

			runningTree = runningTree.left

		} else {

			if runningTree.right == nil {
				return runningTree.data
			}

			runningTree = runningTree.right
		}

		runningIdx++
	}

}

// getBit gets the value of the bit from the byte array data at the given index.
func getBit(data []byte, index uint) bool {
	byteIdx := index / 8
	return getFromByte(data[byteIdx], uint8(index%8))
}

// getFromByte gets the value of the bit from within a given byte.
func getFromByte(data byte, index uint8) bool {
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
