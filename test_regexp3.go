package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/juju/errors"
)

type CRCCheckRecord struct {
	ID        int64  `json:"id"`
	Timestamp string `json:"timestamp"`
	Msg       string `json:"msg"`
}

func (record *CRCCheckRecord) String() string {
	return fmt.Sprintf("CRCCheckRecord[id=%d, Timestamp = %s, Msg = %s]", record.ID, record.Timestamp, record.Msg)
}

var (
	cephLogFile      = "/Users/xsky/go/src/github.com/microyahoo/go-exercises/ceph.log"
	repairCharacters = "trying to repair"
	//2019-05-06 11:58:47.845831 osd.0 osd.0 10.252.3.100:6801/30011 772 : cluster [ERR] oid 2:515bf91a:::xbd_data.503e456f.0000000000000040:0:head crc mismatch, trying to repair
	pattern = regexp.MustCompile(`\A((?P<timestamp>\d{4,}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{1,})) osd.*\[ERR\].*xbd_data\.(?P<imageID>\w+)\.(?P<msg>.*)\z`)
)

func main() {
	test()
}

func test() error {
	file, err := os.Open(cephLogFile)
	if err != nil {
		log.Printf("Cannot open ceph log file: %s, err: [%v]", cephLogFile, err)
		return errors.Trace(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, repairCharacters) {
			var imageID string
			record := new(CRCCheckRecord)
			if pattern.MatchString(line) {
				matches := pattern.FindStringSubmatch(line)
				for i, name := range pattern.SubexpNames() {
					value := matches[i]
					if i == 0 || name == "" || value == "" {
						continue
					}
					switch name {
					case "timestamp":
						record.Timestamp = value
					case "imageID":
						imageID = value
					case "msg":
						record.Msg = value
					default:
						log.Printf("Failed to parse")
					}
				}
				// if err = job.m.CreateCRCCheck(record); err != nil {
				// 	return errors.Trace(err)
				// }
			} else {
				log.Printf("The crc check info %s is invalid", line)
			}
			if len(imageID) > 0 {
				log.Printf("Imageid = %s", imageID)
				// volume, err := job.m.GetVolumeByField("ImageID", imageID)
				// if _, err := job.m.CheckCreateAlertByEvent(models.ResourceBlockVolume, volume,
				// 	models.AlertEventVolumeCRCCheckRepair); err != nil {
				// 	return errors.Trace(err)
				// }
			}
			log.Printf("record = %v", record)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Cannot scanner text file: %s, err: [%v]", cephLogFile, err)
		return errors.Trace(err)
	}
	return nil
}
