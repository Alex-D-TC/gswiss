package node

import (
	"fmt"
	"net"

	"github.com/alex-d-tc/gswiss/netDHT"
	"github.com/alex-d-tc/gswiss/netDHT/ds"
)

type NodeState struct {
	k       uint16
	table   *ds.RouteTable
	address netDHT.DhtAddr
}

func ListenUDP(address net.UDPAddr, handler func([]byte)) {
	conn, err := net.ListenUDP("udp", &address)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {

		var payload [4096]byte

		n, err := conn.Read(payload[:])
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handler(payload[:n])
	}
}

func SendMessage(source netDHT.DhtAddr, dest *netDHT.NeighData, msg NetMessage) bool {

	udpAddr, err := net.ResolveUDPAddr("udp", dest.NetAddress.String())
	if err != nil {
		fmt.Println(err)
		return false
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer conn.Close()

	var message []byte

	message = append(message, byte(msg.MsgType))
	message = append(message, source.Bytes()[:]...)
	message = append(message, msg.Payload...)

	written, err := conn.Write(message)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println(fmt.Sprintf("Written %d bytes", written))
	return true
}

func receiveMessage(conn net.Conn) {
	var message []byte
	conn.Read(message)

}

const (
	locateMsg = 0
	sendMsg   = 1
	joinMsg   = 2
)

type NetMessage struct {
	MsgType uint8
	Payload []byte
	Source  [netDHT.BitSize * 8]byte
}
