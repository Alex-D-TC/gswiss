package ds

import (
	"math/big"

	"github.com/alex-d-tc/gswiss/netDHT"
)

type FingerTable struct {
	siteAddr netDHT.DhtAddr
	bitCount uint16
	table    map[*big.Int]*netDHT.NeighData
}

func MakeFingerTable(bitCount uint16, siteAddr netDHT.DhtAddr) FingerTable {
	return FingerTable{
		siteAddr: siteAddr,
		bitCount: bitCount,
		table:    emptyTable(bitCount, siteAddr),
	}
}

func emptyTable(bitCount uint16, siteAddr netDHT.DhtAddr) map[*big.Int]*netDHT.NeighData {
	pow := big.NewInt(1)
	var table map[*big.Int]*netDHT.NeighData

	var i uint16
	for {
		if i >= bitCount {
			break
		}

		table[siteAddr.Addr.Add(siteAddr.Addr, pow)] = nil

		pow = pow.Mul(pow, pow)
		i++
	}

	return table
}

func (table *FingerTable) GetClosest(siteAddr netDHT.DhtAddr) *netDHT.NeighData {

	pow := big.NewInt(1)

	var i uint16
	for {
		if i >= table.siteAddr.BitSize {
			break
		}

		runningAddr := table.siteAddr.Addr.Add(table.siteAddr.Addr, pow)
		if runningAddr.Cmp(siteAddr.Addr) == 1 {
			return table.table[runningAddr]
		}

		pow = pow.Mul(pow, pow)
		i++
	}

	// default to returning the successor
	return table.table[table.siteAddr.Addr.Add(table.siteAddr.Addr, big.NewInt(1))]
}

func (table *FingerTable) Update(tableIdx *big.Int, newData *netDHT.NeighData) {
	if _, ok := table.table[tableIdx]; ok {
		table.table[tableIdx] = newData
	}
}
