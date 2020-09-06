package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/MichaelZalla/gophercises/07-task/task/todo"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task on your TODO list as completed",
	RunE: func(cmd *cobra.Command, args []string) error {

		target, err := strconv.ParseInt(args[0], 10, 0)

		if err != nil {
			return fmt.Errorf("failed to parse task index from '%s'", args[0])
		}

		var result todo.Todo

		db, err := newWriter()

		if err != nil {
			log.Fatal(err)
		}

		err = db.Update(func(tx *bolt.Tx) error {

			var todo todo.Todo

			b := tx.Bucket(tasksBucket)

			c := b.Cursor()

			k, v := c.First()

			for index := int64(0); index < target; index++ {
				k, v = c.Next()
			}

			if k == nil {
				return fmt.Errorf("no task at index '%d'", target)
			}

			err := json.Unmarshal(v, &todo)

			if err != nil {
				return err
			}

			result = todo

			if todo.IsComplete() {
				return fmt.Errorf("task '%s' has already been done", todo.Description)
			}

			todo.Completed = time.Now().Unix()

			v, err = json.Marshal(todo)

			if err != nil {
				return err
			}

			return b.Put(k, v)

		})

		if err != nil {
			return err
		}

		fmt.Printf("You have completed the \"%s\" task.\n", result.Description)

		return nil

	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
