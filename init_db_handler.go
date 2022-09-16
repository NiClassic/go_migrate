package main

import (
	"database/sql"
	"errors"
	"flag"
)

const (
	createStatement string = `CREATE TABLE IF NOT EXISTS migrations(
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(200) NOT NULL,
		batch SMALLINT UNSIGNED NOT NULL
	)`
)

func NewInitDbCommand() *Command {
	c := NewCommand("init-db",
		"",
		func(args []string) error {
			set := flag.NewFlagSet("init-db", flag.ContinueOnError)
			envPtr := set.String("env", ".env", "Path to the .env file containing database credentials")
			usrPtr := set.String("u", "", "Username for the database user")
			pwPtr := set.String("p", "", "Password for the database user")
			dbPtr := set.String("db", "", "Name of the database")
			portPtr := set.String("port", "3306", "Port of the database connection")
			hostPtr := set.String("host", "localhost", "Host of the database connection")
			err := set.Parse(args)
			if err != nil {
				if errors.Is(err, flag.ErrHelp) {
					return nil
				}
				return err
			}
			c, err := CredentialsFromPointers(envPtr, usrPtr, pwPtr, dbPtr, portPtr, hostPtr)
			if err != nil {
				return err
			}
			conn, err := sql.Open("mysql", c.ToDSN())
			if err != nil {
				return err
			}
			err = conn.Ping()
			if err != nil {
				return err
			}
			_, err = conn.Exec(createStatement)
			return err
		})
	return c
}
