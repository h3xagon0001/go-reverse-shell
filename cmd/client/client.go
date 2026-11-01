package main

import (
	"fmt"
	"log"
	"net"
)

const addr string = "localhost:8888"

var buf []byte

func main() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Dail error: %v", err)
	}

	fmt.Printf("Host: %v\n", conn.LocalAddr().String())

	for {
		buf = make([]byte, 256)
		_, err := conn.Read(buf)
		if err != nil {
			log.Fatalf("Read error: %v", err)
		}
		fmt.Printf("Msg: %s\n", buf)
	}
	// possible reconnection idea
	// start the reader in a goroutine and have main do polling checks for connectivity
	// if connection unavailable, stop goroutine and reconnect
}
