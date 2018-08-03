package main

import (
	"net"
)

func main() {
}

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
	dhtAddress []byte
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
	dhtAddr := data.dhtAddress

	if !getBit(dhtAddr, 0) {
		runningTree = tree.left
	} else {
		runningTree = tree.right
	}

	for {
		runningBit := getBit(dhtAddr, runningIdx)

		if !runningBit && runningTree.left == nil ||
			runningBit && runningTree.right == nil {
			runningTree.data = append(runningTree.data, data)

			// TODO: Add recursive branching if the K-Bucket is full
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

func (tree *KTree) LocateClosest(data []byte) []NodeData {
	var runningTree *innerKTree
	var runningIdx uint = 1

	if !getBit(data, 0) {
		runningTree = tree.left
	} else {
		runningTree = tree.right
	}

	for {
		runningBit := getBit(data, runningIdx)

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
