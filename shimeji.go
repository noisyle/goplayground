package main

import (
	"shimeji/daemon"
	"shimeji/gui"
)

func main() {
	daemon.Start()
	gui.Start()
}
