package address

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Address) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0001 uint32
	zb0001, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zb0001 != 0 {
		err = msgp.ArrayError{Wanted: 0, Got: zb0001}
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Address) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 0
	err = en.Append(0x90)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Address) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 0
	o = append(o, 0x90)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Address) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zb0001 != 0 {
		err = msgp.ArrayError{Wanted: 0, Got: zb0001}
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Address) Msgsize() (s int) {
	s = 1
	return
}
