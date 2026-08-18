package f_mgmt

import (
	"fmt"
	"os"
)

func CheckIfFileExists(filePath string) (bool, error) {
	/*
		os.Lstat() is like os.Stat() but it doesn't follow symlinks.
		It returns a [FileInfo] describing the named file.
		If the file is a symbolic link, the returned FileInfo describes the symbolic link.
		os.Lstat() makes no attempt to follow the symlink.
		If there is an error, it will be of type [*PathError].
	*/
	_, err := os.Lstat(filePath)
	if err != nil {
		if os.IsNotExist(err) { // checking if the error is because the file does not exist or something else
			return false, nil
		}

		return false, fmt.Errorf("Error checking if file exists: %w", err)
	}

	return true, nil
}

func DeleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("File does not exist: %w", err)
		}

		return fmt.Errorf("Error deleting file: %w", err)
	}
	return nil
}
