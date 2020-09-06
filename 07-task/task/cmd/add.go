package cmd

import (
	"fmt"
	"strings"

	"github.com/MichaelZalla/gophercises/07-task/task/todo"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to your TODO list",
	RunE: func(cmd *cobra.Command, args []string) error {

		desc := strings.Join(args, " ")

		todos, err := todo.Create(desc)

		if err != nil {
			return err
		}

		fmt.Printf("Added \"%s\" to your task list.\n", todos[0].Description)

		return nil

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
