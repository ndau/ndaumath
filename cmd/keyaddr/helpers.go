package main

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/pkg/errors"
)

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
	// The expression below (expected+1) makes this check ignore the callback
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

// LogEntry represents a log message for json marshalling.
type LogEntry struct {
	Level     string `json:"lvl"`
	Message   string `json:"msg"`
	Source    string `json:"src"`
	Timestamp int    `json:"ts"`
}

// log marshalls a message to json and outputs to the console
func log(level string, msg string) {
	l := LogEntry{
		Source:    "KEYADDR",
		Level:     level,
		Message:   msg,
		Timestamp: js.Global().Get("Date").Call("now").Int(),
	}
	logJSON, _ := json.Marshal(l)
	js.Global().Get("console").Call("log", fmt.Sprintf("%s", string(logJSON)))
}

// `levels` values are also available in the JS environment as KeyaddrLogLevelDebug, KeyaddrLogLevelInfo, KeyaddrLogLevelError
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

func checkLevel(level string) bool {
	return levels[js.Global().Get("KeyaddrLogLevel").String()] <= levels[level]
}

// logInfo logs with a 'debug' level
func logDebug(msg string) {
	if checkLevel(levelDebug) {
		log(levelDebug, msg)
	}
}

// logError logs with an 'error' level
func logError(msg string) {
	if checkLevel(levelError) {
		log(levelError, msg)
	}
}

// logInfo logs with an 'info' level
func logInfo(msg string) {
	if checkLevel(levelInfo) {
		log(levelInfo, msg)
	}
}

// jsLogReject returns a JS error to the callback and logs the error
func jsLogReject(cb js.Value, str string, sprintfArgs ...interface{}) {
	msg := fmt.Sprintf(str, sprintfArgs...)
	jsErr := js.Global().Get("Error").Invoke(msg)
	logError(msg)
	cb.Invoke(jsErr, nil)
}
