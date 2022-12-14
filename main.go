package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/iancanderson/spookypaths/game"
	"github.com/iancanderson/spookypaths/game/config"
)

func main() {
	ebiten.SetWindowTitle("Spookypaths")

	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowSizeLimits(300, 200, -1, -1)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
