package ds

import (
	"github.com/alex-d-tc/gswiss/netDHT"
	"github.com/alex-d-tc/gswiss/util"
)

type RouteTable struct {
	k           uint16
	sourceAddr  netDHT.DhtAddr
	prefixTable [netDHT.BitSize][]*netDHT.NeighData
}

func MakeRouteTable(k uint16, sourceAddr netDHT.DhtAddr) *RouteTable {

	var table [netDHT.BitSize][]*netDHT.NeighData
	var i uint

	for {
		if i >= netDHT.BitSize {
			break
		}

		table[i] = make([]*netDHT.NeighData, k)

		i++
	}

	return &RouteTable{
		prefixTable: table,
		k:           k,
		sourceAddr:  sourceAddr,
	}
}

func (rTable *RouteTable) GetNearestK(addr netDHT.DhtAddr) []*netDHT.NeighData {
	rightmostIdx, _ := util.GetCommonPrefixSize(rTable.sourceAddr.Bytes(), addr.Bytes())
	var result []*netDHT.NeighData

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

func (rTable *RouteTable) Insert(data *netDHT.NeighData) {

	idx, _ := util.GetCommonPrefixSize(rTable.sourceAddr.Bytes(), data.DhtAddress.Bytes())
	table := rTable.prefixTable

	if uint16(len(table[idx])) == rTable.k {
		// TODO: Add eviction policy
	} else {
		table[idx] = append(table[idx], data)
	}
}
