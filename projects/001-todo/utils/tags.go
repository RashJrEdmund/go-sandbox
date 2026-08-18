package utils

import (
	"fmt"
	"slices"
)

func removeTag(args *[]string, tag string) {
	_tag, __tag := fmt.Sprintf("-%s", tag), fmt.Sprintf("--%s", tag)

	if slices.Contains(*args, _tag) {
		*args = slices.Delete(*args, slices.Index(*args, _tag), slices.Index(*args, _tag)+1)
	}

	if slices.Contains(*args, __tag) {
		*args = slices.Delete(*args, slices.Index(*args, __tag), slices.Index(*args, __tag)+1)
	}
}

func containsTag(args *[]string, tag string) bool {
	_tag, __tag := fmt.Sprintf("-%s", tag), fmt.Sprintf("--%s", tag)

	if slices.Contains(*args, _tag) || slices.Contains(*args, __tag) {
		return true
	}

	return false
}

func ParseInput(args []string) (command string, input string, tags []string) { // input could be a title or an ID
	if len(args) < 1 {
		return command, input, tags // zero values
	}

	// fmt.Println(args, len(args))

	allowedTags := []struct {
		short string
		long  string
	}{
		{short: "p", long: "print"},
		{short: "l", long: "list"}, // list and print flags do the same thing
		{short: "c", long: "completed"},
		{short: "h", long: "help"},
		{short: "v", long: "version"},
	}

	for _, tag := range allowedTags {
		if containsTag(&args, tag.short) || containsTag(&args, tag.long) {
			tags = append(tags, tag.long)
			removeTag(&args, tag.short)
			removeTag(&args, tag.long)
		}
	}

	// if containsTag(&args, "p") || containsTag(&args, "print") {
	// 	tags = append(tags, "print")
	// 	removeTag(&args, "p")
	// 	removeTag(&args, "print")
	// }

	// if containsTag(&args, "c") || containsTag(&args, "completed") {
	// 	tags = append(tags, "completed")
	// 	removeTag(&args, "c")
	// 	removeTag(&args, "completed")
	// }

	// if containsTag(&args, "h") || containsTag(&args, "help") {
	// 	tags = append(tags, "help")
	// 	removeTag(&args, "h")
	// 	removeTag(&args, "help")

	// 	return command, input, tags
	// }

	if len(args) >= 1 {
		command = args[0]
	}

	if len(args) >= 2 {
		input = args[1]
		// for update, input should be `input = []string{args[1], args[2]}`
	}

	return command, input, tags
}
