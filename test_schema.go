package main

import (
	"fmt"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/xeipuuv/gojsonschema"
)

// Validate validates schema with document
func Validate(schema, document string) (result *gojsonschema.Result, err error) {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	documentLoader := gojsonschema.NewStringLoader(document)
	return gojsonschema.Validate(schemaLoader, documentLoader)
}

const (
	enabled = `{
		"type": "boolean"
	}`

	retention = `{
		"type": "string",
		"minLength": 2,
        "pattern": "^(-)?(([0-9\\.]+)w)?(([0-9\\.]+)d)?(([0-9\\.]+)h)?(([0-9\\.]+)m)?(([0-9\\.]+)s)?$"
	}`

	TrashUpdate = `{
		"type": "object",
		"properties": {
			"trash": {
				"type": "object",
				"properties": {
					"retention": ` + retention + `,
					"enabled": ` + enabled + `
				},
				"additionalProperties": false
			}
		},
		"required": ["trash"]
	}`
)

func main() {
	document := `{
		"trash": {
			"retention": "-0.30s",
			"enabled": false
		}
	}`
	document1 := `{
		"trash": {
			"retention": "-0.30s"
		}
	}`
	document2 := `{
		"trash": {
			"enabled": false
		}
	}`
	document3 := `{
		"trash": {
			"retention": "-0.30s",
			"enabled": false,
			"xx": 1
		}
	}`
	document4 := `{
		"trash": {
			"retention": "-x0.30s",
			"enabled": false
		}
	}`
	document5 := `{
		"trash": {
			"enabled": "xx"
		}
	}`
	document6 := `{
	}`
	document7 := `{
		"trash": {
			"retention": ""
		}
	}`
	// fmt.Println(TrashUpdate)
	_, err := simplejson.NewJson([]byte(document))
	if err != nil {
		fmt.Printf("failed to parse document json, err: %v", err)
		return
	}
	_, err = simplejson.NewJson([]byte(TrashUpdate))
	if err != nil {
		fmt.Printf("failed to parse TrashUpdate json, err: %v", err)
		return
	}
	documents := []string{document, document1, document2, document3,
		document4, document5, document6, document7}
	for _, doc := range documents {
		result, err := Validate(TrashUpdate, doc)
		if err != nil {
			fmt.Printf("Failed to validate schema, err: %v", err)
			return
		}
		if !result.Valid() {
			s := "invalid parameters:\n"
			fmt.Println(doc)
			for _, err := range result.Errors() {
				s += fmt.Sprintf("%s\n", err)
			}
			fmt.Printf(s)
		}
	}
}
