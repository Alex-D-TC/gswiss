package ds

import (
	"github.com/alex-d-tc/gswiss/net"
	"github.com/golang-collections/go-datastructures/bitarray"
)

type KTree struct {
	left  *innerKTree
	right *innerKTree
	k     uint
}

type innerKTree struct {
	left  *innerKTree
	right *innerKTree
	data  []net.NodeData
}

func MakeKTree(k uint) KTree {
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

func (tree *KTree) Insert(data net.NodeData) {
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

				// The address space is fully explored. Cannot recursively split
				if runningIdx == uint((dhtAddrByteCount*8)-1) {

				} else {
					// TODO: If the split would result in two K-Buckets which cannot hold at least K items, attempt eviction instead
					// TODO: Prefer eviction over splitting if plausible, to avoid imbalancing

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

func (tree *KTree) LocateClosest(data [net.DhtAddrByteCount]byte) []net.NodeData {
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

// splitByBit takes the node data and splits it into two slices, depending on the value of the bit at the bitIdx location in the pastry address.
// Nodes which have the bit set are part of the rightData slice and the rest are part of the leftData slice.
func splitByBit(bitIdx uint, data []net.NodeData) (leftData []net.NodeData, rightData []net.NodeData) {
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
