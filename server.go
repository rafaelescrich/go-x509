package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	notify := make(chan error)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				notify <- err
				return
			}
			if n > 0 {
				fmt.Printf("unexpected data: %x", buf[:n])
			}
		}
	}()

	for {
		select {
		case err := <-notify:
			if io.EOF == err {
				fmt.Println("connection dropped message", err)
				return
			}
		case <-time.After(time.Second * 1):
			fmt.Println("timeout 1, still alive")
		}
	}
}

// InitServer initalizes the server
func InitServer(port string) {

	_, _ = readKeyPair("server")

	nb := Nonce()

	tb := time.Now().Format(time.RFC850)
	fmt.Println(tb)

	listen, err := net.Listen("tcp4", ":"+port)
	defer listen.Close()
	if err != nil {
		log.Fatalf("Socket listen port %s failed,%s", port, err)
		os.Exit(1)
	}
	log.Printf("Begin listen port: %s", port)

	for {
		conn, _ := listen.Accept()
		go handleConnection(conn)
	}
}
