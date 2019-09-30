package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

type LogEntry struct {
	Level   string `json:"lvl"`
	Message string `json:"msg"`
	Source  string `json:"src"`
}

// validCallback returns true if the last argument is a function.
// Also dispatches an error to the error handler.
func validCallback(args []js.Value) bool {
	obj := getCallback(args)
	if obj.Type() == js.TypeFunction {
		return true
	}
	dispatchError("Last argument must be a callback function.")
	return false
}

// getCallback returns the last argument of the argument array.
func getCallback(args []js.Value) js.Value {
	return args[len(args)-1]
}

// dispatchError passes an error message to the javascript overrideable error handler.
func dispatchError(msg string) {
	js.Global().Call("KeyaddrErrorHandler", msg)
}

func log(l LogEntry) {
	logJSON, _ := json.Marshal(l)
	fmt.Printf("%s\n", string(logJSON))
}

func logDebug(msg string) {
	log(LogEntry{
		Source:  "KEYADDR",
		Level:   "D",
		Message: msg,
	})
}
func logError(msg string) {
	log(LogEntry{
		Source:  "KEYADDR",
		Level:   "E",
		Message: msg,
	})
}
func logInfo(msg string) {
	log(LogEntry{
		Source:  "KEYADDR",
		Level:   "I",
		Message: msg,
	})
}
