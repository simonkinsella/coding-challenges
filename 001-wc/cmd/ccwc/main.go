package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	ExitCodeMissingParamFile = 1
	ExitCodeFileNotFound     = 2
)

func main() {
	bytes := flag.Bool("c", false, "The number of bytes in each input file is written to the standard output")
	lines := flag.Bool("l", false, "The number of lines in each input file is written to the standard output")
	words := flag.Bool("w", false, "The number of words in each input file is written to the standard output")
	chars := flag.Bool("m", false, "The number of bytes in each input file is written to the standard output")
	flag.Parse()
	file := flag.Arg(0)

	var err error
	var input *os.File

	switch flag.NArg() {
	case 0:
		input = os.Stdin
	case 1:
		input, err = os.Open(file)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(ExitCodeFileNotFound)
		}
	default:
		fmt.Println("Usage: -c <file>")
		os.Exit(ExitCodeMissingParamFile)
	}

	combinationMode := flag.NFlag() == 0
	var numWords, numChars, numLines int

	if *chars {
		reader := bufio.NewReader(input)
		for {
			_, _, err = reader.ReadRune()
			if err == io.EOF {
				break
			}
			numChars++
		}
	}

	if *bytes || combinationMode {
		input.Seek(0, io.SeekStart)
		reader := bufio.NewReader(input)
		for {
			n, err := reader.Discard(4096)
			numChars += n
			if err != nil {
				break
			}
		}
	}

	if *lines || combinationMode {
		input.Seek(0, io.SeekStart)
		reader := bufio.NewScanner(input)
		for reader.Scan() {
			numLines++
		}
	}

	if *words || combinationMode {
		input.Seek(0, io.SeekStart)
		reader := bufio.NewScanner(input)
		reader.Split(bufio.ScanWords)
		for reader.Scan() {
			numWords++
		}
	}
	if combinationMode {
		fmt.Printf("%8d %8d %8d %s\n", numLines, numWords, numChars, file)
	} else {
		fmt.Printf("%8d %s\n", numLines+numWords+numChars, file)
	}
}
