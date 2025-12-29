# Kahn
Terminal-based kanban task management application built with Go.

## Features
- Project management with custom names and descriptions
- Three-column kanban board (Not Started, In Progress, Done)
- Task prioritization with Low/Medium/High levels
- SQLite persistence with WAL mode for performance
- Clean terminal UI with keyboard navigation
- Configuration via file, environment variables, or flags

## Requirements
- Go 1.25.4 or higher

## Installation

### Build from Source
```bash
git clone github.com/mattlake/kahn
cd kahn
go build -o kahn .
```

## Quick Start

```bash
# Run the application
./kahn

# The app creates a default project on first run
# Use 'n' to create your first task
```

## Usage

### Navigation
- `h/l` - Navigate between kanban columns
- `j/k` - Navigate tasks within a column
- `space` - Move selected task to next status
- `backspace` - Move selected task to previous status

### Task Management
- `n` - Create new task
- `e` - Edit selected task
- `d` - Delete selected task

### Project Management
- `p` - Switch between projects
- `p` → `n` - Create new project
- `p` → `d` - Delete current project

### Other
- `q` - Quit application
- `esc` - Cancel dialogs/forms

## Configuration

Kahn uses Viper for configuration. Create `~/.kahn/config.toml`:

```toml
[database]
path = "~/.kahn/kahn.db"
journal_mode = "WAL"
busy_timeout = 5000
```

Configuration sources (priority order):
1. Command-line flags
2. Environment variables (`KAHN_` prefix)
3. Config files
4. Default values

## License

This project is licensed under the Non-Profit Open Software License 3.0 (NPOSL-3.0).

This license permits:
- ✅ Free usage and modification
- ✅ Redistribution for non-commercial purposes
- ✅ Private and educational use

This license prohibits:
- ❌ Commercial selling or profit-making
- ❌ Commercial distribution
- ❌ Commercial exploitation

The software is provided **without warranty**. See the [LICENSE](LICENSE) file for full details.
