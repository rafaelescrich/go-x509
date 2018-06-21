package main

import (
	"bufio"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/golang/crypto/argon2"
)

// Protocol defines the types of data to communicate
type Protocol struct {
	Nonce              []byte
	Timestamp          string
	CipheredSessionKey []byte
	SignedMsg          []byte
}

// Reply defines the response
type Reply struct {
	TimestampB         string
	NonceB             []byte
	NonceA             []byte
	CipheredSessionKey []byte
	SignedMsg          []byte
}

// KeyPair is the public and private key pair
type KeyPair struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// Salt is hardcoded
const Salt = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// GenerateMasterKey is the method to generate a key from the salt and
// password
func GenerateMasterKey(password []byte) []byte {
	// The draft RFC recommends time=3, and memory=32*1024
	// is a sensible number. If using that amount of memory (32 MB) is
	// not possible in some contexts then the time parameter can be increased
	//  to compensate.
	// Key(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32)
	return argon2.Key(password, []byte(Salt), 3, 32*1024, 4, 32)

}

// Nonce returns new nonce
func Nonce() []byte {
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		panic(err.Error())
	}
	return nonce
}

// EncryptAESGCM encrypt plaintext with the mk
func EncryptAESGCM(key []byte, nonce []byte, plaintext []byte) ([]byte, error) {
	var ciphertext []byte

	block, err := aes.NewCipher(key)
	if err != nil {
		return ciphertext, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return ciphertext, err
	}

	ciphertext = aesgcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptAESGCM decrypts a ciphertext with a key
func DecryptAESGCM(key []byte, nonce []byte, ciphertext []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
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

// EncryptPubKey encrypt message with public key
func EncryptPubKey(msg []byte, pubKey *rsa.PublicKey) ([]byte, error) {
	var enc []byte
	rng := rand.Reader
	enc, err := rsa.EncryptPKCS1v15(rng, pubKey, msg)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	return enc, err
}

// DecryptPrivKey decrypts message
func DecryptPrivKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	var msg []byte
	rng := rand.Reader
	msg, err := rsa.DecryptPKCS1v15(rng, priv, ciphertext)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	return msg, err
}

// Sign a message
func Sign(msg []byte, privKey *rsa.PrivateKey) ([]byte, error) {
	var signature []byte
	rng := rand.Reader

	hashed := sha256.Sum256(msg)

	signature, err := rsa.SignPKCS1v15(rng, privKey, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return signature, err
	}
	fmt.Printf("Signature: %x\n", signature)
	return signature, nil
}

// VerifySig verifies a signature
func VerifySig(msg []byte, signature []byte, pubKey *rsa.PublicKey) error {
	hashed := sha256.Sum256(msg)

	err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
		return err
	}
	return nil
}
