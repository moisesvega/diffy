package cmd

import (
	"fmt"
	"runtime"
)

type VersionCmd struct {
	version string
}

func (c *VersionCmd) Run() error {
	fmt.Printf("diffy version %s %s/%s\n", c.version, runtime.GOOS, runtime.GOARCH)
	return nil
}
