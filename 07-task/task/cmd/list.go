package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/MichaelZalla/gophercises/07-task/task/todo"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your incomplete tasks",
	RunE: func(cmd *cobra.Command, args []string) error {

		var todos []todo.Todo

		incompleteOnly, err := cmd.Flags().GetBool("incomplete")

		if err != nil {
			return err
		}

		completedOnly, err := cmd.Flags().GetBool("complete")

		if err != nil {
			return err
		}

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

				// (!A && !B) || (A && true) || (B && false)

				if (!completedOnly && !incompleteOnly) ||
					(completedOnly && todo.IsComplete()) ||
					(incompleteOnly && !todo.IsComplete()) {
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
			fmt.Println("You have no tasks at this time.")
			return nil
		}

		fmt.Println("You have the following tasks:")

		for i, t := range todos {
			fmt.Printf("%d. %s\n", i, t)
		}

		return nil

	},
}

func init() {

	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("complete", "c", false, "List only completed tasks")
	listCmd.Flags().BoolP("incomplete", "i", false, "List only incompleted tasks")

}
