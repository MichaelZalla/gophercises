package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/MichaelZalla/gophercises/07-task/task/todo"
	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task on your TODO list as completed",
	RunE: func(cmd *cobra.Command, args []string) error {

		key, err := strconv.Atoi(args[0])

		if err != nil {
			return fmt.Errorf("failed to parse key from '%s'", args[0])
		}

		markAsCompleted := func(t *todo.Todo) error {

			if t.IsComplete() {
				return fmt.Errorf("task '%s' has already been done", t.Description)
			}

			t.Completed = time.Now().Unix()

			return nil

		}

		todos, err := todo.Update(key, markAsCompleted)

		if err != nil {
			return err
		}

		fmt.Printf("You have completed the \"%s\" task.\n", todos[0].Description)

		return nil

	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
