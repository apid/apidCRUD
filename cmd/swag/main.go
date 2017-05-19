package main

import (
	"os"
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"gopkg.in/yaml.v2"
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

func emitHeader() {
	fmt.Printf("%s\n", "// GENERATED FILE - DO NOT EDIT")
	fmt.Printf("%s\n", "package apidCRUD")
	fmt.Printf("%s\n", `import "net/http"`)
}

func emitSwaggerJson(cont []byte) {
	out := contToJson(cont)
	fmt.Printf("const swaggerJSON = `%s`\n", out)
}

func mapcon(obj interface{}) map[string]interface{} {
	switch xobj := obj.(type) {
	case map[string]interface{}:
		return xobj
	case []interface{}:
		return mapcon(xobj[0])
	default:
		fmt.Fprintf(os.Stderr,
			"Error: mapcon conversion error on (%T) %v\n", obj, obj)
		os.Exit(2)
	}
	// NOTREACHED
	return nil
}

func strcon(obj interface{}) string {
	ret, ok := obj.(string)
	if !ok {
		return ""
	}
	return ret
}

func toMethod(str string) string {
	switch {
	case strings.EqualFold(http.MethodGet, str):
		return "http.MethodGet"
	case strings.EqualFold(http.MethodPost, str):
		return "http.MethodPost"
	case strings.EqualFold(http.MethodPut, str):
		return "http.MethodPut"
	case strings.EqualFold(http.MethodPatch, str):
		return "http.MethodPatch"
	case strings.EqualFold(http.MethodDelete, str):
		return "http.MethodDelete"
	default:
		return ""
	}
}

func emitApiTable(cont []byte) {
	obj := map[interface{}]interface{} {}
	err := yaml.Unmarshal(cont, &obj)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error: %s - could not Unmarshal\n", err)
		os.Exit(1)
	}

	cvobj := convert(obj)

	m1 := mapcon(cvobj)
	basePath := m1["basePath"]

	fmt.Printf("\n// basePath = \"%s\"\n", basePath)

	fmt.Printf("\n%s\n", `var apiTable = []apiDesc {`)
	paths := mapcon(m1["paths"])

	// iterate over paths
	for path, val1 := range(paths) {
		m2 := mapcon(val1)
		// iterate over verbs
		for verb, val2 := range(m2) {
			m3 := mapcon(val2)
			method := toMethod(verb)
			if method == "" {
				continue
			}
			handler := strcon(m3["operationId"]) + "Handler"
			fmt.Printf("\t{\"%s\", %s, %s},\n",
				path, method, handler)
		}
	}

	fmt.Printf("%s\n", `}`)
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

	emitHeader()
	emitSwaggerJson(cont)
	emitApiTable(cont)
}
