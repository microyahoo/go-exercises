package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type FileStats struct {
	Total     int64
	Dirs      int64
	Files     int64
	NotExists int64
	Others    int64
	Errors    int64
}

type FileResult struct {
	Path     string
	Type     string // "dir", "file", "not_exists", "other", "error"
	ErrorMsg string
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("用法: %s <文件列表文件>\n", filepath.Base(os.Args[0]))
		fmt.Println("示例: ./file-checker filelist.txt baseDir")
		os.Exit(1)
	}

	listFile := os.Args[1]
	// baseDir := "/mnt/csi-data-aly/shared/public/jiangweichen/"
	baseDir := os.Args[2]

	// 确保基础目录以 / 结尾
	if !strings.HasSuffix(baseDir, "/") {
		baseDir += "/"
	}

	// 检查列表文件是否存在
	if _, err := os.Stat(listFile); os.IsNotExist(err) {
		fmt.Printf("错误: 文件 '%s' 不存在\n", listFile)
		os.Exit(1)
	}

	startTime := time.Now()

	// 读取文件列表
	paths, err := readFileList(listFile)
	if err != nil {
		fmt.Printf("读取文件列表失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("开始检查文件列表: %s\n", listFile)
	fmt.Printf("基础目录: %s\n", baseDir)
	fmt.Printf("待检查项目: %d\n", len(paths))
	fmt.Println("========================================")

	// 并发处理
	results := processConcurrently(paths, baseDir, 10) // 10个并发goroutine

	// 输出结果并统计
	stats := FileStats{}
	for _, result := range results {
		atomic.AddInt64(&stats.Total, 1)
		switch result.Type {
		case "dir":
			atomic.AddInt64(&stats.Dirs, 1)
			fmt.Printf("[目录] %s\n", result.Path)
		case "file":
			atomic.AddInt64(&stats.Files, 1)
			// fmt.Printf("[文件] %s\n", result.Path)
		case "not_exists":
			atomic.AddInt64(&stats.NotExists, 1)
			// fmt.Printf("[不存在] %s\n", result.Path)
		case "other":
			atomic.AddInt64(&stats.Others, 1)
			fmt.Printf("[其他] %s\n", result.Path)
		case "error":
			atomic.AddInt64(&stats.Errors, 1)
			fmt.Printf("[错误] %s - %s\n", result.Path, result.ErrorMsg)
		}
	}

	// 计算耗时
	elapsed := time.Since(startTime)

	// 输出统计信息
	fmt.Println("========================================")
	fmt.Printf("统计结果 (耗时: %v):\n", elapsed.Round(time.Millisecond))
	fmt.Printf("总共检查: %d\n", stats.Total)
	fmt.Printf("存在: %d\n", stats.Dirs+stats.Files+stats.Others)
	if stats.Dirs > 0 {
		fmt.Printf("  - 目录: %d\n", stats.Dirs)
	}
	if stats.Files > 0 {
		fmt.Printf("  - 文件: %d\n", stats.Files)
	}
	if stats.Others > 0 {
		fmt.Printf("  - 其他: %d (特殊文件类型)\n", stats.Others)
	}
	fmt.Printf("不存在: %d\n", stats.NotExists)
	if stats.Errors > 0 {
		fmt.Printf("错误: %d\n", stats.Errors)
	}
}

// 读取文件列表，跳过空行和注释
func readFileList(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var paths []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 跳过空行和注释行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		paths = append(paths, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return paths, nil
}

// 检查单个路径
func checkPath(path, baseDir string) FileResult {
	fullPath := filepath.Join(baseDir, path)

	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return FileResult{
				Path: path,
				Type: "not_exists",
			}
		}
		return FileResult{
			Path:     path,
			Type:     "error",
			ErrorMsg: err.Error(),
		}
	}

	mode := fileInfo.Mode()
	if mode.IsDir() {
		return FileResult{
			Path: path,
			Type: "dir",
		}
	} else if mode.IsRegular() {
		return FileResult{
			Path: path,
			Type: "file",
		}
	} else {
		return FileResult{
			Path: path,
			Type: "other",
		}
	}
}

// 并发处理路径检查
func processConcurrently(paths []string, baseDir string, workers int) []FileResult {
	var wg sync.WaitGroup
	results := make([]FileResult, 0, len(paths))
	resultsCh := make(chan FileResult, len(paths))
	pathsCh := make(chan string, len(paths))

	// 启动 worker goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathsCh {
				resultsCh <- checkPath(path, baseDir)
			}
		}()
	}

	// 发送路径到 channel
	for _, path := range paths {
		pathsCh <- path
	}
	close(pathsCh)

	// 等待所有 worker 完成
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// 收集结果
	for result := range resultsCh {
		results = append(results, result)
	}

	return results
}
