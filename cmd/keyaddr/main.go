package main

/*

Notes:

* goroutines are required in each of the handlers to prevent Go from detecting deadlocks.
* This WASM module loads into the javascriprt engine like a separate application, rather like a server. It does not import a library, like a normal javascript module would. It
* Panics will shut down the wasm application and it will not automatically restart. @TODO, add a line to report a panic.
* There is a global javascript function KeyaddrErrorHandler that can be overriden before this WASM module is loaded.

*/

import (
	"syscall/js"
)

// c is used to control keeping the process open, as well as closing on command.
var c chan bool

func main() {
	js.Global().Get("console").Call("log", "WASM Keyaddr starting")
	c := make(chan bool)

	// put go functions in a javascript object
	obj := make(map[string]interface{})
	obj["newKey"] = js.FuncOf(newKey)
	obj["wordsToBytes"] = js.FuncOf(wordsToBytes)
	obj["deriveFrom"] = js.FuncOf(deriveFrom)
	obj["ndauAddress"] = js.FuncOf(ndauAddress)
	obj["toPublic"] = js.FuncOf(toPublic)
	obj["child"] = js.FuncOf(child)
	obj["sign"] = js.FuncOf(sign)
	obj["hardenedChild"] = js.FuncOf(hardenedChild)
	obj["newKey"] = js.FuncOf(newKey)
	obj["exit"] = js.FuncOf(exit)

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

	<-c // wait indefinitely

	js.Global().Get("console").Call("log", "WASM Keyadder exiting")
}
