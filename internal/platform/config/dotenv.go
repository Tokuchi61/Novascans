package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func LoadDefaultEnvFile() error {
	return LoadEnvFile(".env")
}

func LoadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return err
		}

		return fmt.Errorf("open env file %q: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for lineNo := 1; scanner.Scan(); lineNo++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("parse env file %q line %d: missing '=' separator", path, lineNo)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			return fmt.Errorf("parse env file %q line %d: empty key", path, lineNo)
		}

		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		if unquoted, err := unquoteEnvValue(value); err == nil {
			value = unquoted
		}

		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set env value for %q: %w", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan env file %q: %w", path, err)
	}

	return nil
}

func unquoteEnvValue(value string) (string, error) {
	if len(value) < 2 {
		return value, nil
	}

	quote := value[0]
	if quote != '"' && quote != '\'' {
		return value, nil
	}

	if value[len(value)-1] != quote {
		return "", fmt.Errorf("unterminated quoted value")
	}

	return value[1 : len(value)-1], nil
}
