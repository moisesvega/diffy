package main

import (
	"os"

	"github.com/moisesvega/diffy/internal/cmd"
)

func main() {
	kong := cmd.Main()
	kctx, err := kong.Parse(os.Args[1:])
	kong.FatalIfErrorf(err)
	err = kctx.Run()
	kong.FatalIfErrorf(err)
}
