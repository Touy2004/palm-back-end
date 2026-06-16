package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type CryptoService struct {
	key []byte
}

func NewCryptoService(base64Key string) (*CryptoService, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes")
	}

	return &CryptoService{key: key}, nil
}

func (s *CryptoService) Encrypt(plain []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	cipherText := gcm.Seal(nil, nonce, plain, nil)

	return cipherText, nonce, nil
}

func (s *CryptoService) Decrypt(cipherText []byte, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(nonce) != gcm.NonceSize() {
		return nil, fmt.Errorf("incorrect nonce length: expected %d, got %d", gcm.NonceSize(), len(nonce))
	}

	plain, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plain, nil
}
