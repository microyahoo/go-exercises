package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	file "github.com/microyahoo/go-exercises/file"
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

	method, err := file.ParseCompareMethod(*compare)
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
	absFileList, err := filepath.Abs(*fileListPath)
	if err != nil {
		log.Fatalf("abs filelist: %v", err)
	}

	cfg := file.MigrateConfig{
		Concurrency:   *concurrency,
		CompareMethod: method,
	}
	migrator := file.NewMigrator(cfg)

	taskCh := make(chan file.FileTask, cfg.Concurrency*2)
	go func() {
		defer close(taskCh)
		if err := file.StreamFileList(absFileList, func(rel string) error {
			taskCh <- file.FileTask{
				SrcPath: filepath.Join(absSrc, rel),
				DstPath: filepath.Join(absDst, rel),
			}
			return nil
		}); err != nil {
			log.Fatalf("read file list: %v", err)
		}
	}()

	log.Printf("starting migration: concurrency=%d, compare=%s", cfg.Concurrency, method)

	var stats file.MigrateStats
	start := time.Now()

	stopProgress := file.StartMigrateProgressReporter(&stats, start, 30*time.Second)
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
		file.FormatBytes(totalBytes), file.FormatBytes(file.Throughput(totalBytes, elapsed)))
}
