package encryption

import (
	"testing"
)	

var ciphertext []byte

func TestEncrypt(t *testing.T) {
	var err error
	plaintext := []byte("Hello")
	key := []byte("1234567890123456")
	ciphertext, err = Encrypt(plaintext, key)
	if err != nil {
		t.Error(err)
	}

	if len(ciphertext) == 0 {
		t.Error("ciphertext is empty")
	}	
}

func TestDecrypt(t *testing.T) {
	key := []byte("1234567890123456")
	plaintext, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Error(err)
	}

	if len(plaintext) == 0 {
		t.Error("plaintext is empty")
	}

	if string(plaintext) != "Hello" {
		t.Error("plaintext is not equal to 'Hello'")
	}
}