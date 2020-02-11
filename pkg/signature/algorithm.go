package signature

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	"io"
	"reflect"

	"github.com/pkg/errors"
)

// Algorithm abstracts over a variety of signature algorithms
//
// The required methods here are low-level, for simplicity of external
// implementation. Consumers should consider generating their keys using
// the `Generate` function and then interacting with the keys and signatures
// using the high-level interface.
type Algorithm interface {
	// PublicKeySize is the size in bytes of this algorithm's public keys
	PublicKeySize() int
	// PrivateKeySize is the size in bytes of this algorithm's private keys
	PrivateKeySize() int
	// SignatureSize is the size in bytes of this algorithm's signatures
	SignatureSize() int

	// Public generates a public key when given a private key
	Public(private []byte) []byte

	// Generate creates a new keypair
	Generate(rand io.Reader) (public, private []byte, err error)
	// Sign signs the message with privateKey and returns a signature
	Sign(private, message []byte) []byte
	// Verify verifies a message's signature
	//
	// Return true if the signature is valid
	Verify(public, message, sig []byte) bool
}

// SameAlgorithm returns true when two algorithms are in fact
// the same algorithm, even if they are not the same instance.
//
// Unknown algorithms are never the same.
func SameAlgorithm(a1 Algorithm, a2 Algorithm) bool {
	id1, err := idOf(a1)
	if err != nil {
		return false
	}
	id2, err := idOf(a2)
	if err != nil {
		return false
	}
	return id1 == id2
}

// Marshal the given data into a serialized binary format which includes
// a type byte for the algorithm.
func marshal(al Algorithm, data []byte) (serialized []byte, err error) {
	id, err := idOf(al)
	if err != nil {
		return nil, err
	}
	container := IdentifiedData{
		Data:      data,
		Algorithm: id,
	}
	return container.MarshalMsg(nil)
}

// shallow copy an interface from an example struct
// https://stackoverflow.com/a/22948379/504550
func cloneAl(original Algorithm) Algorithm {
	val := reflect.ValueOf(original)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	return reflect.New(val.Type()).Interface().(Algorithm)
}

// Unmarshal the serialized binary data into an Algorithm instance and
// the originally supplied data.
func unmarshal(serialized []byte) (al Algorithm, data []byte, err error) {
	container := IdentifiedData{}
	leftovers, err := container.UnmarshalMsg(serialized)
	if err != nil {
		return nil, nil, err
	}
	if len(leftovers) > 0 {
		return nil, nil, errors.New("Leftovers present after deserialization")
	}
	return cloneAl(idMap[container.Algorithm]), container.Data, nil
}

func unmarshalWithLeftovers(serialized []byte) (al Algorithm, data, leftovers []byte, err error) {
	container := IdentifiedData{}
	leftovers, err = container.UnmarshalMsg(serialized)
	if err != nil {
		return nil, nil, nil, err
	}
	return cloneAl(idMap[container.Algorithm]), container.Data, leftovers, nil
}

// Generate a high-level keypair
func Generate(al Algorithm, rdr io.Reader) (public PublicKey, private PrivateKey, err error) {
	pubBytes, privBytes, err := al.Generate(rdr)
	if err == nil {
		public = PublicKey{keyBase{algorithm: al, key: pubBytes, extra: []byte{}}}
		private = PrivateKey{keyBase{algorithm: al, key: privBytes, extra: []byte{}}}
	}
	return
}
