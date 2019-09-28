package main

import (
	"fmt"
	"syscall/js"
)

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

func log(level, msg string) {
	fmt.Printf("%s: %s\n", level, msg)
}

func logDebug(msg string) {
	log("KEYADDR DEBUG", msg)
}
func logError(msg string) {
	log("KEYADDR ERROR", msg)
}
func logInfo(msg string) {
	log("KEYADDR INFO", msg)
}
