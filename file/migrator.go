package file

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

type CompareMethod int

const (
	CompareNone CompareMethod = iota
	CompareSize
	CompareMD5
)

func (m CompareMethod) String() string {
	switch m {
	case CompareNone:
		return "none"
	case CompareSize:
		return "size"
	case CompareMD5:
		return "md5"
	default:
		return "unknown"
	}
}

func ParseCompareMethod(s string) (CompareMethod, error) {
	switch s {
	case "none":
		return CompareNone, nil
	case "size":
		return CompareSize, nil
	case "md5":
		return CompareMD5, nil
	default:
		return CompareNone, fmt.Errorf("unknown compare method: %s", s)
	}
}

// MigrateConfig holds all configurable options for file migration.
type MigrateConfig struct {
	Concurrency   int
	CompareMethod CompareMethod
	// BufSize is the buffer size for file copy, default 4MB.
	BufSize int
}

func DefaultMigrateConfig() MigrateConfig {
	return MigrateConfig{
		Concurrency:   32,
		CompareMethod: CompareSize,
		BufSize:       4 * 1024 * 1024,
	}
}

// FileTask represents a single file migration task.
type FileTask struct {
	SrcPath string
	DstPath string
}

type SkipReason int

const (
	NotSkipped    SkipReason = iota
	SkipSame                        // destination file matches source (by size or md5)
	SkipNotFound                    // source file does not exist
)

// MigrateStats aggregates migration statistics using atomic counters,
// safe for concurrent updates from multiple worker goroutines.
type MigrateStats struct {
	Succeeded      atomic.Int64
	SkippedSame    atomic.Int64
	SkippedNoSrc   atomic.Int64
	Failed         atomic.Int64
	Bytes          atomic.Int64
}

// Migrator performs concurrent file migration.
type Migrator struct {
	cfg MigrateConfig
}

func NewMigrator(cfg MigrateConfig) *Migrator {
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 1
	}
	if cfg.BufSize <= 0 {
		cfg.BufSize = 4 * 1024 * 1024
	}
	return &Migrator{cfg: cfg}
}

// Run reads tasks from taskCh, migrates them concurrently,
// and aggregates results into stats. It blocks until taskCh is closed
// and all workers finish.
func (m *Migrator) Run(taskCh <-chan FileTask, stats *MigrateStats) {
	var wg sync.WaitGroup
	for i := 0; i < m.cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := make([]byte, m.cfg.BufSize)
			for task := range taskCh {
				reason, n, err := m.migrate(task, buf)
				switch {
				case err != nil:
					stats.Failed.Add(1)
					fmt.Fprintf(os.Stderr, "[FAIL] %s -> %s: %v\n",
						task.SrcPath, task.DstPath, err)
				case reason == SkipSame:
					stats.SkippedSame.Add(1)
				case reason == SkipNotFound:
					stats.SkippedNoSrc.Add(1)
				default:
					stats.Succeeded.Add(1)
					stats.Bytes.Add(n)
				}
			}
		}()
	}
	wg.Wait()
}

func (m *Migrator) migrate(task FileTask, buf []byte) (reason SkipReason, n int64, err error) {
	srcInfo, err := os.Stat(task.SrcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return SkipNotFound, 0, nil
		}
		return NotSkipped, 0, fmt.Errorf("stat source %s: %w", task.SrcPath, err)
	}
	if !srcInfo.Mode().IsRegular() {
		return NotSkipped, 0, fmt.Errorf("%s is not a regular file", task.SrcPath)
	}

	dstDir := filepath.Dir(task.DstPath)
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return NotSkipped, 0, fmt.Errorf("mkdir %s: %w", dstDir, err)
	}

	if m.cfg.CompareMethod != CompareNone {
		same, err := m.filesEqual(task.SrcPath, task.DstPath, srcInfo)
		if err != nil {
			return NotSkipped, 0, err
		}
		if same {
			return SkipSame, 0, nil
		}
	}

	n, err = m.copyFile(task.SrcPath, task.DstPath, srcInfo.Mode(), buf)
	return NotSkipped, n, err
}

func (m *Migrator) filesEqual(src, dst string, srcInfo os.FileInfo) (bool, error) {
	dstInfo, err := os.Stat(dst)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("stat dest %s: %w", dst, err)
	}

	switch m.cfg.CompareMethod {
	case CompareSize:
		return srcInfo.Size() == dstInfo.Size(), nil
	case CompareMD5:
		srcHash, err := fileMD5(src)
		if err != nil {
			return false, err
		}
		dstHash, err := fileMD5(dst)
		if err != nil {
			return false, err
		}
		return srcHash == dstHash, nil
	default:
		return false, nil
	}
}

func fileMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open %s for md5: %w", path, err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("read %s for md5: %w", path, err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (m *Migrator) copyFile(src, dst string, perm os.FileMode, buf []byte) (written int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, fmt.Errorf("open source %s: %w", src, err)
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return 0, fmt.Errorf("create dest %s: %w", dst, err)
	}
	defer func() {
		if cerr := dstFile.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close dest %s: %w", dst, cerr)
		}
	}()

	written, err = io.CopyBuffer(dstFile, srcFile, buf)
	if err != nil {
		return written, fmt.Errorf("copy %s -> %s: %w", src, dst, err)
	}
	return written, nil
}
