package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/oneiro-ndev/ndaumath/pkg/key"
)

func hdstr(k key.ExtendedKey) string {
	text, err := k.MarshalText()
	check(err)
	return string(text)
}

func hdparse(s string) *key.ExtendedKey {
	k := new(key.ExtendedKey)
	err := k.UnmarshalText([]byte(s))
	check(err)
	return k
}

func cmdHDNew(cmd *cli.Cmd) {
	cmd.Action = func() {
		seed, err := key.GenerateSeed(key.RecommendedSeedLen)
		check(err)
		k, err := key.NewMaster(seed)
		check(err)
		fmt.Println(hdstr(*k))
	}
}

func cmdHDPublic(cmd *cli.Cmd) {
	pvtS := cmd.StringArg("PRIVATE", "", "private key from which to make a public key")
	cmd.Action = func() {
		pvt := hdparse(*pvtS)
		pub, err := pvt.Public()
		check(err)
		fmt.Println(hdstr(*pub))
	}
}

func cmdHDChild(cmd *cli.Cmd) {
	keyS := cmd.StringArg("KEY", "", "key from which to derive a child")
	pathS := cmd.StringArg("PATH", "", "derivation path for child key")
	cmd.Action = func() {
		key := hdparse(*keyS)
		key, err := key.DeriveFrom("/", *pathS)
		check(err)
		fmt.Println(hdstr(*key))
	}
}

func cmdHDConvert(cmd *cli.Cmd) {
	keyS := cmd.StringArg("KEY", "", "old-format key to convert")

	cmd.Action = func() {
		k, err := key.FromOldSerialization(*keyS)
		check(err)
		fmt.Println(hdstr(*k))
	}
}

func cmdHDTruncate(cmd *cli.Cmd) {
	keyS := cmd.StringArg("KEY", "", "key to truncate")

	cmd.Action = func() {
		key := hdparse(*keyS)
		var keyB []byte
		var err error
		if key.IsPrivate() {
			skey, err := key.SPrivKey()
			check(err)
			skey.Truncate()
			keyB, err = skey.MarshalText()
		} else {
			skey, err := key.SPubKey()
			check(err)
			skey.Truncate()
			keyB, err = skey.MarshalText()
		}
		check(err)
		fmt.Println(string(keyB))
	}
}
