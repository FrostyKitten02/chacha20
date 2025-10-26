package main

import (
	"ChaCha20/gui"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(6)
	gui.ShowGui()
}
