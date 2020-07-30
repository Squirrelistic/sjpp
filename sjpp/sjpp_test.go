package main

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const TestFilesFolder = "../test-files"

func TestJsonTemplateProcessing(t *testing.T) {
	files, err := ioutil.ReadDir(TestFilesFolder)
	check(err)

	testCaseRegex := regexp.MustCompile(`^\d\d\.\D.+\.sjt$`)
	for _, file := range files {
		if !file.IsDir() && testCaseRegex.FindStringIndex(file.Name()) != nil {
			var testName = strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			t.Log("Testing file: ", testName)

			paramsMap := loadParamsFromFile(filepath.Join(TestFilesFolder, testName+".params.json"))
			actualContent := processTemplate(filepath.Join(TestFilesFolder, file.Name()), paramsMap)

			expectedContent, err := ioutil.ReadFile(filepath.Join(TestFilesFolder, testName+".expected.json"))
			check(err)

			expectedResult := prettyPrint(string(expectedContent))
			actualResult := prettyPrint(actualContent)
			if string(expectedResult) != actualResult {
				t.Errorf("Test %s failed\nExpected:\n%s\nActual:\n%s", testName, expectedResult, actualResult)
			}
		}
	}
}
