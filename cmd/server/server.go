package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var port string = ":8888"

func main() {
	config := readConfig("server.cfg")
	for _, kv := range config {
		key, value := kv[0], kv[1]
		switch key {
		case "port":
			port = value
		}
	}

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

	fmt.Printf("Server started on port [%v]\n", port)
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

func readConfig(name string) (config [][]string) {
	data, err := os.ReadFile(name)
	if err != nil {
		log.Fatalf("read error: %v\n", err)
	}

	lines := strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n")

	for _, line := range lines {
		keyValuePair := strings.Split(line, "=")
		if len(keyValuePair) == 2 {
			config = append(config, keyValuePair)
		}
	}

	return config
}
