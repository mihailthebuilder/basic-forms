package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
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
	internalUserId, err := decrypt(userId, d.Secret)

	if err != nil {
		return fmt.Errorf("unable to decrypt userId: %s", err)
	}

	filePath := filepath.Join(".", "users", internalUserId.String(), fmt.Sprintf("%s.txt", origin))

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("can't open file: %s", err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("can't write to file: %s", err)
	}

	// Write a newline character
	_, err = file.Write([]byte("\n"))
	if err != nil {
		return fmt.Errorf("can't write new line to file: %s", err)
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
	hashedKey := sha256.Sum256([]byte(secret))
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
	ciphertext := gcm.Seal(nonce, nonce, internalId[:], nil)

	// Return the encrypted data as a hex string
	return hex.EncodeToString(ciphertext), nil
}

func decrypt(encryptedData string, secret string) (uuid.UUID, error) {
	// Convert the secret to a byte array
	hashedKey := sha256.Sum256([]byte(secret))
	key := hashedKey[:]

	// Create a new AES cipher using the secret key
	block, err := aes.NewCipher(key)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create a new GCM (Galois/Counter Mode) cipher mode instance
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decode the hex string back to bytes
	data, err := hex.DecodeString(encryptedData)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	// Extract the nonce size from GCM
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return uuid.Nil, fmt.Errorf("ciphertext too short")
	}

	// Extract the nonce and ciphertext
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt the data using GCM with the nonce
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	// Convert the decrypted plaintext back to UUID
	decryptedUUID, err := uuid.FromBytes(plaintext)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to convert decrypted data to UUID: %w", err)
	}

	// Return the decrypted UUID
	return decryptedUUID, nil
}
