package main

// See: https://www.hackerrank.com/challenges/camelcase/problem?isFullScreen=true

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Complete the camelcase function below.
func camelcase(s string) int32 {

	if len(s) == 0 || len(s) == 1 {
		return int32(len(s))
	}

	wc := 1

	for i := 1; i < len(s); i++ {
		if strings.IndexByte(uppercase, s[i]) > -1 {
			// Rune on the right-hand side (i+1) begins a new word
			wc++
		}
	}

	return int32(wc)

}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 1024*1024)

	s := readLine(reader)

	result := camelcase(s)

	fmt.Fprintf(writer, "%d\n", result)

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
