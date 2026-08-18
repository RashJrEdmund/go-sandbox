package f_mgmt

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Todo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

/*
This version with os.Open() allows for more granular control of the file
as we have access to the file object and can manipulate it as needed.
*/
func ReadWithOpen(filePath string) ([]Todo, error) {
	/*
		os.Open() opens a file for reading.
		If the file does not exist, it returns an error.
	*/
	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("File does not exist: %w", err)
		}

		return nil, fmt.Errorf("Error opening file: %w", err)
	}
	defer file.Close()

	var todoList []Todo

	decoder := json.NewDecoder(file) // just like reading a request body
	if err := decoder.Decode(&todoList); err != nil {
		return nil, fmt.Errorf("Error decoding todos: %w", err)
	}

	return todoList, nil
}

/*
This version with os.ReadFile() is simpler and more concise.
*/
func ReadWithReadFile(filePath string) ([]Todo, error) { // ✅ i'd use this one for simpler cases
	/*
		os.ReadFile() reads the entire file into memory as a byte slice.
		If the file does not exist, it returns an error.
	*/
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// same as errors.Is(err, os.ErrNotExist)
			// but os.IsNotExist() is restricted to only errors returned by the os package.
			return nil, fmt.Errorf("File does not exist: %w", err)
		}
		return nil, fmt.Errorf("Error reading file: %w", err)
	}

	var todoList []Todo

	if err := json.Unmarshal(data, &todoList); err != nil {
		return nil, fmt.Errorf("Error unmarshalling todos: %w", err)
	}

	return todoList, nil
}
