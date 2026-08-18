package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func loadUsers(filePath string) []User {
	var res []User

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err) // will throw an error if the file does not exist
		return res
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&res); err != nil {
		fmt.Println("Error decoding users:", err)
		return res
	}

	return res
}

func saveUsers(users []User, filePath string) {
	userJsonBytes, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Error marshalling users:", err)
		return
	}

	// 0644 sets the file mode to -rw-r--r-- (owner can read/write, group and others can read)
	// Other common file modes:
	// 0666 -rw-rw-rw- (everyone can read/write)
	// 0600 -rw------- (only owner can read/write)
	// 0755 -rwxr-xr-x (owner can read/write/execute, group and others can read/execute)
	er := os.WriteFile(filePath, userJsonBytes, 0644)
	if er != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Users written to", filePath)
}

var filePath = "/tmp/r_apps_file-writing.json" // writing to a temporary file

func main() {
	users := loadUsers(filePath)

	id := 0

	if len(users) > 0 {
		id = users[len(users)-1].ID + 1
	}

	newUser := User{
		ID:   id,
		Name: "John Doe",
	}

	users = append(users, newUser)

	saveUsers(users, filePath)

	fmt.Println(users)
}
