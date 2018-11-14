package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	cli "github.com/jawher/mow.cli"
	"github.com/oneiro-ndev/ndaumath/pkg/key"
	"github.com/oneiro-ndev/ndaumath/pkg/signature"
)

// ktype should always be "", "PUB", or "PVT"
func keytype(ktype string) string {
	return strings.ToUpper(strings.TrimSpace(ktype)) + "KEY"
}

// TODO: doesn't work if we need to read multiple keys in one subcmd
func getKeySpec(ktype string) string {
	return fmt.Sprintf("(%s | --stdin)", keytype(ktype))
}

func getKeyClosure(cmd *cli.Cmd, ktype string, desc string) func() signature.Key {
	key := cmd.StringArg(keytype(ktype), "", desc)
	stdin := cmd.BoolOpt("stdin", false, "if set, read the key from stdin")

	return func() signature.Key {
		var keys string
		if stdin != nil && *stdin {
			in := bufio.NewScanner(os.Stdin)
			if !in.Scan() {
				check(errors.New("stdin selected but empty"))
			}
			check(in.Err())
			keys = in.Text()
		} else if key != nil && len(*key) > 0 {
			keys = *key
		} else {
			check(errors.New("no or multiple keys input--this should be unreachable"))
		}

		switch {
		case signature.MaybePrivate(keys):
			var pk signature.PrivateKey
			err := pk.UnmarshalText([]byte(keys))
			check(err)
			return signature.Key(pk)
		case signature.MaybePublic(keys):
			var pk signature.PublicKey
			err := pk.UnmarshalText([]byte(keys))
			check(err)
			return signature.Key(pk)
		default:
			check(errors.New("provided data is not an ndau key"))
		}

		// UNREACHABLE
		return signature.Key{}
	}
}

func getKeyClosureHD(cmd *cli.Cmd, ktype string, desc string) func() *key.ExtendedKey {
	keyi := cmd.StringArg(keytype(ktype), "", desc)
	stdin := cmd.BoolOpt("stdin", false, "if set, read the key from stdin")

	return func() *key.ExtendedKey {
		var keys string
		if stdin != nil && *stdin {
			in := bufio.NewScanner(os.Stdin)
			if !in.Scan() {
				check(errors.New("stdin selected but empty"))
			}
			check(in.Err())
			keys = in.Text()
		} else if keyi != nil && len(*keyi) > 0 {
			keys = *keyi
		} else {
			check(errors.New("no or multiple keys input--this should be unreachable"))
		}

		ek := new(key.ExtendedKey)
		err := ek.UnmarshalText([]byte(keys))
		check(err)
		return ek
	}
}
