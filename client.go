package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// InitClient initializes the client
func InitClient(port string) {

	_, _ = readKeyPair("client")

	na := Nonce()

	ta := time.Now().Format(time.RFC850)
	fmt.Println(ta)

	addr := "localhost:" + port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	conn.Write(n1)
	conn.Write([]byte("EOF"))
	log.Printf("Send: %x", na)

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])
}
