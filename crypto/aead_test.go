package crypto

import "testing"

func TestAEADKeySize(t *testing.T) {
	if AEADKeySize != 32 {
		t.Error("AEADKeySize != 32")
	}
}

func TestAEADNonceSize(t *testing.T) {
	if AEADNonceSize != 24 {
		t.Error("AEADNonceSize != 24")
	}
}

func TestNewAEADKey(t *testing.T) {
	key, err := NewAEADKey()
	if err != nil {
		t.Error(err)
	}

	if len(key) != AEADKeySize {
		t.Errorf("generated key have bad size (%d)", len(key))
	}
}

func TestNewAEADNonce(t *testing.T) {
	nonce, err := NewAEADNonce()
	if err != nil {
		t.Error(err)
	}

	if len(nonce) != AEADNonceSize {
		t.Errorf("generated nonce have bad size (%d)", len(nonce))
	}
}

func TestAEADSealDstNil(t *testing.T) {
	data := []byte("this is a plaintext message")
	ad := []byte("additional data")

	nonce, err := NewAEADNonce()
	if err != nil {
		t.Error(err)
	}
	key, err := NewAEADKey()
	if err != nil {
		t.Error(err)
	}
	cipher, err := NewAEAD(key)
	if err != nil {
		t.Error(err)
	}

	ciphertext := cipher.Seal(nil, nonce, data, ad)
	if ciphertext == nil {
		t.Error("ciphertext is nil")
	}
}

func TestEncryptDecryptWithNonce(t *testing.T) {
	plaintext := []byte("this is a plaintext message")
	ad := []byte("additional data")

	key, err := NewAEADKey()
	if err != nil {
		t.Error(err)
	}

	cipherText, nonce, err := EncryptWithNonce(key, plaintext, ad)
	if err != nil {
		t.Error(err)
	}

	plaintext2, err := DecryptWithNonce(key, nonce, cipherText, ad)
	if err != nil {
		t.Error(err)
	}

	if string(plaintext) != string(plaintext2) {
		t.Errorf("bad plaintext while decrypting: %s, expedting: %s", string(plaintext2), string(plaintext))
	}
}

func TestEncryptDecrypt(t *testing.T) {
	plaintext := []byte("this is a plaintext message")
	ad := []byte("additional data")

	key, err := NewAEADKey()
	if err != nil {
		t.Error(err)
	}

	cipherText, err := Encrypt(key, plaintext, ad)
	if err != nil {
		t.Error(err)
	}

	plaintext2, err := Decrypt(key, cipherText, ad)
	if err != nil {
		t.Error(err)
	}

	if string(plaintext) != string(plaintext2) {
		t.Errorf("bad plaintext while decrypting: %s, expedting: %s", string(plaintext2), string(plaintext))
	}
}
