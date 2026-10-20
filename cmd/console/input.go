package main

import (
	"bufio"
	"io"
	"strings"
)

func readInput(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		if scanner.Err() != nil {
			return "", scanner.Err()
		}
		return "", io.EOF
	}
	return strings.TrimSpace(scanner.Text()), nil
}
