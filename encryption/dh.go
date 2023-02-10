package encryption

import (
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/curve25519"
)

const dhLen = 32

type KeyPair struct {
	PrivateKey [32]byte
	PublicKey  [32]byte
}

func GenerateKeypair() (*KeyPair, error) {
	var keyPair KeyPair

	if _, err := rand.Read(keyPair.PrivateKey[:]); err != nil {
		err = errors.New("Error generating keypair")
		return nil, err
	}

	curve25519.ScalarBaseMult(&keyPair.PublicKey, &keyPair.PrivateKey)

	return &keyPair, nil
}

func DH(privateKey, publicKey []byte) ([]byte, error) {
	sharedKey, err := curve25519.X25519(privateKey, publicKey)
	if err != nil {
		err = errors.New("Error generating shared key")
		return nil, err
	}

	return sharedKey, err
}

