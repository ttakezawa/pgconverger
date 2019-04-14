package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ttakezawa/pgconverger/lexer"
	"github.com/ttakezawa/pgconverger/parser"
)

func main() {
	var (
		source  = flag.String("source", "", "compare {SOURCE} to {DESIRED}")
		desired = flag.String("desired", "", "compare {SOURCE} to {DESIRED}")
	)
	flag.Parse()

	log.Printf("source:  %s", *source)
	log.Printf("desired: %s", *desired)

	desiredFile, err := os.Open(*desired)
	if err != nil {
		log.Fatal(err.Error())
	}

	desiredFileData, err := ioutil.ReadAll(desiredFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	p := parser.New(lexer.Lex(string(desiredFileData)))
	dataDefinition := p.ParseDataDefinition()

	if errors := p.Errors(); errors != nil {
		for _, err := range errors {
			log.Print(err.Error())
		}
	}

	var builder strings.Builder
	dataDefinition.Source(&builder)
	log.Printf("\n%s", builder.String())
}
