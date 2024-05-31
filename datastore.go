package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Datastore struct{}

func (d Datastore) NewUser(userId string) error {
	dirPath := filepath.Join(".", "users", userId)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("can't create user directory: %s", err)
	}
	return nil
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
