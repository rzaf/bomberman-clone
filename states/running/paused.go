package running

import (
	"bomberman/core"
	"bomberman/game"

	ray "github.com/gen2brain/raylib-go/raylib"
)

var (
	pausedText            *core.Text
	pauseMenuTexts        []*core.Text
	currentPauseMenuIndex int
)

type Paused struct{}

func (Paused) OnEnter() {
	if pauseMenuTexts == nil {
		currentPauseMenuIndex = 0
		pausedText = core.NewText("PAUSED", ray.GetFontDefault(), ray.NewVector2(float32(ray.GetScreenWidth())/2, 25), 25, 3, ray.White)
		pauseMenuTexts = append(pauseMenuTexts, core.NewText("RESUME", ray.GetFontDefault(), ray.NewVector2(float32(ray.GetScreenWidth())/2, 200), 32, 5, ray.Red))
		pauseMenuTexts = append(pauseMenuTexts, core.NewText("RESTART", ray.GetFontDefault(), ray.NewVector2(float32(ray.GetScreenWidth())/2, 250), 32, 5, ray.White))
		pauseMenuTexts = append(pauseMenuTexts, core.NewText("MAIN MENU", ray.GetFontDefault(), ray.NewVector2(float32(ray.GetScreenWidth())/2, 300), 32, 5, ray.White))
		pauseMenuTexts = append(pauseMenuTexts, core.NewText("EXIT", ray.GetFontDefault(), ray.NewVector2(float32(ray.GetScreenWidth())/2, 350), 32, 5, ray.White))
	}

	pausedText.Pos = ray.NewVector2(float32(ray.GetScreenWidth())/2, 25)
	for _, t := range pauseMenuTexts {
		t.Pos.X = float32(ray.GetScreenWidth()) / 2
	}
}

func (Paused) OnExit() {
	if game.State.Get() != game.OFFLINE_BATTLE && game.State.Get() != game.ONLINE_BATTLE {
		core.GetSound("battle1").Stop()
	}
	if game.State.Get() == game.ONLINE_BATTLE {
		return
	}
	if isHost {
		DisconnectServer()
	} else {
		CloseClient()
	}
}

func changePauseMenuColor() {
	for i, t := range pauseMenuTexts {
		if i == currentPauseMenuIndex {
			t.Color = ray.Red
		} else {
			t.Color = ray.White
		}
	}
}

func (Paused) OnWindowResized() {
	fitCamera()
	pausedText.Pos = ray.NewVector2(float32(ray.GetScreenWidth())/2, 25)
	for _, t := range pauseMenuTexts {
		t.Pos.X = float32(ray.GetScreenWidth()) / 2
	}
}

func (Paused) Update() {
	if ray.IsKeyPressed(ray.KeyEscape) {
		game.State.Change(game.LastState.Get())
		return
	}
	if game.LastState.Get() == game.ONLINE_BATTLE {
		updateBattle()
	}

	if ray.IsKeyReleased(ray.KeyUp) {
		if currentPauseMenuIndex == 0 {
			currentPauseMenuIndex = len(pauseMenuTexts) - 1
		} else {
			currentPauseMenuIndex -= 1
		}
		changePauseMenuColor()
	} else if ray.IsKeyReleased(ray.KeyDown) {
		currentPauseMenuIndex = (currentPauseMenuIndex + 1) % len(pauseMenuTexts)
		changePauseMenuColor()
	}
	if ray.IsKeyPressed(ray.KeyEnter) {
		switch currentPauseMenuIndex {
		case 0:
			game.State.Change(game.LastState.Get())
		case 1:
			if game.LastState.Get() == game.OFFLINE_BATTLE {
				game.State.Change(game.BATTLE_MENU)
			} else {
				game.State.Change(game.ONLINE_MENU)
			}
		case 2:
			game.State.Change(game.MENU)
		case 3:
			game.State.Change(game.QUIT)
		}
	}
}

func (Paused) Draw() {
	drawBattle()
	ray.DrawRectangle(0, 0, int32(ray.GetScreenWidth()), int32(ray.GetScreenHeight()), ray.NewColor(0, 0, 0, 200))
	pausedText.DrawCentered()
	for _, t := range pauseMenuTexts {
		t.DrawCentered()
	}
}
