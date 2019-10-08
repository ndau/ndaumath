package signature

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignRoundtripEd25519(t *testing.T) {
	message := make([]byte, 256)
	rand.Read(message)

	public, private, err := Generate(Ed25519, nil)
	require.NoError(t, err)

	signature := private.Sign(message)
	t.Logf("1 data:        %x", signature.data)

	signatureBytes, err := signature.Marshal()
	require.NoError(t, err)
	t.Logf("bytes: %x", signatureBytes)

	signature2 := Signature{}
	err = (&signature2).Unmarshal(signatureBytes)
	require.NoError(t, err)

	t.Logf("2 data:        %x", signature2.data)

	// method 1 to verify a signature
	require.True(t, signature2.Verify(message, public))
	// method 2 to verify a signature
	require.True(t, public.Verify(message, signature2))
}

func TestSignRoundtripEd25519Text(t *testing.T) {
	message := make([]byte, 256)
	rand.Read(message)

	public, private, err := Generate(Ed25519, nil)
	require.NoError(t, err)

	signature := private.Sign(message)
	t.Logf("1 data:        %x\n", signature.data)

	signatureBytes, err := signature.MarshalText()
	require.NoError(t, err)
	t.Logf("bytes: %s\n", signatureBytes)

	signature2 := Signature{}
	err = (&signature2).UnmarshalText(signatureBytes)
	require.NoError(t, err)

	t.Logf("2 data:        %x\n", signature2.data)

	// method 1 to verify a signature
	require.True(t, signature2.Verify(message, public))
	// method 2 to verify a signature
	require.True(t, public.Verify(message, signature2))
}

func TestSignRoundtripSecp256k1(t *testing.T) {
	message := make([]byte, 256)
	rand.Read(message)

	public, private, err := Generate(Secp256k1, nil)
	require.NoError(t, err)

	signature := private.Sign(message)
	t.Logf("1 data:        %x", signature.data)

	signatureBytes, err := signature.Marshal()
	require.NoError(t, err)
	t.Logf("bytes: %x", signatureBytes)

	signature2 := Signature{}
	err = (&signature2).Unmarshal(signatureBytes)
	require.NoError(t, err)

	t.Logf("2 data:        %x", signature2.data)

	// method 1 to verify a signature
	require.True(t, signature2.Verify(message, public))
	// method 2 to verify a signature
	require.True(t, public.Verify(message, signature2))
}

func TestSignRoundtripSecp256k1Text(t *testing.T) {
	message := make([]byte, 256)
	rand.Read(message)

	public, private, err := Generate(Secp256k1, nil)
	require.NoError(t, err)

	signature := private.Sign(message)
	t.Logf("1 data:        %x\n", signature.data)

	signatureBytes, err := signature.MarshalText()
	require.NoError(t, err)
	t.Logf("bytes: %s\n", signatureBytes)

	signature2 := Signature{}
	err = (&signature2).UnmarshalText(signatureBytes)
	require.NoError(t, err)

	t.Logf("2 data:        %x\n", signature2.data)

	// method 1 to verify a signature
	require.True(t, signature2.Verify(message, public))
	// method 2 to verify a signature
	require.True(t, public.Verify(message, signature2))
}

func TestMarshalNull(t *testing.T) {
	var s Signature

	ser, err := s.Marshal()
	require.Nil(t, err)

	var s2 Signature
	err = s2.Unmarshal(ser)
	require.Nil(t, err)
}

func TestMarshalMsgNull(t *testing.T) {
	var s Signature

	var b []byte
	ser, err := s.MarshalMsg(b)
	require.Nil(t, err)

	var s2 Signature
	leftover, err := s2.UnmarshalMsg(ser)
	require.Nil(t, err)
	require.Zero(t, len(leftover))
}

func TestUnmarshal(t *testing.T) {
	pubkbytes := "npuba8jadtbbebmmi9j8838464z7u7vxgpzyfebhuhyrqcnuz97mitidqytia3pu7nbe43pn2m6x"
	var k PublicKey
	err := k.UnmarshalText([]byte(pubkbytes))
	fmt.Println(err)

	pvtkbytes := "npvtayjadtcbibvri9mw8awpwks773tj5nwfz93xzbi98gaqyajxek3q923jq8xt6xxwrw9rn9pqpm83q34vg55cuav3d5hzbgjm98xwiwbzmiwany3qd9tth2qh"
	var k2 PrivateKey
	err = k2.UnmarshalText([]byte(pvtkbytes))
	fmt.Println(err)
}
