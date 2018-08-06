package main

import (
	"fmt"

	"github.com/alex-d-tc/gswiss/netDHT"
	"github.com/alex-d-tc/gswiss/util/netUtil"

	"github.com/alex-d-tc/gswiss/netDHT/node"
)

func main() {

	state, err := node.InitCluster(netUtil.NetAddr{Address: ":8080", ConnType: "udp"}, netDHT.BitSize)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(state)
	for {
		// just run
	}
}
