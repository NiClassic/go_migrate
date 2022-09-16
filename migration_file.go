package main

import (
	"bufio"
	"io"
	"strings"
)

// Function to read sql statements from an io.Reader. It splits statements using a semicolon as a separator.
func ReadStatementsFromReader(r io.Reader) ([]string, error) {
	s := bufio.NewScanner(r)
	s.Split(SplitSQLStatements)
	var statements []string
	for s.Scan() {
		stmt := strings.TrimSpace(s.Text())
		statements = append(statements, stmt)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return statements, nil
}
