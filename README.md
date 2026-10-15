# OpenTask

OpenTask is a fast, local command-line task manager for one person. It stores your tasks in SQLite and provides both simple commands and an optional interactive dashboard.

## Features

- Create, view, update, complete, toggle, and delete tasks
- Set low, medium, or high priority
- Filter by status or priority
- Search titles and descriptions
- Styled terminal tables with readable status and priority indicators
- Optional keyboard-driven dashboard
- Automatic SQLite database initialization
- Works on Linux, macOS, and Windows
- No account, login, network connection, or authentication required

## Requirements

- Go 1.26 or newer
- A terminal

## Quick start

Clone or download the project, then run these commands from the project directory:

```sh
go run ./cmd/task add "Set up OpenTask"
go run ./cmd/task list
go run ./cmd/task done 1
```

The first command creates the database automatically. You do not need to push the project to GitHub or install a published version to use it locally.

## Configuration

OpenTask reads the database location from a `.env` file in the directory where you run the command. Create it by copying the example:

### Linux and macOS

```sh
cp .env.template .env
```

### Windows PowerShell

```powershell
Copy-Item .env.template .env
```

Then edit `.env`:

```dotenv
TASK_DB_PATH=data/opentask.db
```

The `data` directory is created automatically. Relative paths are resolved from the directory containing `.env`, so the same configuration works regardless of the current terminal directory. Absolute Windows paths such as `C:\Users\you\Documents\opentask.db` are also supported.

Configuration precedence is:

1. `--db PATH` command-line option
2. `TASK_DB_PATH` environment variable
3. `TASK_DB_PATH` in `.env`
4. `tasks.db` in the current directory

Examples:

```sh
go run ./cmd/task list
go run ./cmd/task --db ./temporary.db list
```

Keep `.env` private. It is ignored by Git; `.env.template` is the shareable template.

## Commands

### Add tasks

```sh
task add "Write project plan"
task add "Review release notes" --description "Check open questions" --priority high
```

When using the source checkout, replace `task` with `go run ./cmd/task`.

### List, filter, and search

```sh
task list
task list --status pending
task list --status done
task list --priority high
task list --search release
```

### Update and complete

```sh
task update 1 --title "Write implementation plan"
task update 1 --description "Include database decisions"
task update 1 --priority medium
task update 1 --clear-description
task done 1
task toggle 1
```

### Delete

Deletion requires an explicit confirmation flag:

```sh
task delete 1 --force
```

### Dashboard

```sh
task dashboard
```

Use `↑`/`↓` or `j`/`k` to navigate, `Enter` or `d` to toggle a task, `r` to refresh, and `q` or `Ctrl+C` to exit.

For plain output without color:

```sh
task --no-color list
```

## Local development

Run the automated checks from the project directory:

```sh
go test ./...
go vet ./...
go build -o bin/task ./cmd/task
```

On Windows PowerShell, use this build command instead:

```powershell
go build -o bin\task.exe .\cmd\task
```

Run the locally built program:

```sh
./bin/task list
```

To use the shorter `task` command without installing it, add the project's `bin` directory to your current shell's `PATH`:

```sh
export PATH="$PWD/bin:$PATH"
task list
```

This lasts only for the current terminal session. To make it permanent, add the same `export PATH=...` line to your shell startup file, such as `~/.bashrc` or `~/.zshrc`, then open a new terminal.

```powershell
.\bin\task.exe list
```

In Windows PowerShell, use the local executable directly:

```powershell
.\bin\task.exe list
```

To use `task` in the current PowerShell session:

```powershell
$env:Path = "$PWD\bin;$env:Path"
task.exe list
```

For a permanent Windows PATH entry, add the full project `bin` directory through **System Settings → Environment Variables → Path**, then open a new terminal.

## Install and uninstall

To install the current checkout into Go's local binary directory:

```sh
go install ./cmd/task
```

Make sure the directory printed by `go env GOBIN` or `go env GOPATH` plus `/bin` is on your `PATH`.

For zsh on Linux/macOS, the usual setup is:

```sh
export PATH="$(go env GOPATH)/bin:$PATH"
rehash
task --help
```

To make this permanent for future zsh sessions:

```sh
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
rehash
```

Verify the installation location with:

```sh
command -v task
ls "$(go env GOPATH)/bin/task"
```

After that, you can run the program from any directory:

```sh
task list
```

On Windows PowerShell:

```powershell
task.exe list
```

On Linux and macOS, remove the installed program with:

```sh
bin="$(go env GOBIN)"; if [ -z "$bin" ]; then bin="$(go env GOPATH)/bin"; fi; rm -f "$bin/task"
```

On Windows PowerShell, remove it with:

```powershell
$bin = go env GOBIN
if ([string]::IsNullOrWhiteSpace($bin)) { $bin = Join-Path (go env GOPATH) "bin" }
Remove-Item (Join-Path $bin "task.exe") -ErrorAction SilentlyContinue
```

## Project layout

```text
cmd/task/              Application entry point
internal/cli/          Cobra commands and flag handling
internal/config/       .env and database-path resolution
internal/storage/      SQLite connection and task repository
internal/task/         Domain models, validation, and services
internal/ui/           Styled output and dashboard
db/migrations/         Database migration source
```

## License

See [LICENSE](LICENSE).
