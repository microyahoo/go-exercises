package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

func main() {
	fmt.Println(strings.Split("network__addresses", "__"))
	fmt.Println(0x7FFFFFF)

	rot13 := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			r += 'a' - 'A'
			return r
			// return 'A' + (r-'A'+13)%26
		case r >= 'a' && r <= 'z':
			r -= 'a' - 'A'
			return r
			// return 'a' + (r-'a'+13)%26
		}
		return r
	}
	fmt.Println(strings.Map(rot13, "'Twas brillig and the slithy gopher..."))

	ids := make([]string, 0, 2)
	str_map := make(map[int]string)
	str_map[1] = "a"
	str_map[2] = "b"
	str_map[3] = "c"
	for _, id := range str_map {
		ids = append(ids, id)
	}
	fmt.Println(ids)

	// GuessingGame()

	imageInfo := "xbd_data.2995801c"
	fmt.Println(strings.Split(imageInfo, "."))
	fmt.Println(time.Now().UTC().Format(time.RFC3339Nano))

	test_bytes()
	test_bytes2()
	test_string()
	test_string(nil)
	x := "xxx"
	y := "yyy"
	test_string(&x)
	test_string(&x, &y)
	fmt.Printf("print percent: %d%%\n", 100)
	test_string_field()
}

func GuessingGame() {
	var s string
	fmt.Printf("Pick an integer from 0 to 100.\n")
	answer := sort.Search(100, func(i int) bool {
		fmt.Printf("Is your number <= %d? ", i)
		fmt.Scanf("%s", &s)
		fmt.Println(s)
		return s != "" && s[0] == 'y'
	})
	fmt.Printf("Your number is %d.\n", answer)
}

func test_bytes() {
	s := `forward_rbd_trash_list cmd is {"poolname":"pool-25ae576490fe48eea5df5f7e6b50bc96","fsid":"1e3257bb-b0a4-4ade-a95c-cfb255614933","mon_hosts":["10.255.101.73:6789"]}
[{"id":"2391db5","name":"volume-4019909f17e14958b34fa95aeac65cd5","size":"12288"}]`
	index := strings.LastIndex(s, "[")
	fmt.Printf("the index is %v\n", index)
	fmt.Printf("rbd trash ls: %s\n", s[index:])
	// outBytes := []byte(newOut[index:])

	var results RbdTrashLsResult
	outBytes := []byte(s[index:])
	// if err := json.Unmarshal(out.Bytes(), &results); err != nil {
	if err := json.Unmarshal(outBytes, &results); err != nil {
		fmt.Printf("Failed to unmarshal bytes, err: %v", err)
	}
	fmt.Println(results)
}

type RbdTrashLsResult []*RbdTrashVolumeInfo

type RbdTrashVolumeInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size string `json:"size"`
	// Size int64  `json:"size"`
}

type RbdImage struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Format  int64  `json:"format"`
	Version int64  `json:"version"`
	// BlockNamePrefix string `json:"block_name_prefix"`
	ImageID string `json:"block_name_prefix"`
	// Parent          *RbdImageParentInfo `json:"parent"`
}

func test_bytes2() {
	var results RbdImage
	cmdResp := `{
    "block_name_prefix": "xbd_data.2cea9ae4",
    "flags": [],
    "format": 128,
    "name": "volume1",
    "object_size": 4194304,
    "objects": 25600,
    "order": 22,
    "size": 107374182400
}`
	if err := json.Unmarshal([]byte(cmdResp), &results); err != nil {
		fmt.Printf("Failed to unmarshal bytes, err: %v", err)
	}
	fmt.Println(results)
}

func test_string(test ...*string) {
	fmt.Println(len(test))
	fmt.Println(test)
	fmt.Println(test == nil)
	// if len(test) == 1 {
	//  // error
	// 	s, ok := test[0].(*string)
	// 	fmt.Println(s, ok)
	// }
}

func test_string_field() {
	var s = "Cron TrashExpiredScanCron @every 30d"
	a := strings.Fields(s)
	if len(a) < 3 {
		fmt.Println("the length of string is less than 3")
	} else if len(a) > 3 {
		tmp := a[2:]
		a[2] = strings.Join(tmp, " ")
		fmt.Println(a[2])
	}
	fmt.Println(a)
}
