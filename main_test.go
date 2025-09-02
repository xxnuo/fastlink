package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoadConfig(t *testing.T) {
	// Test loading config when file doesn't exist
	t.Run("config file doesn't exist", func(t *testing.T) {
		// Create a temporary home directory
		tmpDir := t.TempDir()

		// Override user home directory for testing
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", tmpDir)
		defer os.Setenv("HOME", originalHome)

		config, err := loadConfig()
		if err != nil {
			t.Fatalf("Expected no error when config doesn't exist, got: %v", err)
		}

		if config.Destination != "" || config.Keep != false || config.Move != false {
			t.Errorf("Expected empty config, got: %+v", config)
		}
	})

	t.Run("valid config file", func(t *testing.T) {
		// Create a temporary home directory
		tmpDir := t.TempDir()

		// Override user home directory for testing
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", tmpDir)
		defer os.Setenv("HOME", originalHome)

		// Create config directory
		configDir := filepath.Join(tmpDir, ".config", "fastlink")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatalf("Failed to create config directory: %v", err)
		}

		// Create config file
		configPath := filepath.Join(configDir, "config.yaml")
		configData := Config{
			Destination: "/tmp/fastlink",
			Keep:        true,
			Move:        false,
		}

		data, err := yaml.Marshal(configData)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		if err := ioutil.WriteFile(configPath, data, 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		config, err := loadConfig()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if config.Destination != "/tmp/fastlink" || config.Keep != true || config.Move != false {
			t.Errorf("Expected config {Destination: '/tmp/fastlink', Keep: true, Move: false}, got: %+v", config)
		}
	})

	t.Run("invalid config file", func(t *testing.T) {
		// Create a temporary home directory
		tmpDir := t.TempDir()

		// Override user home directory for testing
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", tmpDir)
		defer os.Setenv("HOME", originalHome)

		// Create config directory
		configDir := filepath.Join(tmpDir, ".config", "fastlink")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatalf("Failed to create config directory: %v", err)
		}

		// Create invalid config file
		configPath := filepath.Join(configDir, "config.yaml")
		invalidData := "invalid: yaml: content: ["

		if err := ioutil.WriteFile(configPath, []byte(invalidData), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		_, err := loadConfig()
		if err == nil {
			t.Error("Expected error for invalid config file, got none")
		}
	})
}

func TestIsSubPath(t *testing.T) {
	tests := []struct {
		name     string
		parent   string
		child    string
		expected bool
	}{
		{
			name:     "child is subdirectory",
			parent:   "/home/user",
			child:    "/home/user/documents",
			expected: true,
		},
		{
			name:     "child is same path",
			parent:   "/home/user",
			child:    "/home/user",
			expected: true,
		},
		{
			name:     "child is parent directory",
			parent:   "/home/user/documents",
			child:    "/home/user",
			expected: false,
		},
		{
			name:     "completely different paths",
			parent:   "/home/user1",
			child:    "/home/user2",
			expected: false,
		},
		{
			name:     "child with .. in path",
			parent:   "/home/user",
			child:    "/home/user/documents/../..",
			expected: false,
		},
		{
			name:     "nested subdirectory",
			parent:   "/home/user",
			child:    "/home/user/documents/projects/test",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSubPath(tt.parent, tt.child)
			if result != tt.expected {
				t.Errorf("isSubPath(%q, %q) = %v, want %v", tt.parent, tt.child, result, tt.expected)
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source file
	sourceFile := filepath.Join(tmpDir, "source.txt")
	sourceContent := "Hello, World!\nThis is a test file."

	if err := ioutil.WriteFile(sourceFile, []byte(sourceContent), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Test successful copy
	t.Run("successful copy", func(t *testing.T) {
		destFile := filepath.Join(tmpDir, "dest.txt")

		err := copyFile(sourceFile, destFile)
		if err != nil {
			t.Fatalf("copyFile failed: %v", err)
		}

		// Check if destination file exists
		if _, err := os.Stat(destFile); os.IsNotExist(err) {
			t.Error("Destination file was not created")
		}

		// Check file content
		destContent, err := ioutil.ReadFile(destFile)
		if err != nil {
			t.Fatalf("Failed to read destination file: %v", err)
		}

		if string(destContent) != sourceContent {
			t.Errorf("File content mismatch. Expected: %q, got: %q", sourceContent, string(destContent))
		}

		// Check file permissions
		sourceInfo, _ := os.Stat(sourceFile)
		destInfo, _ := os.Stat(destFile)

		if sourceInfo.Mode() != destInfo.Mode() {
			t.Errorf("File permissions mismatch. Expected: %v, got: %v", sourceInfo.Mode(), destInfo.Mode())
		}
	})

	// Test copy to non-existent directory
	t.Run("copy to non-existent directory", func(t *testing.T) {
		destFile := filepath.Join(tmpDir, "nonexistent", "dest.txt")

		err := copyFile(sourceFile, destFile)
		if err == nil {
			t.Error("Expected error when copying to non-existent directory, got none")
		}
	})

	// Test copy non-existent source
	t.Run("copy non-existent source", func(t *testing.T) {
		nonExistentSource := filepath.Join(tmpDir, "nonexistent.txt")
		destFile := filepath.Join(tmpDir, "dest2.txt")

		err := copyFile(nonExistentSource, destFile)
		if err == nil {
			t.Error("Expected error when copying non-existent source, got none")
		}
	})
}

func TestCopyDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source directory structure
	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(filepath.Join(sourceDir, "subdir"), 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create files in source directory
	file1 := filepath.Join(sourceDir, "file1.txt")
	file2 := filepath.Join(sourceDir, "subdir", "file2.txt")

	if err := ioutil.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}

	if err := ioutil.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	// Test successful copy
	t.Run("successful directory copy", func(t *testing.T) {
		destDir := filepath.Join(tmpDir, "dest")

		err := copyDir(sourceDir, destDir, false)
		if err != nil {
			t.Fatalf("copyDir failed: %v", err)
		}

		// Check if destination directory exists
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			t.Error("Destination directory was not created")
		}

		// Check if files were copied
		destFile1 := filepath.Join(destDir, "file1.txt")
		destFile2 := filepath.Join(destDir, "subdir", "file2.txt")

		if _, err := os.Stat(destFile1); os.IsNotExist(err) {
			t.Error("file1.txt was not copied")
		}

		if _, err := os.Stat(destFile2); os.IsNotExist(err) {
			t.Error("subdir/file2.txt was not copied")
		}

		// Check file contents
		content1, _ := ioutil.ReadFile(destFile1)
		content2, _ := ioutil.ReadFile(destFile2)

		if string(content1) != "content1" {
			t.Errorf("file1.txt content mismatch. Expected: 'content1', got: %q", string(content1))
		}

		if string(content2) != "content2" {
			t.Errorf("file2.txt content mismatch. Expected: 'content2', got: %q", string(content2))
		}
	})
}

func TestFastCopy(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("copy file", func(t *testing.T) {
		// Create source file
		sourceFile := filepath.Join(tmpDir, "source_file.txt")
		destFile := filepath.Join(tmpDir, "dest_file.txt")

		if err := ioutil.WriteFile(sourceFile, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		err := fastCopy(sourceFile, destFile, false)
		if err != nil {
			t.Fatalf("fastCopy failed: %v", err)
		}

		// Check if destination exists
		if _, err := os.Stat(destFile); os.IsNotExist(err) {
			t.Error("Destination file was not created")
		}

		// Check if source still exists (copy mode)
		if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
			t.Error("Source file was removed in copy mode")
		}
	})

	t.Run("move file", func(t *testing.T) {
		// Create source file
		sourceFile := filepath.Join(tmpDir, "source_move.txt")
		destFile := filepath.Join(tmpDir, "dest_move.txt")

		if err := ioutil.WriteFile(sourceFile, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		err := fastCopy(sourceFile, destFile, true)
		if err != nil {
			t.Fatalf("fastCopy failed: %v", err)
		}

		// Check if destination exists
		if _, err := os.Stat(destFile); os.IsNotExist(err) {
			t.Error("Destination file was not created")
		}

		// Check if source was removed (move mode)
		if _, err := os.Stat(sourceFile); !os.IsNotExist(err) {
			t.Error("Source file was not removed in move mode")
		}
	})

	t.Run("copy directory", func(t *testing.T) {
		// Create source directory
		sourceDir := filepath.Join(tmpDir, "source_dir")
		destDir := filepath.Join(tmpDir, "dest_dir")

		if err := os.MkdirAll(sourceDir, 0755); err != nil {
			t.Fatalf("Failed to create source directory: %v", err)
		}

		// Create a file in source directory
		testFile := filepath.Join(sourceDir, "test.txt")
		if err := ioutil.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err := fastCopy(sourceDir, destDir, false)
		if err != nil {
			t.Fatalf("fastCopy failed: %v", err)
		}

		// Check if destination directory exists
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			t.Error("Destination directory was not created")
		}

		// Check if file was copied
		destTestFile := filepath.Join(destDir, "test.txt")
		if _, err := os.Stat(destTestFile); os.IsNotExist(err) {
			t.Error("File in directory was not copied")
		}
	})

	t.Run("destination already exists", func(t *testing.T) {
		// Create source and destination files
		sourceFile := filepath.Join(tmpDir, "source_exists.txt")
		destFile := filepath.Join(tmpDir, "dest_exists.txt")

		if err := ioutil.WriteFile(sourceFile, []byte("source"), 0644); err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}

		if err := ioutil.WriteFile(destFile, []byte("dest"), 0644); err != nil {
			t.Fatalf("Failed to create dest file: %v", err)
		}

		err := fastCopy(sourceFile, destFile, false)
		if err == nil {
			t.Error("Expected error when destination already exists, got none")
		}

		if !strings.Contains(err.Error(), "destination already exists") {
			t.Errorf("Expected 'destination already exists' error, got: %v", err)
		}
	})
}

// Benchmark tests
func BenchmarkCopyFile(b *testing.B) {
	tmpDir := b.TempDir()

	// Create a source file with some content
	sourceFile := filepath.Join(tmpDir, "benchmark_source.txt")
	content := strings.Repeat("Hello, World!\n", 1000) // ~13KB

	if err := ioutil.WriteFile(sourceFile, []byte(content), 0644); err != nil {
		b.Fatalf("Failed to create source file: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		destFile := filepath.Join(tmpDir, "benchmark_dest_"+string(rune(i))+".txt")
		if err := copyFile(sourceFile, destFile); err != nil {
			b.Fatalf("copyFile failed: %v", err)
		}
		// Clean up
		os.Remove(destFile)
	}
}

func BenchmarkIsSubPath(b *testing.B) {
	parent := "/home/user/documents/projects"
	child := "/home/user/documents/projects/myproject/src/main.go"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		isSubPath(parent, child)
	}
}
