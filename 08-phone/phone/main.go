package main

import (
	"bytes"
	"context"
	"log"
	"regexp"

	phonedb "github.com/MichaelZalla/gophercises/08-phone/phone/db"
	_ "github.com/mattn/go-sqlite3"
)

const driverName = "sqlite3"
const dataSource = "./phone.db"

var ctx context.Context

var numbers []string

func main() {

	must(phonedb.Reset(driverName, dataSource))

	// Reset our database and phone numbers table

	must(phonedb.Migrate(driverName, dataSource))

	// Connect to our database

	db, err := phonedb.Open(driverName, dataSource)

	must(err)

	defer db.Close()

	err = db.Seed()

	must(err)

	// Fetch all phone number entries in the databas

	phones, err := db.GetPhones()

	must(err)

	for _, p := range phones {

		normalized := normalizeRegex(p.Number)

		if normalized != p.Number {

			// Check whether the normalized version of p.Number already exists

			match, err := db.FindPhone(normalized)

			must(err)

			if match != nil {

				// Normalized version has already been stored; delete this record

				log.Printf("Existing normalization found: %s\n", match.Number)
				log.Printf("Deleting...\n")

				must(db.DeletePhone(p.ID))

			} else {

				// Normalized version has not been stored; update this record

				log.Printf("New normalization found: %s\n", normalized)
				log.Printf("Updating...\n")

				p.Number = normalized

				must(db.UpdatePhone(&p))

			}

		}

	}

}

func normalizeIterative(number string) string {

	var buf bytes.Buffer

	for _, r := range number {
		if r >= '0' && r <= '9' {
			buf.WriteRune(r)
		}
	}

	return buf.String()

}

func normalizeRegex(number string) string {

	re := regexp.MustCompile("\\D")

	return re.ReplaceAllString(number, "")

}

func must(err error) error {

	if err != nil {
		log.Fatal(err)
	}

	return err

}
