package main

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/orashus/rtodo/store"
	"github.com/orashus/rtodo/types"
	"github.com/orashus/rtodo/utils"
)

func listHandler(todoList *[]types.Todo, tags []string) {
	if slices.Contains(tags, TAGS.COMPLETED) {
		res := []types.Todo{}

		for _, todo := range *todoList {
			if todo.Completed {
				res = append(res, todo)
			}
		}

		PrintTodos(&res)
		return
	}

	PrintTodos(todoList)
}

func removeHandler(todoList *[]types.Todo, id string, shouldPrint bool) {
	hasTodo := false

	*todoList = slices.DeleteFunc(*todoList, func(todo types.Todo) bool {
		if todo.ID == id {
			hasTodo = true
		}
		return todo.ID == id
	})

	if !hasTodo {
		fmt.Println("\n", utils.Delimiter)
		fmt.Println("Todo not found")
		if shouldPrint {
			PrintTodos(todoList)
		}
		return
	}

	fmt.Println("\n", utils.Delimiter)
	fmt.Printf("Todo with Id '%s' removed successfully\n", id)
	if shouldPrint {
		PrintTodos(todoList)
	}
}

func addHandler(todoList *[]types.Todo, title string, shouldPrint bool) {
	id, err := utils.GenerateUniqueId()
	if err != nil {
		fmt.Println("\n", utils.Delimiter)
		fmt.Println("Error generating unique ID:", err)
		return
	}

	newTodo := types.Todo{
		ID:        id,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	*todoList = append(*todoList, newTodo)
	if shouldPrint {
		PrintTodos(todoList)
	}
}

const (
	version  = "1.0.0"
	appName  = "rtodo"
	filePath = "todos.json"
)

func main() {
	var TodoList, err = store.LoadTodos(filePath)

	if err != nil {
		fmt.Println("Error loading todos:", err)
		return
	}

	if len(os.Args) < 2 {
		fmt.Printf("%s  version %s\n", appName, version)
		fmt.Println("Please provide a command")
		return
	}

	var shouldPrint bool

	command, input, tags := utils.ParseInput(os.Args[1:])

	if slices.Contains(tags, "print") {
		shouldPrint = true
	}

	if slices.Contains(tags, "version") {
		fmt.Printf("%s  version %s\n", appName, version)
		return
	}

	switch command {
	case COMMANDS.LIST:
		listHandler(&TodoList, tags)
	case COMMANDS.DELETE:
		fallthrough
	case COMMANDS.RM:
		fallthrough
	case COMMANDS.REMOVE:
		if input == "" {
			fmt.Println("Please provide an ID to remove")
			return
		}

		removeHandler(&TodoList, input, shouldPrint)
	case COMMANDS.ADD:
		if input == "" {
			fmt.Println("Please provide a title to add")
			return
		}

		addHandler(&TodoList, input, shouldPrint)
	case COMMANDS.UPDATE:
		// WRITE UPDATE HANDLER
		// Update needs 2 inputs,
	default:
		fmt.Println("Invalid command")
		fmt.Println(utils.Delimiter)

		if shouldPrint {
			PrintTodos(&TodoList)
		}
		return
	}
}

// TOdo, fix persistence problem
