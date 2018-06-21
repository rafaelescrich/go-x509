package main

import (
	"bufio"
	"crypto/rsa"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var publicClientKey *rsa.PublicKey
var sessionKey []byte
var serverKP KeyPair

// HandleFunc is a function that handles an incoming command.
type HandleFunc func(*bufio.ReadWriter)

// Endpoint provides an endpoint to other processess
// that they can send data to.
type Endpoint struct {
	listener net.Listener
	handler  map[string]HandleFunc

	// Maps are not threadsafe, so we need a mutex to control access.
	m sync.RWMutex
}

// NewEndpoint creates a new endpoint. Too keep things simple,
// the endpoint listens on a fixed port number.
func NewEndpoint() *Endpoint {
	// Create a new Endpoint with an empty list of handler funcs.
	return &Endpoint{
		handler: map[string]HandleFunc{},
	}
}

// AddHandleFunc adds a new function for handling incoming data.
func (e *Endpoint) AddHandleFunc(name string, f HandleFunc) {
	e.m.Lock()
	e.handler[name] = f
	e.m.Unlock()
}

// Listen starts listening on the endpoint port on all interfaces.
// At least one handler function must have been added
// through AddHandleFunc() before.
func (e *Endpoint) Listen(port string) error {
	var err error
	e.listener, err = net.Listen("tcp", port)
	if err != nil {
		return errors.Wrapf(err, "Unable to listen on port %s\n", port)
	}
	log.Println("Listen on", e.listener.Addr().String())
	for {
		log.Println("Accept a connection request.")
		conn, err := e.listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		log.Println("Handle incoming messages.")
		e.handleMessages(conn)
	}
}

// handleMessages reads the connection up to the first newline.
// Based on this string, it calls the appropriate HandleFunc.
func (e *Endpoint) handleMessages(conn net.Conn) {
	// Wrap the connection into a buffered reader for easier reading.
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	// Read from the connection until EOF. Expect a command name as the
	// next input. Call the handler that is registered for this command.
	for {
		log.Print("Receive command '")
		cmd, err := rw.ReadString('\n')
		switch {
		case err == io.EOF:
			log.Println("Reached EOF - close this connection.\n   ---")
			return
		case err != nil:
			log.Println("\nError reading command. Got: '"+cmd+"'\n", err)
			return
		}
		// Trim the request string - ReadString does not strip any newlines.
		cmd = strings.Trim(cmd, "\n ")
		log.Println(cmd + "'")

		// Fetch the appropriate handler function from the 'handler' map and call it.
		e.m.RLock()
		handleCommand, ok := e.handler[cmd]
		e.m.RUnlock()
		if !ok {
			log.Println("Command '" + cmd + "' is not registered.")
			return
		}
		handleCommand(rw)
	}
}

// handleStrings handles the "STRING" request.
func handleStrings(rw *bufio.ReadWriter) {
	// Receive a string.
	log.Print("Receive STRING message:")
	s, err := rw.ReadString('\n')
	if err != nil {
		log.Println("Cannot read from connection.\n", err)
	}
	s = strings.Trim(s, "\n ")
	log.Println(s)
	_, err = rw.WriteString("Thank you.\n")
	if err != nil {
		log.Println("Cannot write to connection.\n", err)
	}
	err = rw.Flush()
	if err != nil {
		log.Println("Flush failed.", err)
	}
}

// handleGob handles the "GOB" request. It decodes the received GOB data
// into a struct.
func handleGob(rw *bufio.ReadWriter) {
	log.Print("Receive Protocol data:")
	var data Protocol
	// Create a decoder that decodes directly into a struct variable.
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&data)
	if err != nil {
		log.Println("Error decoding GOB data:", err)
		return
	}

	// // Print the complexData struct and the nested one, too, to prove
	// // that both travelled across the wire.
	log.Printf("Outer complexData struct: \n%#v\n", data)

	// Replying message
	na := data.Nonce
	msg := append(data.Nonce, data.Timestamp...)
	msg = append(msg, data.CipheredSessionKey...)
	err = VerifySig(msg, data.SignedMsg, publicClientKey)
	if err != nil {
		log.Println("Error verifying message:", err)
		return
	}
	sessionKey, err = DecryptPrivKey(data.CipheredSessionKey, serverKP.privateKey)
	// Now we can reply
	nb := Nonce()
	tb := time.Now().Format(time.RFC850)
	sSK, err := EncryptPubKey(sessionKey, publicClientKey)

	resp := append(nb, tb...)
	resp = append(msg, sSK...)

	sm, err := Sign(resp, serverKP.privateKey)
	if err != nil {
		fmt.Println(err)
	}
	reply := Reply{
		NonceB:             nb,
		NonceA:             na,
		TimestampB:         tb,
		CipheredSessionKey: sSK,
		SignedMsg:          sm,
	}
	enc := gob.NewEncoder(rw)
	n, err := rw.WriteString("GOB\n")
	if err != nil {
		fmt.Println(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)")
	}
	err = enc.Encode(reply)
	if err != nil {
		fmt.Println(err, "Encode failed for struct: %#v", reply)
	}
	err = rw.Flush()
	if err != nil {
		fmt.Println(err, "Flush failed.")
	}
}

// InitServer initalizes the server
func InitServer(port string) {

	publicClientKey = new(rsa.PublicKey)
	publicClientKey.N = big.NewInt(0)
	publicClientKey.N.SetBytes([]byte("29439234235147834624400106512350547158011030442832601322153468696324894192687132514357408532936147318236762955403181573699071216172231455769712275184435320099154582974822642760501066771874468050439146958798642517591941662554845991426625213455717261854031816578015209820463429808070353329037406237874208798850018644138706527756636302935154226853894506677547248237440150967630363864947704469855899873459413047449535629575893869062764724732809870070228795939455341871885281507205679866762058558155480584542658591248464177314020994728543478969515862596440757390666404159966689110277124077529957124797448578029298193299009"))
	publicClientKey.E = 65537
	fmt.Println(publicClientKey.N)

	serverKP, _ = readKeyPair("server")

	endpoint := NewEndpoint()

	// Add the handle funcs.
	endpoint.AddHandleFunc("STRING", handleStrings)
	endpoint.AddHandleFunc("GOB", handleGob)

	addr := "localhost:" + port
	// Start listening.
	err := endpoint.Listen(addr)
	if err != nil {
		fmt.Println(err)
	}
}
