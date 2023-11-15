package editor

import (
	"bomberman/core"
	"bomberman/game"
	"fmt"
	"log"
	"os"
	"strings"

	ray "github.com/gen2brain/raylib-go/raylib"
)

var (
	gameMaps     []*game.GameMap
	names        []*core.Text
	currentIndex int

	pickingTiles []*core.Texture
	selectedTile int

	xSlider *core.Slider
	ySlider *core.Slider

	camera        ray.Camera2D
	restartButton *core.Button
	saveButton    *core.Button
	exitButton    *core.Button
)

type Editor struct{}

func (Editor) OnEnter() {
	core.LoadTexture("assets/characters.png", "anims", ray.NewRectangle(0, 0, 100, 100))
	core.LoadTexture("assets/tiles.png", "tiles", ray.NewRectangle(0, 0, 1060, 680))
	gameMaps = nil
	names = nil
	pickingTiles = nil
	game.TileManager.GameMap = nil
	dirs, err := os.ReadDir("assets/maps/")
	if err != nil {
		log.Fatalln(err)
	}
	names = append(names, core.NewText("New Map", ray.GetFontDefault(), ray.NewVector2(100, 50+float32(len(names)+1)*35), 18, 3, ray.White))
	gameMaps = append(gameMaps, &game.GameMap{
		XCount: 5,
		YCount: 5,
	})
	for i, dir := range dirs {
		fmt.Println(dir)
		newMap := game.GameMap{}
		newMap.LoadFromFile("assets/maps/" + dir.Name())
		newMap.Name = dir.Name()
		gameMaps = append(gameMaps, &newMap)
		names = append(names, core.NewText(dir.Name(), ray.GetFontDefault(), ray.NewVector2(100, 50+float32(i+2)*35), 18, 3, ray.White))

	}

	pickingTiles = append(pickingTiles, core.GetTexture("tiles").Crop(ray.NewRectangle(153, 0, 16, 16))) // d wall
	pickingTiles[0].Dest = ray.NewRectangle(250+200, 5, 50, 50)
	pickingTiles = append(pickingTiles, core.GetTexture("tiles").Crop(ray.NewRectangle(170, 0, 16, 16))) // normal wall
	pickingTiles[1].Dest = ray.NewRectangle(250+250+20, 5, 50, 50)
	pickingTiles = append(pickingTiles, core.GetTexture("tiles").Crop(ray.NewRectangle(204, 0, 16, 16))) // floor
	pickingTiles[2].Dest = ray.NewRectangle(250+300+40, 5, 50, 50)
	pickingTiles = append(pickingTiles, core.GetTexture("anims").Crop(ray.NewRectangle(64, 0, 17, 25))) // p1
	pickingTiles[3].Dest = ray.NewRectangle(250+350+40, 5, 40, 50)
	pickingTiles = append(pickingTiles, core.GetTexture("anims").Crop(ray.NewRectangle(167, 0, 17, 25))) // p2
	pickingTiles[4].Dest = ray.NewRectangle(250+400+40, 5, 40, 50)

	game.TileManager.GameMap = gameMaps[0]
	currentIndex = 0
	game.TileManager.Init()
	camera = ray.NewCamera2D(ray.NewVector2(400, 200), ray.NewVector2(0, 0), 0, 1)

	xSlider = core.NewSlider(210, 10, 200, 8, 5, 25, 5)
	ySlider = core.NewSlider(210, 40, 200, 8, 5, 25, 5)

	restartButton = core.NewTextButton(core.NewText("RESET", ray.GetFontDefault(), ray.NewVector2(70, float32(ray.GetScreenHeight())-20), 18, 2, ray.White), ray.NewRectangle(0, float32(ray.GetScreenHeight())-30, 100, 100))
	saveButton = core.NewTextButton(core.NewText("SAVE", ray.GetFontDefault(), ray.NewVector2(150, float32(ray.GetScreenHeight())-20), 18, 2, ray.White), ray.NewRectangle(100, float32(ray.GetScreenHeight())-30, 100, 100))
	exitButton = core.NewTextButton(core.NewText("EXIT", ray.GetFontDefault(), ray.NewVector2(20, 10), 18, 2, ray.White), ray.NewRectangle(0, 0, 100, 20))

	exitButton.OnClick = func() {
		game.State.Change(game.MENU)
	}
	restartButton.OnClick = func() {
		if currentIndex == 0 {
			game.TileManager.Init()
		} else {
			game.TileManager.GameMap = gameMaps[currentIndex].LoadFromFile("assets/maps/" + names[currentIndex].Text)
		}
		xSlider.SetValue(float32(game.TileManager.XCount))
		ySlider.SetValue(float32(game.TileManager.YCount))
	}
	saveButton.OnClick = func() {
		if names[0].Text == "New Map" || strings.TrimSpace(names[0].Text) == "" {
			return
		}
		game.TileManager.Name = strings.TrimSpace(names[0].Text)
		fmt.Printf("trying to save map:%d,%d\n", game.TileManager.XCount, game.TileManager.YCount)
		game.TileManager.SaveInto("assets/maps/")
		names = append(names, core.NewText(game.TileManager.Name, ray.GetFontDefault(), ray.NewVector2(100, 50+float32((len(names)+1)*35)), 18, 3, ray.White))
		var createdMap *game.GameMap = &game.GameMap{}
		game.CopyGameMap(createdMap, game.TileManager.GameMap)
		gameMaps = append(gameMaps, createdMap)

		gameMaps[0].Init()
		currentIndex = len(names) - 1
		game.TileManager.GameMap = gameMaps[currentIndex]
		names[0].Text = "New Map"
	}
}

func (Editor) OnExit() {}

func (Editor) OnWindowResized() {
	restartButton.Text.Pos = ray.NewVector2(70, float32(ray.GetScreenHeight())-20)
	restartButton.Boundary = ray.NewRectangle(0, float32(ray.GetScreenHeight())-30, 100, 100)
	saveButton.Text.Pos = ray.NewVector2(150, float32(ray.GetScreenHeight())-20)
	saveButton.Boundary = ray.NewRectangle(100, float32(ray.GetScreenHeight())-30, 100, 100)
}

func (Editor) Update() {
	restartButton.Update()
	saveButton.Update()
	exitButton.Update()

	if currentIndex == 0 {
		k := ray.GetKeyPressed()
		validKey := false
		if k >= ray.KeyA && k <= ray.KeyZ {
			validKey = true
		}
		switch k {
		case
			ray.KeyPeriod,
			ray.KeySlash,
			ray.KeyApostrophe,
			ray.KeyComma,
			ray.KeyZero,
			ray.KeyOne,
			ray.KeyTwo,
			ray.KeyThree,
			ray.KeyFour,
			ray.KeyFive,
			ray.KeySix,
			ray.KeySeven,
			ray.KeyEight,
			ray.KeyNine,
			ray.KeySemicolon,
			ray.KeyEqual:
			validKey = true
		}
		if validKey {
			if names[0].Text == "New Map" {
				names[0].Text = strings.ToLower(string(k))
			} else {
				names[0].Text += strings.ToLower(string(k))
			}
		}
		if k == ray.KeyBackspace {
			l := len(names[0].Text)
			names[0].Text = names[0].Text[:l-1]
		}
	}
	if ray.IsKeyPressed(ray.KeyUp) {
		if currentIndex == 0 {
			currentIndex = len(names) - 1
		} else {
			currentIndex -= 1
		}

		game.TileManager.GameMap = gameMaps[currentIndex]
		xSlider.SetValue(float32(game.TileManager.XCount))
		ySlider.SetValue(float32(game.TileManager.YCount))
	} else if ray.IsKeyPressed(ray.KeyDown) {
		currentIndex = (currentIndex + 1) % len(names)
		game.TileManager.GameMap = gameMaps[currentIndex]
		xSlider.SetValue(float32(game.TileManager.XCount))
		ySlider.SetValue(float32(game.TileManager.YCount))
	}

	for i := 0; i < len(names); i++ {
		if i == currentIndex {
			names[i].FontSize = 24
			names[i].Color = ray.Red
		} else {
			names[i].FontSize = 18
			names[i].Color = ray.White
		}
		names[i].Measure()
	}
	p := ray.GetMousePosition()
	if p.X > 200 && p.Y > 60 {
		if ray.IsMouseButtonPressed(ray.MouseMiddleButton) {
			ray.SetMouseCursor(ray.MouseCursorCrosshair)
		}
		if ray.IsMouseButtonReleased(ray.MouseMiddleButton) {
			ray.SetMouseCursor(ray.MouseCursorArrow)
		}
		if ray.IsMouseButtonDown(ray.MouseMiddleButton) {
			p2 := ray.GetMouseDelta()
			camera.Offset.X += p2.X
			camera.Offset.Y += p2.Y
		}
		m := ray.GetMouseWheelMove()
		camera.Zoom += m * 0.1
		if ray.IsMouseButtonDown(ray.MouseLeftButton) {
			p2 := ray.GetScreenToWorld2D(p, camera)
			i := int(p2.X / game.TILE_LENGTH)
			j := int(p2.Y / game.TILE_LENGTH)
			fmt.Println(i, j)
			if i > 0 && j > 0 && i < game.TileManager.XCount-1 && j < game.TileManager.YCount-1 {
				var t game.TileInterface
				switch selectedTile {
				case 0:
					t = game.NewWallTile(i, j, true)
				case 1:
					t = game.NewWallTile(i, j, false)
				case 2:
					t = game.NewFloorTile(i, j)
				case 3:
					game.TileManager.P1Pos = ray.NewVector2(float32(i)*game.TILE_LENGTH+game.TILE_LENGTH/2, float32(j)*game.TILE_LENGTH+game.TILE_LENGTH/2)
				case 4:
					game.TileManager.P2Pos = ray.NewVector2(float32(i)*game.TILE_LENGTH+game.TILE_LENGTH/2, float32(j)*game.TILE_LENGTH+game.TILE_LENGTH/2)
				}
				if t != nil {
					game.TileManager.Add(t)
				}
			}
		}
	} else {
		if ray.IsMouseButtonPressed(ray.MouseLeftButton) {
			if ray.CheckCollisionPointRec(p, pickingTiles[0].Dest) {
				selectedTile = 0
			} else if ray.CheckCollisionPointRec(p, pickingTiles[1].Dest) {
				selectedTile = 1
			} else if ray.CheckCollisionPointRec(p, pickingTiles[2].Dest) {
				selectedTile = 2
			} else if ray.CheckCollisionPointRec(p, pickingTiles[3].Dest) {
				selectedTile = 3
			} else if ray.CheckCollisionPointRec(p, pickingTiles[4].Dest) {
				selectedTile = 4
			} else {
				for i := 0; i < len(names); i++ {
					if ray.CheckCollisionPointRec(p, ray.NewRectangle(0, 50+float32(i+1)*35-names[i].Size.Y/2, 200, names[i].Size.Y)) {
						currentIndex = i
						game.TileManager.GameMap = gameMaps[i]
						xSlider.SetValue(float32(gameMaps[i].XCount))
						ySlider.SetValue(float32(gameMaps[i].YCount))
						break
					}
				}
			}
		}
	}

	xSlider.Update()
	ySlider.Update()

	if game.TileManager.XCount != xSlider.GetValueInt() || game.TileManager.YCount != ySlider.GetValueInt() {
		game.TileManager.Change(&game.GameMap{
			XCount: xSlider.GetValueInt(),
			YCount: ySlider.GetValueInt(),
			P1Pos:  game.TileManager.P1Pos,
			P2Pos:  game.TileManager.P2Pos,
		})
		gameMaps[currentIndex] = game.TileManager.GameMap
	}
}

func (Editor) Draw() {
	ray.BeginMode2D(camera)
	game.TileManager.DrawInEditor()
	ray.EndMode2D()

	ray.DrawRectangle(0, 0, 200, int32(ray.GetScreenHeight()), ray.Black)
	ray.DrawRectangle(200, 0, int32(ray.GetScreenWidth()), 60, ray.DarkGray)
	for _, m := range names {
		m.DrawCentered()
	}
	for i, t := range pickingTiles {
		t.Draw()
		if i == selectedTile {
			ray.DrawRectangleLinesEx(t.Dest, 1.5, ray.Red)
		}
	}

	xSlider.Draw()
	ray.DrawText(fmt.Sprintf("%d", xSlider.GetValueInt()), 420, 4, 20, ray.Black)
	ySlider.Draw()
	ray.DrawText(fmt.Sprintf("%d", ySlider.GetValueInt()), 420, 38, 20, ray.Black)

	restartButton.Text.DrawCentered()
	saveButton.Text.DrawCentered()
	exitButton.Draw()
}
