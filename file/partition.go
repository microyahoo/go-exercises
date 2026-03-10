package file

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
)

// PartitionWriter distributes relative paths into N on-disk bucket files
// using consistent hashing so that each bucket can later be loaded into
// memory independently for exact set operations.
//
// For 600M paths with N=256, each bucket holds ~2.3M entries (~230 MB),
// well within a typical server's available memory.
type PartitionWriter struct {
	dir     string
	prefix  string
	n       int
	files   []*os.File
	writers []*bufio.Writer
}

// NewPartitionWriter creates N bucket files in dir named <prefix>_0000 …
// <prefix>_NNNN. Caller must call Close when done writing.
func NewPartitionWriter(dir, prefix string, n int) (*PartitionWriter, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create partition dir %s: %w", dir, err)
	}
	pw := &PartitionWriter{
		dir:     dir,
		prefix:  prefix,
		n:       n,
		files:   make([]*os.File, n),
		writers: make([]*bufio.Writer, n),
	}
	for i := 0; i < n; i++ {
		p := filepath.Join(dir, fmt.Sprintf("%s_%04d", prefix, i))
		f, err := os.Create(p)
		if err != nil {
			pw.Close()
			return nil, fmt.Errorf("create partition file %s: %w", p, err)
		}
		pw.files[i] = f
		pw.writers[i] = bufio.NewWriterSize(f, 64*1024)
	}
	return pw, nil
}

// Add writes relPath to the appropriate bucket file.
func (pw *PartitionWriter) Add(relPath string) error {
	idx := pathPartition(relPath, pw.n)
	if _, err := pw.writers[idx].WriteString(relPath); err != nil {
		return err
	}
	_, err := pw.writers[idx].WriteString("\n")
	return err
}

// PartitionPath returns the filesystem path for bucket idx.
func (pw *PartitionWriter) PartitionPath(idx int) string {
	return filepath.Join(pw.dir, fmt.Sprintf("%s_%04d", pw.prefix, idx))
}

// Close flushes all writers and closes all files.
func (pw *PartitionWriter) Close() error {
	var firstErr error
	for i, w := range pw.writers {
		if w != nil {
			if err := w.Flush(); err != nil && firstErr == nil {
				firstErr = fmt.Errorf("flush partition %d: %w", i, err)
			}
		}
	}
	for i, f := range pw.files {
		if f != nil {
			if err := f.Close(); err != nil && firstErr == nil {
				firstErr = fmt.Errorf("close partition %d: %w", i, err)
			}
			pw.files[i] = nil
		}
	}
	return firstErr
}

// pathPartition returns the bucket index for the given path.
func pathPartition(path string, n int) int {
	h := fnv.New32a()
	h.Write([]byte(path))
	return int(h.Sum32() % uint32(n))
}

// LoadPartitionSet reads all lines from a single bucket file into a
// hash set for O(1) lookups.
func LoadPartitionSet(path string) (map[string]struct{}, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("open partition %s: %w", path, err)
	}
	defer f.Close()

	set := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			set[line] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan partition %s: %w", path, err)
	}
	return set, nil
}
