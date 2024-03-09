package main

import (
	"flag"
	"fmt"
	"os"

	"taskbranch/tasks"
	"taskbranch/tree"

	"github.com/mips171/beetea"
)

var printFlag bool

func init() {
	flag.BoolVar(&printFlag, "print", false, "Print the task graph and exit")
}

func main() {
	taskFilePath := flag.String("task-file", "taskbranch.json", "Path to the task file")
	flag.Parse()

	container, err := tasks.LoadTasks(*taskFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s not found. Please make sure the file exists.\n", *taskFilePath)
		} else {
			fmt.Printf("Failed to load tasks: %v\n", err)
		}
		return
	}

	taskTree := tree.BuildBehaviorTree(container.Tasks, container.DryRun)

	if printFlag {
		tree.PrintTasks(container)
		return
	}

	status := taskTree.Tick()
	if status == beetea.Success {
		fmt.Println("All tasks executed successfully.")
	} else {
		fmt.Println("Some tasks failed.")
	}
}
