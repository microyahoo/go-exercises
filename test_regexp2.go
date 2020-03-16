package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/juju/errors"
)

var blankRegexp = regexp.MustCompile("\\s+")

// DiskDeviceLogEntry ...
type DiskDeviceLogEntry struct {
	Type  string
	Name  string
	Value string
}

func (entry *DiskDeviceLogEntry) String() string {
	return fmt.Sprintf("DiskDeviceLogEntry[Type=%s, Name= %s, Value= %s]", entry.Type, entry.Name, entry.Value)
}

var (
	diskSmartFile = "/Users/xsky/go/src/github.com/microyahoo/go-exercises/disk_smart.txt"
)

func main() {
	file, err := os.Open(diskSmartFile)
	if err != nil {
		log.Printf("Cannot open disk smart file: %s, err: [%v]", diskSmartFile, err)
		return
	}
	defer file.Close()
	logs, err := parseDevstatDeviceLogsOfSata(file)
	log.Printf("logs = %v, err = %v", logs, err)

	s := regexp.MustCompile("a*").Split("abaabaccadaaae", 5)
	log.Printf("s = %s", strings.Join(s, "/"))
	// s: ["", "b", "b", "c", "cadaaae"]
}

func parseDevstatDeviceLogsOfSata(reader io.Reader) (
	logs []*DiskDeviceLogEntry, err error) {

	scanner := bufio.NewScanner(reader)
	startParse := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		log.Printf("**** line = %s", line)
		if len(line) == 0 || "" == line {
			continue
		}
		if strings.HasPrefix(line, "Page") {
			startParse = true
			continue
		}
		if !startParse {
			continue
		}

		parts := blankRegexp.Split(line, 5)
		log.Printf("**** parts = %s", strings.Join(parts, "/"))
		if len(parts) != 5 {
			return nil, errors.Errorf("Invalid line: %s", line)
		}
		if parts[3] == "=" {
			// This is a page separate line.
			continue
		}
		entry := &DiskDeviceLogEntry{
			Type:  "dev_stats",
			Name:  strings.TrimSpace(parts[4]),
			Value: strings.TrimSpace(parts[3]),
		}
		logs = append(logs, entry)
	}

	return logs, nil
}
