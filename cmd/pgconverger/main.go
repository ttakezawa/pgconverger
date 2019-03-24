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
		from = flag.String("from", "", "compare {FROM} to {TO}")
		to   = flag.String("to", "", "compare {FROM} to {TO}")
	)
	flag.Parse()

	log.Printf("from: %s", *from)
	log.Printf("to:   %s", *to)

	toFile, err := os.Open(*to)
	if err != nil {
		log.Fatal(err.Error())
	}

	toFileData, err := ioutil.ReadAll(toFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	p := parser.New(lexer.Lex(string(toFileData)))
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
