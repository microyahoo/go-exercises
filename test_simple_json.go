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

type ClientCodesInfo struct {
	ClientList []*ClientCode `json:"clientlist"`
}

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
	fmt.Println(len(nil))
}
