package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	srcDir := flag.String("src", "", "source directory (required)")
	dstDir := flag.String("dst", "", "destination directory (required)")
	fileListPath := flag.String("filelist", "", "path to file list (file: one relative path per line; directory: each file in it is a file list) (required)")
	concurrency := flag.Int("c", 32, "number of concurrent workers")
	compare := flag.String("compare", "size", "file comparison method: none, size, md5")
	flag.Parse()

	if *srcDir == "" || *dstDir == "" || *fileListPath == "" {
		fmt.Fprintln(os.Stderr, "error: -src, -dst and -filelist are all required")
		flag.Usage()
		os.Exit(1)
	}

	method, err := ParseCompareMethod(*compare)
	if err != nil {
		log.Fatalf("invalid compare method: %v", err)
	}

	absSrc, err := filepath.Abs(*srcDir)
	if err != nil {
		log.Fatalf("abs src: %v", err)
	}
	absDst, err := filepath.Abs(*dstDir)
	if err != nil {
		log.Fatalf("abs dst: %v", err)
	}

	cfg := MigrateConfig{
		Concurrency:   *concurrency,
		CompareMethod: method,
	}
	migrator := NewMigrator(cfg)

	taskCh := make(chan FileTask, cfg.Concurrency*2)

	// Stream-read file list in a separate goroutine so memory usage
	// stays bounded regardless of how many files are in the list.
	go func() {
		defer close(taskCh)
		if err := streamTasks(absSrc, absDst, *fileListPath, taskCh); err != nil {
			log.Fatalf("read file list: %v", err)
		}
	}()

	log.Printf("starting migration: concurrency=%d, compare=%s", cfg.Concurrency, method)

	var stats MigrateStats
	start := time.Now()

	stopProgress := startProgressReporter(&stats, start, 30*time.Second)
	migrator.Run(taskCh, &stats)
	stopProgress()

	elapsed := time.Since(start)

	succeeded := stats.Succeeded.Load()
	skippedSame := stats.SkippedSame.Load()
	skippedNoSrc := stats.SkippedNoSrc.Load()
	failed := stats.Failed.Load()
	totalBytes := stats.Bytes.Load()
	total := succeeded + skippedSame + skippedNoSrc + failed

	log.Printf("migration completed in %s: total=%d succeeded=%d skipped_same=%d skipped_not_found=%d failed=%d",
		elapsed.Round(time.Millisecond), total,
		succeeded, skippedSame, skippedNoSrc, failed)
	log.Printf("total migrated: %s, throughput: %s/s",
		formatBytes(totalBytes), formatBytes(throughput(totalBytes, elapsed)))
}

func streamTasks(srcDir, dstDir, fileListPath string, ch chan<- FileTask) error {
	info, err := os.Stat(fileListPath)
	if err != nil {
		return fmt.Errorf("stat filelist %s: %w", fileListPath, err)
	}

	if !info.IsDir() {
		return streamTasksFromFile(srcDir, dstDir, fileListPath, ch)
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
		if err := streamTasksFromFile(srcDir, dstDir, path, ch); err != nil {
			return err
		}
	}
	return nil
}

func streamTasksFromFile(srcDir, dstDir, path string, ch chan<- FileTask) error {
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
		ch <- FileTask{
			SrcPath: filepath.Join(srcDir, line),
			DstPath: filepath.Join(dstDir, line),
		}
	}
	return scanner.Err()
}

func startProgressReporter(stats *MigrateStats, start time.Time, interval time.Duration) func() {
	ticker := time.NewTicker(interval)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				elapsed := time.Since(start)
				succeeded := stats.Succeeded.Load()
				skippedSame := stats.SkippedSame.Load()
				skippedNoSrc := stats.SkippedNoSrc.Load()
				failed := stats.Failed.Load()
				bytes := stats.Bytes.Load()
				log.Printf("[progress] elapsed=%s succeeded=%d skipped_same=%d skipped_not_found=%d failed=%d migrated=%s throughput=%s/s",
					elapsed.Round(time.Second), succeeded, skippedSame, skippedNoSrc, failed,
					formatBytes(bytes), formatBytes(throughput(bytes, elapsed)))
			case <-done:
				return
			}
		}
	}()
	return func() {
		ticker.Stop()
		close(done)
	}
}

func formatBytes(b int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)
	switch {
	case b >= TB:
		return fmt.Sprintf("%.2f TB", float64(b)/float64(TB))
	case b >= GB:
		return fmt.Sprintf("%.2f GB", float64(b)/float64(GB))
	case b >= MB:
		return fmt.Sprintf("%.2f MB", float64(b)/float64(MB))
	case b >= KB:
		return fmt.Sprintf("%.2f KB", float64(b)/float64(KB))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func throughput(bytes int64, d time.Duration) int64 {
	if d <= 0 {
		return 0
	}
	return int64(float64(bytes) / d.Seconds())
}
