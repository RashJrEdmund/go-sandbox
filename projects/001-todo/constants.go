package main

type tags struct {
	PRINT     string
	LIST      string
	COMPLETED string
	HELP      string
	VERSION   string
}

type commands struct {
	LIST string
	//
	ADD string
	// delete and remove have the same functionality
	DELETE string
	RM     string
	REMOVE string
	//
	RMC string // remove completed
	//
	CLEAR string
	//
	COMPLETE string
	DONE     string
	FINISH   string
	CHECK    string
	MARK     string // complete, done, finished, check, mark all mean the same thing
	//
	UPDATE string
}

var TAGS = tags{
	PRINT:     "print",
	LIST:      "list",
	COMPLETED: "completed",
	HELP:      "help",
	VERSION:   "version",
}

var COMMANDS = commands{
	LIST: "list",
	//
	ADD: "add",
	//
	DELETE: "delete",
	REMOVE: "remove",
	RM:     "rm",
	//
	RMC: "rmc", // remove completed
	//
	CLEAR: "clear",
	//
	COMPLETE: "complete",
	DONE:     "done",
	FINISH:   "finish",
	CHECK:    "check",
	MARK:     "mark",
	//
	UPDATE: "update",
}
