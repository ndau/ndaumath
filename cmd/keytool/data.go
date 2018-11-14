package main

import (
	"encoding/base64"
	"io/ioutil"

	cli "github.com/jawher/mow.cli"
)

func getDataSpec() string {
	return "(--file=<path> | --b64=<base64-encoded data>)"
}

func getDataClosure(cmd *cli.Cmd) func() []byte {
	var (
		filei = cmd.StringOpt("f file", "", "path to file containing applicable data")
		b64i  = cmd.StringOpt("b b64", "", "base64-encoded data with no padding")
	)

	return func() []byte {
		switch {
		case filei != nil && len(*filei) > 0:
			data, err := ioutil.ReadFile(*filei)
			check(err)
			return data
		case b64i != nil && len(*b64i) > 0:
			data, err := base64.RawStdEncoding.DecodeString(*b64i)
			check(err)
			return data
		}
		// unreachable
		return []byte{}
	}
}
