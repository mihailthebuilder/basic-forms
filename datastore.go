package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Datastore struct {
	Secret string
}

type User struct {
	InternalId string `json:"internalId"`
	ExternalId string `json:"externalId"`
}

func (d Datastore) NewUser() (User, error) {
	user := User{}

	internalId := uuid.New()

	user.InternalId = internalId.String()

	// Create the directory if it doesn't exist
	dirPath := filepath.Join(".", "users", user.InternalId)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return user, fmt.Errorf("can't create user directory: %s", err)
	}

	externalId, err := encrypt(internalId, d.Secret)
	if err != nil {
		return user, fmt.Errorf("can't encrypt internalId: %s", err)
	}

	user.ExternalId = externalId
	return user, nil
}

func (d Datastore) AddSubmission(userId string, origin string, content []byte) error {
	filePath := filepath.Join(".", "users", userId, fmt.Sprintf("%s.txt", origin))

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("can't open file: %s", err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("can't write to file: %s", err)
	}

	return nil
}

func (d Datastore) GetSubmissions(userId string, origin string) ([]byte, error) {
	filePath := filepath.Join(".", "users", userId, fmt.Sprintf("%s.txt", origin))

	// check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil
	}

	content, err := os.ReadFile(filePath)
	return content, err
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
