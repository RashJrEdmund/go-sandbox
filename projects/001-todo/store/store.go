package store

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/orashus/rtodo/types"
)

func LoadTodos(filePath string) ([]types.Todo, error) {
	var todoList []types.Todo

	file, err := os.Open(filePath)
	if err != nil {
		return todoList, fmt.Errorf("Error opening file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file) // just like reading a request body
	if err := decoder.Decode(&todoList); err != nil {
		return todoList, fmt.Errorf("Error decoding todos: %w", err)
	}

	return todoList, nil
}

func SaveTodos(filePath string, todos *[]types.Todo) error {
	todoJsonBytes, err := json.Marshal(*todos)
	if err != nil {
		return fmt.Errorf("Error marshalling todos: %w", err)
	}

	// os.WriteFile will create file if not exists or truncate it if it exists
	// truncating the file means deleting all the content and writing the new content
	er := os.WriteFile(filePath, todoJsonBytes, 0644)
	if er != nil {
		return fmt.Errorf("Error writing file: %w", er)
	}

	return nil
}
