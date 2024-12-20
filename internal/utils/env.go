package utils

import (
	"bufio"
	"os"
	"strings"
)

func LoadEnvFile(filePath string) error {
	var	file	*os.File
	var	scanner	*bufio.Scanner
	var	line	string
	var	parts	[]string
	var	key		string
	var	value	string
	var	err		error

	if file, err = os.Open(filePath); err != nil {
		return err
	}
	defer file.Close()
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts = strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key = strings.TrimSpace(parts[0])
		value = strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)
		os.Setenv(key, value)
	}
	return scanner.Err()
}
