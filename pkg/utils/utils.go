package utils

import (
	"bufio"
	"os"
	"strings"
)

// SplitDelimiterList splits a list of items by an additional delimiter
func SplitDelimiterList(items []string, delim string) (map[string]string, error) {

	data := map[string]string{}
	for _, item := range items {
		pieces := strings.Split(item, delim)
		if len(pieces) > 1 {
			data[pieces[0]] = pieces[1]
		} else {
			data[pieces[0]] = ""
		}
	}
	return data, nil
}

// ParseConfigFile parses a simple configuration file, with newlines for each thing,
// and a starting prefix to determine comma, and some other delimiter to determine
// key value pairs
func ParseConfigFile(path, comment, delim string) (map[string]string, error) {

	// We will return key value pairs
	data := map[string]string{}

	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	// Read through with a scanner, this looks to be a parse file
	s := bufio.NewScanner(fd)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Text()

		// Skip over lines that are commented out
		if strings.HasPrefix(line, comment) {
			continue
		}

		// Parse lines (splitting by delim) that are not
		parts := strings.Split(line, delim)
		if len(parts) >= 2 {
			data[parts[0]] = parts[1]
		}
	}
	return data, nil
}
