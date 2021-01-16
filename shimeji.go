package main

import (
	"shimeji/daemon"
	"shimeji/gui"
)

func main() {
	daemon.Init()
	gui.Start()
}
