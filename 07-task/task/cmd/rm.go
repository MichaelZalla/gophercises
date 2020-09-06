package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/MichaelZalla/gophercises/07-task/task/todo"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a new task to your TODO list",
	RunE: func(cmd *cobra.Command, args []string) error {

		desc := strings.Join(args, " ")

		db, err := newWriter()

		if err != nil {
			log.Fatal(err)
		}

		err = db.Update(func(tx *bolt.Tx) error {

			var todo todo.Todo

			b := tx.Bucket(tasksBucket)

			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {

				err := json.Unmarshal(v, &todo)

				if err != nil {
					return err
				}

				if todo.Description == desc {
					return b.Delete(k)
				}

			}

			return fmt.Errorf("task '%s' does not exist", desc)

		})

		if err != nil {
			return err
		}

		fmt.Printf("You have removed the \"%s\" task.\n", desc)

		return nil

	},
}

func init() {

	rootCmd.AddCommand(rmCmd)

}
