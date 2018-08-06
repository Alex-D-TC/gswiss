package netDHT

import (
	"math/big"
	"net"
)

const (
	// TODO: DO NOT FORGET TO CHANGE IT BACK TO 20 AFTER TESTING YA DINGUS
	BitSize = 8
)

type NeighData struct {
	DhtAddress DhtAddr
	NetAddress net.Addr
}

type DhtAddr struct {
	Addr    *big.Int
	BitSize uint16
	Mod     *big.Int
}

func MakeDhtAddr(addr []byte, bitSize uint16) DhtAddr {
	mod := big.NewInt(2)
	mod = mod.Exp(mod, big.NewInt(int64(bitSize)), nil)

	address := big.NewInt(0)
	address = address.SetBytes(addr)
	address = address.Mod(address, mod)

	return DhtAddr{
		Addr:    address,
		BitSize: bitSize,
		Mod:     mod,
	}
}

func (addr DhtAddr) Bytes() []byte {

	result := make([]byte, BitSize/8)

	copyIdx := (BitSize / 8) - len(result)
	copy(result[copyIdx:], addr.Addr.Bytes())

	return result
}

func (addr DhtAddr) RawIncrementBy(add *big.Int) *big.Int {
	return addr.Addr.Mod(addr.Addr.Add(addr.Addr, add), addr.Mod)
}
