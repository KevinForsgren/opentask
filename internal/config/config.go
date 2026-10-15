package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const DefaultDatabasePath = "tasks.db"

// ResolveDatabasePath applies precedence in this order:
// --db, process environment, .env, then the default path.
func ResolveDatabasePath(flagValue string) (string, error) {
	envFile, err := loadDotEnv(".env")
	if err != nil {
		return "", err
	}

	path := strings.TrimSpace(flagValue)
	baseDir := ""
	if path == "" {
		path = strings.TrimSpace(os.Getenv("TASK_DB_PATH"))
	}
	if path == "" {
		path = strings.TrimSpace(envFile["TASK_DB_PATH"])
		if path != "" {
			baseDir = filepath.Dir(envFilePath(".env"))
		}
	}
	if path == "" {
		path = DefaultDatabasePath
	}
	if !filepath.IsAbs(path) && baseDir != "" {
		path = filepath.Join(baseDir, path)
	}
	return filepath.Clean(path), nil
}

func loadDotEnv(name string) (map[string]string, error) {
	path := envFilePath(name)
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	}
	defer file.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(key) == "" {
			return nil, fmt.Errorf("%s:%d: expected KEY=VALUE", name, lineNumber)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		values[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read %s: %w", name, err)
	}
	return values, nil
}

func envFilePath(name string) string {
	if configured := strings.TrimSpace(os.Getenv("OPENTASK_ENV_FILE")); configured != "" {
		return configured
	}
	return name
}
