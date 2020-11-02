package main

import (
	"fmt"
)

type Migration struct {
	Name    string
	Options *MigrationOptions
}

func (m *Migration) InitMigrationCmd() error {
	fmt.Printf("%#v\n", m)
	fmt.Printf("%#v\n", m.Options)
	return nil
}

type MigrationOptions struct {
	FilesFrom string `json:"filesFrom"`
}

func NewMigration() *Migration {
	return &Migration{
		Options: &MigrationOptions{},
	}
}

func main() {
	migration := NewMigration()

	migration.Name = "zhucan"
	migration.Options = &MigrationOptions{
		FilesFrom: "canbi.avi",
	}

	copyData := func() error {
		err := migration.InitMigrationCmd()
		if err != nil {
			fmt.Println(err)
		}
		return nil
	}

	copyData2 := func(m *Migration) error {
		err := m.InitMigrationCmd()
		if err != nil {
			fmt.Println(err)
		}
		return nil
	}

	copyData3 := func(m Migration) error {
		err := m.InitMigrationCmd()
		if err != nil {
			fmt.Println(err)
		}
		return nil
	}

	migration.Options.FilesFrom = "houbi.avi"
	copyData()

	m := *migration
	m.Options.FilesFrom = "doubi.avi"
	copyData2(&m)

	m3 := *migration
	m3.Options.FilesFrom = "doubi2.avi"
	copyData3(m3)

	count := 1000
	var i int64
	totals := []int64{}
	for i = 0; i < int64(count); i++ {
		totals = append(totals, i)
	}
	tdArray := SplitArray(totals, 6)
	for i := range tdArray {
		fmt.Println(tdArray[i], "\n")
	}
}

func SplitArray(arr []int64, num int64) [][]int64 {
	max := int64(len(arr))
	if max < num {
		return nil
	}
	var segmens = make([][]int64, 0)
	quantity := max / num
	end := int64(0)
	for i := int64(1); i <= num; i++ {
		qu := i * quantity
		if i != num {
			segmens = append(segmens, arr[i-1+end:qu])
		} else {
			segmens = append(segmens, arr[i-1+end:])
		}
		end = qu - i
	}
	return segmens
}
