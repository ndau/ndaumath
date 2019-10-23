package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/pkg/errors"
)

type LogEntry struct {
	Level   string `json:"lvl"`
	Message string `json:"msg"`
	Source  string `json:"src"`
}

// handleArgs returns a callback and remaining arguments.
// It also prints and returns error messages.
// when this function returns err != nil, the handler should end asap.
func handleArgs(args []js.Value, expected int, source string) (js.Value, []js.Value, error) {

	// check that callback is a function
	cb := args[len(args)-1]
	if cb.Type() != js.TypeFunction {
		msg := fmt.Sprintf("couldn't parse %s arguments: last argument must be a callback function.", source)
		dispatchError(msg)
		return js.Value{}, []js.Value{}, errors.New(msg)
	}

	// check argument length
	// ignores callback
	if len(args) != expected+1 {
		msg := fmt.Sprintf("couldn't parse %s arguments: incorrect amount of arguments", source)
		cb.Invoke(msg, nil)
		return js.Value{}, []js.Value{}, errors.New(msg)
	}

	remainder := args[:len(args)-1]
	return cb, remainder, nil
}

// dispatchError passes an error message to the javascript overrideable error handler.
func dispatchError(msg string) {
	js.Global().Call("KeyaddrErrorHandler", msg)
}

func log(l LogEntry) {
	logJSON, _ := json.Marshal(l)
	fmt.Printf("%s\n", string(logJSON))
}

// TODO: add log levels
var levels = map[string]int{
	"D": 0,
	"I": 1,
	"E": 2,
	"":  0, // default to 0
}

const (
	levelDebug = "D"
	levelInfo  = "I"
	levelError = "E"
)

func logDebug(msg string) {
	if levels[js.Global().Get("KeyaddrLogLevel").String()] <= levels[levelDebug] {
		log(LogEntry{
			Source:  "KEYADDR",
			Level:   levelDebug,
			Message: msg,
		})
	}
}
func logError(msg string) {
	if levels[js.Global().Get("KeyaddrLogLevel").String()] <= levels[levelError] {
		log(LogEntry{
			Source:  "KEYADDR",
			Level:   levelError,
			Message: msg,
		})
	}
}
func logInfo(msg string) {
	if levels[js.Global().Get("KeyaddrLogLevel").String()] <= levels[levelInfo] {
		log(LogEntry{
			Source:  "KEYADDR",
			Level:   levelInfo,
			Message: msg,
		})
	}
}
