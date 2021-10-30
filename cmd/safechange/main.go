package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"

	"github.com/xaionaro-go/safechange"
)

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func syntaxExit() {
	_, _ = fmt.Fprintf(os.Stderr, "syntax: safechange <go file A> <go file B>\n")
	os.Exit(2)
}

func main() {

	// init

	flag.Parse()
	if flag.NArg() != 2 {
		syntaxExit()
	}
	filePathA := flag.Arg(0)
	filePathB := flag.Arg(1)

	// read

	fileDataA, err := ioutil.ReadFile(filePathA)
	assertNoError(err)

	fileDataB, err := ioutil.ReadFile(filePathB)
	assertNoError(err)

	// parse

	fSetA := token.NewFileSet()
	astA, err := parser.ParseFile(fSetA, "", fileDataA, 0)
	assertNoError(err)
	fileA := safechange.NewAstFile(astA)

	fSetB := token.NewFileSet()
	astB, err := parser.ParseFile(fSetB, "", fileDataB, 0)
	assertNoError(err)
	fileB := safechange.NewAstFile(astB)

	// result

	if fileA.EquivalentTo(fileB) {
		os.Exit(0)
	}

	os.Exit(1)
}
