package cmd

import (
	"fmt"

	"github.com/MichaelZalla/gophercises/07-task/task/todo"
	"github.com/spf13/cobra"
)

var completeFilter = func(t todo.Todo) bool {
	return t.IsComplete()
}

var incompleteFilter = func(t todo.Todo) bool {
	return !t.IsComplete()
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your incomplete tasks",
	RunE: func(cmd *cobra.Command, args []string) error {

		incompleteOnly, err := cmd.Flags().GetBool("incomplete")

		if err != nil {
			return err
		}

		completedOnly, err := cmd.Flags().GetBool("complete")

		if err != nil {
			return err
		}

		if err := todo.Init(); err != nil {
			return err
		}

		var todos []todo.Todo

		if incompleteOnly {
			todos, err = todo.GetFiltered(incompleteFilter)
		} else if completedOnly {
			todos, err = todo.GetFiltered(completeFilter)
		} else {
			todos, err = todo.GetAll()
		}

		if err != nil {
			return err
		}

		if len(todos) == 0 {
			if completedOnly {
				fmt.Println("You have not completed any tasks.")
			} else {
				fmt.Println("You have no tasks at this time.")
			}
			return nil
		}

		if completedOnly {
			fmt.Printf("You have completed %d task(s):\n", len(todos))
		} else {
			fmt.Printf("You have %d task(s):\n", len(todos))
		}

		for _, t := range todos {
			fmt.Printf("%d. %s\n", t.Key, t)
		}

		return nil

	},
}

func init() {

	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("complete", "c", false, "List only completed tasks")

	listCmd.Flags().BoolP("incomplete", "i", false, "List only incompleted tasks")

}
