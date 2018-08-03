package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"net"
)

func main() {
}

func startServer() {

	sv, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Started server")
		for {
			conn, err := sv.Accept()
			if err != nil {
				fmt.Println(err)
			}
			go handleConn(conn)
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		result, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(result)
	}

}

func computeID(data []byte) [sha256.Size]byte {
	return sha256.Sum256(data)
}
