package tasks

import (
	"github.com/mips171/beetea"
)

type ExecutionStrategy interface {
	Execute(task Task, dryRun bool) beetea.Status
	CheckCondition(condition Condition, dryRun bool) bool
}
