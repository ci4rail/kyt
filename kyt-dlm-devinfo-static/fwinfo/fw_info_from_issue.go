package fwinfo

import (
	"bufio"
	"errors"
	"os"
)

const fileName = "/etc/issue"

// Read reads the firmware information from /etc/issue file
func Read() (string, error) {
	return readFile(fileName)
}

// take the firmware version from 2nd line of f
// return content of 2nd line
func readFile(f string) (string, error) {
	fwinfo := ""
	file, err := os.Open(f)

	if err != nil {
		return fwinfo, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	line := 1
	for scanner.Scan() {
		txt := scanner.Text()
		if line == 2 {
			fwinfo = txt
		}
		line++
	}
	if fwinfo == "" {
		err = errors.New("No fwinfo line found")
	}
	return fwinfo, err
}
