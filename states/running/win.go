package running

import (
	"bomberman/core"
	"bomberman/game"

	ray "github.com/gen2brain/raylib-go/raylib"
)

var (
	winText             *core.Text
	winMenuTexts        []*core.Text
	currentWinMenuIndex int
)

type Win struct{}

func (Win) OnEnter() {
	core.GetSound("stage-clear").Play()
	if winMenuTexts == nil {
		currentWinMenuIndex = 0
		winText = core.NewText("", ray.GetFontDefault(), ray.NewVector2(float32(ray.GetScreenWidth())/2, 100), 40, 5, ray.White)
		winMenuTexts = append(winMenuTexts, core.NewText("RESTART", ray.GetFontDefault(), ray.NewVector2(float32(ray.GetScreenWidth())/2, 250), 32, 5, ray.Red))
		winMenuTexts = append(winMenuTexts, core.NewText("MAIN MENU", ray.GetFontDefault(), ray.NewVector2(float32(ray.GetScreenWidth())/2, 300), 32, 5, ray.White))
		winMenuTexts = append(winMenuTexts, core.NewText("EXIT", ray.GetFontDefault(), ray.NewVector2(float32(ray.GetScreenWidth())/2, 350), 32, 5, ray.White))
	}
	p1.Lock()
	p2.Lock()
	if p1.Wins > p2.Wins {
		winText.SetText("player 1 won")
	} else {
		winText.SetText("player 2 won")
	}
	p1.Unlock()
	p2.Unlock()
	winText.Measure()
	currentWinMenuIndex = 0
	winText.Pos = ray.NewVector2(float32(ray.GetScreenWidth())/2, 25)
	for _, t := range winMenuTexts {
		t.Pos.X = float32(ray.GetScreenWidth()) / 2
	}
	changeWinMenuColor()
}

func (Win) OnExit() {}

func changeWinMenuColor() {
	for i, t := range winMenuTexts {
		if i == currentWinMenuIndex {
			t.Color = ray.Red
		} else {
			t.Color = ray.White
		}
	}
}

func (Win) OnWindowResized() {
	fitCamera()
	winText.Pos = ray.NewVector2(float32(ray.GetScreenWidth())/2, 25)
	for _, t := range winMenuTexts {
		t.Pos.X = float32(ray.GetScreenWidth()) / 2
	}
}

func (Win) Update() {
	if ray.IsKeyReleased(ray.KeyUp) {
		if currentWinMenuIndex == 0 {
			currentWinMenuIndex = len(winMenuTexts) - 1
		} else {
			currentWinMenuIndex -= 1
		}
		changeWinMenuColor()
	} else if ray.IsKeyReleased(ray.KeyDown) {
		currentWinMenuIndex = (currentWinMenuIndex + 1) % len(winMenuTexts)
		changeWinMenuColor()
	}

	if ray.IsKeyPressed(ray.KeyEnter) {
		switch currentWinMenuIndex {
		case 0:
			if game.LastState.Get() == game.OFFLINE_BATTLE {
				game.State.Change(game.BATTLE_MENU)
			} else {
				game.State.Change(game.ONLINE_MENU)
			}
		case 1:
			game.State.Change(game.MENU)
		case 2:
			game.State.Change(game.QUIT)
		}
	}
}

func (Win) Draw() {
	drawBattle()
	ray.DrawRectangle(0, 0, int32(ray.GetScreenWidth()), int32(ray.GetScreenHeight()), ray.NewColor(0, 0, 0, 200))
	winText.DrawCentered()
	for _, t := range winMenuTexts {
		t.DrawCentered()
	}
}
