package running

import (
	"github.com/rzaf/bomberman-clone/core"
	"github.com/rzaf/bomberman-clone/game"

	ray "github.com/gen2brain/raylib-go/raylib"
)

var (
	winText *core.Text
	winMenu *game.Menu
)

type Win struct{}

func (Win) OnEnter() {
	core.GetSound("stage-clear").Play()
	if winMenu == nil {
		winText = core.NewText("", ray.GetFontDefault(), ray.NewVector2(float32(ray.GetScreenWidth())/2, 100), 40, 5, ray.White)
		restartItem := game.NewMenuItem("RESTART", func() {
			switch game.LastState.Get() {
			case game.ONLINE_BATTLE:
				game.State.Change(game.ONLINE_MENU)
			case game.PAUSED: // bug
				game.State.Change(game.MENU)
			default:
				game.State.Change(game.BATTLE_MENU)
			}
		})
		mainMenuItem := game.NewMenuItem("MAIN MENU", func() {
			game.State.Change(game.MENU)
		})
		exitItem := game.NewMenuItem("EXIT", func() {
			game.State.Change(game.QUIT)
		})
		winMenu = game.NewMenu(restartItem, mainMenuItem, exitItem)
		winMenu.FontSize = 25
		winMenu.SelectedFontSize = 25
		winMenu.Padding = 50
		winMenu.Pos.Y = 200
	}
	winMenu.Pos.X = float32(game.Width) / 2
	winMenu.SetIndex(0)
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
}

func (Win) OnExit() {}

func (Win) OnWindowResized() {
	fitCamera()
	winText.Pos = ray.NewVector2(float32(game.Width)/2, 25)
	winText.Measure()
	winMenu.Pos.X = float32(game.Width) / 2
	winMenu.Refresh()
}

func (Win) Update() {
	winMenu.Update()
}

func (Win) Draw() {
	drawBattle()
	ray.DrawRectangle(0, 0, int32(ray.GetScreenWidth()), int32(ray.GetScreenHeight()), ray.NewColor(0, 0, 0, 200))
	winText.DrawCentered()
	winMenu.Draw()
}
