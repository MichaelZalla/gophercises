package db

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
)

// DB wraps a sql.DB
type DB struct {
	db *sql.DB
}

// Phone represents a phone number record in our database
type Phone struct {
	ID     int
	Number string
}

func (p Phone) String() string {
	return fmt.Sprintf("<id=%d, number=%s>", p.ID, p.Number)
}

// Open retrieves a handle to the database given by the data source
func Open(driverName, dataSource string) (*DB, error) {

	db, err := sql.Open(driverName, dataSource)

	if err != nil {
		return nil, err
	}

	return &DB{db}, nil

}

// Close releases a handle to the database
func (db *DB) Close() error {

	return db.db.Close()

}

// Seed populates the phone_numbers table with initial, unnormalized data
func (db *DB) Seed() error {

	// Read a list of unnormalized phone numbers from a data file

	unnormalizedPhoneNumbers, err := readLines("./data/numbers.txt")

	must(err)

	// Insert each number into our phone numbers table

	for _, number := range unnormalizedPhoneNumbers {

		err := insertPhone(db.db, number)

		if err != nil {
			return err
		}

	}

	return nil

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

func insertPhone(db *sql.DB, number string) error {

	_, err := db.Exec("insert into phone_numbers(number) values(?)", number)

	if err != nil {
		return err
	}

	return nil

}

// GetPhones returns a Phone for each entry in the phone_numbers table
func (db *DB) GetPhones() ([]Phone, error) {

	var ps []Phone

	rows, err := db.db.Query("select id, number from phone_numbers")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var p Phone

		err := rows.Scan(&p.ID, &p.Number)

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

func getPhone(db *sql.DB, id int) (string, error) {

	var number string

	err := db.QueryRow("select * from phone_numbers where id=?", id).Scan(&id, &number)

	if err != nil {
		return "", nil
	}

	return number, nil

}

// FindPhone locates a Phone record for the given number, if one exists
func (db *DB) FindPhone(number string) (*Phone, error) {

	var p Phone

	row := db.db.QueryRow("select id, number from phone_numbers where number=?", number)
	err := row.Scan(&p.ID, &p.Number)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &p, nil

}

// UpdatePhone writes updates to an existing Phone record
func (db *DB) UpdatePhone(p *Phone) error {

	statement := "update phone_numbers set number=? where id=?"

	_, err := db.db.Exec(statement, p.Number, p.ID)

	return err

}

// DeletePhone removes an existing Phone record
func (db *DB) DeletePhone(id int) error {

	statement := "delete from phone_numbers where id=?"

	_, err := db.db.Exec(statement, id)

	return err

}

// Migrate sets up a new copy of the database at the given data source
func Migrate(driverName, dataSource string) error {

	db, err := sql.Open(driverName, dataSource)

	if err != nil {
		return err
	}

	err = createPhoneNumbersTable(db)

	if err != nil {
		return err
	}

	return db.Close()

}

func createPhoneNumbersTable(db *sql.DB) error {

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

// Reset deletes the database if it exists and creates a new one
func Reset(driverName, dataSource string) error {

	// Open a new connection to our sqlite3 DB process

	db, err := sql.Open("sqlite3", dataSource)

	if err != nil {
		return err
	}

	err = resetDB(db)

	if err != nil {
		return err
	}

	return db.Close()

}

func resetDB(db *sql.DB) error {

	// _, err := db.Exec("drop database if exists " + name)

	// if err != nil {
	// 	return err
	// }

	// ctx = context.Background()

	// err = db.PingContext(ctx)

	// if err != nil {
	// 	return err
	// }

	return createDB(db)

}

func createDB(db *sql.DB) error {

	// _, err := db.Exec("create database " + name)

	// return err

	return nil

}

func must(err error) error {
	if err != nil {
		log.Fatal(err)
	}
	return err
}
