package ds

import (
	"github.com/alex-d-tc/gswiss/netDHT"
	"github.com/alex-d-tc/gswiss/util"
)

type RouteTable struct {
	k           uint16
	sourceAddr  [netDHT.DhtAddrByteCount]byte
	prefixTable [netDHT.DhtAddrByteCount * 8][]*netDHT.NodeData
}

func MakeRouteTable(k uint16, sourceAddr [netDHT.DhtAddrByteCount]byte) *RouteTable {

	var table [netDHT.DhtAddrByteCount * 8][]*netDHT.NodeData
	var i uint

	for {
		if i >= netDHT.DhtAddrByteCount*8 {
			break
		}

		table[i] = make([]*netDHT.NodeData, k)

		i++
	}

	return &RouteTable{
		prefixTable: table,
		k:           k,
		sourceAddr:  sourceAddr,
	}
}

func (rTable *RouteTable) GetNearestK(addr [netDHT.DhtAddrByteCount]byte) []*netDHT.NodeData {
	rightmostIdx, _ := util.GetCommonPrefixSize(rTable.sourceAddr[:], addr[:])
	var result []*netDHT.NodeData

	table := rTable.prefixTable

	for {
		if rightmostIdx < 0 || uint16(len(result)) >= rTable.k {
			break
		}

		if uint16(len(result)+len(table[rightmostIdx])) > rTable.k {
			// TODO: Elect to only get the nodes which are 'closest' to me physically or by some other measure apart from just the prefix differences
			result = append(result, table[rightmostIdx][:rTable.k-uint16(len(result))]...)
		} else {
			result = append(result, table[rightmostIdx]...)
		}

		rightmostIdx--
	}

	return result
}

func (rTable *RouteTable) Insert(data *netDHT.NodeData) {

	idx, _ := util.GetCommonPrefixSize(rTable.sourceAddr[:], data.DhtAddress[:])
	table := rTable.prefixTable

	if uint16(len(table[idx])) == rTable.k {
		// TODO: Add eviction policy
	} else {
		table[idx] = append(table[idx], data)
	}
}
