package main

import (
	"fmt"
	"os"
	"slices"
	"time"
)

var TodoList []Todo

type Todo struct {
	id        string
	title     string
	completed bool
	createdAt time.Time
}

func printTodos(todos []Todo) {
	fmt.Printf("%-30s %-30s %-10s %-50s", "ID", "Created At", "Completed", "Title")
	fmt.Println(Delimiter)
	for _, todo := range todos {
		fmt.Printf("%-30s %-30s %-10v %-50s", todo.id, todo.createdAt.Format(time.RFC3339), todo.completed, todo.title)
	}
}

func finishedHandler() {
	res := []Todo{}

	for _, todo := range TodoList {
		if todo.completed {
			res = append(res, todo)
		}
	}

	printTodos(res)
}

func removeHandler(id string, shouldPrint bool) {
	hasTodo := false

	TodoList = slices.DeleteFunc(TodoList, func(todo Todo) bool {
		if todo.id == id {
			hasTodo = true
		}
		return todo.id == id
	})

	if !hasTodo {
		fmt.Println("\n", Delimiter)
		fmt.Println("Todo not found")
		if shouldPrint {
			printTodos(TodoList)
		}
		return
	}

	fmt.Println("\n", Delimiter)
	fmt.Printf("Todo with Id '%s' removed successfully\n", id)
	if shouldPrint {
		printTodos(TodoList)
	}
}

func addHandler(title string, shouldPrint bool) {
	id, err := GenerateUniqueId()
	if err != nil {
		fmt.Println("\n", Delimiter)
		fmt.Println("Error generating unique ID:", err)
		return
	}

	newTodo := Todo{
		id:        id,
		title:     title,
		completed: false,
		createdAt: time.Now(),
	}

	TodoList = append(TodoList, newTodo)
	if shouldPrint {
		printTodos(TodoList)
	}
}

const version = "1.0.0"

const appName = "rtodo"

func main() {
	var shouldPrint bool

	if len(os.Args) < 2 {
		fmt.Printf("%s  version %s\n", appName, version)
		fmt.Println("Please provide a command")
		return
	}

	command, input, tags := ParseInput(os.Args[1:])

	if slices.Contains(tags, "print") {
		shouldPrint = true
	}

	switch command {
	case "list":
		printTodos(TodoList)
	case "done":
		fallthrough
	case "finish":
		finishedHandler()
	case "delete":
		fallthrough
	case "remove":
		if input == "" {
			fmt.Println("Please provide an ID to remove")
			return
		}

		removeHandler(input, shouldPrint)
	case "add":
		if input == "" {
			fmt.Println("Please provide a title to add")
			return
		}

		addHandler(input, shouldPrint)
	default:
		fmt.Println("Invalid command")
		fmt.Println(Delimiter)

		if shouldPrint {
			printTodos(TodoList)
		}
		return
	}
}

// TOdo, fix persistence problem
