package main

import (
	"fmt"
	"net"

	"github.com/alex-d-tc/gswiss/netDHT/node"
)

func main() {

	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println(err)
		panic("")
	}

	node.ListenUDP(*addr, handleConn)
}

func handleConn(payload []byte) {

	fmt.Println("UDP packet received")
	fmt.Println(payload)
}
