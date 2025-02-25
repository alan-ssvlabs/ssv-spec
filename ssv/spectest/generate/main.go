package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec/ssv/spectest/tests"
	"github.com/ssvlabs/ssv-spec/types"
	comparable2 "github.com/ssvlabs/ssv-spec/types/testingutils/comparable"

	"github.com/ssvlabs/ssv-spec/ssv/spectest"
)

//go:generate go run main.go

func main() {
	clearStateComparisonFolder()

	all := map[string]tests.SpecTest{}
	for _, testF := range spectest.AllTests {
		test := testF()

		// write json test
		n := reflect.TypeOf(test).String() + "_" + test.TestName()
		if all[n] != nil {
			panic(fmt.Sprintf("duplicate test: %s\n", n))
		}
		all[n] = test
	}

	byts, err := json.Marshal(all)
	if err != nil {
		panic(err.Error())
	}

	if len(all) != len(spectest.AllTests) {
		panic("did not generate all tests\n")
	}

	log.Printf("found %d tests\n", len(all))
	writeJson(byts)

	for _, testF := range spectest.AllTests {
		test := testF()

		// generate post state comparison
		post, err := test.GetPostState()
		if err != nil {
			err = errors.Wrapf(err, "failed to get post state for test: %s", test.TestName())
			panic(err.Error())
		}
		writeJsonStateComparison(test.TestName(), reflect.TypeOf(test).String(), post)
	}
}

func clearStateComparisonFolder() {
	_, basedir, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller info")
	}
	dir := filepath.Join(strings.TrimSuffix(basedir, "main.go"), "state_comparison")

	if err := os.RemoveAll(dir); err != nil {
		panic(err.Error())
	}

	if err := os.Mkdir(dir, 0700); err != nil {
		panic(err.Error())
	}
}

func writeJsonStateComparison(name, testType string, post interface{}) {
	postMap, ok := post.(map[string]types.Root)

	if !ok {
		writeSingleSCJson(name, testType, post)
		return
	}
	name = strings.ReplaceAll(name, " ", "_")
	for subTestName, postState := range postMap {
		writeSingleSCJson(subTestName, filepath.Join(testType, name), postState)
	}
}

func writeSingleSCJson(path string, testType string, post interface{}) {
	if post == nil { // If nil, test not supporting post state comparison yet
		log.Printf("skipping state comparison json, not supported: %s\n", path)
		return
	}
	byts, err := json.MarshalIndent(post, "", "		")
	if err != nil {
		panic(err.Error())
	}

	scDir := scDir(testType)
	if err != nil {
		panic(err.Error())
	}

	file := filepath.Join(scDir, fmt.Sprintf("%s.json", path))
	// try to create directory if it doesn't exist
	if err := os.MkdirAll(scDir, 0700); err != nil && !os.IsExist(err) {
		panic(err.Error())
	}

	log.Printf("writing state comparison json: %s\n", file)
	if err := os.WriteFile(file, byts, 0400); err != nil {
		panic(err.Error())
	}
}

func scDir(testType string) string {
	_, basedir, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller info")
	}
	basedir = strings.TrimSuffix(basedir, "main.go")
	scDir := comparable2.GetSCDir(basedir, testType)
	return scDir
}

func writeJson(data []byte) {
	_, basedir, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller info")
	}
	basedir = strings.TrimSuffix(basedir, "main.go")

	// try to create directory if it doesn't exist
	_ = os.Mkdir(basedir, os.ModeDir)

	file := filepath.Join(basedir, "tests.json.gz")

	// Create a buffer to write the gzipped file to
	buf := new(bytes.Buffer)

	// Create a new gzip writer
	gzipWriter := gzip.NewWriter(buf)

	// Write the JSON data to the gzip writer
	if _, err := gzipWriter.Write(data); err != nil {
		panic(err.Error())
	}

	// Close the gzip writer
	if err := gzipWriter.Close(); err != nil {
		panic(err.Error())
	}

	// Write the gzipped data to a file
	if err := os.WriteFile(file, buf.Bytes(), 0400); err != nil {
		panic(err.Error())
	}
}
