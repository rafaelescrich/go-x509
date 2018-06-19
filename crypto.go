package main

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"

	"github.com/golang/crypto/argon2"
)

// KeyPair is the public and private key pair
type KeyPair struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// Salt is hardcoded
const Salt = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// GenerateMasterKey is the method to generate a key from the salt and
// password
func GenerateMasterKey(password string) []byte {
	// The draft RFC recommends time=3, and memory=32*1024
	// is a sensible number. If using that amount of memory (32 MB) is
	// not possible in some contexts then the time parameter can be increased
	//  to compensate.
	// Key(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32)
	return argon2.Key([]byte(password), []byte(Salt), 3, 32*1024, 4, 32)

}

func genKeys(filename string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	publicKey := &privateKey.PublicKey
	fmt.Println("Private Key: ", privateKey)
	fmt.Println("Public key: ", publicKey)

	pemPrivateFile, err := os.Create("certs/" + filename + ".pem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var pemPrivateBlock = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	err = pem.Encode(pemPrivateFile, pemPrivateBlock)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pemPrivateFile.Close()

}

func readKeyPair(from string) (kp KeyPair, err error) {
	privateKeyFile, err := os.Open("certs/" + from + ".pem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	size := pemfileinfo.Size()
	pembytes := make([]byte, size)
	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)
	data, _ := pem.Decode([]byte(pembytes))
	privateKeyFile.Close()

	privateKey, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	publicKey := &privateKey.PublicKey
	fmt.Println("Private Key: ", privateKey)
	fmt.Println("Public key: ", publicKey)

	kp = KeyPair{privateKey, publicKey}
	return kp, err
}

func readPubKey(from string) {
	pubKey, err := os.Open("certs/" + from + ".pem")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pemfileinfo, _ := pubKey.Stat()
	size := pemfileinfo.Size()
	pembytes := make([]byte, size)
	buffer := bufio.NewReader(pubKey)
	_, err = buffer.Read(pembytes)
	block, _ := pem.Decode([]byte(pembytes))
	pubKey.Close()

	cert, _ := x509.ParseCertificate(block.Bytes)
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	fmt.Printf("This is the "+from+" public key: %x\n", rsaPublicKey.N)
}

func pkSign(hash []byte, privkbytes []byte) (r, s *big.Int, err error) {
	zero := big.NewInt(0)
	privkey, err := x509.ParseECPrivateKey(privkbytes)
	if err != nil {
		return zero, zero, err
	}

	r, s, err = ecdsa.Sign(rand.Reader, privkey, hash)
	if err != nil {
		return zero, zero, err
	}
	return r, s, nil
}

func pkVerify(hash []byte, pubkbytes []byte, r *big.Int, s *big.Int) (result bool) {
	pubk, err := x509.ParsePKIXPublicKey(pubkbytes)
	if err != nil {
		return false
	}

	switch pubk := pubk.(type) {
	case *ecdsa.PublicKey:
		return ecdsa.Verify(pubk, hash, r, s)
	default:
		return false
	}
}
