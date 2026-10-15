# OpenTask

OpenTask is a single-user, local task tracker with a fast command-line workflow and an optional interactive dashboard. It stores tasks in SQLite and requires no account, authentication, or network service.

## Install

For local testing, you do not need to push the code to GitHub or install `task@latest`.
Run the CLI directly from this checkout:

```sh
go run ./cmd/task --help
go run ./cmd/task --db /tmp/opentask-test.db add "Test local task" --priority high
go run ./cmd/task --db /tmp/opentask-test.db list
```

The database path can be changed for each test run. Using a temporary path keeps test data out of the repository.

To build and run a local binary:

```sh
mkdir -p bin
go build -o bin/task ./cmd/task
./bin/task --help
./bin/task --db /tmp/opentask-test.db list
```

Run the automated checks from the repository root:

```sh
go test ./...
go vet ./...
```

To test the interactive dashboard locally:

```sh
go run ./cmd/task --db /tmp/opentask-test.db dashboard
```

Press `q` or `Ctrl+C` to exit. Use `r` to refresh, `↑`/`↓` (or `j`/`k`) to navigate, and `Enter` or `d` to toggle the selected task.

## Install a local binary

```sh
go install github.com/kevinforsgren/opentask/cmd/task@latest
```

For the current checkout, install without GitHub by using:

```sh
go install ./cmd/task
```

The installed executable is placed in the directory returned by `go env GOBIN`, or in `$(go env GOPATH)/bin` when `GOBIN` is not set. Make sure that directory is on your `PATH`.

Remove the installed program with this command. It handles both the default Go bin directory and a custom `GOBIN`:

```sh
bin="$(go env GOBIN)"; if [ -z "$bin" ]; then bin="$(go env GOPATH)/bin"; fi; rm -f "$bin/task"
```

To remove only the local build artifacts and test database:

```sh
rm -f ./bin/task /tmp/opentask-test.db
```

## Usage

```sh
task add "Write project plan" --priority high
task add "Review notes" --description "Check open questions" --priority medium
task list
task list --status pending --priority high
task list --search project
task update 1 --title "Write implementation plan"
task done 1
task toggle 2
task delete 2 --force
task dashboard
```

The database is initialized automatically on first use. By default it is `./tasks.db`. Set `TASK_DB_PATH` or pass `--db PATH` to use another location:

```sh
TASK_DB_PATH="$HOME/.local/share/opentask/tasks.db" task list
task --db ./work/tasks.db list
```

Completed tasks are dimmed in terminal output. Color and Unicode are optional presentation enhancements; use `--no-color` for plain output and scripting environments.

## Development

```sh
go test ./...
go vet ./...
```

The application deliberately excludes authentication, accounts, synchronization, remote APIs, recurring tasks, reminders, and multi-user permissions.
