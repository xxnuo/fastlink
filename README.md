# fastlink

A fast and efficient tool to move files/directories to a storage location and create symbolic links at the original location, helping you save disk space while maintaining file accessibility.

[简体中文](README.zh.md)

## Features

- 🚀 **Fast file/directory operations** using [github.com/spf13/fileflow](https://github.com/spf13/fileflow)
- 🔗 **Automatic symbolic link creation** after moving files
- 📁 **Batch operations** on multiple files and directories
- ⚙️ **Flexible configuration** via YAML config file
- 🛡️ **Safety checks** to prevent recursive operations and data loss
- 📋 **Copy mode** option to keep original files

## Installation

### From Source

```bash
git clone https://github.com/xxnuo/fastlink.git
cd fastlink
make install
```

### Manual Installation

```bash
go build -o fastlink main.go
sudo cp fastlink /usr/local/bin/fastlink
chmod +x /usr/local/bin/fastlink
```

## Quick Start

1. **Create configuration file** (optional but recommended):
   ```bash
   mkdir -p ~/.config/fastlink
   cp config.sample.yaml ~/.config/fastlink/config.yaml
   # Edit the config file to set your default destination
   ```

2. **Move a file and create symbolic link**:
   ```bash
   fastlink /path/to/large-file.zip
   ```

3. **Copy mode (keep original file)**:
   ```bash
   fastlink -k /path/to/important-file.pdf
   ```

## Usage

```bash
fastlink [-k|--keep] <source> [<destination>]
```

### Options

- `-k, --keep`: Keep original file mode. Copy the file to destination instead of moving it, and don't create a symbolic link at the original location.

### Arguments

- `<source>`: The file or directory to process (required)
- `<destination>`: Target location (optional if configured in config file)

## Configuration

Create a configuration file at `~/.config/fastlink/config.yaml`:

```yaml
# Default destination directory
# Files will be moved here when no destination is specified
destination: "/mnt/storage/fastlink"

# Keep original files by default
# true:  Copy files, keep originals (no symbolic links created)
# false: Move files, create symbolic links (default behavior)
keep: false
```

### Configuration Priority

1. Command line arguments (highest priority)
2. Configuration file settings
3. Default values (lowest priority)

## Examples

### Basic Usage

**Move file to configured destination:**
```bash
fastlink ~/Downloads/large-video.mp4
# Result: File moved to /mnt/storage/fastlink/large-video.mp4
#         Symbolic link created at ~/Downloads/large-video.mp4
```

**Move directory to specific location:**
```bash
fastlink ~/Documents/old-projects /backup/archives/
# Result: Directory moved to /backup/archives/old-projects
#         Symbolic link created at ~/Documents/old-projects
```

**Copy mode (preserve original):**
```bash
fastlink --keep ~/important-document.pdf ~/backup/
# Result: File copied to ~/backup/important-document.pdf
#         Original file remains untouched at ~/important-document.pdf
```

### Advanced Examples

**Archive large directories:**
```bash
# Move multiple large directories to external storage
fastlink ~/Videos/raw-footage /mnt/external/storage/
fastlink ~/Development/old-projects /mnt/external/storage/
fastlink ~/Downloads/iso-files /mnt/external/storage/
```

**Backup important files:**
```bash
# Copy important files while keeping originals
fastlink -k ~/.ssh /backup/ssh-keys/
fastlink -k ~/Documents/contracts /backup/documents/
```

**Free up space while maintaining access:**
```bash
# Move large files but keep them accessible via symbolic links
fastlink ~/Downloads/ubuntu-22.04.iso
fastlink ~/Videos/family-vacation-2023
fastlink ~/.cache/large-app-cache
```

### Error Cases

**Recursive move prevention:**
```bash
fastlink /home/user/documents /home/user/documents/backup
# Error: recursive move is not allowed
```

**Missing configuration:**
```bash
fastlink /some/file
# Error: destination not provided and not found in config
```

**Destination already exists:**
```bash
fastlink file.txt /backup/
# If /backup/file.txt already exists:
# Error: destination already exists: /backup/file.txt
```

## How It Works

1. **Safety Checks**: Validates source exists and prevents recursive operations
2. **Path Resolution**: Converts all paths to absolute paths for reliability
3. **Destination Preparation**: Creates destination directories if needed
4. **File Operations**: 
   - **Normal mode**: Copies file → Removes original → Creates symbolic link
   - **Keep mode**: Copies file only
5. **Symbolic Link Creation**: Creates link pointing to the new location

### Symbolic Links Explained

When fastlink moves a file, it creates a symbolic link at the original location that points to the new location. This means:

- ✅ Applications can still access the file using the original path
- ✅ File appears to be in the original location
- ✅ Actual file data is stored in the destination location
- ✅ You save disk space on the original location

Example:
```bash
# Before
/home/user/large-file.zip (1GB file)

# After: fastlink /home/user/large-file.zip
/home/user/large-file.zip -> /mnt/storage/fastlink/large-file.zip
# (symbolic link)              (actual 1GB file)
```

## Safety Features

- **Recursive operation prevention**: Cannot move a directory into itself
- **Existing file protection**: Won't overwrite existing files at destination
- **Symbolic link handling**: Skips processing of symbolic links to prevent loops
- **Path validation**: Ensures all paths are valid and accessible

## Troubleshooting

### Common Issues

**Permission denied errors:**
```bash
# Ensure you have write permissions to both source and destination
ls -la /path/to/source
ls -la /path/to/destination
```

**Symbolic link not working:**
```bash
# Check if the symbolic link exists and points to the right location
ls -la /original/path
readlink /original/path
```

**Configuration not found:**
```bash
# Verify config file exists and is readable
ls -la ~/.config/fastlink/config.yaml
cat ~/.config/fastlink/config.yaml
```

### Getting Help

Run `fastlink` without arguments to see usage information:
```bash
fastlink
```

## Development

### Building from Source

```bash
git clone https://github.com/xxnuo/fastlink.git
cd fastlink
go mod download
go build -o fastlink main.go
```

### Running Tests

```bash
go test -v ./...
```

### Available Make Targets

```bash
make build     # Build the binary
make test      # Run tests
make install   # Build and install to /usr/local/bin
make uninstall # Remove from /usr/local/bin
```

## Dependencies

- [github.com/spf13/fileflow](https://github.com/spf13/fileflow) - Fast file operations
- [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3) - YAML configuration parsing

## License

This project is open source. Please check the repository for license details.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.
