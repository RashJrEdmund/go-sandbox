package main

type tags struct {
	PRINT     string
	COMPLETED string
	HELP      string
}

type commands struct {
	LIST string
	ADD  string
	// delete and remove have the same functionality
	DELETE string
	RM     string
	REMOVE string
	//
	UPDATE string
}

var TAGS = tags{
	PRINT:     "print",
	COMPLETED: "completed",
	HELP:      "help",
}

var COMMANDS = commands{
	LIST:   "list",
	ADD:    "add",
	DELETE: "delete",
	REMOVE: "remove",
	RM:     "rm",
	UPDATE: "update",
}
