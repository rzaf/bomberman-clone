package running

import (
	"github.com/rzaf/bomberman-clone/core"
	"github.com/rzaf/bomberman-clone/game"

	ray "github.com/gen2brain/raylib-go/raylib"
)

var (
	pausedText *core.Text
	pauseMenu  *game.Menu
)

type Paused struct{}

func (Paused) OnEnter() {
	if pauseMenu == nil {
		pausedText = core.NewText("PAUSED", ray.GetFontDefault(), ray.NewVector2(float32(ray.GetScreenWidth())/2, 25), 25, 3, ray.White)
		resumeItem := game.NewMenuItem("RESUME", func() {
			game.State.Change(game.LastState.Get())
		})
		restartItem := game.NewMenuItem("RESTART", func() {
			if game.LastState.Get() == game.OFFLINE_BATTLE {
				game.State.Change(game.BATTLE_MENU)
			} else {
				game.State.Change(game.ONLINE_MENU)
			}
		})
		mainMenuItem := game.NewMenuItem("MAIN MENU", func() {
			game.State.Change(game.MENU)
		})
		quitItem := game.NewMenuItem("EXIT", func() {
			game.State.Change(game.QUIT)
		})
		pauseMenu = game.NewMenu(resumeItem, restartItem, mainMenuItem, quitItem)
		pauseMenu.FontSize = 25
		pauseMenu.SelectedFontSize = 25
		pauseMenu.Padding = 50
		pauseMenu.Pos.Y = 200
	}

	pauseMenu.Pos.X = float32(game.Width) / 2
	pauseMenu.SetIndex(0)

	pausedText.Pos = ray.NewVector2(float32(ray.GetScreenWidth())/2, 25)
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

func (Paused) OnWindowResized() {
	fitCamera()
	pausedText.Pos.X = float32(game.Width) / 2
	pausedText.Measure()
	pauseMenu.Pos.X = float32(game.Width) / 2
	pauseMenu.Refresh()
}

func (Paused) Update() {
	if game.IsKeyPressed("pause") {
		game.State.Change(game.LastState.Get())
		return
	}
	if game.LastState.Get() == game.ONLINE_BATTLE {
		updateBattle()
	}

	pauseMenu.Update()
}

func (Paused) Draw() {
	drawBattle()
	ray.DrawRectangle(0, 0, int32(ray.GetScreenWidth()), int32(ray.GetScreenHeight()), ray.NewColor(0, 0, 0, 200))
	pausedText.DrawCentered()
	pauseMenu.Draw()
}
