package domain

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

type aesCrypto struct {
	gcm   cipher.AEAD
	nonce []byte
}

var cryptPassphrase string
var aesCrypt *aesCrypto

func SetCryptPassphrase(passphrase string) error {
	cryptPassphrase = createHash(passphrase)
	return loadAesCrypt()
}

func init() {
	cryptPassphrase = createHash("defaultPassphrase")
	err := loadAesCrypt()
	if err != nil {
		panic(err)
	}
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func loadAesCrypt() error {
	block, err := aes.NewCipher([]byte(cryptPassphrase))
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	aesCrypt = &aesCrypto{gcm: gcm, nonce: nonce}
	return nil
}

func encrypt(data string) (string, error) {
	return aesCrypt.encrypt(data)
}

func decrypt(data string) (string, error) {
	return aesCrypt.decrypt(data)
}

func (crypto *aesCrypto) encrypt(data string) (string, error) {
	ciphertext := crypto.gcm.Seal(nil, crypto.nonce, []byte(data), nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

func (crypto *aesCrypto) decrypt(data string) (string, error) {
	decoded, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	plaintext, err := crypto.gcm.Open(nil, crypto.nonce, decoded, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
