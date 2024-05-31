package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

type Encryption struct {
	Secret string
}

func (e Encryption) Encrypt(input string) (string, error) {
	// Convert the secret to a byte array
	hashedKey := sha256.Sum256([]byte(e.Secret))
	key := hashedKey[:]

	// Create a new AES cipher using the secret key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create a new GCM (Galois/Counter Mode) cipher mode instance
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate a nonce using GCM's standard nonce size
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the data using GCM with the nonce
	ciphertext := gcm.Seal(nonce, nonce, []byte(input), nil)

	// Return the encrypted data as a hex string
	return hex.EncodeToString(ciphertext), nil
}

func (e Encryption) Decrypt(input string) (string, error) {
	// Convert the secret to a byte array
	hashedKey := sha256.Sum256([]byte(input))
	key := hashedKey[:]

	// Create a new AES cipher using the secret key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create a new GCM (Galois/Counter Mode) cipher mode instance
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decode the hex string back to bytes
	data, err := hex.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex string: %w", err)
	}

	// Extract the nonce size from GCM
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Extract the nonce and ciphertext
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt the data using GCM with the nonce
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}
