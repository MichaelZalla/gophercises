package cmd

import (
	"fmt"
	"strconv"

	"github.com/MichaelZalla/gophercises/07-task/task/todo"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a new task to your TODO list",
	RunE: func(cmd *cobra.Command, args []string) error {

		key, err := strconv.Atoi(args[0])

		if err != nil {
			return fmt.Errorf("failed to parse key from '%s'", args[0])
		}

		err = todo.Delete(key)

		if err != nil {
			return err
		}

		fmt.Printf("You have removed the task with key %d.\n", key)

		return nil

	},
}

func init() {

	rootCmd.AddCommand(rmCmd)

}
