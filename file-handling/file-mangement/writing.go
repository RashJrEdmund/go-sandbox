package f_mgmt

import (
	"encoding/json"
	"fmt"
	"os"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func CreateFile(filePath string) error {
	/*
		os.Create() creates a file if it doesn't exist and opens it for writing.
		If the file already exists, it is truncated.
	*/
	file, err := os.Create(filePath)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("File already exists: %w", err)
		}

		return fmt.Errorf("Error creating file: %w", err)
	}
	defer file.Close()

	return nil
}

func WriteToFile(users []User, filePath string) error {
	userJsonBytes, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("Error marshalling users: %w", err)
	}

	/*
		os.WriteFile will create file if not exists or truncate it if it exists
		truncating the file means deleting all the content and writing the new content

		Permissions:
			0644 sets the file mode to -rw-r--r-- (owner can read/write, group and others can read)
			Other common file modes:
			0666 -rw-rw-rw- (everyone can read/write)
			0600 -rw------- (only owner can read/write)
			0755 -rwxr-xr-x (owner can read/write/execute, group and others can read/execute)
	*/
	er := os.WriteFile(filePath, userJsonBytes, 0644)
	if er != nil {
		return fmt.Errorf("Error writing file: %w", err)
	}

	fmt.Println("Users written to", filePath)
	return nil
}
