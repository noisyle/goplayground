package main

import (
	"shimeji/daemon"
	"shimeji/gui"
)

func main() {
	d := daemon.Start()
	gui.Start(d)
}
