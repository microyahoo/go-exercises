package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-ini/ini"
)

// Embeded ...
type Embeded struct {
	Dates  []time.Time `delim:"|" comment:"Time data"`
	Places []string    `ini:"places,omitempty"`
	None   []int       `ini:",omitempty"`
}

// IP ...
type IP struct {
	Value []string `ini:"ips,omitempty,allowshadow"`
}

// Author ...
type Author struct {
	Name      string `ini:"NAME"`
	Male      bool
	Age       int `comment:"Author's age"`
	GPA       float64
	NeverMind string `ini:"-"`
	*Embeded  `comment:"Embeded section"`
	IP        `comment:"Author's ips"`
}

func main() {
	a := &Author{"Unknwon", true, 21, 2.8, "",
		&Embeded{
			[]time.Time{time.Now(), time.Now()},
			[]string{"HangZhou", "Boston"},
			[]int{},
		},
		IP{[]string{"10.255.101.74", "10.255.101.75"}},
	}
	cfg := ini.Empty(ini.LoadOptions{
		AllowShadows: true,
	})
	err := ini.ReflectFrom(cfg, a)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = cfg.SaveTo("test.ini")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
