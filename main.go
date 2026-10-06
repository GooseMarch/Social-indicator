package main

import (
	"os"
	"lab2_bubble/cli"
	"lab2_bubble/tui"
)

func main() {
	if len(os.Args) < 2 {
		tui.Run()
		return
	}
	cli.RunCLI(os.Args[1:])
}