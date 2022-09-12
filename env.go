package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

// Function to load environment variables from a .env file from the path provided.
func LoadCredentialsCustomEnvPath(path string) (*Credentials, error) {
	if !strings.HasSuffix(path, ".env") {
		return nil, errors.New("path does not point to a .env file")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	pairs := scanPairsFromFile(file)
	err = file.Close()
	if err != nil {
		return nil, err
	}
	return credentialsFromMap(pairs)

}

// Function to load pairs in the form of KEY=VALUE from a reader into a map
func scanPairsFromFile(reader io.Reader) map[string]string {
	scanner := bufio.NewScanner(reader)
	pairs := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			continue
		}
		name, value := strings.ToLower(strings.TrimSpace(tokens[0])), strings.TrimSpace(tokens[1])
		pairs[name] = value
	}
	return pairs
}

func credentialsFromMap(pairs map[string]string) (*Credentials, error) {
	var c Credentials
	if pairs["username"] == "" {
		return nil, errors.New("found no value for `username` in .env file")
	}
	c.Username = pairs["username"]
	if pairs["password"] == "" {
		return nil, errors.New("found no value for `password` in .env file")
	}
	c.Password = pairs["password"]
	if pairs["dbname"] == "" {
		return nil, errors.New("found no value for `dbname` in .env file")
	}
	c.DatabaseName = pairs["dbname"]
	if pairs["host"] == "" {
		return nil, errors.New("found no value for `host` in .env file")
	}
	c.Host = pairs["host"]
	if pairs["port"] == "" {
		return nil, errors.New("found no value for `port` in .env file")
	}
	c.Port = pairs["port"]
	return &c, nil
}
