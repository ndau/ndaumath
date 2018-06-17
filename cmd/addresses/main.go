package main

//go:generate gopherjs build --minify

// This is an experiment to see if gopherjs can reasonably generate js code from go source
// so that we can have a single-source solution for keys and addresses.

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oneiro-ndev/ndaumath/pkg/address"
)

func main() {
	// validate accepts an ndau key (like those returned by generate) and either
	// returns the same key, or an empty string if the key was not valid
	// function validate(key: string): string
	js.Global.Set("validate", func(addr string) string {
		a, err := address.Validate(addr)
		if err != nil {
			return ""
		}
		return a.String()
	})

	// generate creates a key of the appropriate kind (which must be one of a, n, e, or x)
	// and uses data (which must be at least 32 bytes long) to generate a new ndau key
	// if anything goes wrong, the result is an empty string
	// function generate(kind: string, data: string) : string
	js.Global.Set("generate", func(kind string, data string) string {
		kinds := map[string]address.Kind{
			string(address.KindUser):      address.KindUser,
			string(address.KindNdau):      address.KindNdau,
			string(address.KindExchange):  address.KindExchange,
			string(address.KindEndowment): address.KindEndowment,
		}
		k, ok := kinds[kind]
		if !ok {
			return ""
		}
		a, err := address.Generate(k, []byte(data))
		if err != nil {
			return ""
		}
		return a.String()
	})
}
