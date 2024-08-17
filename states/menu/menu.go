package menu

import (
	"fmt"
	"github.com/rzaf/bomberman-clone/core"
	"github.com/rzaf/bomberman-clone/game"

	ray "github.com/gen2brain/raylib-go/raylib"
)

var (
	menu              *game.Menu
	gameCamera        ray.Camera2D = ray.NewCamera2D(ray.NewVector2(0, 0), ray.NewVector2(0, 0), 0, 1)
	githubProfileLink *core.Button
	gitVersion        *core.Text
)

type MainMenu struct{}

func (m MainMenu) OnEnter() {
	fmt.Println("*** entering main menu state")
	core.LoadTexture("assets/characters.png", "anims", ray.NewRectangle(0, 0, 100, 100))
	core.LoadTexture("assets/tiles.png", "tiles", ray.NewRectangle(0, 0, 1060, 680))
	game.TileManager.GameMap = &game.GameMap{}
	game.TileManager.LoadFromFile("assets/menuMap.txt")
	if githubProfileLink == nil {
		gitVersion = core.NewText("v "+game.VersionString, ray.GetFontDefault(), ray.NewVector2(float32(game.Width)-70, float32(game.Height)-20), 15, 2, ray.White)
		t2 := core.NewText("github.com/rzaf/bomberman-clone", ray.GetFontDefault(), ray.NewVector2(90, float32(game.Height)-20), 15, 2, ray.White)
		githubProfileLink = core.NewTextButton(
			t2,
			ray.NewRectangle(t2.Pos.X, t2.Pos.Y, t2.Size.X, t2.Size.Y),
		)

		githubProfileLink.OnClick = func() {
			ray.OpenURL("https://github.com/rzaf/bomberman-clone/")
		}
	}
	if menu == nil {
		battleItem := game.NewMenuItem("OFFLINE BATTLE", func() {
			game.State.Change(game.BATTLE_MENU)
		})
		onlineBattleItem := game.NewMenuItem("ONLINE BATTLE", func() {
			game.State.Change(game.ONLINE_MENU)
		})
		editorItem := game.NewMenuItem("LEVEL EDITOR", func() {
			game.State.Change(game.EDITOR)
		})
		settingItem := game.NewMenuItem("SETTING", func() {
			game.State.Change(game.SETTING)
		})
		quitItem := game.NewMenuItem("QUIT", func() {
			game.State.Change(game.QUIT)
		})
		menu = game.NewMenu(battleItem, onlineBattleItem, editorItem, settingItem, quitItem)
		menu.FontSize = 42
		menu.SelectedFontSize = 52
		menu.Padding = 100
	}

	i, j := game.TileManager.Length()
	gameCamera.Target.X = float32(game.TILE_LENGTH*i) / 2
	gameCamera.Target.Y = float32(game.TILE_LENGTH*j) / 2
	m.OnWindowResized()
}

func (MainMenu) OnExit() {}

func (MainMenu) OnWindowResized() {
	menu.Pos.X = float32(game.Width) / 2
	menu.Pos.Y = float32(game.Height)/2 - 200
	menu.Refresh()
	fitCamera()
	gitVersion.Pos = ray.NewVector2(float32(game.Width)-70, float32(game.Height)-20)
	githubProfileLink.Text.Pos = ray.NewVector2(90, float32(game.Height)-20)
	githubProfileLink.Boundary.Y = githubProfileLink.Text.Pos.Y
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

func (MainMenu) Update() {
	menu.Update()
	githubProfileLink.Update()
}

func (MainMenu) Draw() {
	ray.BeginMode2D(gameCamera)
	game.TileManager.Draw()
	ray.EndMode2D()

	ray.DrawRectangle(0, 0, int32(ray.GetScreenWidth()), int32(ray.GetScreenHeight()), ray.NewColor(0, 0, 0, 150))
	menu.Draw()
	ray.DrawText("Source Code:", 5, int32(game.Height)-18, 12, ray.White)
	gitVersion.Draw()
	githubProfileLink.Draw()
}
