package cmd

import (
	"fmt"
	"time"

	"github.com/MichaelZalla/gophercises/07-task/task/todo"
	"github.com/spf13/cobra"
)

var recentlyCompleteFilter = func(t todo.Todo) bool {

	since := time.Since(time.Unix(t.Completed, 0))

	return (since.Hours() <= 24)

}

// completedCmd represents the completed command
var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "List all of your completed tasks",
	RunE: func(cmd *cobra.Command, args []string) error {

		todos, err := todo.GetFiltered(recentlyCompleteFilter)

		if err != nil {
			return err
		}

		if len(todos) == 0 {
			fmt.Println("You haven't compeleted any tasks today.")
			return nil
		}

		fmt.Printf("You have finished %d task(s) today:\n", len(todos))

		for _, t := range todos {
			fmt.Printf("%d. %s\n", t.Key, t)
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(completedCmd)
}
