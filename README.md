# Kahn
Terminal-based kanban task management application built with Go.

## Features
- Project management with custom names and descriptions
- Three-column kanban board (Not Started, In Progress, Done)
- Task prioritization with Low/Medium/High levels
- SQLite persistence with WAL mode for performance
- Clean terminal UI with keyboard navigation
- Multiplatform support (Linux, macOS, Windows)
- Flexible configuration via file, environment variables, or flags
- Configurable database location

## Requirements
- Go 1.25.4 or higher
- C compiler (GCC/Clang on Unix, MinGW/TDM-GCC on Windows)

## Quick Start

### Linux / macOS / WSL2
```bash
# Clone, build, run
git clone https://github.com/mattlake/kahn
cd kahn
go build -o kahn .
./kahn
```

### Windows (Choose One)

**Native Windows Build**
```powershell
# Clone, build, run
git clone https://github.com/mattlake/kahn
cd kahn
go build -o kahn.exe
.\kahn.exe
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

### Path Format Guidelines

**Recommended Practices:**
- **Windows**: Use forward slashes `"C:/Users/Name/tasks.db"`
- **Cross-platform**: Tilde expansion works everywhere `~/.kahn/`
- **Absolute paths**: Use platform-appropriate format
- **Relative paths**: Work relative to execution directory `./tasks.db`

### Config File Locations
Search order: `./config.toml` → `~/.kahn/config.toml` → `/etc/kahn/config.toml`

### Configuration Priority (High to Low)
1. **Command-line flags** (highest priority)
2. **Environment variables** (`KAHN_` prefix)
3. **Config files**
4. **Default values** (lowest priority)

## License

This project is licensed under the MIT License.

This license permits:
- ✅ Free commercial use and modification
- ✅ Distribution and private use
- ✅ Sublicensing and patent rights
- ✅ No liability warranty

The only requirement is including the original copyright notice in redistributions.

See the [LICENSE](LICENSE) file for full details.
