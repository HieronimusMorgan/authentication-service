package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// Encryption interface defines the methods for encryption, decryption, and hashing
type Encryption interface {
	Encrypt(text string) (string, error)
	Decrypt(encryptedText string) (string, error)
	HashPhoneNumber(phone string) string
	HashPassword(password string) (string, error)
	CheckPassword(hash, password string) error
}

// encryption struct contains AES key and IV
type encryption struct {
	key []byte
	iv  []byte
}

// NewEncryption initializes encryption with a 32-byte AES key and 16-byte IV
func NewEncryption(key string, iv string) Encryption {
	hashedKey := sha256.Sum256([]byte(key)) // Ensure 32-byte key for AES-256
	hashedIV := sha256.Sum256([]byte(iv))   // Ensure IV is at least 16 bytes
	return &encryption{
		key: hashedKey[:],
		iv:  hashedIV[:aes.BlockSize], // Use first 16 bytes of the hash
	}
}

// Encrypt encrypts text using AES-256 with a fixed IV (Deterministic)
func (a *encryption) Encrypt(text string) (string, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", errors.New("failed to create AES cipher: " + err.Error())
	}

	cipherText := make([]byte, len(text))

	// Use fixed IV for deterministic encryption
	stream := cipher.NewCFBEncrypter(block, a.iv)
	stream.XORKeyStream(cipherText, []byte(text))

	// Encode as Base64 for safe storage
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts an AES-256 encrypted string using the fixed IV
func (a *encryption) Decrypt(encryptedText string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", errors.New("failed to decode base64 data: " + err.Error())
	}

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", errors.New("failed to create AES cipher: " + err.Error())
	}

	// Use the same fixed IV to decrypt
	stream := cipher.NewCFBDecrypter(block, a.iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

// HashPhoneNumber hashes a phone number using SHA-256 for deterministic lookups
func (a *encryption) HashPhoneNumber(phone string) string {
	hash := sha256.Sum256([]byte(phone))
	return hex.EncodeToString(hash[:])
}

// HashPassword securely hashes a password using bcrypt
func (a *encryption) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password: " + err.Error())
	}
	return string(hashed), nil
}

// CheckPassword verifies a bcrypt hashed password
func (a *encryption) CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
