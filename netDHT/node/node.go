package node

import (
	"crypto/sha256"
	"fmt"
	"net"

	"github.com/alex-d-tc/gswiss/util/netUtil"

	"github.com/alex-d-tc/gswiss/netDHT"
	"github.com/alex-d-tc/gswiss/netDHT/ds"
)

type NodeState struct {
	successor   *netDHT.NeighData
	predecessor *netDHT.NeighData
	table       *ds.FingerTable
	address     netDHT.DhtAddr
}

func makeAddress(entropySource []byte, byteCount uint16) []byte {
	hash := sha256.Sum256(entropySource)
	return hash[:byteCount]
}

func InitCluster(listenAddr netUtil.NetAddr, bitCount uint16) (*NodeState, error) {

	rawDhtAddr := makeAddress([]byte(listenAddr.String()), bitCount/8)
	siteAddr := netDHT.MakeDhtAddr(rawDhtAddr, bitCount)
	fingerTable := ds.MakeFingerTable(bitCount, siteAddr)

	addr, err := net.ResolveUDPAddr("udp", listenAddr.String())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	go ListenUDP(addr, handleMsg)

	fmt.Println("Cluster initialised")

	return &NodeState{
		successor:   nil,
		predecessor: nil,
		table:       &fingerTable,
		address:     siteAddr,
	}, nil
}

func BootStrap(bootStrapAddr netUtil.NetAddr) (*NodeState, error) {
	return nil, nil
}

func handleMsg(payload []byte) {
	fmt.Println("UDP packet received")
	fmt.Println(payload)
}

func ListenUDP(address *net.UDPAddr, handler func([]byte)) {
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {

		fmt.Println(fmt.Sprintf("Listening on interface %s", address.String()))

		var payload [4096]byte

		n, err := conn.Read(payload[:])
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handler(payload[:n])
	}
}

func SendMessageUDP(source netDHT.DhtAddr, dest *netDHT.NeighData, msg NetMessage) bool {

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
