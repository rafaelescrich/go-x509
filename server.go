// package main

// import (
// 	"crypto/rand"
// 	"crypto/tls"
// 	"crypto/x509"
// 	"log"
// 	"net"
// )

// func main() {
// 	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
// 	if err != nil {
// 		log.Fatalf("server: loadkeys: %s", err)
// 	}
// 	config := tls.Config{Certificates: []tls.Certificate{cert}}
// 	config.Rand = rand.Reader
// 	service := "0.0.0.0:8000"
// 	listener, err := tls.Listen("tcp", service, &config)
// 	if err != nil {
// 		log.Fatalf("server: listen: %s", err)
// 	}
// 	log.Print("server: listening")
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Printf("server: accept: %s", err)
// 			break
// 		}
// 		defer conn.Close()
// 		log.Printf("server: accepted from %s", conn.RemoteAddr())
// 		tlscon, ok := conn.(*tls.Conn)
// 		if ok {
// 			log.Print("ok=true")
// 			state := tlscon.ConnectionState()
// 			for _, v := range state.PeerCertificates {
// 				log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
// 			}
// 		}
// 		go handleClient(conn)
// 	}
// }

// func handleClient(conn net.Conn) {
// 	defer conn.Close()
// 	buf := make([]byte, 512)
// 	for {
// 		log.Print("server: conn: waiting")
// 		n, err := conn.Read(buf)
// 		if err != nil {
// 			if err != nil {
// 				log.Printf("server: conn: read: %s", err)
// 			}
// 			break
// 		}
// 		log.Printf("server: conn: echo %q\n", string(buf[:n]))
// 		n, err = conn.Write(buf[:n])

// 		n, err = conn.Write(buf[:n])
// 		log.Printf("server: conn: wrote %d bytes", n)

// 		if err != nil {
// 			log.Printf("server: write: %s", err)
// 			break
// 		}
// 	}
// 	log.Println("server: conn: closed")
// }
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
				fmt.Println("unexpected data: %s", buf[:n])
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

func main() {
	listen, err := net.Listen("tcp4", ":8000")
	defer listen.Close()
	if err != nil {
		log.Fatalf("Socket listen port %s failed,%s", "8000", err)
		os.Exit(1)
	}
	log.Printf("Begin listen port: %s", "8000")

	for {
		conn, _ := listen.Accept()
		go handleConnection(conn)
	}
}
