package tasks

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mips171/beetea"
)

// executes tasks and checks conditions using OS commands
type OSCommandStrategy struct{}

func (s *OSCommandStrategy) Execute(task Task, dryRun bool) beetea.Status {
	if dryRun {
		fmt.Printf("[DRY RUN] Would execute: %s\n", task.Command)
		return beetea.Success
	}

	fmt.Println("Executing:", task.Command)
	cmd := exec.Command("sh", "-c", task.Command)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error executing task %s: %v\n", task.ID, err)
		return beetea.Failure
	}
	return beetea.Success
}

func (s *OSCommandStrategy) CheckCondition(condition Condition, dryRun bool) bool {
    if dryRun {
        fmt.Println("[DRY RUN] Checking condition:", condition.CheckCommand)
        return true
    }

    fmt.Println("Checking condition:", condition.CheckCommand)
    cmd := exec.Command("sh", "-c", condition.CheckCommand)
    output, err := cmd.CombinedOutput()
    trimmedOutput := strings.TrimSpace(string(output))

    if err != nil {
        fmt.Printf("Condition check failed: %s, Output: %s\n", err, trimmedOutput)
        return false
    }

    return trimmedOutput == condition.ExpectedOutcome
}
