package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"encoding/json"
)

// borrowed from:
// http://stackoverflow.com/questions/40737122/convert-yaml-to-json-without-struct-golang
func convert(i interface{}) interface{} {
    switch x := i.(type) {
    case map[interface{}]interface{}:
        m2 := map[string]interface{}{}
        for k, v := range x {
            m2[k.(string)] = convert(v)
        }
        return m2
    case []interface{}:
        for i, v := range x {
            x[i] = convert(v)
        }
    }
    return i
}

func fileToJson(cont []byte) []byte {
	obj := map[interface{}]interface{} {}
	err := yaml.Unmarshal(cont, &obj)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error: %s - could not Unmarshal\n", err)
		os.Exit(1)
	}

	cvobj := convert(obj)

	jsonStr, err := json.MarshalIndent(cvobj, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error: %s - json Marshal\n", err)
		os.Exit(1)
	}
	return jsonStr
}

func main() {
	var cont []byte
	var err error
	if len(os.Args) < 2 {
		cont, err = ioutil.ReadAll(os.Stdin)
	} else {
		cont, err = ioutil.ReadFile(os.Args[1])
	}
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error: %s - could not read input\n", err)
		os.Exit(1)
	}
	out := fileToJson(cont)
	fmt.Printf("%s\n", out)
}
