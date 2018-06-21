package main

import (
	"bufio"
	"crypto/rsa"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var publicServerKey *rsa.PublicKey

// Open a connection to a listening server
func Open(addr string) (*bufio.ReadWriter, error) {

	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

// InitClient initializes the client
func InitClient(port string) {

	na := Nonce()

	ta := time.Now().Format(time.RFC850)

	publicServerKey = new(rsa.PublicKey)
	publicServerKey.N = big.NewInt(0)
	publicServerKey.N.SetBytes([]byte("25051951314883923742917198599529893300770518822415211040991740181632958960962635078102875145607825968301915170902286354113671683016267722874229586470863762841128785490324207250806021365714673967692587295217586223182246350596856570245398605001881889738202791572639458781465096571800834170496806213181467221238557421074289617717044191701133749784730601791755859685341464781226429030426262337766566507887258578986832662340005054608773498915680887555688337820331951717998363693120373895850380432254423060642025902405330940080700898045904675145999350668552583934370139984494401758998971809529970529308560476801938964941467"))
	publicServerKey.E = 65537
	fmt.Println(publicServerKey.N)

	sessionKey := GenerateMasterKey(na)
	cSK, err := EncryptPubKey(sessionKey, publicServerKey)

	clientKP, err := readKeyPair("client")
	if err != nil {
		fmt.Println(err)
	}

	msg := append(na, ta...)
	msg = append(msg, cSK...)

	sm, err := Sign(msg, clientKP.privateKey)
	if err != nil {
		fmt.Println(err)
	}

	testStruct := Protocol{
		Nonce:              na,
		Timestamp:          ta,
		CipheredSessionKey: cSK,
		SignedMsg:          sm,
	}

	addr := "localhost:" + port

	// Open a connection to the server.
	rw, err := Open(addr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	log.Println("Send a struct as GOB:")
	log.Printf("Outer complexData struct: \n%#v\n", testStruct)
	enc := gob.NewEncoder(rw)
	n, err := rw.WriteString("GOB\n")
	if err != nil {
		fmt.Println(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)")
	}
	err = enc.Encode(testStruct)
	if err != nil {
		fmt.Println(err, "Encode failed for struct: %#v", testStruct)
	}
	err = rw.Flush()
	if err != nil {
		fmt.Println(err, "Flush failed.")
	}

	log.Print("Receive Protocol data:")
	var data Reply
	// Create a decoder that decodes directly into a struct variable.
	dec := gob.NewDecoder(rw)
	err = dec.Decode(&data)
	if err != nil {
		log.Println("Error decoding GOB data:", err)
		return
	}
	// // Print the complexData struct and the nested one, too, to prove
	// // that both travelled across the wire.
	log.Printf("Outer complexData struct: \n%#v\n", data)

	log.Println("Send the string request.")
	n, err = rw.WriteString("STRING\n")
	if err != nil {
		fmt.Println(err, "Could not send the STRING request ("+strconv.Itoa(n)+" bytes written)")
	}
	n, err = rw.WriteString("Additional data.\n")
	if err != nil {
		fmt.Println(err, "Could not send additional STRING data ("+strconv.Itoa(n)+" bytes written)")
	}
	log.Println("Flush the buffer.")
	err = rw.Flush()
	if err != nil {
		fmt.Println(err, "Flush failed.")
	}

	// Read the reply.
	log.Println("Read the reply.")
	response, err := rw.ReadString('\n')
	if err != nil {
		fmt.Println(err, "Client: Failed to read the reply: '"+response+"'")
	}

	log.Println("STRING request: got a response:", response)
}
