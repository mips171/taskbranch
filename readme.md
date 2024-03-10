## TaskBranch: Behavior-tree-oriented system administration tool

### 🚧 Under construction 🚧

### Documentation

TaskBranch simplifies task management, especially for workflows requiring tasks to be executed in a specific order based on dependencies. By leveraging [behavior trees](https://en.wikipedia.org/wiki/Behavior_tree_(artificial_intelligence,_robotics_and_control)), TaskBranch ensures tasks are executed in accordance with their defined dependencies, and can react to conditions, making it ideal for automating complex workflows. These workflows can include system administration, but can also be used to automate workflows or tasks on any computer system. Thanks to cross-compilation support by Go, you can run TaskBranch natively on your architecture.

### Configuration

TaskBranch operates based on a JSON configuration file that defines tasks, their respective commands, and dependencies. Each task is identified by a unique id, contains the command to be executed, and lists its dependsOn dependencies (other task IDs).

JSON Configuration Structure
```json
{
  "tasks": [
    {
      "id": "unique_task_id",
      "command": "shell_command_to_execute",
      "dependsOn": ["id_of_another_task", "..."]
    }
    // Additional tasks
  ]
}
```

### Simple Use Case
In a simple use case, you might have a two-step process where the second task should only run after the first one completes successfully. For instance, upgrading a system followed by a reboot:

```json
{
  "tasks": [
    {
      "id": "task1",
      "command": "echo 'Performing task 1'",
      "dependsOn": []
    },
    {
      "id": "task2",
      "command": "echo 'Performing task 2'",
      "dependsOn": ["task1"]
    }
  ],
  "dryRun": false
}
``` 

This configuration specifies that the reboot task will only be executed after the upgrade task has completed successfully.

### More Complex Use Case

In a more complex scenario, you might have multiple dependencies and tasks that need to be executed in a specific order:

``` json
{
  "tasks": [
    {
      "id": "download",
      "command": "echo downloading software >> test.out",
      "dependsOn": []
    },
    {
      "id": "install",
      "command": "echo installing software >> test.out",
      "dependsOn": ["download"]
    },
    {
      "id": "configure",
      "command": "echo configuring software >> test.out",
      "dependsOn": ["install"]
    },
    {
      "id": "test",
      "command": "echo testing software >> test.out",
      "dependsOn": ["configure"]
    },
    {
      "id": "deploy",
      "command": "echo deploying software >> test.out",
      "dependsOn": ["test"]
    }
  ]
}
```

This configuration ensures a software deployment process follows a specific order: downloading, installing, configuring, testing, and finally deploying the software. Each step depends on the successful completion of the previous one.

## Running TaskBranch
To run TaskBranch, use the following command, specifying the path to your configuration file as needed:

```sh
taskbranch -file path/to/your/tasks.json
```
If the `-file` flag is omitted, TaskBranch will default to using `taskbranch.json` in the current directory.

```sh
taskbranch -print
```
If the `-print` flag is given, TaskBranch will print a tree structure of the tasks.

```
Tasks:
└─── ID: checkConfigExists
        Condition:
            Check Command: test -f /path/to/config.file
            Expected Outcome: 
└─── ID: upgrade
        Command: echo upgrading >> test.out
        Depends On: checkConfigExists
    └─── ID: checkConfigExists
            Condition:
                Check Command: test -f /path/to/config.file
                Expected Outcome: 
└─── ID: reboot
        Command: echo rebooting >> test.out
        Depends On: upgrade
    └─── ID: upgrade
            Command: echo upgrading >> test.out
            Depends On: checkConfigExists
        └─── ID: checkConfigExists
                Condition:
                    Check Command: test -f /path/to/config.file
                    Expected Outcome: 
```
