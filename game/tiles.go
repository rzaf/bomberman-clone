package game

import (
	"bomberman/core"
	"fmt"

	ray "github.com/gen2brain/raylib-go/raylib"
)

var (
	DropUpgrades bool = true
)

type TileInterface interface {
	tile()
	Draw()
	Index() (int, int)
}

type Tile struct {
	I, J         int
	X, Y         int
	Width        int
	Height       int
	texture      *core.Texture
	HasCollision bool
}

func NewTile(i, j int, texture *core.Texture) *Tile {
	t := &Tile{i, j, i * TILE_LENGTH, j * TILE_LENGTH, TILE_LENGTH, TILE_LENGTH, texture, false}
	CollisionManager.Add(Collision(t))
	return t
}

func (t *Tile) tile() {}

func (t *Tile) Position() ray.Vector2 {
	return ray.NewVector2(float32(t.X), float32(t.Y))
}

func (t *Tile) String() string {
	return fmt.Sprintf("x:%d,y:%d,w:%d,h:%d", t.X, t.Y, t.Width, t.Height)
}

func (t *Tile) Draw() {
	t.texture.DrawAt(ray.NewRectangle(float32(t.X), float32(t.Y), float32(t.Width), float32(t.Height)))
	if ShowCollsions && t.HasCollision {
		ray.DrawRectangleLinesEx(t.GetCollision(), 2, ray.Blue)
	}
}

func (t *Tile) Index() (int, int) {
	return t.I, t.J
}

func (t *Tile) GetCollision() ray.Rectangle {
	rec := ray.NewRectangle(float32(t.X), float32(t.Y), float32(t.Width), float32(t.Height))
	return rec
}

func (t *Tile) Enabled() bool {
	return t.HasCollision
}

type WallTile struct {
	Tile
	isDestoryable       bool
	destroyingAnimation *core.Animation
	isDestroyed         bool
}

func NewWallTile(i, j int, isDestoryable bool) *WallTile {
	t := &WallTile{
		Tile:                Tile{i, j, i * TILE_LENGTH, j * TILE_LENGTH, TILE_LENGTH, TILE_LENGTH, nil, true},
		isDestoryable:       isDestoryable,
		destroyingAnimation: nil,
		isDestroyed:         false,
	}
	if isDestoryable {
		t.destroyingAnimation = core.NewAnimation(
			core.NewSpriteSheet(core.GetTexture("tiles"), 6, 6, 1, 16, 16, 1, 0, ray.NewRectangle(153, 17, 101, 16), false, false),
			80,
			false,
		)
		t.texture = core.GetTexture("tiles").Crop(ray.NewRectangle(153, 0, 16, 16))
	} else {
		t.texture = core.GetTexture("tiles").Crop(ray.NewRectangle(170, 0, 16, 16))
	}
	CollisionManager.Add(Collision(t))
	return t
}

func (w *WallTile) Destroy() {
	if w.isDestoryable {
		fmt.Printf("destroying wall %d,%d \n", w.I, w.J)
		w.isDestroyed = true
		w.destroyingAnimation.Play()
		// time.AfterFunc(250*time.Millisecond, w.remove)
	}
}

func (w *WallTile) Remove() {
	fmt.Printf("removing wall %d,%d \n", w.I, w.J)
	CollisionManager.Remove(w)
	// TileManager.Add(NewFloorTile(w.I, w.J))
	TileManager.Remove(w.I, w.J)
	if DropUpgrades {
		AddUpgradeTile(w.I, w.J, GetRandomUpgrade())
	} else {
		AddUpgradeTile(w.I, w.J, UPGRADE_NONE)
	}
}

func (w *WallTile) Update() {
	if w.isDestoryable && w.isDestroyed {
		w.destroyingAnimation.Update()
		if !w.destroyingAnimation.IsPlaying() {
			w.Remove()
		}
	}
}

func (w *WallTile) Draw() {
	if w.isDestoryable && w.isDestroyed {
		w.destroyingAnimation.DrawAt(ray.NewRectangle(float32(w.X), float32(w.Y), float32(w.Width), float32(w.Height)))
	} else {
		w.texture.DrawAt(ray.NewRectangle(float32(w.X), float32(w.Y), float32(w.Width), float32(w.Height)))
	}
	// if ShowCollsions && w.HasCollision {
	// 	ray.DrawRectangleLinesEx(w.GetCollision(), 2, ray.Blue)
	// }
}

type FloorTile struct {
	Tile
}

func NewFloorTile(i, j int) *FloorTile {
	t := &FloorTile{
		Tile: Tile{i, j, i * TILE_LENGTH, j * TILE_LENGTH, TILE_LENGTH, TILE_LENGTH, nil, true},
	}
	if j == 1 {
		t.texture = core.GetTexture("tiles").Crop(ray.NewRectangle(221, 0, 16, 16))
	} else {
		t.texture = core.GetTexture("tiles").Crop(ray.NewRectangle(204, 0, 16, 16))
	}
	return t
}

type FallingTile struct {
	Tile
}

func NewFallingTile(i, j int) *FallingTile {
	t := &FallingTile{
		Tile: Tile{i, j, i * TILE_LENGTH, j * TILE_LENGTH, TILE_LENGTH, TILE_LENGTH, nil, true},
	}
	t.texture = core.GetTexture("tiles").Crop(ray.NewRectangle(238, 0, 16, 16))
	CollisionManager.Add(Collision(t))
	return t
}
