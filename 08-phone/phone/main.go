package main

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

type phone struct {
	id     int
	number string
}

func (p phone) String() string {
	return fmt.Sprintf("<id=%d, number=%s>", p.id, p.number)
}

const databaseName = "phone.db"

var ctx context.Context

var numbers []string

func main() {

	// Reset our database and phone numbers table

	db, err := resetDatabase()

	must(err)

	must(resetPhoneNumbersTable(db))

	// Read a list of unnormalized phone numbers from a data file

	unnormalizedPhoneNumbers, err := readLines("./data/numbers.txt")

	must(err)

	_, err = populatePhoneNumbersTable(db, unnormalizedPhoneNumbers)

	must(err)

	phones, err := getPhones(db)

	must(err)

	for _, p := range phones {

		normalized := normalizeRegex(p.number)

		if normalized != p.number {

			// Check whether the normalized version of p.number already exists

			match, err := findPhone(db, normalized)

			must(err)

			if match != nil {

				// Normalized version has already been stored; delete this record

				log.Printf("Existing normalization found: %s\n", match.number)
				log.Printf("Deleting...\n")

				must(deletePhone(db, p.id))

			} else {

				// Normalized version has not been stored; update this record

				log.Printf("New normalization found: %s\n", normalized)
				log.Printf("Updating...\n")

				p.number = normalized

				must(updatePhone(db, p))

			}

		}

	}

	// doMoreDatabaseStuff()

}

func must(err error) error {
	if err != nil {
		log.Fatal(err)
	}
	return err
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

	re := regexp.MustCompile("\\D") // [^0-9]

	return re.ReplaceAllString(number, "")

}

func readLines(filepath string) ([]string, error) {

	lines := []string{}

	file, err := os.Open(filepath)

	must(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	must(scanner.Err())

	return lines, nil

}

func insertPhone(db *sql.DB, number string) (int, error) {

	result, err := db.Exec("insert into phone_numbers(number) values(?)", number) //.Scan(&id) /* returning id */

	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	return int(id), nil

}

func getPhone(db *sql.DB, id int) (*phone, error) {

	var p *phone

	err := db.QueryRow("select id, number from phone_numbers where id=?", id).Scan(p.id, p.number)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return p, nil

}

func findPhone(db *sql.DB, number string) (*phone, error) {

	var p phone

	err := db.QueryRow("select id, number from phone_numbers where number=?", number).Scan(&p.id, &p.number)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &p, nil

}

func updatePhone(db *sql.DB, p phone) error {

	statement := "update phone_numbers set number=? where id=?"

	_, err := db.Exec(statement, p.number, p.id)

	return err

}

func deletePhone(db *sql.DB, id int) error {

	statement := "delete from phone_numbers where id=?"

	_, err := db.Exec(statement, id)

	return err

}

func getPhones(db *sql.DB) ([]phone, error) {

	var ps []phone

	rows, err := db.Query("select id, number from phone_numbers")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var p phone

		err := rows.Scan(&p.id, &p.number)

		if err != nil {
			return nil, err
		}

		ps = append(ps, p)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ps, nil

}

func resetDatabase() (*sql.DB, error) {

	// Clean up existing DB

	os.Remove(fmt.Sprintf("./%s", databaseName))

	// Open a new connection to our sqlite3 DB process

	db, err := sql.Open("sqlite3", fmt.Sprintf("./%s", databaseName))

	if err != nil {
		return nil, err
	}

	// Ping the database process to make sure we can actually connect

	ctx = context.Background()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil

}

func resetPhoneNumbersTable(db *sql.DB) error {

	// Execute a query to create our phone_numbers table, if necessary; if the
	// table already exists, drop all rows so we can start with an empty table;

	// @TODO(mzalla) Add an index on number column

	q := `
		create table if not exists phone_numbers (id integer not null primary key, number text);
		delete from phone_numbers;
		`

	_, err := db.Exec(q)

	if err != nil {
		return err
	}

	return nil

}

func populatePhoneNumbersTable(db *sql.DB, numbers []string) ([]int, error) {

	ids := make([]int, len(numbers))

	// Insert each number into our phone numbers table

	for i, number := range numbers {

		id, err := insertPhone(db, number)

		if err != nil {
			return nil, err
		}

		ids[i] = id

	}

	return ids, nil

}

func doMoreDatabaseStuff(db *sql.DB) {

	// Read back the data we just wrote

	rows, err := db.Query("select id, number from phone_numbers")

	must(err)

	defer rows.Close()

	for rows.Next() {

		var id int
		var number string

		must(rows.Scan(&id, &number))

		fmt.Println(id, number)

	}

	err = rows.Err()

	must(err)

	// Query a single row from our table

	st, err := db.Prepare("select number from phone_numbers where id=?")

	must(err)

	defer st.Close()

	var number string

	err = st.QueryRow("3").Scan(&number)

	must(err)

	fmt.Println(number)

	// Delete all numbers from our table

	_, err = db.Exec("delete from phone_numbers")

	must(err)

	// Re-populate our table

	_, err = db.Exec("insert into phone_numbers(id, number) values(1, '123-456-7894'), (2, '123-456-7890'), (3, '1234567892')")

	must(err)

}
