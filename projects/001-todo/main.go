package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/orashus/rtodo/store"
	"github.com/orashus/rtodo/types"
	"github.com/orashus/rtodo/utils"
)

func listHandler(tags []string) {
	todoList, err := store.LoadTodos(FILE_PATH)
	if err != nil {
		PrintDelimiter()
		fmt.Println("Error loading todos ->", err)
		return
	}

	if slices.Contains(tags, TAGS.COMPLETED) {
		res := []types.Todo{}

		for _, todo := range todoList {
			if todo.Completed {
				res = append(res, todo)
			}
		}

		PrintTodos(&res)
		return
	}

	PrintTodos(&todoList)
}

func removeHandler(id string, shouldPrint bool) {
	todoList, err := store.LoadTodos(FILE_PATH)
	if err != nil {
		PrintDelimiter()
		fmt.Println("Error loading todos ->", err)
		return
	}

	hasTodo := false

	todoList = slices.DeleteFunc(todoList, func(todo types.Todo) bool {
		if todo.ID == id {
			hasTodo = true
		}
		return todo.ID == id
	})

	if !hasTodo {
		PrintDelimiter()
		fmt.Println("Todo not found")
		if shouldPrint {
			PrintTodos(&todoList)
		}
		return
	}

	if err := store.SaveTodos(FILE_PATH, &todoList); err != nil {
		PrintDelimiter()
		fmt.Println("Error saving todos:", err)
		return
	}

	PrintDelimiter()
	fmt.Printf("Todo with Id '%s' removed successfully\n", id)
	if shouldPrint {
		PrintTodos(&todoList)
	}
}

func clearHandler(shouldPrint bool) {
	err := os.Remove(FILE_PATH)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) { // or I could do os.IsNotExist(err)
			/*
				This just checks if the error is known to report that a file or directory does not exist.
				os.IsNotExist() does the same thing but is restricted to only errors returned by the os package.
					It is satisfied by [ErrNotExist] as well as some syscall errors.
			*/
			PrintDelimiter()
			fmt.Println("No todos to clear")
			return
		}

		PrintDelimiter()
		fmt.Println("Error clearing todos:", err)
		return
	}

	PrintDelimiter()
	fmt.Println("Todos cleared successfully")
	if shouldPrint {
		PrintTodos(&[]types.Todo{})
	}
}

func addHandler(title string, shouldPrint bool) {
	todoList, err := store.LoadTodos(FILE_PATH)
	if err != nil {
		PrintDelimiter()
		fmt.Println("Error loading todos ->", err)
		return
	}

	id, err := utils.GenerateUniqueId(
		Map(todoList, func(todo types.Todo, _ int) string {
			return todo.ID
		}),
	)

	if err != nil {
		PrintDelimiter()
		fmt.Println("Error generating unique ID:", err)
		return
	}

	newTodo := types.Todo{
		ID:        id,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	todoList = append(todoList, newTodo)
	if err := store.SaveTodos(FILE_PATH, &todoList); err != nil {
		PrintDelimiter()
		fmt.Println("Error saving todos:", err)
		return
	}

	PrintDelimiter()
	fmt.Println("Todo added successfully")
	if shouldPrint {
		PrintTodos(&todoList)
	}
}

const (
	version   = "1.0.0"
	appName   = "rtodo"
	FILE_PATH = "/tmp/r_apps_rtodo.json"
	// FILE_PATH = "test-todos.json"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%s  version %s\n", appName, version)
		fmt.Println("Please provide a command")
		return
	}

	var shouldPrint bool

	command, input, tags := utils.ParseInput(os.Args[1:])

	if slices.Contains(tags, "print") || slices.Contains(tags, "list") {
		shouldPrint = true
	}

	if slices.Contains(tags, "version") {
		fmt.Printf("%s  version %s\n", appName, version)
		return
	}

	switch command {
	case COMMANDS.LIST:
		listHandler(tags)
	case COMMANDS.DELETE:
		fallthrough
	case COMMANDS.RM:
		fallthrough
	case COMMANDS.REMOVE:
		if input == "" {
			fmt.Println("Please provide an ID to remove")
			return
		}

		removeHandler(input, shouldPrint)
	case COMMANDS.CLEAR:
		clearHandler(shouldPrint)
	case COMMANDS.ADD:
		if input == "" {
			fmt.Println("Please provide a title to add")
			return
		}

		addHandler(input, shouldPrint)
	case COMMANDS.UPDATE:
		// WRITE UPDATE HANDLER
		// Update needs 2 inputs,
	default:
		fmt.Println("Invalid command")
		PrintDelimiter()

		if shouldPrint {
			todoList, _ := store.LoadTodos(FILE_PATH)
			PrintTodos(&todoList)
		}
		return
	}
}
