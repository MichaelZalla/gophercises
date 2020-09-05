package main

// See: https://www.hackerrank.com/challenges/caesar-cipher-1/problem?isFullScreen=true

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const upper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var lower = strings.ToLower(upper)

func rotate(r byte, base, delta int) byte {
	offset := int(r) - base
	offset = (offset + delta) % 26
	return byte(base + offset)
}

func cipher(r byte, k int) byte {

	if r >= 'A' && r <= 'Z' {
		return rotate(r, 'A', k)
	}

	if r >= 'a' && r <= 'z' {
		return rotate(r, 'a', k)
	}

	return r

}

// Complete the caesarCipher function below.
func caesarCipher(s string, k int32) string {

	var builder strings.Builder

	for _, b := range s {
		err := builder.WriteByte(cipher(byte(b), int(k)))
		if err != nil {
			panic(err)
		}
	}

	return builder.String()

}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 1024*1024)

	nTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)
	n := int32(nTemp)

	_ = n

	s := readLine(reader)

	kTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)
	k := int32(kTemp)

	result := caesarCipher(s, k)

	fmt.Fprintf(writer, "%s\n", result)

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
