package main

import (
	"fmt"
	"log"
	"net"
)

const port string = ":8888"

func main() {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}

	fmt.Println("Server started!")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Accept error: %v", err)
		}

		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	var addr string = conn.RemoteAddr().String()
	var msg string

	fmt.Printf("Host connected: %v\n", addr)

	for {
		fmt.Printf("%v> ", addr)
		n, err := fmt.Scanln(&msg)
		if n != 0 {
			if err != nil {
				log.Fatalf("Scan error: %v", err)
			}
			conn.Write([]byte(msg))
		}
	}
}
