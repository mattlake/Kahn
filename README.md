# Kahn

[![Go Version](https://img.shields.io/badge/Go-1.25.5+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/mattlake/kahn/actions)

A fast, elegant terminal-based kanban task management application built with Go, providing efficient task management without leaving your command line.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Screenshots

![Kahn kanban board interface showing three columns with sample tasks](https://github.com/user-attachments/assets/2fb93e7a-c419-4247-b31f-d4e6758c1905)

## Features
- Project management with custom names and descriptions
- Three-column kanban board (Not Started, In Progress, Done)
- Task prioritization with Low/Medium/High levels
- Clean terminal UI with keyboard navigation
- Multiplatform support (Linux, macOS, Windows)
- Flexible configuration via file, environment variables, or flags

## Requirements (when built from source)
- Go 1.25.5 or higher
- C compiler (GCC/Clang on Unix, MinGW/TDM-GCC on Windows)

## Installation

### From Source
```bash
git clone https://github.com/mattlake/kahn
cd kahn
go install .
```

### From Releases
Download pre-compiled binaries from the [Releases](https://github.com/mattlake/kahn/releases) page for your platform.

## Usage

### Navigation
| Key(s) | Action |
|--------|--------|
| `h` / `l` | Navigate between kanban columns |
| `j` / `k` | Navigate tasks within a column |
| `space` | Move selected task to next status |
| `backspace` | Move selected task to previous status |

### Task Management
| Key(s) | Action |
|--------|--------|
| `n` | Create new task |
| `e` | Edit selected task |
| `d` | Delete selected task |

### Project Management
| Key(s) | Action |
|--------|--------|
| `p` | Switch between projects |
| `p` → `n` | Create new project |
| `p` → `d` | Delete current project |

### Other
- `q` - Quit application
- `esc` - Cancel dialogs/forms

## Configuration

### Database Location

**Default Database Path:**
- **Windows**: `%USERPROFILE%\.kahn\kahn.db`
- **Linux/macOS/WSL2**: `~/.kahn/kahn.db`

### Database Path Configuration

You can customize the database location as shown below:

**Linux/macOS/WSL2 (`~/.kahn/config.toml`):**
```toml
[database]
# Home directory location
path = "~/.kahn/kahn.db"

# Custom absolute path
path = "/home/yourname/tasks/my-kahn.db"

# Using shell expansion
path = "~/Dropbox/tasks/kahn.db"

# System-wide shared database
path = "/opt/kahn/shared.db"
```

**Windows (`%USERPROFILE%\.kahn\config.toml`):**
```toml
[database]
# Default Windows path (auto-expanded)
path = "~/.kahn/kahn.db"

# Custom path with forward slashes
path = "C:/Users/YourName/Documents/Tasks/kahn.db"

# Different drive
path = "D:/Work/Project Management/kahn.db"
```

### Config File Locations
Search order: `./config.toml` → `~/.kahn/config.toml` → `/etc/kahn/config.toml`

### Configuration Priority (High to Low)
1. **Command-line flags** (highest priority)
2. **Environment variables** (`KAHN_` prefix)
3. **Config files**
4. **Default values** (lowest priority)

### Environment Variables

You can use environment variables with the `KAHN_` prefix:

```bash
export KAHN_DATABASE_PATH="/custom/path/kahn.db"
export KAHN_DATABASE_BUSY_TIMEOUT="3000"
```

## Contributing

We welcome contributions! Please follow these guidelines:

1. **Fork** the repository and create a feature branch
2. **Add tests** for new functionality
3. **Ensure all tests pass** with `go test ./...`
4. **Format code** with `go fmt ./...`
5. **Submit a pull request** with a clear description

### Development Guidelines
- Follow the existing code style and patterns
- Write comprehensive tests for new features
- Use table-driven tests for multiple scenarios
- Document public APIs and complex logic

## Troubleshooting

### Common Issues

**Database Permission Errors**
```bash
# Ensure the ~/.kahn directory exists and has proper permissions
mkdir -p ~/.kahn
chmod 755 ~/.kahn
```

### Having Issues?

- Open an issue on [GitHub](https://github.com/mattlake/kahn/issues)
- Check existing issues for solutions

## License

This project is licensed under the MIT License.

This license permits:
- ✅ Free commercial use and modification
- ✅ Distribution and private use
- ✅ Sublicensing and patent rights
- ✅ No liability warranty

The only requirement is including the original copyright notice in redistributions.

See the [LICENSE](LICENSE) file for full details.
