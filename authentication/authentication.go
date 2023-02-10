package authentication

import (
	"crypto"
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha256"
)

// GenerateKeyPair generates a public and private key pair.
func GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, nil
}

func SignMessage(privKey *rsa.PrivateKey, message []byte) ([]byte, error) {
	hashedMessage := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashedMessage[:])
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func VerifySignature(pubKey *rsa.PublicKey, message []byte, signature []byte) error {
	hashedMessage := sha256.Sum256(message)
	err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashedMessage[:], signature)
	if err != nil {
		return err
	}

	return nil
}