package main

import (
	"fmt"
)

// Struct that holds information about the database authentication
type Credentials struct {
	Username     string
	Password     string
	DatabaseName string
	Port         string
	Host         string
}

func NewCredentials(username, password, databaseName, port, host string) *Credentials {
	return &Credentials{
		Username:     username,
		Password:     password,
		DatabaseName: databaseName,
		Port:         port,
		Host:         host,
	}
}

// Method to format the credential information into a valid mysql-dsn string,
// e.g. root:root@tcp(127.0.0.1:3306)/demo.
// See: https://github.com/go-sql-driver/mysql#dsn-data-source-name
func (c *Credentials) ToDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.Host, c.Port, c.DatabaseName)
}
