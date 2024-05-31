package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/google/uuid"
)

type Datastore struct {
	Secret string
}

type User struct {
	InternalId string `json:"internalId"`
	ExternalId string `json:"externalId"`
}

func (d Datastore) newUser() (User, error) {
	user := User{}

	internalId := uuid.New()

	user.InternalId = internalId.String()

	externalId, err := encrypt(internalId, d.Secret)
	if err != nil {
		return user, fmt.Errorf("can't encrypt internalId: %s", err)
	}

	user.ExternalId = externalId
	return user, nil
}

func encrypt(internalId uuid.UUID, secret string) (string, error) {
	// Convert the secret to a byte array
	key := []byte(secret)

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
	ciphertext := gcm.Seal(nonce, nonce, internalId[:], nil)

	// Return the encrypted data as a hex string
	return hex.EncodeToString(ciphertext), nil
}
