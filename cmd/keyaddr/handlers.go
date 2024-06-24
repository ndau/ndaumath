package main

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import (
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"syscall/js"

	"github.com/ndau/ndaumath/pkg/address"
	"github.com/ndau/ndaumath/pkg/keyaddr"
	"github.com/ndau/ndaumath/pkg/signature"
)

// exit will quit the go program, causing the application to no longer respond to function calls.
func exit(this js.Value, args []js.Value) interface{} {
	go func() {
		logInfo("keyaddr exiting...")

		// clean args
		callback, _, err := handleArgs(args, 0, "exit") // _ are remaining args but we're not expecting any
		if err != nil {
			return
		}
		callback.Invoke(nil, "keyaddr exiting...")
		// cause main to be unblocked waiting for channel input and exit
		close(waitChannel)
		return

	}()
	return nil
}

// errorHandler is a javascript overrideable error handler that prints to the console.
func errorHandler(this js.Value, args []js.Value) interface{} {
	msg := args[0]
	js.Global().Get("console").Call("log", fmt.Sprintf("WASM Keyaddr Error: %s", msg))
	return nil
}

// JS Usage: newKey(recoveryBytes, cb)
func newKey(this js.Value, args []js.Value) interface{} {
	go func(args []js.Value) {
		logDebug("newKey")

		// clean args
		callback, remainder, err := handleArgs(args, 1, "newKey")
		if err != nil {
			return
		}

		recoveryBytes := remainder[0].String()

		// do work
		key, err := keyaddr.NewKey(recoveryBytes)
		if err != nil {
			jsLogReject(callback, "error creating new key: %s", err)
			return
		}

		// return result
		callback.Invoke(nil, key.Key)
		return
	}(args)
	return nil
}

// JS Usage: newEdKey(cb) --> no args
func newEdKey(this js.Value, args []js.Value) interface{} {
	go func(args []js.Value) {
		logDebug("newEdKey")

		// clean args
		callback, _, err := handleArgs(args, 0, "newEdKey")
		if err != nil {
			return
		}

		// do work
		public, private, err := signature.Generate(signature.Ed25519, nil)
		if err != nil {
			jsLogReject(callback, "error creating new ed key: %s", err)
			return
		}

		pubkeyText, err := public.MarshalText()
		if err != nil {
			jsLogReject(callback, "error marshalling new public ed key: %s", err)
			return
		}

		privkeyText, err := private.MarshalText()
		if err != nil {
			jsLogReject(callback, "error marshalling new private ed key: %s", err)
			return
		}

		// Make an map[string]interface{} so syscall will turn it into a js object.
		obj := make(map[string]interface{})
		obj["pubkey"] = string(pubkeyText)
		obj["privkey"] = string(privkeyText)

		// return result
		callback.Invoke(nil, js.ValueOf(obj))
		return
	}(args)
	return nil
}

// JS Usage: newEdKey(seed, cb)
func newEdKeyFromSeed(this js.Value, args []js.Value) interface{} {
	go func(args []js.Value) {
		logDebug("newEdKeyFromSeed")

		// clean args
		callback, remainder, err := handleArgs(args, 1, "newEdKeyFromSeed")
		if err != nil {
			return
		}

		seedString := remainder[0].String()
		seed, err := hex.DecodeString(seedString)
		if err != nil {
			jsLogReject(callback, "error decoding data from hex: %s", err)
			return
		}

		// do work
		public, private, err := signature.GenerateFromSeed(signature.Ed25519, seed)
		if err != nil {
			jsLogReject(callback, "error creating new ed key: %s", err)
			return
		}

		pubkeyText, err := public.MarshalText()
		if err != nil {
			jsLogReject(callback, "error marshalling new public ed key: %s", err)
			return
		}

		privkeyText, err := private.MarshalText()
		if err != nil {
			jsLogReject(callback, "error marshalling new private ed key: %s", err)
			return
		}

		// Make an map[string]interface{} so syscall will turn it into a js object.
		obj := make(map[string]interface{})
		obj["pubkey"] = string(pubkeyText)
		obj["privkey"] = string(privkeyText)

		// return result
		callback.Invoke(nil, js.ValueOf(obj))
		return
	}(args)
	return nil
}

// JS Usage: addrFromPublicKey(pubKey, cb)
func addrFromPublicKey(this js.Value, args []js.Value) interface{} {
	go func(args []js.Value) {
		logDebug("addrFromPublicKey")

		// clean args
		callback, remainder, err := handleArgs(args, 1, "addrFromPublicKey")
		if err != nil {
			return
		}

		pubKeyString := remainder[0].String()

		// do work
		key, err := signature.ParseKey(pubKeyString)
		if err != nil {
			jsLogReject(callback, "error parsing public key: %s", err)
			return
		}

		_, ok := key.(*signature.PublicKey)
		if !ok {
			jsLogReject(callback, "addresses can only be generated from public keys")
			return
		}

		addr, err := address.Generate(address.KindExchange, key.KeyBytes())
		if err != nil {
			jsLogReject(callback, "error generating address from key: %s", err)
			return
		}

		// return result
		callback.Invoke(nil, addr.String())
		return
	}(args)
	return nil
}

// JS Usage: wordsToBytes(language, words, cb)
// language defaults to en if not specified.
func wordsToBytes(this js.Value, args []js.Value) interface{} {
	go func(args []js.Value) {
		logDebug("wordsToBytes")
		// clean args
		callback, remainder, err := handleArgs(args, 2, "wordsToBytes")
		if err != nil {
			return
		}

		lang := "en" // default to english if language not specified
		if remainder[0].Type() != js.TypeUndefined {
			lang = remainder[0].String()
		}
		words := remainder[1].String()

		// checks for a space character to count words
		re := regexp.MustCompile(" ")
		wc := len(re.FindAllStringIndex(words, -1)) + 1

		logDebug(fmt.Sprintf("number of words:%v lang:%s", wc, lang))

		// do work
		bs, err := keyaddr.WordsToBytes(lang, words)
		if err != nil {

			jsLogReject(callback, "error converting words to bytes: %s", err)
			return
		}

		// return result
		callback.Invoke(nil, bs)
		return
	}(args)
	return nil
}

// JS Usage: deriveFrom(parentKey, parentPath, childPath, cb)
func deriveFrom(this js.Value, args []js.Value) interface{} {
	go func(args []js.Value) {
		logDebug("deriveFrom")
		// clean args
		callback, remainder, err := handleArgs(args, 3, "deriveFrom")
		if err != nil {
			return
		}

		parentKey := remainder[0].String()
		parentPath := remainder[1].String()
		childPath := remainder[2].String()

		// do work
		der, err := keyaddr.DeriveFrom(parentKey, parentPath, childPath)
		if err != nil {
			jsLogReject(callback, "error deriving new key: %s", err)
			return
		}

		// return result
		callback.Invoke(nil, der.Key)

		return
	}(args)
	return nil
}

// JS Usage: wordsFromPrefix(lang, prefix, max, cb)
func wordsFromPrefix(this js.Value, args []js.Value) interface{} {
	go func(args []js.Value) {
		logDebug("wordsFromPrefix")
		// clean args
		callback, remainder, err := handleArgs(args, 3, "wordsFromPrefix")
		if err != nil {
			return
		}

		lang := remainder[0].String()
		prefix := remainder[1].String()
		max := remainder[2].Int()

		// do work
		words := keyaddr.WordsFromPrefix(lang, prefix, max)

		// return result
		callback.Invoke(nil, words)

		return
	}(args)
	return nil
}

// JS Usage: isPrivate(key, cb)
func isPrivate(this js.Value, args []js.Value) interface{} {
	go func(args []js.Value) {
		logDebug("isPrivate")
		// clean args
		callback, remainder, err := handleArgs(args, 1, "isPrivate")
		if err != nil {
			return
		}

		k := &keyaddr.Key{
			Key: remainder[0].String(),
		}

		// do work
		isPrivateResult, err := k.IsPrivate()
		if err != nil {
			jsLogReject(callback, "error testing key type: %s", err)
			return
		}

		// return result
		callback.Invoke(nil, isPrivateResult)

		return
	}(args)
	return nil
}

// constructs a key from a string. Possibly to check for validity?
// JS Usage: fromString(key, cb)
func fromString(this js.Value, args []js.Value) interface{} {
	go func(args []js.Value) {
		logDebug("fromString")
		// clean args
		callback, remainder, err := handleArgs(args, 1, "fromString")
		if err != nil {
			return
		}

		str := remainder[0].String()

		// do work
		key, err := keyaddr.FromString(str)
		if err != nil {
			jsLogReject(callback, "error constructing a key from a string: %s", err)
			return
		}

		// Make an map[string]interface{} so syscall will turn it into a js object.
		obj := make(map[string]interface{})
		obj["key"] = key.Key

		// return result
		callback.Invoke(nil, js.ValueOf(obj))

		return
	}(args)
	return nil
}

// JS Usage: wordsFromBytes(lang, base64bytes, cb)
func wordsFromBytes(this js.Value, args []js.Value) interface{} {
	go func(args []js.Value) {
		logDebug("wordsFromBytes")
		// clean args
		callback, remainder, err := handleArgs(args, 2, "wordsFromBytes")
		if err != nil {
			return
		}

		lang := remainder[0].String()
		bs := remainder[1].String()

		// do work
		words, err := keyaddr.WordsFromBytes(lang, bs)
		if err != nil {
			jsLogReject(callback, "error converting bytes to words: %s", err)
			return
		}

		// return result
		callback.Invoke(nil, words)

		return
	}(args)
	return nil
}

// JS Usage: ndauAddress(privateKey, cb)
func ndauAddress(this js.Value, args []js.Value) interface{} {
	go func(args []js.Value) {
		logDebug("ndauAddress")
		// clean args
		callback, remainder, err := handleArgs(args, 1, "ndauAddress")
		if err != nil {
			return
		}

		k := &keyaddr.Key{
			Key: remainder[0].String(),
		}

		// do work
		addr, err := k.NdauAddress()
		if err != nil {
			jsLogReject(callback, "error getting ndau address: %s", err)
			return
		}
		// return result
		callback.Invoke(nil, addr.Address)

		return
	}(args)
	return nil
}

// JS Usage: toPublic(privateKey, cb)
func toPublic(this js.Value, args []js.Value) interface{} {
	go func(args []js.Value) {
		logDebug("toPublic")
		// clean args
		callback, remainder, err := handleArgs(args, 1, "toPublic")
		if err != nil {
			return
		}

		k := &keyaddr.Key{
			Key: remainder[0].String(),
		}

		// do work
		pub, err := k.ToPublic()
		if err != nil {
			jsLogReject(callback, "error converting to public key: %s", err)
			return
		}

		// return result
		callback.Invoke(nil, pub.Key)

		return
	}(args)

	return nil
}

// JS Usage: child(privateKey, n, cb)
func child(this js.Value, args []js.Value) interface{} {
	go func(args []js.Value) {
		logDebug("child")
		// clean args
		callback, remainder, err := handleArgs(args, 2, "child")
		if err != nil {
			return
		}

		k := &keyaddr.Key{
			Key: remainder[0].String(),
		}

		if remainder[1].Type() != js.TypeNumber {
			jsLogReject(callback, "n must be of type Number")
			return
		}

		n := remainder[1].Int()
		if n < math.MinInt32 || n > math.MaxInt32 {
			jsLogReject(callback, "n must not overflow int32")
			return
		}

		n32 := int32(n)

		// do work
		key, err := k.Child(int32(n32))
		if err != nil {
			jsLogReject(callback, "error creating child key: %s", err)
			return
		}

		// return result
		callback.Invoke(nil, key.Key)

		return
	}(args)
	return nil
}

// JS Usage: sign(privateKey, base64Message, cb)
func sign(this js.Value, args []js.Value) interface{} {
	go func(args []js.Value) {
		logDebug("sign")
		// clean args
		callback, remainder, err := handleArgs(args, 2, "sign")
		if err != nil {
			return
		}

		k := keyaddr.Key{
			Key: remainder[0].String(),
		}

		msg := remainder[1].String()

		// do work
		sig, err := k.SignSecP(msg)
		if err != nil {
			logError(fmt.Sprintf("error creating signature: key length: %s, msg: %s, err: %s", len(k.Key), msg, err.Error()))
			jsLogReject(callback, "error creating signature: %s", err)
			return
		}

		// return result
		callback.Invoke(nil, sig.Signature)
		return
	}(args)
	return nil
}

// JS Usage: signEdB64(privateKey, base64Message, cb)
func signEdB64(this js.Value, args []js.Value) interface{} {
	js.Global().Get("console").Call("log", "sign starting")
	go func(args []js.Value) {
		logDebug("signEdB64")
		// clean args
		callback, remainder, err := handleArgs(args, 2, "signEdB64")
		if err != nil {
			return
		}

		k := keyaddr.Key{
			Key: remainder[0].String(),
		}

		msg := remainder[1].String()

		// do work
		sig, err := k.SignEdB64(msg)
		if err != nil {
			logError(fmt.Sprintf("error creating signature: key length: %s, msg: %s, err: %s", len(k.Key), msg, err.Error()))
			jsLogReject(callback, "error creating signature: %s", err)
			return
		}

		// return result
		callback.Invoke(nil, sig.Signature)
		return
	}(args)
	return nil
}

// JS Usage: signEdText(privateKey, textMessage, cb)
func signEdText(this js.Value, args []js.Value) interface{} {
	js.Global().Get("console").Call("log", "sign starting")
	go func(args []js.Value) {
		logDebug("signEdText")
		// clean args
		callback, remainder, err := handleArgs(args, 2, "signEdText")
		if err != nil {
			return
		}

		k := keyaddr.Key{
			Key: remainder[0].String(),
		}

		msg := remainder[1].String()

		// do work
		sig, err := k.SignEdText(msg)
		if err != nil {
			logError(fmt.Sprintf("error creating signature: key length: %s, msg: %s, err: %s", len(k.Key), msg, err.Error()))
			jsLogReject(callback, "error creating signature: %s", err)
			return
		}

		// return result
		callback.Invoke(nil, sig.Signature)
		return
	}(args)
	return nil
}

// JS Usage: hardenedChild(privateKey, n, cb)
func hardenedChild(this js.Value, args []js.Value) interface{} {
	go func(args []js.Value) {
		logDebug("hardenedChild")
		// clean args
		callback, remainder, err := handleArgs(args, 2, "hardenedChild")
		if err != nil {
			return
		}

		k := &keyaddr.Key{
			Key: remainder[0].String(),
		}

		if remainder[1].Type() != js.TypeNumber {
			jsLogReject(callback, "n must be of type Number")
			return
		}

		n := remainder[1].Int()
		if n < math.MinInt32 || n > math.MaxInt32 {
			jsLogReject(callback, "n must not overflow int32")
			return
		}

		n32 := int32(n)

		// do work
		key, err := k.HardenedChild(n32)
		if err != nil {
			jsLogReject(callback, "error hardening child: %s", err)
			return
		}
		// return result
		callback.Invoke(nil, key.Key)

		return
	}(args)
	return nil
}
