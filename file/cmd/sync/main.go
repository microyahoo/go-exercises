package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	file "github.com/microyahoo/go-exercises/file"
)

func main() {
	dstDir := flag.String("dst", "", "destination directory (required)")
	fileListPath := flag.String("filelist", "", "path to source file list (required)")
	concurrency := flag.Int("c", 32, "number of concurrent workers")
	mode := flag.String("mode", "approx", "sync mode: exact (partitioned hash set, zero false-positive) or approx (bloom filter, faster but ~fp-rate miss)")
	dryRun := flag.Bool("dry-run", true, "preview deletions without removing files")

	// exact-mode options
	partitions := flag.Int("partitions", 256, "number of hash partitions (each bucket ≈ total_files/N in memory)")
	tempDir := flag.String("temp-dir", "", "directory for partition temp files; auto-created if empty")

	// approx-mode options
	expectedFiles := flag.Uint64("expected-files", 700000000, "expected file count for bloom filter sizing")
	fpRate := flag.Float64("fp-rate", 0.01, "bloom filter false positive rate")

	flag.Parse()

	if *dstDir == "" || *fileListPath == "" {
		fmt.Fprintln(os.Stderr, "error: -dst and -filelist are required")
		flag.Usage()
		os.Exit(1)
	}
	if *mode != "exact" && *mode != "approx" {
		fmt.Fprintln(os.Stderr, "error: -mode must be 'exact' or 'approx'")
		os.Exit(1)
	}

	absDst, err := filepath.Abs(*dstDir)
	if err != nil {
		log.Fatalf("abs dst: %v", err)
	}
	absFileList, err := filepath.Abs(*fileListPath)
	if err != nil {
		log.Fatalf("abs filelist: %v", err)
	}

	var dstats file.DeleteStats
	totalStart := time.Now()

	if *mode == "approx" {
		syncApprox(absDst, absFileList, *concurrency, *dryRun, *expectedFiles, *fpRate, &dstats)
	} else {
		syncExact(absDst, absFileList, *concurrency, *dryRun, *partitions, *tempDir, &dstats)
	}

	totalElapsed := time.Since(totalStart)
	log.Printf("sync completed in %s: scanned=%d deleted=%d failed=%d dirs_removed=%d",
		totalElapsed.Round(time.Millisecond),
		dstats.Scanned.Load(), dstats.Deleted.Load(),
		dstats.Failed.Load(), dstats.DirsRemoved.Load())
}

// syncExact uses partitioned hash sets for zero-false-positive deletion.
func syncExact(absDst, fileListPath string, concurrency int, dryRun bool, partitions int, tempDir string, stats *file.DeleteStats) {
	// Prepare temp directory.
	td := tempDir
	var autoTempDir string
	if td == "" {
		var err error
		td, err = os.MkdirTemp("", "filesync-*")
		if err != nil {
			log.Fatalf("create temp dir: %v", err)
		}
		autoTempDir = td
	}
	if autoTempDir != "" {
		defer func() {
			log.Printf("removing temp directory: %s", autoTempDir)
			os.RemoveAll(autoTempDir)
		}()
	}

	// Phase 1: build source partitions from filelist.
	log.Printf("phase 1 – building source index: %d partitions, temp_dir=%s", partitions, td)
	srcPW, err := file.NewPartitionWriter(td, "src", partitions)
	if err != nil {
		log.Fatalf("create source partition writer: %v", err)
	}

	var srcCount int64
	if err := file.StreamFileList(fileListPath, func(rel string) error {
		srcCount++
		return srcPW.Add(rel)
	}); err != nil {
		log.Fatalf("read file list: %v", err)
	}
	if err := srcPW.Close(); err != nil {
		log.Fatalf("close source partitions: %v", err)
	}
	log.Printf("phase 1 done: indexed %d source paths", srcCount)

	// Phase 2: walk destination into partition files.
	log.Println("phase 2 – scanning destination into partition files...")
	dstPW, err := file.NewPartitionWriter(td, "dst", partitions)
	if err != nil {
		log.Fatalf("create dest partition writer: %v", err)
	}

	var walkScanned atomic.Int64
	dirs, walkErr := file.WalkDestinationPaths(absDst, dstPW, &walkScanned)
	if walkErr != nil {
		log.Printf("phase 2 walk error: %v", walkErr)
	}
	if err := dstPW.Close(); err != nil {
		log.Fatalf("close dest partitions: %v", err)
	}
	log.Printf("phase 2 done: scanned %d destination files", walkScanned.Load())

	// Phase 3: diff partitions and delete extras.
	if dryRun {
		log.Println("phase 3 – diffing partitions (DRY-RUN)...")
	} else {
		log.Println("phase 3 – diffing partitions and deleting extras...")
	}
	deleteStart := time.Now()
	stopProgress := file.StartDeleteProgressReporter(stats, deleteStart, 30*time.Second)

	if diffErr := file.DeleteExtrasExact(absDst, srcPW, dstPW, partitions, concurrency, dryRun, stats); diffErr != nil {
		log.Printf("phase 3 error: %v", diffErr)
	}
	stopProgress()

	deleteElapsed := time.Since(deleteStart)
	log.Printf("phase 3 done in %s: deleted=%d failed=%d",
		deleteElapsed.Round(time.Millisecond), stats.Deleted.Load(), stats.Failed.Load())

	// Phase 4: clean up empty directories.
	if !dryRun && stats.Deleted.Load() > 0 && len(dirs) > 0 {
		log.Println("cleaning up empty directories...")
		removed := file.CleanEmptyDirs(dirs)
		stats.DirsRemoved.Store(removed)
		log.Printf("removed %d empty directories", removed)
	}
}

// syncApprox uses a bloom filter for fast approximate deletion.
func syncApprox(absDst, fileListPath string, concurrency int, dryRun bool, expectedFiles uint64, fpRate float64, stats *file.DeleteStats) {
	// Phase 1: build bloom filter from filelist.
	bf := file.NewBloomFilter(expectedFiles, fpRate)
	log.Printf("phase 1 – building bloom filter: %.2f MB, %d hash functions, expected=%d, fp_rate=%.4f",
		bf.MemoryUsageMB(), bf.NumHashFunctions(), expectedFiles, fpRate)

	var srcCount int64
	if err := file.StreamFileList(fileListPath, func(rel string) error {
		srcCount++
		bf.Add([]byte(rel))
		return nil
	}); err != nil {
		log.Fatalf("read file list: %v", err)
	}
	log.Printf("phase 1 done: indexed %d source paths", srcCount)

	// Phase 2: walk destination and delete extras.
	if dryRun {
		log.Println("phase 2 – scanning destination (DRY-RUN)...")
	} else {
		log.Println("phase 2 – scanning destination and deleting extras...")
	}
	deleteStart := time.Now()
	stopProgress := file.StartDeleteProgressReporter(stats, deleteStart, 30*time.Second)

	dirs, walkErr := file.DeleteExtrasApprox(absDst, bf, concurrency, dryRun, stats)
	stopProgress()

	if walkErr != nil {
		log.Printf("phase 2 walk error: %v", walkErr)
	}
	deleteElapsed := time.Since(deleteStart)
	log.Printf("phase 2 done in %s: scanned=%d deleted=%d failed=%d",
		deleteElapsed.Round(time.Millisecond),
		stats.Scanned.Load(), stats.Deleted.Load(), stats.Failed.Load())

	// Phase 3: clean up empty directories.
	if !dryRun && stats.Deleted.Load() > 0 && len(dirs) > 0 {
		log.Println("cleaning up empty directories...")
		removed := file.CleanEmptyDirs(dirs)
		stats.DirsRemoved.Store(removed)
		log.Printf("removed %d empty directories", removed)
	}
}
