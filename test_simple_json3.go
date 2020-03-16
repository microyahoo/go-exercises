package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/patrickmn/go-cache"
	"os"
	"reflect"
	"strings"
	"time"
)

//ClientCodesInfo ...
type ClientCodesInfo struct {
	ClientList []*ClientCode `json:"clientlist"`
}

// ClientCode ...
type ClientCode struct {
	Iqn *string `json:"iqn"`
	Wwn *string `json:"wwn"`
}

func (cc *ClientCode) String() string {
	return fmt.Sprintf("iqn = %v, wwn = %v", *cc.Iqn, *cc.Wwn)
}

func main() {
	out := "{ \"clientlist\": [ { \"iqn\": \"iqn.1994-05.com.redhat:3261d96c63a\" }, { \"wwn\": \"iqn.1994-05.com.redhat:3261d96c63a\" }] }"
	data, _ := simplejson.NewJson([]byte(out))
	fmt.Println(data)
	iqnArr, _ := data.Get("clientlist").Array()
	fmt.Println(iqnArr)
	for i, a := range iqnArr {
		fmt.Println(i)
		fmt.Println(a)
		if a, ok := a.(map[string]interface{}); ok {
			if code, ok := a["iqn"]; ok {
				fmt.Println(code)
				fmt.Printf("%v %T\n", code, code)
			}
			if code, ok := a["wwn"]; ok {
				fmt.Println(code)
				fmt.Printf("%v %T\n", code, code)
			}
		}
	}

	fmt.Println("----------1---------")
	m := make(map[int]string)
	m[2] = "a"
	m[4] = "b"
	for x := range m {
		fmt.Println(x)
	}

	fmt.Println("----------2---------")
	var hostMap map[string]int
	hostMap = make(map[string]int)
	hostMap["a"] = 1
	hostMap["b"] = 2
	fmt.Println(hostMap)

	fmt.Println("----------3---------")
	out = "{ \"clientlist\": [] }"
	// out = "{}"
	// var clientCodes *ClientCodesInfo
	clientCodes := new(ClientCodesInfo)
	// clientCodes := make(map[string][]*ClientCode)
	if err := json.Unmarshal([]byte(out), clientCodes); err != nil {
		fmt.Println("Unmarshal failed.")
	}
	fmt.Println(clientCodes)
	for _, cc := range clientCodes.ClientList {
		// for _, cc := range clientCodes["clientlist"] {
		if cc.Iqn != nil {
			fmt.Println(*cc.Iqn)
		}
		if cc.Wwn != nil {
			fmt.Println(*cc.Wwn)
		}
	}
	fmt.Println("----------4---------")
	fmt.Println(len(hostMap))
	wwn1 := "11:22"
	wwn2 := "33:44"
	iqn1 := "iqn-1"
	iqn2 := "iqn-2"
	clientCodes = &ClientCodesInfo{ClientList: []*ClientCode{{Wwn: &wwn1, Iqn: &iqn1}, {Wwn: &wwn2, Iqn: &iqn2}}}
	fmt.Println(clientCodes)
	ccBytes, err := json.Marshal(clientCodes)
	if err != nil {
		fmt.Println("Marshal failed")
	} else {
		fmt.Println(ccBytes)
		var out bytes.Buffer
		json.Indent(&out, ccBytes, "", "\t")
		out.WriteTo(os.Stdout)
	}
	bs := `bs`
	ss := "string"
	fmt.Printf("\n%T %T\n", bs, ss)

	fmt.Println("----------5---------")
	const jsonStream = `
        {"name":"ethancai", "fansCount": 9223372036854775807}
    `

	decoder := json.NewDecoder(strings.NewReader(jsonStream))
	decoder.UseNumber() // UseNumber causes the Decoder to unmarshal a number into an interface{} as a Number instead of as a float64.

	var user interface{}
	if err := decoder.Decode(&user); err != nil {
		fmt.Println("error:", err)
		return
	}

	mm := user.(map[string]interface{})
	fansCount := mm["fansCount"]
	fmt.Printf("%+v \n", reflect.TypeOf(fansCount).PkgPath()+"."+reflect.TypeOf(fansCount).Name())

	v, err := fansCount.(json.Number).Int64()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("%+v \n", v)

	fmt.Println("----------6---------")
	b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
	var f interface{}
	err = json.Unmarshal(b, &f)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(f)
	fmt.Println("----------7---------")
	cac := cache.New(time.Second*2, 5*time.Second)
	cac.Set("x", "mmm", cache.NoExpiration)
	cac.Set("y", "nnn", cache.DefaultExpiration)
	foo, found := cac.Get("x")
	if found {
		fmt.Println(foo)
	}
	bar, found := cac.Get("y")
	if found {
		fmt.Println(bar)
	}
	time.Sleep(time.Second * 5)
	foo, found = cac.Get("x")
	if found {
		fmt.Println(foo)
	} else {
		fmt.Println("x not found")
	}
	bar, found = cac.Get("y")
	if found {
		fmt.Println(bar)
	} else {
		fmt.Println("y not found")
	}
	time.Sleep(time.Second * 5)
	foo, found = cac.Get("x")
	if found {
		fmt.Println(foo)
	} else {
		fmt.Println("x not found")
	}
	fmt.Println(cac.Items())
	// fmt.Println(len(nil))

	fmt.Println("----------8---------")
	j := simplejson.New()
	j.SetPath([]string{"ResourceTrashResource", "resource_id"}, 1)
	j.SetPath([]string{"ResourceTrashResource", "resource_type"}, "block_volume")
	fmt.Println(j)

	metadata := `
{
    "id": 1,
    "addition_hb_addr": "10.255.101.75:6808/28436",
    "arch": "x86_64",
    "back_addr": "10.255.101.75:6807/28436",
    "ceph_version": "ceph version SDS_4.2.000.0.191128.2-7-gd2ac60c (d2ac60ce3cb6701cbd9e6a0af44875780bf74940) mimic (dev)",
    "cpu": "Intel(R) Xeon(R) CPU E5-2680 0 @ 2.70GHz",
    "default_device_class": "hdd",
    "distro": "centos",
    "distro_description": "CentOS Linux 7 (Core)",
    "distro_version": "7",
    "front_addr": "10.255.101.75:6806/28436",
    "hb_back_addr": "10.255.101.75:6809/28436",
    "hb_front_addr": "10.255.101.75:6810/28436",
    "hostname": "ceph-3",
    "kernel_description": "#1 SMP Tue Aug 22 21:09:27 UTC 2017",
    "kernel_version": "3.10.0-693.el7.x86_64",
    "mem_swap_kb": "0",
    "mem_total_kb": "16268092",
    "os": "Linux",
    "osd_data": "/var/lib/ceph/osd/ceph-1",
    "osd_objectstore": "xstore2",
    "pid": 28436,
    "rotational": "1",
    "uuid": "a7dac8de-9e63-45b3-8730-e0c72959771f",
    "up": true,
    "in": true,
    "exit_code": 0
}
`
	var osdMetadata OsdMetadata
	err = json.Unmarshal([]byte(metadata), &osdMetadata)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("metadata = %#v", osdMetadata)

}

// OsdMetadata defines osd metadata info
type OsdMetadata struct {
	ID             uint64 `json:"id"`
	Up             bool   `json:"up"`
	In             bool   `json:"in"`
	ExitCode       int    `json:"exit_code"`
	HostName       string `json:"hostname"`
	OS             string `json:"os"`
	OsdData        string `json:"osd_data"`
	OsdObjectStore string `json:"osd_objectstore"`
	Arch           string `json:"arch"`
	AdditionHbAddr string `json:"addition_hb_addr"`
}
