package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

var tasksBucket = []byte("tasks")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "task is a CLI for managing your TODOs.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newReader() (*bolt.DB, error) {

	db, err := bolt.Open("tasks.db", 0666, &bolt.Options{ReadOnly: true})

	if err != nil {
		return nil, err
	}

	return db, nil

}

func newWriter() (*bolt.DB, error) {

	db, err := bolt.Open("tasks.db", 0666, nil)

	if err != nil {
		return nil, err
	}

	return db, nil

}

func init() {

	// Initialize tha Bolt database to store our tasks

	db, err := newWriter()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Guarantee that the tasks bucket exists

	err = db.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists(tasksBucket)

		if err != nil {
			return fmt.Errorf("called CreateBucketIfNotExists() with name '%s'", tasksBucket)
		}

		return nil

	})

	if err != nil {
		log.Fatal(err)
	}

}
