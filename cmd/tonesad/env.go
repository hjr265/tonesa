package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func LoadEnvFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		ln, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		parts := strings.SplitN(ln, "=", 2)
		if len(parts) != 2 {
			continue
		}
		err = os.Setenv(parts[0], strings.TrimSpace(parts[1]))
		if err != nil {
			return err
		}
	}
	return nil
}
