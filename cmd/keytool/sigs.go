package main

import (
	"errors"
	"fmt"

	"github.com/oneiro-ndev/ndaumath/pkg/signature"

	cli "github.com/jawher/mow.cli"
)

func cmdSign(cmd *cli.Cmd) {
	cmd.Spec = fmt.Sprintf(
		"%s %s",
		getKeySpec("PVT"),
		getDataSpec(),
	)

	getKey := getKeyClosurePrivate(cmd, "PVT", "sign with this private key")
	getData := getDataClosure(cmd)

	cmd.Action = func() {
		key := getKey()
		data := getData()

		sig := key.Sign(data)
		sigb, err := sig.MarshalText()
		check(err)
		fmt.Println(string(sigb))
	}
}

func cmdVerify(cmd *cli.Cmd) {
	cmd.Spec = fmt.Sprintf(
		"[-v] %s SIGNATURE %s",
		getKeySpec("PUB"),
		getDataSpec(),
	)

	getKey := getKeyClosurePublic(cmd, "PUB", "verify with this public key")
	getData := getDataClosure(cmd)

	verbose := cmd.BoolOpt("v verbose", false, "indicate success or failure on stdout in addition to the return code")
	sigi := cmd.StringArg("SIGNATURE", "", "verify this signature")

	cmd.Action = func() {
		key := getKey()
		data := getData()

		if sigi == nil || len(*sigi) == 0 {
			check(errors.New("signature not specified"))
		}
		var sig signature.Signature
		err := sig.UnmarshalText([]byte(*sigi))
		check(err)

		v := false
		if verbose != nil && *verbose {
			v = true
		}

		if sig.Verify(data, key) {
			if v {
				fmt.Println("OK")
			}
		} else {
			if v {
				fmt.Println("NO MATCH")
			}
			cli.Exit(2)
		}
	}
}
