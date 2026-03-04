package core

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// GenerateSkillHash creates a SHA256 hash representing the content of the entire skill directory
func GenerateSkillHash(skillDir string) (string, error) {
	hash := sha256.New()
	var files []string

	err := filepath.WalkDir(skillDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			rel, _ := filepath.Rel(skillDir, path)
			files = append(files, rel)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	sort.Strings(files)

	for _, rel := range files {
		path := filepath.Join(skillDir, rel)
		f, err := os.Open(path)
		if err != nil {
			return "", err
		}

		// Write filename to hash
		hash.Write([]byte(rel))

		// Write file content to hash
		if _, err := io.Copy(hash, f); err != nil {
			f.Close()
			return "", err
		}
		f.Close()
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// CopyDir recursively copies a directory from src to dst
func CopyDir(src string, dst string) error {
	// First, remove existing directory to ensure a clean copy
	os.RemoveAll(dst)

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// It's a file
		return copyFile(path, targetPath)
	})
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// CreateSymlink removes any existing target and creates a symlink
func CreateSymlink(src, dst string) error {
	os.RemoveAll(dst)

	// Create parent directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	return os.Symlink(src, dst)
}
