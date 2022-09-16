package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

func NewUpCommand() *Command {
	c := NewCommand("up", "", func(args []string) error {
		set := flag.NewFlagSet("up", flag.ContinueOnError)
		dirPtr := set.String("d", "database", "Path to the migration directory")
		envPtr := set.String("env", ".env", "Path to the .env file containing database credentials")
		usrPtr := set.String("u", "", "Username for the database user")
		pwPtr := set.String("p", "", "Password for the database user")
		dbPtr := set.String("db", "", "Name of the database")
		nPtr := set.Int("n", -1, "Number of migrations to run")
		portPtr := set.String("port", "3306", "Port of the database connection")
		hostPtr := set.String("host", "localhost", "Host of the database connection")
		pretendPtr := set.Bool("pretend", false, "Print the sql-statements instead of executing them")
		err := set.Parse(args)
		if err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return nil
			}
			return err
		}
		if *dirPtr == "" {
			return fmt.Errorf("migration directory path cannot be empty")
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
		ranMigrations, err := MigrationsFromDatabase(conn)
		if err != nil {
			return err
		}
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		migrationPath := path.Join(wd, *dirPtr)
		mFiles, err := os.ReadDir(migrationPath)
		unranMigrations := make([]string, 0, len(mFiles))
		if err != nil {
			return err
		}
		for _, file := range mFiles {
			if file.IsDir() || strings.Contains(file.Name(), "down") || !strings.Contains(file.Name(), "up") {
				continue
			}
			name, _, found := strings.Cut(file.Name(), ".")
			if !found {
				return fmt.Errorf("file %s is not a valid sql file", file.Name())
			}
			hasMatch := false
			for _, m := range ranMigrations {
				if strings.EqualFold(m.Name, name) {
					hasMatch = true
				}
			}
			if !hasMatch {
				unranMigrations = append(unranMigrations, file.Name())
			}

		}
		if len(unranMigrations) < 1 {
			fmt.Println("No unran migrations")
		}
		batchId := GetBatchId(&ranMigrations)
		for i, m := range unranMigrations {
			if *nPtr > 0 {
				if i > *nPtr-1 {
					break
				}
			}
			name, _, found := strings.Cut(m, ".")
			if !found {
				continue
			}
			fmt.Printf("Scanning migration file %s\n", name)
			f, err := os.Open(path.Join(migrationPath, m))
			if err != nil {
				return err
			}
			statements, err := ReadStatementsFromReader(f)
			if err != nil {
				return err
			}
			fmt.Printf("Found %d statement(s)\n", len(statements))
			for _, stmt := range statements {
				if *pretendPtr {
					fmt.Printf("%s%s%s\n", Green, stmt, Reset)
					continue
				}
				_, err := conn.Exec(stmt)
				if err != nil {
					return err
				}
			}
			if *pretendPtr {
				continue
			}
			err = f.Close()
			if err != nil {
				return err
			}

			res, err := conn.Exec("INSERT INTO migrations (name, batch) VALUES (?,?)", name, batchId)
			if err != nil {
				return err
			}
			rows, err := res.RowsAffected()
			if err != nil {
				return err
			}
			if rows < 1 {
				return fmt.Errorf("migration %s could not be stored in the database", m)
			}
		}
		return nil
	})
	return c
}
