package signature

//go:generate msgp

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

//msgp:tuple IdentifiedData

// AlgorithmID is an identifier uniquely associated with each supported signature algorithm
type AlgorithmID uint8

// IdentifiedData is a byte slice associated with an algorithm
type IdentifiedData struct {
	Algorithm AlgorithmID
	Data      []byte
}
