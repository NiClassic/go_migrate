package main

import (
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//TODO: reverse migrations

const (
	usage string = `Usage: migrator [SUBCOMMAND] [options] ...
		
Migrator is a simple CLI-tool written in Golang to manage database migrations within a project.
It supports mysql databases.
	
Available SUBCOMMANDS are:
	
	create		Create two timestamped migration files, one for running the migration and one for reversing the migration respectively
	init-db		Create the migration table for the desired database
	`
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(0)
	}
	createCmd := NewCreateCommand()
	initCmd := NewInitDbCommand()
	upCmd := NewUpCommand()
	switch strings.ToLower(os.Args[1]) {
	case createCmd.Name:
		err := createCmd.Handler(os.Args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case initCmd.Name:
		err := initCmd.Handler(os.Args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case upCmd.Name:
		err := upCmd.Handler(os.Args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	default:
		fmt.Fprintf(os.Stderr, "Command %s is not a valid sub-command\n", os.Args[1])
		os.Exit(0)
	}
}
