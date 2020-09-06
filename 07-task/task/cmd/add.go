package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/MichaelZalla/gophercises/07-task/task/convert"
	"github.com/MichaelZalla/gophercises/07-task/task/todo"
	bolt "go.etcd.io/bbolt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to your TODO list",
	RunE: func(cmd *cobra.Command, args []string) error {

		todo := todo.Todo{
			Description: strings.Join(args, " "),
			Completed:   0,
		}

		todoBytes, err := json.Marshal(todo)

		if err != nil {
			return err
		}

		db, err := newWriter()

		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()

		err = db.Update(func(tx *bolt.Tx) error {

			b := tx.Bucket(tasksBucket)

			if b == nil {
				return fmt.Errorf("failed to find bucket '%s'", tasksBucket)
			}

			todoID, err := b.NextSequence()

			if err != nil {
				return err
			}

			return b.Put(convert.IntToBytes(int(todoID)), todoBytes)

		})

		if err != nil {
			return err
		}

		fmt.Printf("Added \"%s\" to your task list.\n", todo.Description)

		return nil

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
