package types

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Lock) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "notice":
			err = z.NoticePeriod.DecodeMsg(dc)
			if err != nil {
				return
			}
		case "unlock":
			if dc.IsNil() {
				err = dc.ReadNil()
				if err != nil {
					return
				}
				z.UnlocksOn = nil
			} else {
				if z.UnlocksOn == nil {
					z.UnlocksOn = new(Timestamp)
				}
				err = z.UnlocksOn.DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Lock) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "notice"
	err = en.Append(0x82, 0xa6, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65)
	if err != nil {
		return
	}
	err = z.NoticePeriod.EncodeMsg(en)
	if err != nil {
		return
	}
	// write "unlock"
	err = en.Append(0xa6, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b)
	if err != nil {
		return
	}
	if z.UnlocksOn == nil {
		err = en.WriteNil()
		if err != nil {
			return
		}
	} else {
		err = z.UnlocksOn.EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Lock) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "notice"
	o = append(o, 0x82, 0xa6, 0x6e, 0x6f, 0x74, 0x69, 0x63, 0x65)
	o, err = z.NoticePeriod.MarshalMsg(o)
	if err != nil {
		return
	}
	// string "unlock"
	o = append(o, 0xa6, 0x75, 0x6e, 0x6c, 0x6f, 0x63, 0x6b)
	if z.UnlocksOn == nil {
		o = msgp.AppendNil(o)
	} else {
		o, err = z.UnlocksOn.MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Lock) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "notice":
			bts, err = z.NoticePeriod.UnmarshalMsg(bts)
			if err != nil {
				return
			}
		case "unlock":
			if msgp.IsNil(bts) {
				bts, err = msgp.ReadNilBytes(bts)
				if err != nil {
					return
				}
				z.UnlocksOn = nil
			} else {
				if z.UnlocksOn == nil {
					z.UnlocksOn = new(Timestamp)
				}
				bts, err = z.UnlocksOn.UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Lock) Msgsize() (s int) {
	s = 1 + 7 + z.NoticePeriod.Msgsize() + 7
	if z.UnlocksOn == nil {
		s += msgp.NilSize
	} else {
		s += z.UnlocksOn.Msgsize()
	}
	return
}
