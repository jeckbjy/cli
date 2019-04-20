package cmds

import (
	"log"

	"github.com/jeckbjy/cli"
)

func onCmdGet(ctx *cli.Context) {
	log.Printf("process get\n")

	var flags struct {
		File string `cli:"file"`
	}

	ctx.Bind(&flags)

	log.Printf("file:%+v\n", flags.File)
}
