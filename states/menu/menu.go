package menu

import (
	"fmt"
	"github.com/rzaf/bomberman-clone/core"
	"github.com/rzaf/bomberman-clone/game"

	ray "github.com/gen2brain/raylib-go/raylib"
)

const (
	MENU_BATTLE uint8 = iota
	MENU_ONLINE
	MENU_EDITOR
	MENU_EXIT
	MENU_COUNT
)

var (
	currentMenu uint8 = MENU_BATTLE

	battleText *core.Text
	onlineText *core.Text
	editorText *core.Text
	exitText   *core.Text
	gameCamera ray.Camera2D = ray.NewCamera2D(ray.NewVector2(0, 0), ray.NewVector2(0, 0), 0, 1)
)

type MainMenu struct{}

func (MainMenu) OnEnter() {
	fmt.Println("*** entering main menu state")
	core.LoadTexture("assets/characters.png", "anims", ray.NewRectangle(0, 0, 100, 100))
	core.LoadTexture("assets/tiles.png", "tiles", ray.NewRectangle(0, 0, 1060, 680))
	game.TileManager.GameMap = &game.GameMap{}
	game.TileManager.LoadFromFile("assets/menuMap.txt")
	if battleText == nil {
		battleText = core.NewText("OFFLINE BATTLE", ray.GetFontDefault(), ray.NewVector2(0, 0), 52, 4, ray.Red)
		onlineText = core.NewText("ONLINE BATTLE", ray.GetFontDefault(), ray.NewVector2(0, 0), 42, 4, ray.White)
		editorText = core.NewText("LEVEL EDITOR", ray.GetFontDefault(), ray.NewVector2(0, 0), 42, 4, ray.White)
		exitText = core.NewText("QUIT", ray.GetFontDefault(), ray.NewVector2(0, 0), 42, 4, ray.White)
	}
	changeTexts()
	resizeTexts()

	i, j := game.TileManager.Length()
	gameCamera.Target.X = float32(game.TILE_LENGTH*i) / 2
	gameCamera.Target.Y = float32(game.TILE_LENGTH*j) / 2
	fitCamera()
}

func (MainMenu) OnExit() {}

func (MainMenu) OnWindowResized() {
	resizeTexts()
	fitCamera()
}

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
	battleText.Color = ray.White
	onlineText.Color = ray.White
	editorText.Color = ray.White
	exitText.Color = ray.White
	battleText.FontSize = 42
	onlineText.FontSize = 42
	editorText.FontSize = 42
	exitText.FontSize = 42
	switch currentMenu {
	case MENU_BATTLE:
		battleText.Color = ray.Red
		battleText.FontSize = 52
	case MENU_EDITOR:
		editorText.Color = ray.Red
		editorText.FontSize = 52
	case MENU_ONLINE:
		onlineText.Color = ray.Red
		onlineText.FontSize = 52
	case MENU_EXIT:
		exitText.Color = ray.Red
		exitText.FontSize = 52
	}
	battleText.Measure()
	onlineText.Measure()
	editorText.Measure()
	exitText.Measure()
}

func resizeTexts() {
	battleText.Pos.X = float32(ray.GetScreenWidth()) / 2
	onlineText.Pos.X = float32(ray.GetScreenWidth()) / 2
	editorText.Pos.X = float32(ray.GetScreenWidth()) / 2
	exitText.Pos.X = float32(ray.GetScreenWidth()) / 2
	battleText.Pos.Y = float32(ray.GetScreenHeight())/2 - 200
	onlineText.Pos.Y = float32(ray.GetScreenHeight())/2 - 100
	editorText.Pos.Y = float32(ray.GetScreenHeight()) / 2
	exitText.Pos.Y = float32(ray.GetScreenHeight())/2 + 100
}

func (MainMenu) Update() {

	if ray.IsKeyPressed(ray.KeyUp) {
		if currentMenu == 0 {
			currentMenu = MENU_COUNT - 1
		} else {
			currentMenu -= 1
		}
		changeTexts()
	} else if ray.IsKeyPressed(ray.KeyDown) {
		currentMenu = (currentMenu + 1) % MENU_COUNT
		changeTexts()
	}

	if ray.IsKeyPressed(ray.KeyEnter) {
		switch currentMenu {
		case MENU_BATTLE:
			game.State.Change(game.BATTLE_MENU)
		case MENU_EDITOR:
			game.State.Change(game.EDITOR)
		case MENU_ONLINE:
			game.State.Change(game.ONLINE_MENU)
		case MENU_EXIT:
			game.State.Change(game.QUIT)
		}
	}
}

func (MainMenu) Draw() {
	ray.BeginMode2D(gameCamera)
	game.TileManager.Draw()
	ray.EndMode2D()

	ray.DrawRectangle(0, 0, int32(ray.GetScreenWidth()), int32(ray.GetScreenHeight()), ray.NewColor(0, 0, 0, 150))
	battleText.DrawCentered()
	editorText.DrawCentered()
	onlineText.DrawCentered()
	exitText.DrawCentered()
}
