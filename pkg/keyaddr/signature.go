package keyaddr

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import "github.com/oneiro-ndev/ndaumath/pkg/signature"

// Signature is the result of signing a block of data with a key.
type Signature struct {
	Signature string
}

// SignatureFrom converts a `signature.Signature` into a `*Signature`
func SignatureFrom(sig signature.Signature) (*Signature, error) {
	sigB, err := sig.MarshalText()
	if err != nil {
		return nil, err
	}

	return &Signature{string(sigB)}, nil
}

// ToSignature converts a `Signature` into a `signature.Signature`
func (s Signature) ToSignature() (signature.Signature, error) {
	sig := signature.Signature{}
	err := sig.UnmarshalText([]byte(s.Signature))
	return sig, err
}
