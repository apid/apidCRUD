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

func contToJson(cont []byte) []byte {
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
	var err error
	F := os.Stdin
	if len(os.Args) > 1 {
		fn := os.Args[1]
		F, err = os.Open(fn)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"Error: %s - could not open %s\n", err, fn)
			os.Exit(2)
		}
	}

	cont, err := ioutil.ReadAll(F)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error: %s - read error\n", err)
		os.Exit(2)
	}
	out := contToJson(cont)
	fmt.Printf("%s\n", out)
}
