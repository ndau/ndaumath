package types

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Duration) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 int64
		zb0001, err = dc.ReadInt64()
		if err != nil {
			return
		}
		(*z) = Duration(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Duration) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteInt64(int64(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Duration) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendInt64(o, int64(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Duration) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 int64
		zb0001, bts, err = msgp.ReadInt64Bytes(bts)
		if err != nil {
			return
		}
		(*z) = Duration(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Duration) Msgsize() (s int) {
	s = msgp.Int64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Timestamp) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zb0001 int64
		zb0001, err = dc.ReadInt64()
		if err != nil {
			return
		}
		(*z) = Timestamp(zb0001)
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Timestamp) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteInt64(int64(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Timestamp) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendInt64(o, int64(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Timestamp) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zb0001 int64
		zb0001, bts, err = msgp.ReadInt64Bytes(bts)
		if err != nil {
			return
		}
		(*z) = Timestamp(zb0001)
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Timestamp) Msgsize() (s int) {
	s = msgp.Int64Size
	return
}
