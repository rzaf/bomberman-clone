package running

import (
	"github.com/rzaf/bomberman-clone/core"
	"github.com/rzaf/bomberman-clone/game"
	// "github.com/rzaf/bomberman-clone/pb"
	// "sync"
	"fmt"

	ray "github.com/gen2brain/raylib-go/raylib"
)

type OnlineBattle struct{}

func (OnlineBattle) OnEnter() {
	fmt.Println("*** entering onlineBattle state")
	core.LoadTexture("assets/characters.png", "anims", ray.NewRectangle(0, 0, 100, 100))
	core.LoadTexture("assets/tiles.png", "tiles", ray.NewRectangle(0, 0, 1060, 680))
	loadUi()
	if !core.GetSound("battle1").IsPlaying() {
		core.GetSound("battle1").Play()
	}

	if !isHost {
		game.DropUpgrades = false
	} else {
		game.SaveUpgrades = true
	}
}

func (OnlineBattle) OnExit() {
	if game.State.Get() == game.PAUSED {
		return
	}

	core.GetSound("battle1").Stop()
	if isHost {
		DisconnectServer()
	} else {
		CloseClient()
	}
}

func (OnlineBattle) OnWindowResized() {
	fitCamera()
}

func (OnlineBattle) Update() {
	updateBattle()
}

func (OnlineBattle) Draw() {
	drawBattle()
}
