package crypto

import (
	"crypto/cipher"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
)

const (
	// AEADKeySize is the size of the key used by this AEAD, in bytes.
	AEADKeySize = chacha20poly1305.KeySize

	// AEADNonceSize is the size of the nonce used with the XChaCha20-Poly1305
	// variant of this AEAD, in bytes.
	AEADNonceSize = chacha20poly1305.NonceSizeX
)

// NewAEADKey generates a new random secret key.
func NewAEADKey() ([]byte, error) {
	return RandBytes(AEADKeySize)
}

// NewAEADNonce generates a new random nonce.
func NewAEADNonce() ([]byte, error) {
	return RandBytes(AEADNonceSize)
}

// NewAEAD returns a XChaCha20-Poly1305 AEAD that uses the given 256-bit key.
//
// XChaCha20-Poly1305 is a ChaCha20-Poly1305 variant that takes a longer nonce, suitable to be
// generated randomly without risk of collisions. It should be preferred when nonce uniqueness cannot
// be trivially ensured, or whenever nonces are randomly generated.
func NewAEAD(key []byte) (cipher.AEAD, error) {
	return chacha20poly1305.NewX(key)
}

// Encrypt is an helper function to symetrically encrypt a piece of data using XChaCha20-Poly1305
// returning the nonce separatly
func EncryptWithNonce(key, plaintext, additionalData []byte) (ciphertext, nonce []byte, err error) {
	nonce, err = NewAEADNonce()
	if err != nil {
		return
	}
	cipher, err := NewAEAD(key)
	if err != nil {
		return
	}
	ciphertext = cipher.Seal(nil, nonce, plaintext, additionalData)
	return
}

// DecryptWithNonce is an helper function to symetrically  decrypt a piece of data using XChaCha20-Poly1305
// taking the nonce as a separate piece of input
func DecryptWithNonce(key, nonce, ciphertext, additionalData []byte) (plaintext []byte, err error) {
	cipher, err := NewAEAD(key)
	if err != nil {
		return
	}
	plaintext, err = cipher.Open(nil, nonce, ciphertext, additionalData)
	return
}

// Encrypt is an helper function to symetrically encrypt a piece of data using XChaCha20-Poly1305
// the nonce is prepended to the ciphertext in the returned buffer
func Encrypt(key, plaintext, additionalData []byte) (ciphertext []byte, err error) {
	nonce, err := NewAEADNonce()
	if err != nil {
		return
	}
	cipher, err := NewAEAD(key)
	if err != nil {
		return
	}
	ciphertext = cipher.Seal(nil, nonce, plaintext, additionalData)
	ciphertext = append(nonce, ciphertext...)
	return
}

// DecryptWithNonce is an helper function to symetrically  decrypt a piece of data using XChaCha20-Poly1305
// The nonce should be at the begining of the ciphertext
func Decrypt(key, ciphertext, additionalData []byte) (plaintext []byte, err error) {
	cipher, err := NewAEAD(key)
	if err != nil {
		return
	}

	if len(ciphertext) < AEADNonceSize {
		err = errors.New("crypto.Decrypt: len(ciphertext) < NonceSize")
		return
	}
	nonce := ciphertext[:AEADNonceSize]
	ciphertext = ciphertext[AEADNonceSize:]

	plaintext, err = cipher.Open(nil, nonce, ciphertext, additionalData)
	return
}
