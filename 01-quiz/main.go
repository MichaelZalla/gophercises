package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func normalize(input string) string {
	return strings.TrimSpace(strings.ToLower(input))
}

func getProblems(entries [][]string) []problem {

	problems := make([]problem, len(entries))

	for i, entry := range entries {
		problems[i] = problem{
			q: entry[0],
			a: normalize(entry[1]),
		}
	}

	return problems

}

func main() {

	// Define CLI flags

	var filepath = flag.String("csv", "./data/problems.csv", "path to your questions CSV file")

	var timeLimit = flag.Int("limit", 30, "time limit for the quiz (default: 30 seconds).")

	// Initialize CLI flags

	flag.Parse()

	// Get a file reader for the CSV file

	file, err := os.Open(*filepath)

	defer file.Close()

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to open CSV file '%v' (%s).", *filepath, err))
	}

	// Get a CSV reader to decode the CSV data into a native Go structure

	reader := csv.NewReader(file)

	// Scan through each record in the CSV file

	entries, err := reader.ReadAll()

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to read CSV file '%v' (%s).", &filepath, err))
	}

	problems := getProblems(entries)

	correct := 0

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	for i, p := range problems {

		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)

		answerCh := make(chan string)

		go func() {

			var answer string

			fmt.Scanf("%s\n", &answer)

			answerCh <- normalize(answer)

		}()

		select {

		case <-timer.C:

			fmt.Printf("You scored %d out of %d.\n", correct, len(problems))

			return

		case answer := <-answerCh:

			if answer == p.a {
				correct++
			}

		}

	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))

}
