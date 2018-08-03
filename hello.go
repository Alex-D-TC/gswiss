package main

import (
	"net"

	"github.com/golang-collections/go-datastructures/bitarray"
)

func main() {
}

const dhtAddrByteCount uint = 20

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

					if getBit(dhtAddr[:], runningIdx+1) {
						prefix.SetBit(uint64(runningIdx + 1))
					}

					leftData, rightData := splitByPrefix(prefix, runningIdx+1, runningTree.data)

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

// TODO: Implement splitting by prefix
func splitByPrefix(prefix bitarray.BitArray, runningIdx uint, data []NodeData) (leftData []NodeData, rightData []NodeData) {
	return nil, nil
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

func getBit(data []byte, index uint) bool {
	byteIdx := index / 8
	return getFromByte(data[byteIdx], uint8(index%8))
}

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
