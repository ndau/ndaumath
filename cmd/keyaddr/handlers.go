package main

import (
	"fmt"
	"regexp"
	"syscall/js"

	"github.com/oneiro-ndev/ndaumath/pkg/keyaddr"
)

// exit will quit the go program, causing the application to no longer respond to function calls.
func exit(this js.Value, args []js.Value) interface{} {
	c <- true
	return nil
}

// errorHandler is a javascript overrideable error handler that prints to the console.
func errorHandler(this js.Value, args []js.Value) interface{} {
	msg := args[0]
	js.Global().Get("console").Call("log", fmt.Sprintf("WASM Keyaddr Error: %s", msg))
	return nil
}

// js usage: newKey(recoveryBytes, cb)
func newKey(this js.Value, args []js.Value) interface{} {

	go func(this js.Value, args []js.Value) interface{} {

		logInfo("newKey")
		// clean args
		if !validCallback(args) {
			return nil
		}
		callback := getCallback(args)

		recoveryBytes := args[0].String()

		// do work
		key, err := keyaddr.NewKey(recoveryBytes)
		if err != nil {
			callback.Invoke(fmt.Sprintf("error creating new key: %s", err), nil)
			return nil
		}

		// return result
		callback.Invoke(nil, fmt.Sprintf("%s", key.Key))
		return nil
	}(this, args)
	return nil
}

// JS Usage: wordsToBytes(language, words, cb)
func wordsToBytes(this js.Value, args []js.Value) interface{} {
	go func(this js.Value, args []js.Value) interface{} {
		logInfo("wordsToBytes")
		// clean args
		if !validCallback(args) {
			return nil
		}
		callback := getCallback(args)

		lang := "en"
		if args[0].Type() != js.TypeUndefined {
			lang = args[0].String()
		}
		words := args[1].String()

		re := regexp.MustCompile(" ")
		matches := re.FindAllStringIndex(words, -1)
		logDebug(fmt.Sprintf("number of words:%s lang:%s", len(matches), lang))

		// do work
		bs, err := keyaddr.WordsToBytes(lang, words)
		if err != nil {
			callback.Invoke(fmt.Sprintf("error converting words to bytes: %s", err), nil)
			return nil
		}

		// return result
		callback.Invoke(nil, bs)
		return nil
	}(this, args)
	return nil
}

// JS Usage: deriveFrom(parentKey, parentPath, childPath, cb)
func deriveFrom(this js.Value, args []js.Value) interface{} {

	go func(this js.Value, args []js.Value) interface{} {
		logInfo("deriveFrom")

		// clean args
		if !validCallback(args) {
			return nil
		}
		callback := getCallback(args)
		parentKey := args[0].String()
		parentPath := args[1].String()
		childPath := args[2].String()

		// do work
		der, err := keyaddr.DeriveFrom(parentKey, parentPath, childPath)
		if err != nil {
			callback.Invoke(fmt.Sprintf("error deriving new key: %s", err), nil)
			return nil
		}

		// return result
		callback.Invoke(nil, der.Key)

		return nil
	}(this, args)
	return nil
}

// JS usage: ndauAddress(privateKey, cb)
func ndauAddress(this js.Value, args []js.Value) interface{} {
	go func(this js.Value, args []js.Value) interface{} {
		logInfo("ndauAddress")

		// clean args
		if !validCallback(args) {
			return nil
		}
		callback := getCallback(args)

		k := &keyaddr.Key{
			Key: args[0].String(),
		}

		// do work
		addr, err := k.NdauAddress()
		if err != nil {
			callback.Invoke(fmt.Sprintf("error getting ndau address: %s", err), nil)
			return nil
		}
		// return result
		callback.Invoke(nil, addr.Address)

		return nil
	}(this, args)
	return nil
}

// JS usage: toPublic(privateKey, cb)
func toPublic(this js.Value, args []js.Value) interface{} {

	go func(this js.Value, args []js.Value) interface{} {
		logInfo("toPublic")

		// clean args
		if !validCallback(args) {
			return nil
		}
		callback := getCallback(args)

		k := &keyaddr.Key{
			Key: args[0].String(),
		}

		// do work
		pub, err := k.ToPublic()
		if err != nil {
			callback.Invoke(fmt.Sprintf("error converting to public key: %s", err), nil)
			return nil
		}

		// return result
		callback.Invoke(nil, pub.Key)

		return nil
	}(this, args)

	return nil
}

// JS usage: child(privateKey, n)
func child(this js.Value, args []js.Value) interface{} {

	go func(this js.Value, args []js.Value) interface{} {
		logInfo("child")

		// clean args
		if !validCallback(args) {
			return nil
		}
		callback := getCallback(args)

		k := &keyaddr.Key{
			Key: args[0].String(),
		}

		if args[1].Type() != js.TypeNumber {
			callback.Invoke("n must be of type Number", nil)
			return nil
		}

		n := args[1].Int()
		if n < -2147483648 || n > 2147483647 {
			callback.Invoke("n must not overflow int32")
			return nil
		}

		n32 := int32(n)

		// do work
		key, err := k.Child(int32(n32))
		if err != nil {
			callback.Invoke(fmt.Sprintf("error creating child key: %s", err), nil)
			return nil
		}

		// return result
		callback.Invoke(nil, key.Key)

		return nil
	}(this, args)
	return nil
}

// js usage: sign(privateKey, base64Message, cb)
func sign(this js.Value, args []js.Value) interface{} {

	go func(this js.Value, args []js.Value) interface{} {
		logInfo("sign")

		// clean args
		if !validCallback(args) {
			return nil
		}
		callback := getCallback(args)

		k := keyaddr.Key{
			Key: args[0].String(),
		}

		msg := args[1].String()

		// do work
		sig, err := k.Sign(msg)
		if err != nil {
			logError(fmt.Sprintf("key length: %s, msg: %s, err: %s", len(k.Key), msg, err.Error()))
			callback.Invoke(fmt.Sprintf("error creating signature: %s", err), nil)
			return nil
		}

		// return result
		callback.Invoke(nil, sig.Signature)
		return nil
	}(this, args)

	return nil
}

// JS Usage: hardenedChild(privateKey, n)
func hardenedChild(this js.Value, args []js.Value) interface{} {
	go func(this js.Value, args []js.Value) interface{} {
		logInfo("hardenedChild")

		// clean args
		if !validCallback(args) {
			return nil
		}
		callback := getCallback(args)

		k := &keyaddr.Key{
			Key: args[0].String(),
		}

		if args[1].Type() != js.TypeNumber {
			callback.Invoke("n must be of type Number", nil)
			return nil
		}

		n := args[1].Int()
		if n < -2147483648 || n > 2147483647 {
			callback.Invoke("n must not overflow int32")
			return nil
		}

		n32 := int32(n)

		// do work
		key, err := k.HardenedChild(n32)
		if err != nil {
			callback.Invoke(fmt.Sprintf("error hardening child: %s", err), nil)
			return nil
		}
		// return result
		callback.Invoke(nil, key.Key)

		return nil
	}(this, args)
	return nil
}
