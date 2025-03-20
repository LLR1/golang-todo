package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const tasksFile = "tasks.json"

var tasks []string

func loadTasks() error {
	// Use os.ReadFile if you want to read entire file at once:
	data, err := os.ReadFile(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			tasks = []string{}
			return nil
		}
		return err
	}
	return json.Unmarshal(data, &tasks)
}

func saveTasks() error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tasksFile, data, 0644)
}

func main() {
	if err := loadTasks(); err != nil {
		fmt.Fprintln(os.Stderr, "Error loading tasks:", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Please specify a task to add.")
			return
		}
		task := os.Args[2]
		tasks = append(tasks, task)
		if err := saveTasks(); err != nil {
			fmt.Println("Error saving tasks:", err)
			return
		}
		fmt.Println("Task added:", task)
	case "list":
		// If no tasks, inform the user
		if len(tasks) == 0 {
			fmt.Println("Task list is empty.")
			return
		}
		// Print all tasks
		fmt.Println("Task list:")
		for i, t := range tasks {
			fmt.Printf("%d. %s\n", i+1, t)
		}
	case "delete":
		// Check if a task number is provided
		if len(os.Args) < 3 {
			fmt.Println("Please specify the task number to delete.")
			return
		}
		// Convert the task number from string to integer
		index, err := strconv.Atoi(os.Args[2])
		if err != nil || index < 1 || index > len(tasks) {
			fmt.Println("Invalid task number.")
			return
		}
		// Remove the task from the slice
		removed := tasks[index-1]
		tasks = append(tasks[:index-1], tasks[index:]...)
		if err := saveTasks(); err != nil {
			fmt.Println("Error saving tasks:", err)
			return
		}
		fmt.Println("Task deleted:", removed)
	default:
		printUsage()
	}
}

// printUsage displays the usage instructions for the CLI tool
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  add <task>    - add a new task")
	fmt.Println("  list          - show all tasks")
	fmt.Println("  delete <num>  - delete a task by its number")
}
