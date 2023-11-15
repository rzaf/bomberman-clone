package running

import (
	"bomberman/core"
	"bomberman/game"
	"fmt"
	"log"
	"os"

	ray "github.com/gen2brain/raylib-go/raylib"
)

const (
	s_MENU_MAP uint8 = iota
	s_MENU_N_ROUNDS
	s_MENU_TIME
	s_MENU_COUNT
)

var (
	mapNames []*core.Text
	gameMaps []*game.GameMap

	currentMapIndex  int
	menuState        uint8
	rounds           int = 1
	roundText        *core.Text
	roundTimeSeconds int = 60
	timeText         *core.Text
)

type BattleMenu struct{}

func (BattleMenu) OnEnter() {
	mapNames = nil
	gameMaps = nil
	dirs, err := os.ReadDir("assets/maps/")
	if err != nil {
		log.Fatalln(err)
	}
	for _, dir := range dirs {
		newMap := game.GameMap{}
		newMap.LoadFromFile("assets/maps/" + dir.Name())
		newMap.Name = dir.Name()
		gameMaps = append(gameMaps, &newMap)
		mapNames = append(mapNames, core.NewText(dir.Name(), ray.GetFontDefault(), ray.NewVector2(0, 0), 22, 3, ray.Red))
	}
	changeMap(0)
	menuState = s_MENU_MAP
	roundText = core.NewText(fmt.Sprintf("%d", rounds), ray.GetFontDefault(), ray.NewVector2(0, 0), 22, 2, ray.White)
	timeText = core.NewText(fmt.Sprintf("%d:%02d", roundTimeSeconds/60, roundTimeSeconds%60), ray.GetFontDefault(), ray.NewVector2(0, 0), 22, 2, ray.White)

}

func (BattleMenu) OnExit() {}

func changeMap(index int) {
	currentMapIndex = index
	game.TileManager.GameMap = gameMaps[currentMapIndex]
	i, j := game.TileManager.Length()
	gameCamera.Target.X = float32(game.TILE_LENGTH*i) / 2
	gameCamera.Target.Y = float32(game.TILE_LENGTH*j) / 2
	fitCamera()
}

func changeMenuColor() {
	switch menuState {
	case s_MENU_MAP:
		for _, n := range mapNames {
			n.Color = ray.Red
		}
		roundText.Color = ray.White
		timeText.Color = ray.White
	case s_MENU_N_ROUNDS:
		for _, n := range mapNames {
			n.Color = ray.White
		}
		roundText.Color = ray.Red
		timeText.Color = ray.White
	case s_MENU_TIME:
		for _, n := range mapNames {
			n.Color = ray.White
		}
		roundText.Color = ray.White
		timeText.Color = ray.Red
	}
}

func (BattleMenu) OnWindowResized() {
	fitCamera()
}

func (BattleMenu) Update() {
	if ray.IsKeyPressed(ray.KeyEnter) {
		loadLevel(mapNames[currentMapIndex].Text)
		fmt.Println(game.LastState, "is last state online_menu:", game.LastState.Get() == game.ONLINE_MENU)
		if game.LastState.Get() == game.ONLINE_MENU && isHost {
			isHostWaiting = true
			game.State.Change(game.ONLINE_MENU)
		} else {
			game.State.Change(game.OFFLINE_BATTLE)
		}
		return
	}
	if ray.IsKeyPressed(ray.KeyUp) {
		if menuState == 0 {
			menuState = s_MENU_COUNT - 1
		} else {
			menuState -= 1
		}
		changeMenuColor()
	} else if ray.IsKeyPressed(ray.KeyDown) {
		menuState = (menuState + 1) % s_MENU_COUNT
		changeMenuColor()
	}

	if ray.IsKeyPressed(ray.KeyRight) {
		switch menuState {
		case s_MENU_MAP:
			changeMap((currentMapIndex + 1) % len(gameMaps))
		case s_MENU_N_ROUNDS:
			if rounds < 5 {
				rounds += 1
			}
			roundText.Text = fmt.Sprintf("%d", rounds)
		case s_MENU_TIME:
			if roundTimeSeconds < 300 {
				roundTimeSeconds += 1
			}
			timeText.Text = fmt.Sprintf("%d:%02d", roundTimeSeconds/60, roundTimeSeconds%60)
		}
	} else if ray.IsKeyPressed(ray.KeyLeft) {
		switch menuState {
		case s_MENU_MAP:
			i := currentMapIndex
			if i == 0 {
				i = len(gameMaps) - 1
			} else {
				i -= 1
			}
			changeMap(i)
		case s_MENU_N_ROUNDS:
			if rounds > 1 {
				rounds -= 1
			}
			roundText.Text = fmt.Sprintf("%d", rounds)
		case s_MENU_TIME:
			if roundTimeSeconds > 30 {
				roundTimeSeconds -= 1
			}
			timeText.Text = fmt.Sprintf("%d:%02d", roundTimeSeconds/60, roundTimeSeconds%60)
		}
	}
}

func (BattleMenu) Draw() {
	ray.BeginMode2D(gameCamera)
	game.TileManager.Draw()
	ray.EndMode2D()

	ray.DrawRectangle(int32(ray.GetScreenWidth())/2-200, int32(ray.GetScreenHeight())-250, 400, 250, ray.DarkGray)
	ray.DrawText("MAP:", int32(ray.GetScreenWidth())/2-180, int32(ray.GetScreenHeight())-200, 20, ray.Black)
	mapNames[currentMapIndex].DrawAt(ray.NewVector2(float32(ray.GetScreenWidth()/2), float32(ray.GetScreenHeight()-200)))
	ray.DrawText(fmt.Sprintf("%d/%d", currentMapIndex+1, len(gameMaps)), int32(ray.GetScreenWidth())/2+120, int32(ray.GetScreenHeight())-200, 20, ray.Black)

	ray.DrawText("ROUNDS:", int32(ray.GetScreenWidth())/2-180, int32(ray.GetScreenHeight())-150, 20, ray.Black)
	roundText.DrawAt(ray.NewVector2(float32(ray.GetScreenWidth())/2, float32(ray.GetScreenHeight())-150))

	ray.DrawText("TIME:", int32(ray.GetScreenWidth())/2-180, int32(ray.GetScreenHeight())-100, 20, ray.Black)
	timeText.DrawAt(ray.NewVector2(float32(ray.GetScreenWidth())/2, float32(ray.GetScreenHeight())-100))
}
