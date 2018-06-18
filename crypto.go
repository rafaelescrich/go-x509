package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"math/big"
)

// KeyPair is the public and private key pair
type KeyPair struct {
	privatekey, publickey []byte
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
