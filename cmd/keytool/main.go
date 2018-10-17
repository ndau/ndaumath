package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("keytool", "manipulate key strings on the command line")

	app.Command("hd", "manipulate HD keys", hd)

	app.Run(os.Args)
}

// hd subcommand
func hd(cmd *cli.Cmd) {
	cmd.Command("new", "create a new HD key", cmdHDNew)
	cmd.Command("public", "create a public key from supplied key", cmdHDPublic)
	cmd.Command("child", "create a child key derived from the supplied key", cmdHDChild)
	cmd.Command("convert", "convert an old-format key into the new format", cmdHDConvert)
	cmd.Command("truncate", "remove HD portions from an HD key leaving only the key itself", cmdHDTruncate)
}
