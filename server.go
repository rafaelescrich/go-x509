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
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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

func readPEM() {
	// needs to change reading from file instead of reading from const
	const rootPEM = `
-----BEGIN CERTIFICATE-----
MIIEBDCCAuygAwIBAgIDAjppMA0GCSqGSIb3DQEBBQUAMEIxCzAJBgNVBAYTAlVT
MRYwFAYDVQQKEw1HZW9UcnVzdCBJbmMuMRswGQYDVQQDExJHZW9UcnVzdCBHbG9i
YWwgQ0EwHhcNMTMwNDA1MTUxNTU1WhcNMTUwNDA0MTUxNTU1WjBJMQswCQYDVQQG
EwJVUzETMBEGA1UEChMKR29vZ2xlIEluYzElMCMGA1UEAxMcR29vZ2xlIEludGVy
bmV0IEF1dGhvcml0eSBHMjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AJwqBHdc2FCROgajguDYUEi8iT/xGXAaiEZ+4I/F8YnOIe5a/mENtzJEiaB0C1NP
VaTOgmKV7utZX8bhBYASxF6UP7xbSDj0U/ck5vuR6RXEz/RTDfRK/J9U3n2+oGtv
h8DQUB8oMANA2ghzUWx//zo8pzcGjr1LEQTrfSTe5vn8MXH7lNVg8y5Kr0LSy+rE
ahqyzFPdFUuLH8gZYR/Nnag+YyuENWllhMgZxUYi+FOVvuOAShDGKuy6lyARxzmZ
EASg8GF6lSWMTlJ14rbtCMoU/M4iarNOz0YDl5cDfsCx3nuvRTPPuj5xt970JSXC
DTWJnZ37DhF5iR43xa+OcmkCAwEAAaOB+zCB+DAfBgNVHSMEGDAWgBTAephojYn7
qwVkDBF9qn1luMrMTjAdBgNVHQ4EFgQUSt0GFhu89mi1dvWBtrtiGrpagS8wEgYD
VR0TAQH/BAgwBgEB/wIBADAOBgNVHQ8BAf8EBAMCAQYwOgYDVR0fBDMwMTAvoC2g
K4YpaHR0cDovL2NybC5nZW90cnVzdC5jb20vY3Jscy9ndGdsb2JhbC5jcmwwPQYI
KwYBBQUHAQEEMTAvMC0GCCsGAQUFBzABhiFodHRwOi8vZ3RnbG9iYWwtb2NzcC5n
ZW90cnVzdC5jb20wFwYDVR0gBBAwDjAMBgorBgEEAdZ5AgUBMA0GCSqGSIb3DQEB
BQUAA4IBAQA21waAESetKhSbOHezI6B1WLuxfoNCunLaHtiONgaX4PCVOzf9G0JY
/iLIa704XtE7JW4S615ndkZAkNoUyHgN7ZVm2o6Gb4ChulYylYbc3GrKBIxbf/a/
zG+FA1jDaFETzf3I93k9mTXwVqO94FntT0QJo544evZG0R0SnU++0ED8Vf4GXjza
HFa9llF7b1cq26KqltyMdMKVvvBulRP/F/A8rLIQjcxz++iPAsbw+zOzlTvjwsto
WHPbqCRiOwY1nQ2pM714A5AuTHhdUDqB1O6gyHA43LL5Z/qHQF1hwFGPa4NrzQU6
yuGnBXj8ytqU0CwIPX4WecigUCAkVDNx
-----END CERTIFICATE-----`

	block, _ := pem.Decode([]byte(rootPEM))
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	fmt.Println(rsaPublicKey.N)
	// fmt.Println(rsaPublicKey.E)
}

// InitServer initalizes the server
func InitServer(port string) {

	// read from file the server's private key

	// read from pem the client's public key
	readPEM()

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
