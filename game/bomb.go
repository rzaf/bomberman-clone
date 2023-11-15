package game

import (
	"bomberman/core"
	// "fmt"
	"time"

	ray "github.com/gen2brain/raylib-go/raylib"
)

const (
	timeOut = 2500
)

type BombTile struct {
	Tile
	bomber      *Player
	startTime   int64
	isDestroyed bool
	radius      int
	rRadius     int
	lRadius     int
	uRadius     int
	dRadius     int
	idle        *core.Animation
	explosionC  *core.Animation
	explosionUD *core.Animation
	explosionLR *core.Animation
	explosionR  *core.Animation
	explosionL  *core.Animation
	explosionU  *core.Animation
	explosionD  *core.Animation
}

func NewBombTile(i, j, explosionRadius int, p *Player) *BombTile {
	b := &BombTile{
		Tile:    Tile{i, j, i * TILE_LENGTH, j * TILE_LENGTH, TILE_LENGTH, TILE_LENGTH, nil, true},
		bomber:  p,
		radius:  explosionRadius,
		rRadius: explosionRadius,
		lRadius: explosionRadius,
		uRadius: explosionRadius,
		dRadius: explosionRadius,
	}
	b.HasCollision = false
	CollisionManager.Add(Collision(b))
	t1 := core.GetTexture("tiles")

	var offset float32 = 0
	if p.Number == 2 {
		offset += 128
	}
	b.idle = core.NewAnimation(
		core.NewSpriteSheet(t1, 4, 4, 1, 16, 16, 0, 0, ray.NewRectangle(0+offset, 173, 64, 16), false, false),
		200,
		true,
	)
	b.explosionC = core.NewAnimation(
		core.NewSpriteSheet(t1, 5, 1, 5, 16, 16, 0, 0, ray.NewRectangle(96+offset, 189, 16, 80), true, true),
		80,
		false,
	)
	b.explosionUD = core.NewAnimation(
		core.NewSpriteSheet(t1, 5, 1, 5, 16, 16, 0, 0, ray.NewRectangle(64+offset, 189, 16, 80), true, true),
		80,
		false,
	)
	b.explosionLR = core.NewAnimation(
		core.NewSpriteSheet(t1, 5, 1, 5, 16, 16, 0, 0, ray.NewRectangle(80+offset, 189, 16, 80), true, true),
		80,
		false,
	)
	b.explosionR = core.NewAnimation(
		core.NewSpriteSheet(t1, 5, 1, 5, 16, 16, 0, 0, ray.NewRectangle(48+offset, 189, 16, 80), true, true),
		80,
		false,
	)
	b.explosionL = core.NewAnimation(
		core.NewSpriteSheet(t1, 5, 1, 5, 16, 16, 0, 0, ray.NewRectangle(32+offset, 189, 16, 80), true, true),
		80,
		false,
	)
	b.explosionD = core.NewAnimation(
		core.NewSpriteSheet(t1, 5, 1, 5, 16, 16, 0, 0, ray.NewRectangle(16+offset, 189, 16, 80), true, true),
		80,
		false,
	)
	b.explosionU = core.NewAnimation(
		core.NewSpriteSheet(t1, 5, 1, 5, 16, 16, 0, 0, ray.NewRectangle(0+offset, 189, 16, 80), true, true),
		80,
		false,
	)
	b.idle.Play()
	b.startTime = time.Now().UnixMilli()
	return b
}

func (b *BombTile) AddToTileSet() *BombTile {
	i, j := b.I, b.J
	switch TileManager.Get(i, j).(type) {
	case *FloorTile:
	case *UpgradeTile:
	default:
		return nil
	}

	for _, t := range TileManager.extraTiles {
		bt, Isbomb := t.(*BombTile)
		if Isbomb {
			if bt.I == i && bt.J == j {
				return nil
			}
		}
	}
	TileManager.Add(b)
	return b
}

func (b *BombTile) checkTileCollision(index int, d Direction) bool {
	col := b.GetCollision()
	var sideRec ray.Rectangle
	switch d {
	case RIGHT:
		sideRec = ray.NewRectangle(col.X+TILE_LENGTH*float32(index+1), col.Y, TILE_LENGTH, TILE_LENGTH)
	case LEFT:
		sideRec = ray.NewRectangle(col.X-TILE_LENGTH*float32(index+1), col.Y, TILE_LENGTH, TILE_LENGTH)
	case UP:
		sideRec = ray.NewRectangle(col.X, col.Y-TILE_LENGTH*float32(index+1), TILE_LENGTH, TILE_LENGTH)
	case DOWN:
		sideRec = ray.NewRectangle(col.X, col.Y+TILE_LENGTH*float32(index+1), TILE_LENGTH, TILE_LENGTH)
	}
	sideCollisions := CollisionManager.CheckCollisions(sideRec)

	for _, c := range sideCollisions {
		switch t := c.(type) {
		case *Player:
			t.Die()
		case *WallTile:
			if t.isDestoryable {
				t.Destroy()
			}
			switch d {
			case RIGHT:
				b.rRadius = index
			case LEFT:
				b.lRadius = index
			case UP:
				b.uRadius = index
			case DOWN:
				b.dRadius = index
			}
			return true
		case *BombTile:
			if !t.isDestroyed {
				t.Destroy()
			}
			switch d {
			case RIGHT:
				b.rRadius = index
			case LEFT:
				b.lRadius = index
			case UP:
				b.uRadius = index
			case DOWN:
				b.dRadius = index
			}
			return true
		case *UpgradeTile:
			timer := time.NewTimer(time.Millisecond * 200)
			go func() {
				<-timer.C
				if !t.isDestroyed.Load() {
					t.Destroy()
				}
			}()
		}
	}
	return false
}

func (b *BombTile) checkExplosionCollisions() {
	col := b.GetCollision()
	centerCollisions := CollisionManager.CheckCollisions(col)
	for _, c := range centerCollisions {
		switch t := c.(type) {
		case *Player:
			t.Die()
		case *UpgradeTile:
			if !t.isDestroyed.Load() {
				t.Destroy()
			}
		}
	}

	for i := 0; i < b.rRadius; i++ {
		if b.checkTileCollision(i, RIGHT) {
			break
		}
	}
	for i := 0; i < b.lRadius; i++ {
		if b.checkTileCollision(i, LEFT) {
			break
		}
	}
	for i := 0; i < b.uRadius; i++ {
		if b.checkTileCollision(i, UP) {
			break
		}
	}
	for i := 0; i < b.dRadius; i++ {
		if b.checkTileCollision(i, DOWN) {
			break
		}
	}
	// fmt.Println(b.lRadius, b.uRadius, b.rRadius, b.dRadius)
}

func (b *BombTile) Destroy() {
	b.isDestroyed = true
	core.GetSound("bomb-explode").Play()
	b.explosionC.Play()
	b.explosionUD.Play()
	b.explosionLR.Play()
	b.explosionR.Play()
	b.explosionD.Play()
	b.explosionU.Play()
	b.explosionL.Play()
	b.checkExplosionCollisions()
}

func (b *BombTile) Update() {
	if b.isDestroyed {
		b.explosionC.Update()
		b.explosionUD.Update()
		b.explosionLR.Update()
		b.explosionR.Update()
		b.explosionL.Update()
		b.explosionU.Update()
		b.explosionD.Update()
		b.checkExplosionCollisions()
	} else {
		if time.Now().UnixMilli()-b.startTime > timeOut {
			b.Destroy()
		}
		b.idle.Update()
		if !b.HasCollision {
			arePlayersOutside := true
			centerCollisions := CollisionManager.CheckCollisions(b.GetCollision())
			for _, c := range centerCollisions {
				switch c.(type) {
				case *Player:
					arePlayersOutside = false
				}
			}
			if arePlayersOutside {
				b.HasCollision = true
			}
		}
	}
}

func (b *BombTile) Draw() {
	if b.isDestroyed {
		if b.explosionC.IsPlaying() {
			b.explosionC.DrawAt(ray.NewRectangle(float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height)))
			for i := 0; i < b.rRadius; i++ {
				if i == b.rRadius-1 {
					b.explosionR.DrawAt(ray.NewRectangle(float32(b.X+TILE_LENGTH*(i+1)), float32(b.Y), float32(b.Width), float32(b.Height)))
				} else {
					b.explosionLR.DrawAt(ray.NewRectangle(float32(b.X+TILE_LENGTH*(i+1)), float32(b.Y), float32(b.Width), float32(b.Height)))
				}
			}
			for i := 0; i < b.lRadius; i++ {
				if i == b.lRadius-1 {
					b.explosionL.DrawAt(ray.NewRectangle(float32(b.X-TILE_LENGTH*(i+1)), float32(b.Y), float32(b.Width), float32(b.Height)))
				} else {
					b.explosionLR.DrawAt(ray.NewRectangle(float32(b.X-TILE_LENGTH*(i+1)), float32(b.Y), float32(b.Width), float32(b.Height)))
				}
			}
			for i := 0; i < b.uRadius; i++ {
				if i == b.uRadius-1 {
					b.explosionU.DrawAt(ray.NewRectangle(float32(b.X), float32(b.Y-TILE_LENGTH*(i+1)), float32(b.Width), float32(b.Height)))
				} else {
					b.explosionUD.DrawAt(ray.NewRectangle(float32(b.X), float32(b.Y-TILE_LENGTH*(i+1)), float32(b.Width), float32(b.Height)))
				}
			}
			for i := 0; i < b.dRadius; i++ {
				if i == b.dRadius-1 {
					b.explosionD.DrawAt(ray.NewRectangle(float32(b.X), float32(b.Y+TILE_LENGTH*(i+1)), float32(b.Width), float32(b.Height)))
				} else {
					b.explosionUD.DrawAt(ray.NewRectangle(float32(b.X), float32(b.Y+TILE_LENGTH*(i+1)), float32(b.Width), float32(b.Height)))
				}
			}
		} else {
			b.Remove()
		}
	} else {
		b.idle.DrawAt(ray.NewRectangle(float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height)))
	}
	if ShowCollsions && b.HasCollision {
		ray.DrawRectangleLinesEx(b.GetCollision(), 2, ray.Blue)
	}
}

func (b *BombTile) Remove() {
	b.HasCollision = false
	b.bomber.upgrades.bombCount++
	for i, t := range TileManager.extraTiles {
		bt, Isbomb := t.(*BombTile)
		if Isbomb && bt == b {
			TileManager.extraTiles[i] = nil
		}
	}
}
