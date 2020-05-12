package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ErrorCodeDesc ...
type ErrorCodeDesc struct {
	Desc string `json:"desc"`
}

// ErrorCodeDescSuggest ...
type ErrorCodeDescSuggest struct {
	Desc    string `json:"desc"`
	Suggest string `json:"suggest"`
	ZHValue string `json:"zh_value"`
}

func main() {
	// errCodes := make(map[string]ErrorCodeDesc)
	errCodeSuggests := make(map[string]ErrorCodeDescSuggest)
	var errCodes map[string]ErrorCodeDesc
	errCodeSuggests2 := make(map[string]ErrorCodeDescSuggest)

	// read file
	data1, err := ioutil.ReadFile("./error_code.json")
	if err != nil {
		fmt.Print(err)
	}
	data2, err := ioutil.ReadFile("./code.json")
	if err != nil {
		fmt.Print(err)
	}

	// unmarshall it
	err = json.Unmarshal(data1, &errCodes)
	if err != nil {
		fmt.Println("error error_code.json: ", err)
	}
	err = json.Unmarshal(data2, &errCodeSuggests)
	if err != nil {
		fmt.Println("error code.json: ", err)
	}
	// fmt.Printf("%#v\n", errCodes)
	// fmt.Printf("%#v\n", errCodeSuggests)

	// compare it
	for key1, val1 := range errCodes {
		if val2, ok := errCodeSuggests[key1]; ok {
			if val1.Desc != val2.Desc {
				fmt.Printf("Not equal of two json files: key=%s, val1 = %s, val2= %s\n", key1, val1.Desc, val2.Desc)
			}
		} else {
			errCodeSuggests2[key1] = ErrorCodeDescSuggest{Desc: val1.Desc}
		}
	}
	var codeBytes []byte
	if codeBytes, err = json.MarshalIndent(errCodeSuggests2, "", "\t"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	destFile, err := os.Create("output.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer destFile.Close()
	if _, err = destFile.Write(codeBytes); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("hello world\n", 123)
}
