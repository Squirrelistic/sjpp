package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func jsonMarshal(i interface{}) string {
	b, err := json.Marshal(i)
	check(err)
	return string(b)
}

func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	check(err)
	s := string(b)
	return s[1 : len(s)-1]
}

func processTemplate(templateFile string, paramsMap map[string]interface{}) string {
	absTemplateFile, err1 := filepath.Abs(templateFile)
	check(err1)

	fmt.Println("Loading template:", absTemplateFile)
	templateContent, err2 := ioutil.ReadFile(absTemplateFile)
	check(err2)

	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"json": func(i interface{}) string {
			return jsonMarshal(i)
		},
		"jsonEscape": func(s string) string {
			return jsonEscape(s)
		},
		"include": func(includeFile string) string {
			if !filepath.IsAbs(includeFile) {
				// non-absolute path are relative to the parent template file location
				includeFile = filepath.Join(filepath.Dir(absTemplateFile), includeFile)
			}
			// TODO: Detect endless recursion (caused by circular include)
			return processTemplate(includeFile, paramsMap)
		},
	}).Option("missingkey=error").Parse(string(templateContent)))

	outputBuffer := new(bytes.Buffer)
	err3 := tmpl.Execute(outputBuffer, paramsMap)
	check(err3)

	return outputBuffer.String()
}

func loadParamsFromFile(paramsFile string) map[string]interface{} {
	absParamsFile, err1 := filepath.Abs(paramsFile)
	check(err1)

	fmt.Println("Loading parameters:", absParamsFile)
	paramsContent, err2 := ioutil.ReadFile(absParamsFile)
	check(err2)

	paramsMap := map[string]interface{}{}
	err3 := json.Unmarshal(paramsContent, &paramsMap)
	check(err3)

	return paramsMap
}

func prettyPrint(jsonText string) string {
	paramsMap := map[string]interface{}{}
	err1 := json.Unmarshal([]byte(jsonText), &paramsMap)
	if err1 != nil {
		fmt.Fprintln(os.Stderr, "Output JSON in not valid:\n"+jsonText)
	}
	check(err1)

	b, err2 := json.MarshalIndent(paramsMap, "", "    ")
	check(err2)

	return string(b)
}

func saveOutputToFile(output string, outputFile string) {
	absOutputFile, err1 := filepath.Abs(outputFile)
	check(err1)

	fmt.Println("Saving output:", absOutputFile)
	outputWriter, err2 := os.Create(absOutputFile)
	check(err2)
	defer outputWriter.Close()

	bufferedOutputWriter := bufio.NewWriter(outputWriter)
	_, err3 := bufferedOutputWriter.WriteString(output)
	check(err3)
	bufferedOutputWriter.Flush()
}

var version = "1.0.0"

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	}()

	templateFileName := flag.String("t", "", "Squirrelistic JSON Template (SJT) file")
	paramsFileName := flag.String("p", "", "Parameters file (JSON)")
	outputFileName := flag.String("o", "", "Output file")

	flag.Parse()

	if *templateFileName == "" || *paramsFileName == "" || *outputFileName == "" {
		fmt.Fprintln(os.Stderr, "Squirrelistic JSON Preprocessor "+version)
		flag.PrintDefaults()
		os.Exit(1)
	}

	paramsMap := loadParamsFromFile(*paramsFileName)
	output := processTemplate(*templateFileName, paramsMap)
	prettyOutput := prettyPrint(output)

	saveOutputToFile(prettyOutput, *outputFileName)
}
