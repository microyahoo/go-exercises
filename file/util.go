package file

import (
	"fmt"
	"log"
	"time"
)

// FormatBytes returns a human-readable byte size string.
func FormatBytes(b int64) string {
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

// Throughput computes bytes per second.
func Throughput(bytes int64, d time.Duration) int64 {
	if d <= 0 {
		return 0
	}
	return int64(float64(bytes) / d.Seconds())
}

// StartMigrateProgressReporter periodically logs migration progress.
func StartMigrateProgressReporter(stats *MigrateStats, start time.Time, interval time.Duration) func() {
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
					FormatBytes(bytes), FormatBytes(Throughput(bytes, elapsed)))
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
