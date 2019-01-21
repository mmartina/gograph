package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MarcGrol/golangAnnotations/generator"
	"github.com/MarcGrol/golangAnnotations/parser"
)

const Version = "0.1"

func main() {
	inputDir, tags := processArgs()

	parsedSources, err := parser.New().ParseSourceDir(inputDir, "^.*.go$", "^"+generator.GenfilePrefix+".*.go$")
	if err != nil {
		log.Fatalf("Error parsing golang sources in %s: %s", inputDir, err)
		os.Exit(-1)
	}

	g := NewGenerator()
	err = g.Generate(inputDir, tags, parsedSources)
	if err != nil {
		log.Fatalf("Error executing %s: %s\n", os.Args[0], err)
		os.Exit(-2)
	}

	os.Exit(0)
}

func processArgs() (string, string) {
	inputDir := flag.String("input-dir", "", "Directory to be examined")
	tags := flag.String("tags", "", "Build tags to be added")
	help := flag.Bool("help", false, "Usage information")
	version := flag.Bool("version", false, "Version information")
	flag.Parse()
	if *version == true {
		fmt.Fprintf(os.Stderr, "Version: %s\n", Version)
		os.Exit(1)
	}
	if *help == true || *inputDir == "" {
		fmt.Fprintf(os.Stderr, "Usage:\n%s [flags]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	return *inputDir, *tags
}
