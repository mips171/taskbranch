package tasks

import (
	"encoding/json"
	"os"
)

type Condition struct {
	CheckCommand    string            `json:"checkCommand"`
	ExpectedOutcome string            `json:"expectedOutcome"`
	Strategy        ExecutionStrategy `json:"-"`
}

type Task struct {
	DependsOn []string          `json:"dependsOn"`
	ID        string            `json:"id"`
	Command   string            `json:"command"`
	ExecuteIf string            `json:"executeIf"`
	Condition *Condition        `json:"condition"`
	Strategy  ExecutionStrategy `json:"-"`
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
	return &container, nil
}
