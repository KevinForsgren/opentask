package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDatabasePathFromDotEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("TASK_DB_PATH=data/tasks.db\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("OPENTASK_ENV_FILE", envPath)
	t.Setenv("TASK_DB_PATH", "")

	path, err := ResolveDatabasePath("")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "data", "tasks.db")
	if path != want {
		t.Fatalf("path = %q, want %q", path, want)
	}
}

func TestResolveDatabasePathFlagWins(t *testing.T) {
	t.Setenv("OPENTASK_ENV_FILE", filepath.Join(t.TempDir(), ".env"))
	t.Setenv("TASK_DB_PATH", "from-environment.db")

	path, err := ResolveDatabasePath("from-flag.db")
	if err != nil {
		t.Fatal(err)
	}
	if path != "from-flag.db" {
		t.Fatalf("path = %q, want %q", path, "from-flag.db")
	}
}
