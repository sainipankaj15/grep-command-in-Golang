package worker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Result struct {
	Line    string
	LineNum int
	Path    string
}

type Results struct {
	Inner []Result
}

func NewResult(line string, lineNum int, path string) Result {
	// return Result{line, lineNum, path}
	return Result{Line: line, LineNum: lineNum, Path: path}
	// Both lines are same : Just a different method to write
}

func FindInFile(path string, find string) *Results {

	file, err := os.Open(path)

	if err != nil {
		fmt.Println("Error while file opening : ", err)
	}

	results := Results{Inner: make([]Result, 0)}

	// will create a scanner using bufio package
	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan(){
		if strings.Contains(scanner.Text(), find){
			foundResult := Result{scanner.Text(), lineNum, path}
			results.Inner = append(results.Inner, foundResult)
		}
		lineNum += 1
	}

	if len(results.Inner) == 0 {
		fmt.Println("Nothing found")
		return nil
	}else{
		return &results
	}
}
