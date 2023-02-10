package encryption

import (
	"bytes"
	"testing"
)

func TestGenerateKeypair(t *testing.T) {
	_, err := GenerateKeypair()
	if err != nil {
		t.Errorf("Error generating keypair: %s", err)
	}
}

func TestDH(t *testing.T) {
	aliceKeyPair, _ := GenerateKeypair()
	bobKeyPair, _ := GenerateKeypair()

	sharedKey1, err := DH(aliceKeyPair.PrivateKey[:], bobKeyPair.PublicKey[:])
	if err != nil {
		t.Errorf("Error generating shared key: %s", err)
	}

	sharedKey2, err := DH(bobKeyPair.PrivateKey[:], aliceKeyPair.PublicKey[:])
	if err != nil {
		t.Errorf("Error generating shared key: %s", err)
	}

	if !bytes.Equal(sharedKey1, sharedKey2) {
		t.Errorf("Shared keys are not equal")
	}

	if len(sharedKey1) != dhLen {
		t.Errorf("Shared key length is not %d bytes", dhLen)
	}
}