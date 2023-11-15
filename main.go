package main

import (
	"bomberman/game"
	"bomberman/states/editor"
	"bomberman/states/menu"
	"bomberman/states/running"
	"fmt"

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
	ray.InitWindow(600, 600, "Bomberman")
	ray.SetWindowState(ray.FlagWindowResizable)
	ray.SetWindowSize(600, 600)
	ray.SetExitKey(0)
	ray.SetMasterVolume(0.5)
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
		if gs != game.State.Get() {
			changeState()
		} else {
			if ray.IsWindowResized() {
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
