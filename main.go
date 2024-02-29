package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  int = 640
	screenHeight int = 480
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("color-bird")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
