package main

import (
	"fmt"
	"github.com/rzaf/bomberman-clone/game"
	"github.com/rzaf/bomberman-clone/states/editor"
	"github.com/rzaf/bomberman-clone/states/menu"
	"github.com/rzaf/bomberman-clone/states/running"
	"github.com/rzaf/bomberman-clone/states/setting"

	ray "github.com/gen2brain/raylib-go/raylib"
)

var (
	stateI game.StateI
	gs     game.GameState
)

func changeState() {
	game.LastState.Change(gs)
	fmt.Printf("state changed from %v to %v \n", game.LastState.Get(), game.State.Get())
	stateI.OnExit()
	switch game.State.Get() {
	case game.OFFLINE_BATTLE:
		stateI = running.OfflineBattle{}
	case game.BATTLE_MENU:
		stateI = running.BattleMenu{}
	case game.PAUSED:
		stateI = running.Paused{}
	case game.WIN:
		stateI = running.Win{}
	case game.ONLINE_BATTLE:
		stateI = running.OnlineBattle{}
	case game.ONLINE_MENU:
		stateI = running.OnlineMenu{}
	case game.MENU:
		stateI = menu.MainMenu{}
	case game.SETTING:
		stateI = setting.Setting{}
	case game.EDITOR:
		stateI = editor.Editor{}
	}
	stateI.OnEnter()
	gs.Change(game.State.Get())
}

func main() {
	ray.SetTraceLogLevel(ray.LogWarning)
	ray.InitAudioDevice()
	running.LoadSounds()
	ray.SetTargetFPS(60)
	ray.InitWindow(int32(game.Width), int32(game.Height), "bomberman-clone")
	ray.SetWindowState(ray.FlagWindowResizable)
	ray.SetWindowSize(game.Width, game.Height)
	ray.SetExitKey(0)
	setting.Load()
	running.LoadTextures()
	g1 := game.MENU
	g2 := game.MENU
	game.LastState = &g1
	game.State = &g2
	gs = game.MENU
	stateI = menu.MainMenu{}
	stateI.OnEnter()
	for !ray.WindowShouldClose() && gs != game.QUIT {
		// ray.DrawFPS(0, 0)
		game.UpdateKeys()
		if gs != game.State.Get() {
			changeState()
		} else {
			if ray.IsWindowResized() {
				game.Width = ray.GetScreenWidth()
				game.Height = ray.GetScreenHeight()
				stateI.OnWindowResized()
			}
			stateI.Update()
		}
		ray.BeginDrawing()
		ray.ClearBackground(ray.Gray)
		if gs != game.State.Get() {
			continue
		}
		stateI.Draw()
		ray.EndDrawing()
	}
	running.Unload()
	ray.CloseWindow()
	ray.CloseAudioDevice()
}
