# gvm-baron

A Go Version Manager for easily installing, managing, and switching between different versions of Go.

## Usage

### Basic Commands

```bash
# Show help
gvm --help

# Display current Go version
gvm current

# List available Go versions
gvm list latest|lts|stable|all

# Install a Go version
gvm install <version>

# Switch to a Go version
gvm use <version>

# Create an alias
gvm alias <source> <target>

# Remove an alias
gvm alias-delete <name>

# Remove an installed version
gvm remove <version>

# Update the versions cache
gvm refresh-versions
```

### Command Details

#### List Versions

```bash
# List stable versions (default)
gvm list
gvm list stable

# List latest versions
gvm list latest

# List LTS versions
gvm list lts

# List all versions
gvm list all
```

#### Install and Use Versions

```bash
# Install a specific version
gvm install 1.21.0
gvm install 1.20.5

# Install using shortcuts
gvm install latest
gvm install lts

# Switch to an installed version
gvm use 1.21.0
gvm use lts
```

#### Aliases

```bash
# Create an alias
gvm alias 1.21.0 goprod
gvm alias 1.20.5 godev
goprod run .

# Remove an alias
gvm alias-delete goprod
```

#### Version Management

```bash
# Remove an installed version
gvm remove 1.20.5

# This will also remove any aliases pointing to that version
```

### Flags

- `-h, --help`: Show help message
- `-d, --debug`: Enable debug output for troubleshooting
- `-n, --no-cache`: Disable caching for fresh version lookups

### Examples

```bash
# Install the latest stable Go version with debug output
gvm install latest -d

# List all versions without using cache
gvm list all -n

# Install and switch to Go 1.21.0
gvm install 1.21.0
gvm use 1.21.0

# Create development and production aliases
gvm alias 1.21.0 prod
gvm alias 1.20.8 dev

# Switch between environments
gvm use prod
gvm use dev
```

## Troubleshooting

### Enable Debug Mode

```bash
gvm install 1.21.0 --debug
```

### Clean Installation

To start fresh, remove the gvm-baron directory:

```bash
# Linux/macOS
rm -rf ~/.gvm-baron

# Windows
Remove-Item -Recurse -Force $env:USERPROFILE\.gvm-baron
```
