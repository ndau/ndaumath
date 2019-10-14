package address

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import "encoding"

//========== I M P O R T A N T ====================
// NOTE: These go generate commands are deliberately triple-commented.
// The address field is private, but we want it to be serialized. This
// requires a manual step -- make the field public (by naming it Addr)
// remove the extra slash from these two lines, and do `go generate`.
// Then change all occurences of Addr back to addr (including in the
// generated code). And recomment those two lines.

// Yes, this is horribly ugly.
//=================================================

///go:generate msgp
///msgp:tuple Address

// An Address is a 48-character string uniquely identifying an Ndau account
//
// For type-safety purposes, it is an opaque struct. This should help make
// it difficult to accidentally pass in a wrong string or something: so long
// as one gets an Address by means of the Generate or Validate functions,
// it is known to be good.
type Address struct {
	addr string
}

var _ encoding.TextMarshaler = (*Address)(nil)
var _ encoding.TextUnmarshaler = (*Address)(nil)

// MarshalText implements encoding.TextMarshaler
func (a Address) MarshalText() ([]byte, error) {
	return []byte(a.addr), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (a *Address) UnmarshalText(text []byte) error {
	s := string(text)
	_, err := Validate(s)
	if err != nil {
		return err
	}
	a.addr = s
	return nil
}
