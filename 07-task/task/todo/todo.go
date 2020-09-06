package todo

import "fmt"

// Todo represents a task and its associated completion status
type Todo struct {
	Description string `json:"desc"`
	Completed   int64  `json:"completed"`
}

// IsComplete indicates whether or not a Todo must still be accomplished
func (t Todo) IsComplete() bool {
	return t.Completed != 0
}

func (t Todo) String() string {

	var pre string

	if t.IsComplete() {
		pre = "[*]"
	} else {
		pre = "[ ]"
	}

	return fmt.Sprintf("%s %s", pre, t.Description)

}
