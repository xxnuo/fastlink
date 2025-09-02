package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Destination string `yaml:"destination"`
	Keep        bool   `yaml:"keep"`
	Move        bool   `yaml:"move"`
}

func main() {
	var keep bool
	var move bool
	flag.BoolVar(&keep, "k", false, "keep original file mode flag")
	flag.BoolVar(&keep, "keep", false, "keep original file mode flag")
	flag.BoolVar(&move, "m", false, "quick move mode flag")
	flag.BoolVar(&move, "move", false, "quick move mode flag")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-m] [--move] [-k] [--keep] <source> [<destination>]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nQuickly move files/directories to a certain location and create soft links at the original location\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	source := args[0]
	var destination string

	// Get destination from args or config
	if len(args) >= 2 {
		destination = args[1]
	} else {
		config, err := loadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		destination = config.Destination
		// Use config keep flag if not overridden by command line
		if !keep {
			keep = config.Keep
		}
		// Use config move flag if not overridden by command line
		if !move {
			move = config.Move
		}
	}

	if destination == "" {
		fmt.Fprintf(os.Stderr, "Error: destination not provided and not found in config\n")
		os.Exit(1)
	}

	// Convert to absolute paths
	sourceAbs, err := filepath.Abs(source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting absolute path for source: %v\n", err)
		os.Exit(1)
	}

	destinationAbs, err := filepath.Abs(destination)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting absolute path for destination: %v\n", err)
		os.Exit(1)
	}

	// Check for recursive move
	if isSubPath(sourceAbs, destinationAbs) {
		fmt.Fprintf(os.Stderr, "Error: recursive move is not allowed\n")
		os.Exit(1)
	}

	// Check if source exists
	sourceInfo, err := os.Stat(sourceAbs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: source does not exist: %v\n", err)
		os.Exit(1)
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destinationAbs, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating destination directory: %v\n", err)
		os.Exit(1)
	}

	// Determine final destination path
	var finalDest string
	if sourceInfo.IsDir() {
		finalDest = filepath.Join(destinationAbs, filepath.Base(sourceAbs))
	} else {
		finalDest = filepath.Join(destinationAbs, filepath.Base(sourceAbs))
	}

	// Copy/move the file or directory
	if err := fastCopy(sourceAbs, finalDest, move); err != nil {
		fmt.Fprintf(os.Stderr, "Error copying file/directory: %v\n", err)
		os.Exit(1)
	}

	if !keep {
		// Create soft link at original location
		if err := os.RemoveAll(sourceAbs); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing original file: %v\n", err)
			os.Exit(1)
		}

		if err := os.Symlink(finalDest, sourceAbs); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating symlink: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Moved %s to %s and created symlink\n", sourceAbs, finalDest)
	} else {
		fmt.Printf("Copied %s to %s (keeping original)\n", sourceAbs, finalDest)
	}
}

func loadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to get user home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".config", "fastlink", "config.yaml")

	// If config file doesn't exist, return empty config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unable to parse config file: %w", err)
	}

	return &config, nil
}

func isSubPath(parent, child string) bool {
	// Clean paths to handle .. and . properly
	parent = filepath.Clean(parent)
	child = filepath.Clean(child)

	// Check if child path starts with parent path
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}

	// If relative path doesn't start with "..", then child is under parent
	return !strings.HasPrefix(rel, "..")
}

func fastCopy(source, dest string, move bool) error {
	// Get source info
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}

	// Check if destination already exists
	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("destination already exists: %s", dest)
	}

	// Handle directory
	if sourceInfo.IsDir() {
		if err := copyDir(source, dest, move); err != nil {
			return fmt.Errorf("failed to copy directory: %w", err)
		}
	} else {
		// Handle file
		if err := copyFile(source, dest); err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}
	}

	// Remove source if this is a move operation
	if move {
		if err := os.RemoveAll(source); err != nil {
			return fmt.Errorf("failed to remove source after move: %w", err)
		}
	}

	return nil
}

// copyFile copies a single file from source to destination
func copyFile(source, dest string) error {
	// Open source file
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get source file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create destination file
	dstFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Set destination file permissions to match source
	if err := dstFile.Chmod(srcInfo.Mode()); err != nil {
		return err
	}

	// Copy file content with buffer for efficiency
	buffer := make([]byte, 64*1024) // 64KB buffer
	if _, err := io.CopyBuffer(dstFile, srcFile, buffer); err != nil {
		return err
	}

	// Sync to ensure data is written to disk
	if err := dstFile.Sync(); err != nil {
		return err
	}

	return nil
}

// copyDir recursively copies a directory from source to destination
func copyDir(source, dest string, move bool) error {
	// Get source directory info
	srcInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// Create destination directory with same permissions
	if err := os.MkdirAll(dest, srcInfo.Mode()); err != nil {
		return err
	}

	// Read source directory entries
	entries, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(source, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := copyDir(srcPath, destPath, false); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}
