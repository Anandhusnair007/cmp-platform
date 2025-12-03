package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encryptor handles encryption at rest
type Encryptor struct {
	key []byte
}

// NewEncryptor creates a new encryptor with a key
func NewEncryptor(key []byte) (*Encryptor, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes (AES-256)")
	}
	return &Encryptor{key: key}, nil
}

// Encrypt encrypts data using AES-256-GCM
func (e *Encryptor) Encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts data using AES-256-GCM
func (e *Encryptor) Decrypt(encryptedData string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// KeyManager manages encryption keys
type KeyManager struct {
	currentKeyID string
	keys         map[string][]byte
}

// NewKeyManager creates a new key manager
func NewKeyManager() *KeyManager {
	return &KeyManager{
		keys: make(map[string][]byte),
	}
}

// GenerateKey generates a new encryption key
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32) // AES-256
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

// SetCurrentKey sets the current encryption key
func (km *KeyManager) SetCurrentKey(keyID string, key []byte) {
	km.currentKeyID = keyID
	km.keys[keyID] = key
}

// GetEncryptor returns an encryptor with the current key
func (km *KeyManager) GetEncryptor() (*Encryptor, error) {
	if km.currentKeyID == "" {
		return nil, errors.New("no current key set")
	}
	key, exists := km.keys[km.currentKeyID]
	if !exists {
		return nil, errors.New("current key not found")
	}
	return NewEncryptor(key)
}
