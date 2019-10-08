package keyaddr

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import (
	"encoding/base64"
	"strings"

	"github.com/oneiro-ndev/ndaumath/pkg/words"
)

// WordsFromBytes takes an array of bytes and converts it to a space-separated list of
// words that act as a mnemonic. A 16-byte input array will generate a list of 12 words.
func WordsFromBytes(lang string, data string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	sa, err := words.FromBytes(lang, b)
	if err != nil {
		return "", err
	}
	return strings.Join(sa, " "), nil
}

// WordsToBytes takes a space-separated list of words and generates the set of bytes
// from which it was generated (or an error). The bytes are encoded as a base64 string
// using standard base64 encoding, as defined in RFC 4648.
func WordsToBytes(lang string, w string) (string, error) {
	wordlist := strings.Split(w, " ")
	b, err := words.ToBytes(lang, wordlist)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// WordsFromPrefix accepts a language and a prefix string and returns a sorted, space-separated list
// of words that match the given prefix. max can be used to limit the size of the returned list
// (if max is 0 then all matches are returned, which could be up to 2K if the prefix is empty).
func WordsFromPrefix(lang string, prefix string, max int) string {
	return words.FromPrefix(lang, prefix, max)
}
