// Command line: concordance.go <inputDoc> [<outputFile>]
//
//  <inputDoc> is a file with the text to be used for the concordance
//  <outputFile> (defaults to stdout) is the file where the resulting concordance should be written
package main

import (
	"io"
	"log"
	"os"

	"github.com/AOrnstein/concordance"
)

func main() {
	var output io.Writer = os.Stdout

	args := os.Args
	if len(args) < 2 {
		log.Fatal("usage: concordance <inputFile> [<outputFile>]\n")
	}

	input := args[1]
	inFile, err := os.Open(input)
	if err != nil {
		log.Fatalf("Invalid argument [1]: input file: \"%s\": %v\n", input, err)
	}
	defer inFile.Close()

	if len(args) > 2 {
		outFilePath := args[2]
		outFile, err := os.OpenFile(outFilePath, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Invalid argument [1]: output file: \"%s\": %v\n", outFilePath, err)
		}
		defer outFile.Close()
		output = outFile

	}

	err = concordance.GenerateConcordance(inFile, output)
	if err != nil {
		log.Fatalf("Failed to create concordance: %v\n", err)
	}
}
