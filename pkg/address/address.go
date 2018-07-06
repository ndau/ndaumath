package address

import "encoding"

///go:generate msgp

///msgp:tuple Address

// An Address is a 48-character string uniquely identifying an Ndau account
//
// For type-safety purposes, it is an opaque struct. This should help make
// it difficult to accidentally pass in a wrong string or something: so long
// as one gets and Address by means of the Generate or Validate functions,
// it is known to be good.
type Address struct {
	addr string
}

var _ encoding.TextMarshaler = (*Address)(nil)
var _ encoding.TextUnmarshaler = (*Address)(nil)

// MarshalText implements encoding.TextMarshaler
func (a *Address) MarshalText() ([]byte, error) {
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
