package main

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

/*

Notes:

* goroutines are required in each of the handlers to prevent Go from flagging the program as deadlocked.
* This WASM module loads into the javascriprt engine like a separate application, rather like a server. It does not import a library, like a normal javascript module would.
* Panics will shut down the wasm application and it will not automatically restart. @TODO, add a line to report a panic.
* There is a global javascript function KeyaddrErrorHandler that can be overriden before this WASM module is loaded.
* There are three different levels of errors here
	* Panics which are avoided by code and should not ever occur. A @todo is mentioned above to handle reporting for this case.
	* Errors that are sent back through the callback. They are given as the first argument in the node-style
		callbacks, which `promisify` can put into promises, which can be handled with chained `.catch()` methods
		or the `try`/`catch` blocks.
	* Application level errors that cannot be sent to a callback. This can be because a callback was not provided, or an error occured that prevented a callback from being used. `dispatchError` and `KeyaddrErrorHandler` is used for these cases.

*/

import (
	"syscall/js"
)

// waitChannel is used to control keeping the process open, as well as closing on command.
var waitChannel chan struct{}

func main() {
	js.Global().Get("console").Call("log", "WASM Keyaddr starting")

	waitChannel = make(chan struct{})

	// put go functions in a javascript object
	obj := map[string]interface{}{
		"newKey":        js.FuncOf(newKey),
		"wordsToBytes":  js.FuncOf(wordsToBytes),
		"deriveFrom":    js.FuncOf(deriveFrom),
		"ndauAddress":   js.FuncOf(ndauAddress),
		"toPublic":      js.FuncOf(toPublic),
		"child":         js.FuncOf(child),
		"sign":          js.FuncOf(sign),
		"hardenedChild": js.FuncOf(hardenedChild),
		"exit":          js.FuncOf(exit),
	}

	// Register all functions globally under KeyaddrNS. Either `window` in browsers, or
	// `global` in node. NS stands for node style and refers to the functions which use
	// callbacks instead of promises. In node, the functions may be easily turned into
	// promisified functions with `util.promisify`.
	js.Global().Set("KeyaddrNS", js.ValueOf(obj))

	// register default error handler if error handler doesn't already exist
	if js.Global().Get("KeyaddrErrorHandler").Truthy() == false {
		js.Global().Get("console").Call("log", "Global function KeyaddrErrorHandler not detected. Adding default error handler function. To override, instantiate KeyaddrErrorHandler early in the JS code.")
		js.Global().Set("KeyaddrErrorHandler", js.FuncOf(errorHandler))
	}

	<-waitChannel // wait indefinitely

	js.Global().Get("console").Call("log", "WASM Keyadder exiting")
}
