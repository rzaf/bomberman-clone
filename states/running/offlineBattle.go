package running

import (
	"bomberman/core"
	"bomberman/game"
	"fmt"

	ray "github.com/gen2brain/raylib-go/raylib"
)

type OfflineBattle struct{}

func (OfflineBattle) OnExit() {
	if game.State.Get() != game.PAUSED {
		core.GetSound("battle1").Stop()
	}
}

func (OfflineBattle) OnEnter() {
	fmt.Println("*** entering OfflineBattle state")
	core.LoadTexture("assets/characters.png", "anims", ray.NewRectangle(0, 0, 100, 100))
	core.LoadTexture("assets/tiles.png", "tiles", ray.NewRectangle(0, 0, 1060, 680))
	loadUi()
	if !core.GetSound("battle1").IsPlaying() {
		core.GetSound("battle1").Play()
	}
	game.LoadDefaultKeys()
}

func (OfflineBattle) OnWindowResized() {
	fitCamera()
}

func (OfflineBattle) Update() {
	updateBattle()
}

func (OfflineBattle) Draw() {
	drawBattle()
}
