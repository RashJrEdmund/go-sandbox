package store

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/orashus/rtodo/types"
)

func LoadTodos(filePath string) ([]types.Todo, error) {
	byteData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []types.Todo{}, nil
		}
		return nil, fmt.Errorf("Error opening file: %w", err)
	}

	var todoList []types.Todo

	if err := json.Unmarshal(byteData, &todoList); err != nil {
		return nil, fmt.Errorf("Error unmarshalling todos: %w", err)
	}

	return todoList, nil
}

func SaveTodos(filePath string, todos *[]types.Todo) error {
	todoJsonBytes, err := json.MarshalIndent(*todos, "", "  ") // 2 spaces for indentation in json file
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
