// package main

// import (
// 	"crypto/tls"
// 	"crypto/x509"
// 	"fmt"
// 	"io"
// 	"log"
// )

// func main() {
// 	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
// 	if err != nil {
// 		log.Fatalf("server: loadkeys: %s", err)
// 	}
// 	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
// 	conn, err := tls.Dial("tcp", "127.0.0.1:8000", &config)
// 	if err != nil {
// 		log.Fatalf("client: dial: %s", err)
// 	}
// 	defer conn.Close()
// 	log.Println("client: connected to: ", conn.RemoteAddr())

// 	state := conn.ConnectionState()
// 	for _, v := range state.PeerCertificates {
// 		fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
// 		fmt.Println(v.Subject)
// 	}
// 	log.Println("client: handshake: ", state.HandshakeComplete)
// 	log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)

// 	message := "Hello\n"
// 	n, err := io.WriteString(conn, message)
// 	if err != nil {
// 		log.Fatalf("client: write: %s", err)
// 	}
// 	log.Printf("client: wrote %q (%d bytes)", message, n)

// 	reply := make([]byte, 256)
// 	n, err = conn.Read(reply)
// 	log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
// 	log.Print("client: exiting")
// }

package main

import (
	"log"
	"net"
	"os"
)

// InitClient initializes the client
func InitClient(port string) {

	// read from file the client's private key

	// read from pem the server's public key

	addr := "localhost:" + port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	conn.Write([]byte("test"))
	conn.Write([]byte("\r\n\r\n"))
	log.Printf("Send: %s", "test")

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])
}
