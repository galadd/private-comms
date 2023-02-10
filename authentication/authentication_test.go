package authentication

import (
	"crypto/rsa"
	"testing"
)

var (
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey

	signature []byte
	message []byte
)

func TestGenerateKeyPair(t *testing.T) {
	var err error
	privKey, pubKey, err = GenerateKeyPair()
	if err != nil {
		t.Error(err)
	}

	if privKey == nil {
		t.Error("privKey is nil")
	}

	if pubKey == nil {
		t.Error("pubKey is nil")
	}

	if privKey.PublicKey != *pubKey {
		t.Error("pubKey does not match privKey")
	}
}

func TestSignMessage(t *testing.T) {
	var err error
	message = []byte("Hello")
	signature, err = SignMessage(privKey, message)
	if err != nil {
		t.Error(err)
	}

	if len(signature) == 0 {
		t.Error("signature is empty")
	}
}

func TestVerifySignature(t *testing.T) {
	err := VerifySignature(pubKey, message, signature)
	if err != nil {
		t.Error(err)
	}
}