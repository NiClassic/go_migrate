package main

import "bytes"

//Split function to retrieve sql using a bufio.Scanner. Scans text until a semicolon is found.
//	r := bufio.NewScanner(f)
// 	r.Split(SplitSQLStatements)
//	for r.Scan() {
//		statement := r.Text()
//	}
func SplitSQLStatements(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, ';'); i >= 0 {
		return i + 1, data[0:i], nil
	}
	return 0, nil, nil
}
