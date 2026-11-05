package main

import (
	"fmt"
	"log"
	"net"
	"bufio"
	"os"
)

const port string = ":8888"

func main() {
	ln := startServer(port)

	conn, err := ln.Accept()
	if err != nil {
		log.Fatalf("Accept error: %v", err)
	}

	handleConnection(conn)
}

func startServer(port string) (ln net.Listener) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}

	fmt.Println("Server started!")
	return ln
}

func handleConnection(conn net.Conn) {
	var addr string = conn.RemoteAddr().String()

	fmt.Printf("Host connected: %v\n", addr)

	for {
		fmt.Printf("%v> ", addr)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		if err := scanner.Err(); err != nil {
			fmt.Printf("Scanner error: %v\n", err)
		}

		conn.Write([]byte(scanner.Text()))
	}
}
