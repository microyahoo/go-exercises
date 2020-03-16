package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// const timeLayout = "2006-01-02 15:04:05"

const (
	timeLayout                = "2006-01-02 15:04:05.000000"
	timeLayoutNumericTimezone = "2006-01-02 15:04:05.000000-0700"
)

type crcCheckInfoMap map[uint64][2]map[string]time.Time

func main() {
	subhealth := `
 {
     "OSD_SLOW_DISK_PERF": {
         "severity": 1,
         "suggest": "Please check the slow OSD detail disk perf",
         "time": "2020-01-07 19:49:56.299539",
         "slow_disks": {
             "osd.1": [
                 "1min:hdd_write_avg_lat [c=54,s=513571,avg=9510] > [c=60,s=14845,avg=247] x 10 at [0B,16383B);"
             ]
         }
     },
     "OSD_SLOW_NET_PERF": {
	     "severity": 1,
         "time": "2020-01-07 18:49:57.299539",
         "slow_ips": {
              "10.252.8.87": [
                 "10.252.8.87:[c=45,s=4527381,avg=100608] > 10.252.8.90:[c=48,s=4808483,avg=100176] x 1 at [0B,16383B);"
              ],
              "10.252.8.90": [
                  "10.252.8.90:[c=91,s=9084606,avg=99830] > 10.252.8.87:[c=86,s=4058198,avg=47188] x 1 at [0B,16383B);",
                 "10.252.8.90:[c=4165,s=418341491,avg=100442] > 10.252.8.87:[c=4256,s=218624808,avg=51368] x 1 at [0B,16383B);"
              ]
         },
         "suggest": "Please use ping/iperf and other network tools to check the slow networks"
     }
 }
`
	// subhealth = `
	// {
	// "OSD_SLOW_NET_PERF": {},
	// "OSD_SLOW_DISK_PERF": {}
	// }
	// `
	var subhealthStatusInfo SubhealthStatusInfo
	err := json.Unmarshal([]byte(subhealth), &subhealthStatusInfo)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(subhealthStatusInfo)
	fmt.Println(subhealthStatusInfo.SlowDiskInfo == nil)
	fmt.Println(subhealthStatusInfo.SlowDiskInfo)
	fmt.Println(subhealthStatusInfo.SlowNetworkInfo == nil)
	fmt.Println(subhealthStatusInfo.SlowNetworkInfo)

	const longForm = "2006-01-02 15:04:05"
	t, _ := time.Parse(longForm, "2020-01-07 18:49:57.299539")
	fmt.Println(t)

	fmt.Println("----------------------------")
	timeMap := make(map[*time.Time]string)
	now := time.Now()
	timeMap[&now] = "abc"
	fmt.Println(timeMap[&now])

	var buffer [256]byte
	fmt.Println(len(buffer))
	fmt.Println(cap(buffer))

	crcInfo := `
{
    "mismatched": {
        "osd.1": [
            {
                "time": "2020-03-02 15:55:27.621342",
                "object": "1:3d7ba35d:::object1:head"
            }
        ]
    },
    "recovered": {
        "osd.1": [
            {
                "time": "2020-03-02 15:55:27.621440",
                "object": "1:3d7ba35d:::object1:head"
            }
        ]
    }
}
`

	// crcInfo = `{}`
	fmt.Println("---------crc info-------------------")
	var info CRCMismatchInfo
	if err = json.Unmarshal([]byte(crcInfo), &info); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", info)
	fmt.Printf("%#v\n", info)
	fmt.Println(info.MismatchedInfo != nil)

	fmt.Println(time.Now().Format(timeLayout))

	t1 := "2020-03-02 15:55:27.621342"
	location := time.Now().Location()
	t2, _ := time.ParseInLocation(timeLayout, t1, location)
	fmt.Println(t2)
	fmt.Println(t2.Format(timeLayoutNumericTimezone))
	fmt.Println(time.Time{}.Format(timeLayoutNumericTimezone))

	fmt.Println("--------------map--------------")
	infoMap := make(map[uint64][]map[string]time.Time)
	fmt.Println(infoMap)
	if _, ok := infoMap[1]; !ok {
		infoMap[1] = []map[string]time.Time{
			make(map[string]time.Time),
			make(map[string]time.Time),
		}
	}
	fmt.Println(infoMap[1][1])
	infoMap[1][1] = map[string]time.Time{
		"a": time.Now(),
		"b": time.Now(),
	}
	var infoArray [2]map[string]time.Time
	fmt.Printf("[2]map[string]time.Time= %T\n", infoArray)
	// fmt.Printf("[2]map[string]time.Time= %T\n", infoMap[1])
	infoArray[0] = map[string]time.Time{
		"a": time.Now(),
		"b": time.Now(),
	}
	fmt.Println(infoMap)
	fmt.Println(infoArray)

}

// CRCMismatchInfo contains the mismatched and recoverd CRC info.
type CRCMismatchInfo struct {
	MismatchedInfo map[string][]*CRCObjectInfo `json:"mismatched"`
	RecoveredInfo  map[string][]*CRCObjectInfo `json:"recovered"`
}

func (info CRCMismatchInfo) String() string {
	output := []string{"MismatchedInfo: ["}
	for osdName, obj := range info.MismatchedInfo {
		output = append(output, fmt.Sprintf("osdInfo: %s, object: %s", osdName, obj))
	}
	output = append(output, "]\nRecoveredInfo: [")
	for osdName, obj := range info.RecoveredInfo {
		output = append(output, fmt.Sprintf("osdInfo: %s, object: %s", osdName, obj))
	}
	output = append(output, "]\n")
	return strings.Join(output, "")
}

// CRCObjectInfo returns the object info of mismatched and recovered CRC.
type CRCObjectInfo struct {
	Time   time.Time `json:"time"`
	Object string    `json:"object"`
}

func (s *CRCObjectInfo) String() string {
	return fmt.Sprintf("CRCObjectInfo[Time=%s, Object=%s]", s.Time, s.Object)
}

// UnmarshalJSON overrides the Unmarshaler.UnmarshalJSON
func (s *CRCObjectInfo) UnmarshalJSON(data []byte) error {
	type Alias CRCObjectInfo
	aux := &struct {
		Time string `json:"time"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Time != "" {
		location := time.Now().Location()
		t, err := time.ParseInLocation(timeLayout, aux.Time, location)
		if err != nil {
			return err
		}
		s.Time = t
	}
	return nil
}

// SubhealthStatusInfo contains data from ceph subhealth status
type SubhealthStatusInfo struct {
	SlowDiskInfo    *OsdSlowDiskPerf    `json:"OSD_SLOW_DISK_PERF"`
	SlowNetworkInfo *OsdSlowNetworkPerf `json:"OSD_SLOW_NET_PERF"`
}

func (info SubhealthStatusInfo) String() string {
	return fmt.Sprintf("Info[SlowDiskInfo=%s, SlowNetworkInfo=%s\n", info.SlowDiskInfo, info.SlowNetworkInfo)
}

// OsdSlowDiskPerf contains osd slow disk info
type OsdSlowDiskPerf struct {
	Severity   uint64              `json:"severity"`
	Time       time.Time           `json:"time"`
	SlowDisks  map[string][]string `json:"slow_disks"`
	Suggestion string              `json:"suggest"`
}

// UnmarshalJSON overrides the Unmarshaler.UnmarshalJSON
func (s *OsdSlowDiskPerf) UnmarshalJSON(data []byte) error {
	type Alias OsdSlowDiskPerf
	aux := &struct {
		Time string `json:"time"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Time != "" {
		location := time.Now().Location()
		t, err := time.ParseInLocation(timeLayout, aux.Time, location)
		if err != nil {
			return err
		}
		s.Time = t
	}
	return nil
}

func (s *OsdSlowDiskPerf) String() string {
	return fmt.Sprintf("SlowDiskPerf[Severity=%d, Time=%v, SlowDisks=%v, Suggestion=%s]", s.Severity, s.Time, s.SlowDisks, s.Suggestion)
}

// OsdSlowNetworkPerf contains osd slow network info
type OsdSlowNetworkPerf struct {
	Severity     uint64              `json:"severity"`
	Time         time.Time           `json:"time"`
	SlowNetworks map[string][]string `json:"slow_ips"`
	Suggestion   string              `json:"suggest"`
}

// UnmarshalJSON overrides the Unmarshaler.UnmarshalJSON
func (s *OsdSlowNetworkPerf) UnmarshalJSON(data []byte) error {
	type Alias OsdSlowNetworkPerf
	aux := &struct {
		Time string `json:"time"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	fmt.Printf("********time=%s***\n", aux.Time)
	if aux.Time != "" {
		location := time.Now().Location()
		t, err := time.ParseInLocation(timeLayout, aux.Time, location)
		if err != nil {
			return err
		}
		s.Time = t
	}
	return nil
}

func (s *OsdSlowNetworkPerf) String() string {
	return fmt.Sprintf("SlowNetworkPerf[Severity=%d, Time=%v, SlowNetworks=%v, Suggestion=%s]", s.Severity, s.Time, s.SlowNetworks, s.Suggestion)
}
