package address

import (
	"encoding"
	"fmt"
)

//go:generate msgp

//msgp:tuple Kind

// Kind indicates the type of address in use; this is an external indication
// designed to help users evaluate their own actions; it may or may not be
// enforced by the blockchain.
type Kind struct {
	k string
}

var _ encoding.TextMarshaler = (*Kind)(nil)
var _ encoding.TextUnmarshaler = (*Kind)(nil)

// MarshalText implements encoding.TextMarshaler
func (kind Kind) MarshalText() ([]byte, error) {
	return []byte(kind.k), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (kind *Kind) UnmarshalText(text []byte) error {
	s := string(text)
	if !IsValidKind(s) {
		return fmt.Errorf("invalid kind: %s", s)
	}
	kind.k = s
	return nil
}
