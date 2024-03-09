package tasks

import (
	"encoding/json"
	"os"
)

type Condition struct {
	CheckCommand    string            `json:"checkCommand"`
	ExpectedOutcome string            `json:"expectedOutcome"`
	StrategyID      string            `json:"strategy"` // Adjusted to match your JSON key for strategy identification
	Strategy        ExecutionStrategy `json:"-"`
}

type Task struct {
	DependsOn  []string          `json:"dependsOn"`
	ID         string            `json:"id"`
	Command    string            `json:"command"`
	ExecuteIf  string            `json:"executeIf"`
	Condition  *Condition        `json:"condition"`
	StrategyID string            `json:"strategy"` // Adjusted to match your JSON key for strategy identification
	Strategy   ExecutionStrategy `json:"-"`
}

type TasksContainer struct {
	Tasks  []Task `json:"tasks"`
	DryRun bool   `json:"dryRun"`
}

func LoadTasks(filename string) (*TasksContainer, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var container TasksContainer
	err = json.Unmarshal(data, &container)
	if err != nil {
		return nil, err
	}

	// Assign strategies to tasks and their conditions based on StrategyID
	for i, task := range container.Tasks {
		// Assign strategy to the task
		if task.Strategy == nil {
			container.Tasks[i].Strategy = StrategyFactory(task.StrategyID)
		}
		// Assign strategy to the condition, if it exists
		if task.Condition != nil && task.Condition.Strategy == nil {
			container.Tasks[i].Condition.Strategy = StrategyFactory(task.Condition.StrategyID)
		}
	}

	return &container, nil
}

func StrategyFactory(strategyID string) ExecutionStrategy {
	switch strategyID {
	case "OSCommandStrategy":
		return &OSCommandStrategy{} // Assuming this is a concrete type implementing ExecutionStrategy.
	default:
		return nil // Return a sensible default or nil if unrecognized.
	}
}
