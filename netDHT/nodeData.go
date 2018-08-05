package netDHT

import (
	"net"
)

const (
	// TODO: DO NOT FORGET TO CHANGE IT BACK TO 20 AFTER TESTING YA DINGUS
	DhtAddrByteCount = 1
)

type NodeData struct {
	DhtAddress [DhtAddrByteCount]byte
	NetAddress net.Addr
	connection net.Conn
}

func (data *NodeData) Cleanup() {
	data.connection.Close()
}
