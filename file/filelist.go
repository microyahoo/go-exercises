package file

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// StreamFileList reads file list(s) from fileListPath (a single file or
// a directory of list files) and calls fn for each cleaned relative
// path. Lines that are empty or start with '#' are skipped.
func StreamFileList(fileListPath string, fn func(relPath string) error) error {
	info, err := os.Stat(fileListPath)
	if err != nil {
		return fmt.Errorf("stat filelist %s: %w", fileListPath, err)
	}

	if !info.IsDir() {
		return streamFileListFrom(fileListPath, fn)
	}

	entries, err := os.ReadDir(fileListPath)
	if err != nil {
		return fmt.Errorf("read filelist dir %s: %w", fileListPath, err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(fileListPath, entry.Name())
		log.Printf("processing file list: %s", path)
		if err := streamFileListFrom(path, fn); err != nil {
			return err
		}
	}
	return nil
}

func streamFileListFrom(path string, fn func(relPath string) error) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file list %s: %w", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if err := fn(filepath.Clean(line)); err != nil {
			return err
		}
	}
	return scanner.Err()
}
