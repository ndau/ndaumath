package ed25519

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

	impl "golang.org/x/crypto/ed25519"
)

// Ed25519 is the eponymous algorithm; see https://ed25519.cr.yp.to/
//
// Never edit this; it would be a const if go were smarter
var Ed25519 = ed25519{}

type ed25519 struct{}

// PublicKeySize implements Algorithm
func (ed25519) PublicKeySize() int {
	return impl.PublicKeySize
}

// PrivateKeySize implements Algorithm
func (ed25519) PrivateKeySize() int {
	return impl.PrivateKeySize
}

// SignatureSize implements Algorithm
func (ed25519) SignatureSize() int {
	return impl.SignatureSize
}

// Generate implements Algorithm
func (e ed25519) Generate(rand io.Reader) (public, private []byte, err error) {
	return impl.GenerateKey(rand)
}

// Sign implements Algorithm
func (e ed25519) Sign(private, message []byte) []byte {
	return impl.Sign(impl.PrivateKey(private), message)

}

// Verify implements Algorithm
func (ed25519) Verify(public, message, sig []byte) bool {
	return impl.Verify(impl.PublicKey(public), message, sig)
}

// Public generates a public key when given a private key
func (ed25519) Public(private []byte) []byte {
	return impl.PrivateKey(private).Public().(impl.PublicKey)
}
