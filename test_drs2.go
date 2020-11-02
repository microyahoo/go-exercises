package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

type t struct {
	A string
}
type mig struct {
	T *t
}

func NewMig() *mig {
	return &mig{
		T: &t{},
	}
}

func (o *mig) initMigrationCmd() (*exec.Cmd, error) {
	var options []string

	if len(o.T.A) > 0 {
		options = append(options, "--files-from")
		options = append(options, o.T.A)
	}
	cmd := exec.Command("rclone", options...)
	return cmd, nil
}

func main() {
	m := NewMig()
	//var cmd *exec.Cmd
	count := 8
	totals := []int{}
	for i := 0; i < count; i++ {
		totals = append(totals, i)
	}

	var wg sync.WaitGroup
	tdArray := splitArray(totals, 5)

	copyData := func(m *mig) error {
		fmt.Println("pre ......", m, m.T.A)
		time.Sleep(time.Second * 2)
		fmt.Println("post ......", m, m.T.A)
		return nil
	}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(m)
	for _, subArray := range tdArray {
		wg.Add(len(subArray))
		for _, num := range subArray {
			m1 := NewMig()
			json.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(m1)
			// m1 := *m
			m1.T.A = fmt.Sprintf("haha+%d", num)
			go func(m *mig) {
				defer wg.Done()
				_ = copyData(m)
			}(m1)
		}
		wg.Wait()
		fmt.Println("---------------------------next---------------------")
	}
}

func splitArray(arr []int, num int) [][]int {
	var segmens = make([][]int, 0)
	max := int(len(arr))
	if max <= num || num <= 0 {
		return append(segmens, arr)
	}
	var i = 0
	for {
		if i <= max {
			if i+num <= max {
				if len(arr[i:i+num]) != 0 {
					segmens = append(segmens, arr[i:i+num])
				}
			} else {
				if len(arr[i:]) != 0 {
					segmens = append(segmens, arr[i:])
				}
				break
			}
		} else {
			if i-num < max {
				if len(arr[i-num:]) != 0 {
					segmens = append(segmens, arr[i-num:])
				}
			}
			break
		}
		i += num
	}
	return segmens
}
