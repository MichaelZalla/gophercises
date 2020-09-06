package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MichaelZalla/gophercises/07-task/task/todo"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

// completedCmd represents the completed command
var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "List all of your completed tasks",
	RunE: func(cmd *cobra.Command, args []string) error {

		var todos []todo.Todo

		db, err := newReader()

		if err != nil {
			return err
		}

		err = db.View(func(tx *bolt.Tx) error {

			results := []todo.Todo{}

			c := tx.Bucket(tasksBucket).Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {

				var todo todo.Todo

				err := json.Unmarshal(v, &todo)

				if err != nil {
					return err
				}

				t := time.Unix(todo.Completed, 0)

				if err != nil {
					log.Print(err)
					continue
				}

				// log.Printf("Todo desc: %s\n", todo.Description)
				// log.Printf("Todo completed: %s\n", todo.Completed)
				// log.Printf("Seconds since completion: %f\n", time.Since(t).Seconds())

				since := time.Since(t)

				if since.Hours() <= 24 {
					results = append(results, todo)
				}

			}

			todos = make([]todo.Todo, len(results))

			copy(todos, results)

			return nil

		})

		if err != nil {
			return err
		}

		if len(todos) == 0 {
			fmt.Println("You haven't compeleted any tasks today.")
			return nil
		}

		fmt.Println("You have finished the following tasks today:")

		for i, t := range todos {
			fmt.Printf("%d. %s\n", i, t)
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(completedCmd)
}
