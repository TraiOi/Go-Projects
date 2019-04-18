package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"os"
	"path/filepath"

	"github.com/h2non/filetype"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ignoreFiles(filepath string) bool {
	buf, _ := ioutil.ReadFile(filepath)
	if filetype.IsImage(buf) {
		return true
	}
	return false
}

func ignoreRegexFiles(filepath string) bool {
	if(regexp.MustCompile(".+\\.js$").MatchString(filepath)) {
		return true
	}
	return false
}

func getListFiles() []string {
	result := []string{};
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if ! info.IsDir() {
			if ! ignoreRegexFiles(path) {
				result = append(result, path)
			}
		}
		return nil
	})
	check(err)
	return result
}

func readFile(filepath string) []byte {
	if ignoreFiles(filepath) == false {
		content, err := ioutil.ReadFile(filepath)
		check(err)
		return content
	}
	return []byte{0}
}

func checkShellInclude(filepath string) bool {
	content := string(readFile(filepath))
	re := regexp.MustCompile("\\@include \\\"(\\\\[0-9]{3}.*)+\\\"\\;")
	result := re.MatchString(content)
	return result
}

func checkShellHexNoEOL(filepath string) bool {
	content := string(readFile(filepath))
	re := regexp.MustCompile("(\\\\x[0-9A-Fa-f]{2}[0-9A-Za-z]*)+.+[^\x0A]$")
	result := re.MatchString(content)
	return result
}

func checkShellROT13NoEOL(filepath string) bool {
	content := string(readFile(filepath))
	re := regexp.MustCompile("str_rot13.+[^\x0A]$")
	result := re.MatchString(content)
	return result
}

func scanShell(filepath string) {
	if checkShellInclude(filepath) {
		shellpath := fmt.Sprintf("[@include] %s", filepath)
		fmt.Println(shellpath)
	} else if checkShellHexNoEOL(filepath) {
		shellpath := fmt.Sprintf("[Hex NoEOL] %s", filepath)
		fmt.Println(shellpath)
	} else if checkShellROT13NoEOL(filepath) {
		shellpath := fmt.Sprintf("[ROT13 NoEOL] %s", filepath)
		fmt.Println(shellpath)
	}
}

func main() {
	listfiles := getListFiles()
	for i := range listfiles {
		scanShell(listfiles[i])
	}
}
