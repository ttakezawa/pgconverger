package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ttakezawa/pgconverger/diff"
)

func main() {
	var (
		source  = flag.String("source", "", "compare {SOURCE} to {DESIRED}")
		desired = flag.String("desired", "", "compare {SOURCE} to {DESIRED}")
	)
	flag.Parse()

	log.Printf("source:  %s", *source)
	log.Printf("desired: %s", *desired)

	sourceFile, err := os.Open(*source)
	if err != nil {
		log.Fatal(err.Error())
	}

	desiredFile, err := os.Open(*desired)
	if err != nil {
		log.Fatal(err.Error())
	}

	err, df := diff.Process(sourceFile, desiredFile)
	if err != nil {
		log.Printf("%+v", err)
		return
	}

	fmt.Println(df)
}
