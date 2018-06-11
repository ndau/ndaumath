package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/oneiro-ndev/signature/pkg/signature"

	arg "github.com/alexflint/go-arg"
)

type args struct {
	Generate   bool     `help:"generate a keypair from the input data"`
	Sign       bool     `help:"sign a block of data"`
	Validate   bool     `help:"validate a signature"`
	Hex        bool     `arg:"-x" help:"interpret input data as hex bytes"`
	Verbose    bool     `arg:"-v" help:"Be verbose"`
	OutputFile string   `arg:"-o" help:"Output filename"`
	Input      string   `arg:"-i" help:"Input filename"`
	Data       []string `arg:"positional"`
	Comment    string   `arg:"-c" help:"Comment for key files"`
	Keyfile    string   `arg:"-k" help:"Key filename prefix"`
}

func (args) Description() string {
	return `
	Generates keypairs compatible with ndau, and signs blocks of data using those keys.

	Examples:
	signtool --generate bytestream -k keyfile
	# treats bytestream as entropy and generates a keypair, writing them to keyfile and keyfile.pub
	# default keyfile is "key"
	# if there is no input data at all, reads from the system entropy source

	signtool --sign -k keyfile -i inputstream
	# reads data from inputstream, hashes it, and signs it with the private key from keyfile;
	# sends the signature to outputfile

	signtool --validate -k keyfile -s sigfile -i inputfile
	# reads data from inputstream and validates the signature (using keyfile.pub)
	`
}

// reads an input stream and extracts anything that looks like pairs of hex characters (bytes)
func readAsHex(in io.Reader) ([]byte, error) {
	all, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	pat := regexp.MustCompile("[0-9a-fA-F]{2}")
	// extract the possible hex values
	hexes := pat.FindAllString(string(all), -1)
	var output []byte
	// now convert them
	for _, h := range hexes {
		b, _ := strconv.ParseUint(h, 16, 8)
		output = append(output, byte(b))
	}
	return output, nil
}

const keyType = "ed25519ndau"

// wrapHex encodes b as hex, wrapping at column w
func wrapHex(b []byte, w int) string {
	enc := hex.EncodeToString(b)
	out := ""
	for i := 0; len(enc) > 0; i += w {
		n := len(enc)
		if n > w {
			n = w
		}
		out += enc[0:n] + "\n"
		enc = enc[n:]
	}
	return out
}

func writePrivateKey(filename string, pvt signature.PrivateKey, note string) error {
	const (
		hdr      = "---- Begin Private Key ----"
		ftr      = "---- End Private Key ----"
		format   = hdr + "\nKey Type: %[1]s\nNote: %[3]s\n\n%[2]s\n" + ftr + "\n"
		maxWidth = 64
	)

	b, err := pvt.Marshal()
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	out := wrapHex(b, maxWidth)

	_, err = fmt.Fprintf(f, format, keyType, out, note)
	return err
}

func writePublicKey(filename string, pub signature.PublicKey, note string) error {
	b, err := pub.Marshal()
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := hex.EncodeToString(b)
	_, err = fmt.Fprintf(f, "%s %s %s\n", keyType, enc, note)
	return err
}

func readPublicKey(filename string) (signature.PublicKey, error) {
	sig := signature.PublicKey{}
	f, err := os.Open(filename)
	if err != nil {
		return sig, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	sp := strings.SplitN(string(data), " ", 3)
	if len(sp) < 2 || sp[0] != keyType {
		return sig, errors.New("Not an ndau public key file")
	}
	b, err := hex.DecodeString(sp[1])
	if err != nil {
		return sig, err
	}
	err = sig.Unmarshal(b)
	return sig, err
}

func readPrivateKey(filename string) (signature.PrivateKey, error) {
	sig := signature.PrivateKey{}
	f, err := os.Open(filename)
	if err != nil {
		return sig, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	sp := strings.Split(string(data), "\n")
	// TODO: Validate key type
	if len(sp) < 4 {
		return sig, errors.New("Not an ndau private key file")
	}
	keydata := strings.Join(sp[3:len(sp)-2], "")
	b, err := hex.DecodeString(keydata)
	if err != nil {
		return sig, err
	}
	err = sig.Unmarshal(b)
	return sig, err
}

func main() {
	var args args
	args.Keyfile = "key"
	arg.MustParse(&args)

	// figure out where we get our input stream from
	var in io.Reader
	in = strings.NewReader(strings.Join(args.Data, " "))
	if args.Input != "" {
		if args.Input == "-" {
			in = os.Stdin
		} else {
			f, err := os.Open(args.Input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			defer f.Close()
			in = f
		}
	}

	var data []byte
	var err error
	if args.Hex {
		data, err = readAsHex(in)
	} else {
		data, err = ioutil.ReadAll(in)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	if args.Verbose {
		fmt.Printf("input data:\n%s\n", hex.Dump(data))
	}

	switch {
	case args.Generate:
		// we're creating a keypair
		var r io.Reader
		if len(data) > 0 {
			r = bytes.NewReader(data)
		}
		pub, pvt, err := signature.Generate(signature.Ed25519, r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s generating key\n", err)
			os.Exit(1)
		}
		err = writePublicKey(args.Keyfile+".pub", pub, args.Comment)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s writing public key\n", err)
			os.Exit(1)
		}
		err = writePrivateKey(args.Keyfile, pvt, args.Comment)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s writing private key\n", err)
			os.Exit(1)
		}
	case args.Sign:
		// we're generating a signature so we need the private key
		if len(data) == 0 {
			fmt.Fprintf(os.Stderr, "we need data to sign", err)
			os.Exit(1)
		}

		pvt, err := readPrivateKey(args.Keyfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s reading private key '%s'\n", err, args.Keyfile)
			os.Exit(1)
		}
		sig := pvt.Sign(data)
		b, err := sig.Marshal()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s marshalling signature\n", err)
			os.Exit(1)
		}

		f := os.Stdout
		if args.OutputFile != "" {
			f, err = os.Create(args.OutputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s creating output file\n", err)
				os.Exit(1)
			}
		}
		f.WriteString(wrapHex(b, 80))
		f.Close()

	case args.Validate:
	}
	os.Exit(0)

}
