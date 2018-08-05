package node

import (
	"crypto/sha256"
	"fmt"
	"net"

	"github.com/alex-d-tc/gswiss/netDHT"
	"github.com/alex-d-tc/gswiss/netDHT/ds"
)

type NodeState struct {
	k       uint16
	table   *ds.RouteTable
	address [netDHT.DhtAddrByteCount]byte
}

func Bootstrap(k uint16, bootStrapperDHTAddress [netDHT.DhtAddrByteCount]byte, bootStrapperNetAddress net.Addr) NodeState {

	conn, err := net.Dial("udp", bootStrapperNetAddress.String())
	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("Connection to peer %s failed on network %s", bootStrapperNetAddress.String(), bootStrapperNetAddress.Network()))
	}

	fmt.Println(conn)

	dhtAddress := makeAddress(bootStrapperDHTAddress[:])

	table := ds.MakeRouteTable(k, dhtAddress)

	return NodeState{
		k:       k,
		address: dhtAddress,
		table:   table,
	}
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

func SendMessage(source [netDHT.DhtAddrByteCount]byte, dest netDHT.NodeData, msg NetMessage) bool {

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
	message = append(message, source[:]...)
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

func makeAddress(entropySource []byte) [netDHT.DhtAddrByteCount]byte {

	hash := sha256.Sum256(entropySource)
	var dhtAddress [netDHT.DhtAddrByteCount]byte
	copy(dhtAddress[:], hash[:netDHT.DhtAddrByteCount])

	return dhtAddress
}

const (
	locateMsg = 0
	sendMsg   = 1
)

type NetMessage struct {
	MsgType uint8
	Payload []byte
	Source  [netDHT.DhtAddrByteCount]byte
}
