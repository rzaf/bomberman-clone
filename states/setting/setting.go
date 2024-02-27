package setting

import (
	"fmt"
	"github.com/rzaf/bomberman-clone/core"
	"github.com/rzaf/bomberman-clone/game"

	ray "github.com/gen2brain/raylib-go/raylib"
)

type Setting struct{}

var (
	isSelected       bool = false
	currentMenuIndex int
	audioText        *core.Text
	keysText         *core.Text
	controllerText   *core.Text
	backText         *core.Text
	tlX              int32
	gameCamera       ray.Camera2D = ray.NewCamera2D(ray.NewVector2(0, 0), ray.NewVector2(0, 0), 0, 1)
)

func fitCamera() {
	gameCamera.Offset.X = float32(ray.GetScreenWidth() / 2)
	gameCamera.Offset.Y = float32(ray.GetScreenHeight() / 2)

	screenRatio := float32(ray.GetScreenHeight()) / float32(ray.GetScreenWidth())
	i, j := game.TileManager.Length()
	gridRatio := float32(game.TILE_LENGTH*j) / float32(game.TILE_LENGTH*i)

	if screenRatio < gridRatio {
		gameCamera.Zoom = float32(ray.GetScreenHeight()) / float32(game.TILE_LENGTH*j)
	} else {
		gameCamera.Zoom = float32(ray.GetScreenWidth()) / float32(game.TILE_LENGTH*i)
	}
}

func changeTexts() {
	keysText.Color = ray.White
	controllerText.Color = ray.White
	audioText.Color = ray.White
	backText.Color = ray.White

	keysText.FontSize = 30
	controllerText.FontSize = 30
	audioText.FontSize = 30
	backText.FontSize = 30

	if !isSelected {
		audioText.Pos.X = 40
		keysText.Pos.X = 40
		controllerText.Pos.X = 40
		backText.Pos.X = 40
	} else {
		audioText.Pos.X = 0
		keysText.Pos.X = 0
		controllerText.Pos.X = 0
		backText.Pos.X = 0
	}

	switch currentMenuIndex {
	case 0:
		audioText.Color = ray.Red
		audioText.FontSize = 45
		audioText.Pos.X += 20
	case 1:
		keysText.Color = ray.Red
		keysText.FontSize = 45
		keysText.Pos.X += 20
	case 2:
		controllerText.Color = ray.Red
		controllerText.FontSize = 45
		controllerText.Pos.X += 20
	case 3:
		backText.Color = ray.Red
		backText.FontSize = 45
		backText.Pos.X += 20
	}
	keysText.Measure()
	controllerText.Measure()
	audioText.Measure()
	backText.Measure()
}

func (Setting) OnWindowResized() {
	tlX = int32(ray.GetScreenWidth()) / 4
	fitCamera()
}

func (Setting) OnEnter() {
	fmt.Println("Entering setting state!")
	if keysText == nil {
		audioText = core.NewText("sounds", ray.GetFontDefault(), ray.NewVector2(0, 50), 30, 4, ray.White)
		keysText = core.NewText("keyboard", ray.GetFontDefault(), ray.NewVector2(0, 100), 30, 4, ray.White)
		controllerText = core.NewText("controller", ray.GetFontDefault(), ray.NewVector2(0, 150), 30, 4, ray.White)
		backText = core.NewText("back", ray.GetFontDefault(), ray.NewVector2(0, 200), 30, 4, ray.White)
	}
	i, j := game.TileManager.Length()
	gameCamera.Target.X = float32(game.TILE_LENGTH*i) / 2
	gameCamera.Target.Y = float32(game.TILE_LENGTH*j) / 2
	fitCamera()
	tlX = int32(ray.GetScreenWidth()) / 4
	currentMenuIndex = 0
	changeTexts()
	isSelected = false
	isWaiting = false
	masterV = game.MasterVolume
	musicV = game.MusicVolume
	effectV = game.EffectVolume
	currentI = 0
	currentJ = 0
	currentGamePadId = 0
}

func (Setting) OnExit() {}

func (Setting) Update() {
	if !isSelected {
		if game.IsKeyPressed("accept") {
			switch currentMenuIndex {
			case 0:
				isSelected = true
			case 1:
				currentI = 0
				currentJ = 0
				isSelected = true
			case 2:
				currentI = 2
				currentJ = 0
				isSelected = true
			case 3:
				game.State.Change(game.MENU)
			}
			changeTexts()
			return
		}

		if game.IsKeyPressed("back") {
			game.State.Change(game.MENU)
		}
		if game.IsKeyPressed("p1-Down") {
			currentMenuIndex = (currentMenuIndex + 1) % 4
			changeTexts()
		} else if game.IsKeyPressed("p1-Up") {
			if currentMenuIndex == 0 {
				currentMenuIndex = 3
			} else {
				currentMenuIndex -= 1
			}
			changeTexts()
		}
	} else {
		switch currentMenuIndex {
		case 0:
			updateAudio()
		case 1:
			updateKeyboardMapping()
		case 2:
			updateControllerMapping()
		}
	}
}

func (Setting) Draw() {
	ray.BeginMode2D(gameCamera)
	game.TileManager.Draw()
	ray.EndMode2D()

	ray.DrawRectangle(0, 0, int32(game.Width), int32(game.Height), ray.NewColor(0, 0, 0, 150))

	audioText.Draw()
	keysText.Draw()
	controllerText.Draw()
	backText.Draw()

	switch currentMenuIndex {
	case 0:
		drawAudio()
	case 1:
		drawKeyboardMapping()
	case 2:
		drawControllerMapping()
	}
}
