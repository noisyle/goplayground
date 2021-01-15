package main

import (
	"shimeji/daemon"
	"shimeji/gui"
)

func main() {
	d := daemon.Init()
	gui.Init(d)
}
