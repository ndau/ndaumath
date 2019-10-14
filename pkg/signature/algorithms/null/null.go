package null

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import (
	"errors"
	"io"
)

// This implements algorithm for the Null signature type; this type
// exists to allow serializing and deserializing empty signatures without error

// Null is a convenience constant.
// Never edit this; it would be a const if go were smarter
var Null = null{}

type null struct{}

// PublicKeySize implements Algorithm
func (null) PublicKeySize() int {
	return 0
}

// PrivateKeySize implements Algorithm
func (null) PrivateKeySize() int {
	return 0
}

// SignatureSize implements Algorithm
func (null) SignatureSize() int {
	return 0
}

// Generate implements Algorithm
func (null) Generate(rand io.Reader) (public, private []byte, err error) {
	return []byte{}, []byte{}, errors.New("generating null keys is not permitted")
}

// Sign implements Algorithm
func (n null) Sign(private, message []byte) []byte {
	return []byte{}
}

// Verify implements Algorithm
func (null) Verify(public, message, sig []byte) bool {
	return false
}

// Public generates a public key when given a private key
func (null) Public(private []byte) []byte {
	return []byte{}
}
