package tasks

import (
	"fmt"

	"github.com/mips171/beetea"
)

type CustomLogicStrategy struct{}

func (c CustomLogicStrategy) Execute(task Task, dryRun bool) beetea.Status {

	fmt.Printf("Performing custom logic for task: %s\n", task.ID)
	return beetea.Success
}
