package running

import (
	"bomberman/core"
	"bomberman/game"
	"fmt"
	"sync/atomic"
	"time"

	ray "github.com/gen2brain/raylib-go/raylib"
)

var (
	isResetting bool = false
	shouldReset atomic.Bool

	p1 *game.Player
	p2 *game.Player

	elapsedTime float32      = 0
	gameCamera  ray.Camera2D = ray.NewCamera2D(ray.NewVector2(0, 0), ray.NewVector2(0, 0), 0, 1)

	fallingU, fallingD, fallingR, fallingL int
	fallingI, fallingJ                     int
	fallingCnt                             int
	fallingDir                             game.Direction
)

func fitCamera() {
	gameCamera.Offset.X = float32(ray.GetScreenWidth() / 2)
	gameCamera.Offset.Y = float32(ray.GetScreenHeight() / 2)

	screenRatio := float32(ray.GetScreenHeight()) / float32(ray.GetScreenWidth())
	i, j := game.TileManager.Length()
	gridRatio := float32(game.TILE_LENGTH*j) / float32(game.TILE_LENGTH*i)

	fmt.Println(screenRatio, gridRatio)

	if screenRatio < gridRatio {
		gameCamera.Zoom = float32(ray.GetScreenHeight()) / float32(game.TILE_LENGTH*j)
	} else {
		gameCamera.Zoom = float32(ray.GetScreenWidth()) / float32(game.TILE_LENGTH*i)
	}
	// fmt.Printf("camera fitted zoom:%f ofsset:(%.2f,%.2f)\n", gameCamera.Zoom, gameCamera.Offset.X, gameCamera.Offset.Y)
}

func LoadSounds() {
	core.LoadSound("assets/audio/Battle 1.mp3", "battle1").SetVolume(0.8)
	core.LoadSound("assets/audio/Place Bomb.wav", "place-bomb")
	core.LoadSound("assets/audio/Bomb Bounce.wav", "bomb-bounce")
	core.LoadSound("assets/audio/Bomb Explodes.wav", "bomb-explode").SetVolume(0.8)
	core.LoadSound("assets/audio/Hurry Up 1.wav", "huryy1")
	core.LoadSound("assets/audio/Hurry Up 2.wav", "huryy2")
	core.LoadSound("assets/audio/Item Bounce.wav", "item-bounce")
	core.LoadSound("assets/audio/Item Get.wav", "item-get")
	core.LoadSound("assets/audio/Stage Clear.wav", "stage-clear")
	core.LoadSound("assets/audio/Stage Start.wav", "stage-start")
}

func LoadTextures() {
	core.LoadTexture("assets/characters.png", "anims", ray.NewRectangle(0, 0, 100, 100))
	core.LoadTexture("assets/tiles.png", "tiles", ray.NewRectangle(0, 0, 1060, 680))
}

func Unload() {
	core.UnloadSound("battle1")
	core.UnloadSound("place-bomb")
	core.UnloadSound("bomb-bounce")
	core.UnloadSound("bomb-explode")
	core.UnloadSound("huryy1")
	core.UnloadSound("huryy2")
	core.UnloadSound("item-bounce")
	core.UnloadSound("item-get")
	core.UnloadSound("stage-clear")
	core.UnloadSound("stage-start")

	core.UnloadTexture("anims")
	core.UnloadTexture("tiles")
}

func checkWinState() {
	if !p1.IsAlive || !p2.IsAlive {
		isResetting = true
		time.AfterFunc(time.Second*3, func() {
			p1.Lock()
			p2.Lock()
			defer p1.Unlock()
			defer p2.Unlock()
			if !p1.IsAlive && !p2.IsAlive {
				fmt.Println("draw")
			} else {
				if !p1.IsAlive {
					fmt.Println("p2 won the round")
					// p2Wins += 1
					p2.Wins += 1
				}
				if !p2.IsAlive {
					fmt.Println("p1 won the round")
					// p1Wins += 1
					p1.Wins += 1
				}
			}
			maxScore := rounds
			if p1.Wins == maxScore {
				game.State.Change(game.WIN)
			} else if p2.Wins == maxScore {
				game.State.Change(game.WIN)
			} else {
				shouldReset.Store(true)
			}
		})
	}
}

func loadLevel(name string) {
	fmt.Println("loading", name)
	game.TileManager.GameMap = &game.GameMap{
		Name: "assets/maps/" + name,
	}
	game.TileManager.LoadFromFile(game.TileManager.Name)
	reset()
	p1.Wins = 0
	p2.Wins = 0
}

func reset() {
	isResetting = false
	shouldReset.Store(false)
	elapsedTime = 0
	game.CollisionManager.Reset()
	if p1 == nil {
		p1 = game.NewPlayer(1)
	}
	if p2 == nil {
		p2 = game.NewPlayer(2)
	}
	game.CollisionManager.Add(p1)
	game.CollisionManager.Add(p2)
	p1.Reset()
	p2.Reset()
	p1.Scale = 1.8
	p2.Scale = 1.8
	game.DropUpgrades = true
	game.SaveUpgrades = false
	game.TileManager.Reload()
	fitCamera()
	p1.Position = game.TileManager.P1Pos
	p2.Position = game.TileManager.P2Pos
	if p1.Position.X == 0 {
		p1.Position = ray.NewVector2(game.TILE_LENGTH+game.TILE_LENGTH/2, game.TILE_LENGTH+game.TILE_LENGTH/2)
	}
	if p2.Position.X == 0 {
		p2.Position = ray.NewVector2(float32(game.TileManager.YCount-3)*game.TILE_LENGTH+game.TILE_LENGTH/2, game.TILE_LENGTH+game.TILE_LENGTH/2)
	}
	if game.State.Get() == game.ONLINE_BATTLE {
		if isHost {
			p2.IsControllable = false
			game.SaveUpgrades = true
		} else {
			game.DropUpgrades = false
			p1.IsControllable = false
		}
	}
	i, j := game.TileManager.Length()
	gameCamera.Target.X = float32(game.TILE_LENGTH*i) / 2
	gameCamera.Target.Y = float32(game.TILE_LENGTH*j) / 2
	fallingU = 1
	fallingL = 1
	fallingD = game.TileManager.YCount - 2
	fallingR = game.TileManager.XCount - 2
	fallingI = 1
	fallingJ = 1
	fallingCnt = 0
	fallingDir = game.RIGHT
}

func fallWalls() {
	if t := int(elapsedTime) - roundTimeSeconds; t > 0 {
		if t > fallingCnt && fallingCnt+1 < (game.TileManager.XCount-2)*(game.TileManager.YCount-2) {
			game.TileManager.Remove(fallingI, fallingJ)
			// ti := game.NewTile(fallingI, fallingJ, core.GetTexture("tiles").Crop(ray.NewRectangle(238, 0, 16, 16)))
			ti := game.NewFallingTile(fallingI, fallingJ)
			game.TileManager.Add(ti)
			fallingCnt++
			if ray.CheckCollisionRecs(p1.GetCollision(), ti.GetCollision()) {
				p1.Die()
			}
			if ray.CheckCollisionRecs(p2.GetCollision(), ti.GetCollision()) {
				p2.Die()
			}
			switch fallingDir {
			case game.RIGHT:
				fallingI++
				if fallingI == fallingR {
					fallingR--
					fallingDir = game.DOWN
				}
			case game.DOWN:
				fallingJ++
				if fallingJ == fallingD {
					fallingD--
					fallingDir = game.LEFT
				}
			case game.LEFT:
				fallingI--
				if fallingI == fallingL {
					fallingL++
					fallingDir = game.UP
				}
			case game.UP:
				fallingJ--
				if fallingJ-1 == fallingU {
					fallingU++
					fallingDir = game.RIGHT
				}
			}
		}
	}
}

func updateBattle() {
	elapsedTime += ray.GetFrameTime()
	fallWalls()
	if !isResetting {
		checkWinState()
	} else if shouldReset.Load() {
		reset()
		return
	}

	if ray.IsKeyPressed(ray.KeyEscape) {
		game.State.Change(game.PAUSED)
		return
	}

	if ray.IsKeyPressed(ray.KeyF1) {
		game.ShowCollsions = !game.ShowCollsions
	}
	p1.Update()
	p2.Update()

	// mouseWheel := ray.GetMouseWheelMove()
	// if mouseWheel > 0 {
	// 	p1.Scale += 0.2
	// } else if mouseWheel < 0 {
	// 	p1.Scale -= 0.2
	// }
	game.TileManager.Update()
	updateUi()
}

func drawBattle() {
	ray.BeginMode2D(gameCamera)
	game.TileManager.Draw()
	if p1.Position.Y > p2.Position.Y {
		p2.Draw()
		p1.Draw()
	} else {
		p1.Draw()
		p2.Draw()
	}
	ray.EndMode2D()
	drawUi()
}
