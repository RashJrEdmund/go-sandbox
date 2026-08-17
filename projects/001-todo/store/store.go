package store

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/orashus/rtodo/types"
)

// const dummyData = `
// [
// 	{
// 		"id": "12345",
// 		"title": "Buy groceries",
// 		"completed": false,
// 		"created_at": "2021-01-01T00:00:00Z"
// 	},
// 	{
// 		"id": "67890",
// 		"title": "Buy a new car",
// 		"completed": true,
// 		"created_at": "2021-01-02T00:00:00Z"
// 	}
// ]`

func LoadTodos(filePath string) ([]types.Todo, error) {
	// if err := json.Unmarshal([]byte(dummyData), &TodoList); err != nil {
	// 	return nil, fmt.Errorf("error unmarshalling todos: %w", err)
	// }

	file, err := os.Open(filePath)

	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	defer file.Close()

	var TodoList []types.Todo

	decoder := json.NewDecoder(file) // just like reading a request body

	if err := decoder.Decode(&TodoList); err != nil {
		return nil, fmt.Errorf("error decoding todos: %w", err)
	}

	return TodoList, nil
}

// func saveTodos(path string, todos []Todo) error {
// 	//
// }
