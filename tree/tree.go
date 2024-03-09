package tree

import (
	"fmt"
	"strings"
	"taskbranch/tasks"

	"github.com/mips171/beetea"
)

// Dynamically build the behavior tree based on task dependencies.
// buildBehaviorTree builds a behavior tree based on the given tasks and dryRun flag.

// It creates nodes for each task, including condition checks if present.
// If a task has dependencies, it creates sequences of nodes that need to be executed in order.
// If no sequences were created (no dependencies), it returns a selector of all actions.
// If there is only one sequence, it returns that sequence. Otherwise, it returns a selector of sequences.

// The behavior tree is built using the beetea package: https://github.com/mips171/beetea
//
// Parameters:
//   - tasks: A slice of Task objects representing the tasks in the behavior tree.
//   - dryRun: A boolean flag indicating whether the behavior tree should be executed in dry run mode.
//
// Returns:
//   - A beetea.Node representing the root node of the behavior tree.
func BuildBehaviorTree(tasks []tasks.Task, dryRun bool) beetea.Node {
	taskNodes := make(map[string]beetea.Node)

	for _, task := range tasks {
		if task.Condition != nil {
			conditionNode := createConditionNode(task.Condition, dryRun)
			actionNode := createTaskNode(task, dryRun)
			seq := beetea.NewSequence(conditionNode, actionNode)
			taskNodes[task.ID] = seq
		} else {
			taskNodes[task.ID] = createTaskNode(task, dryRun)
		}
	}
	var sequences []*beetea.Sequence

	for _, task := range tasks {
		if len(task.DependsOn) > 0 {
			var seqNodes []beetea.Node
			for _, depID := range task.DependsOn {
				depNode, exists := taskNodes[depID]
				if exists {
					seqNodes = append(seqNodes, depNode)
				}
			}
			seqNodes = append(seqNodes, taskNodes[task.ID])
			sequence := beetea.NewSequence(seqNodes...)
			sequences = append(sequences, sequence)
		}
	}

	// If no sequences were created (no dependencies), just return a selector of all actions
	if len(sequences) == 0 {
		var actionNodes []beetea.Node
		for _, node := range taskNodes {
			actionNodes = append(actionNodes, node)
		}
		return beetea.NewSelector(actionNodes...)
	}

	// For simplicity, return the first sequence or a selector of sequences if multiple
	if len(sequences) == 1 {
		return sequences[0]
	} else {
		var seqNodes []beetea.Node
		for _, seq := range sequences {
			seqNodes = append(seqNodes, seq)
		}
		return beetea.NewSelector(seqNodes...)
	}
}

func FindTaskByID(id string, tasks []tasks.Task) *tasks.Task {
	for _, task := range tasks {
		if task.ID == id {
			return &task
		}
	}
	return nil
}

func PrintTasks(container *tasks.TasksContainer) {
	fmt.Println("Tasks:")
	for _, task := range container.Tasks {
		printTask(task, "", true, container.Tasks)
	}
}

func printTask(task tasks.Task, prefix string, isLast bool, tasks []tasks.Task) {
	fmt.Print(prefix)
	if isLast {
		fmt.Print("└─── ")
		prefix += "    "
	} else {
		fmt.Print("├─── ")
		prefix += "│   "
	}
	fmt.Printf("ID: %s\n", task.ID)
	if task.Command != "" {
		fmt.Printf("%s    Command: %s\n", prefix, task.Command)
	}
	if len(task.DependsOn) > 0 {
		fmt.Printf("%s    Depends On: %s\n", prefix, strings.Join(task.DependsOn, ", "))
	}
	if task.Condition != nil {
		fmt.Println(prefix + "    Condition:")
		fmt.Printf("%s        Check Command: %s\n", prefix, task.Condition.CheckCommand)
		fmt.Printf("%s        Expected Outcome: %s\n", prefix, task.Condition.ExpectedOutcome)
	}
	if len(task.DependsOn) > 0 {
		lastIdx := len(task.DependsOn) - 1
		for idx, dep := range task.DependsOn {
			isLast := idx == lastIdx
			depTask := FindTaskByID(dep, tasks)
			if depTask != nil {
				printTask(*depTask, prefix, isLast, tasks)
			}
		}
	}
}

func createTaskNode(task tasks.Task, dryRun bool) *beetea.ActionNode {
	// Ensure there is a default strategy if none is specified
	if task.Strategy == nil {
		task.Strategy = &tasks.OSCommandStrategy{}
	}

	return beetea.NewAction(func() beetea.Status {
		return task.Strategy.Execute(task, dryRun)
	})
}

func createConditionNode(condition *tasks.Condition, dryRun bool) *beetea.ConditionNode {
	// Ensure there is a default strategy if none is specified
	if condition.Strategy == nil {
		condition.Strategy = &tasks.OSCommandStrategy{}
	}

	return beetea.NewCondition(func() bool {
		return condition.Strategy.CheckCondition(*condition, dryRun)
	})
}
