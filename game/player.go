package game

import (
	"fmt"
	"github.com/rzaf/bomberman-clone/core"
	"sync"
	"time"

	ray "github.com/gen2brain/raylib-go/raylib"
)

const (
	defaultSpeed = 90
	maxSpeed     = 250
	maxFire      = 10
	maxBomb      = 10
)

type upgrades struct {
	bombCount   int
	bombRadius  int
	heartsCount int
	passBombs   bool
	passWalls   bool
	kickBomb    bool
	grabBomb    bool
}

type Player struct {
	Number    int
	Position  ray.Vector2
	Velocity  ray.Vector2
	Speed     float32
	Scale     float32
	Direction Direction
	upgrades
	Animations     PlayerAnimations
	IsAlive        bool
	HasCollision   bool
	diedAt         int64
	IsControllable bool
	Wins           int

	LastBombX, LastBombY int
	sync.Mutex
}

func NewPlayer(number int) *Player {
	p := &Player{
		Number:    number,
		Direction: DOWN,
		Speed:     defaultSpeed,
		Scale:     1,
		upgrades: upgrades{
			bombCount:   1,
			bombRadius:  1,
			heartsCount: 1,
			passBombs:   false,
			passWalls:   false,
			kickBomb:    false,
			grabBomb:    false,
		},
		HasCollision:   true,
		IsAlive:        true,
		IsControllable: true,
		Wins:           0,
		LastBombX:      -1,
		LastBombY:      -1,
	}
	p.Animations.load(p.Number)
	p.Animations.WalkingDown.Play()
	p.Animations.WalkingUp.Play()
	p.Animations.WalkingLeft.Play()
	p.Animations.WalkingRight.Play()
	return p
}

func (p *Player) Reset() {
	p.Lock()
	defer p.Unlock()
	p.bombCount = 1
	p.bombRadius = 1
	p.passBombs = false
	p.passWalls = false
	p.kickBomb = false
	p.grabBomb = false

	p.IsAlive = true
	p.diedAt = 0
	p.HasCollision = true
	p.Speed = defaultSpeed
	p.IsControllable = true
	p.Direction = DOWN
}

func (p *Player) getNormalizedDirection() ray.Vector2 {
	var l, r, u, d int
	pN := fmt.Sprintf("p%d-", p.Number)

	if IsKeyDown(pN + "Left") {
		l = 1
	} else {
		l = 0
	}
	if IsKeyDown(pN + "Right") {
		r = 1
	} else {
		r = 0
	}
	if IsKeyDown(pN + "Up") {
		u = 1
	} else {
		u = 0
	}
	if IsKeyDown(pN + "Down") {
		d = 1
	} else {
		d = 0
	}

	norm := ray.NewVector2(float32(r-l), float32(d-u))
	if norm.X != 0 || norm.Y != 0 {
		norm = ray.Vector2Normalize(norm)
	}
	return norm
}

func (p *Player) handleInput() {
	norm := p.getNormalizedDirection()

	if norm.X != 0 {
		if norm.X > 0 {
			p.Direction = RIGHT
			p.Animations.WalkingRight.Update()
			p.Animations.WalkingLeft.Index = 0
		} else {
			p.Direction = LEFT
			p.Animations.WalkingRight.Index = 0
			p.Animations.WalkingLeft.Update()
		}
		// p.position.X += float32(lr) * p.Speed * ray.GetFrameTime()
		p.Velocity.X = float32(norm.X) * p.Speed
	} else {
		p.Velocity.X = 0
		p.Animations.WalkingLeft.Index = 0
		p.Animations.WalkingRight.Index = 0
	}

	if norm.Y != 0 {
		if norm.Y > 0 {
			p.Direction = DOWN
			p.Animations.WalkingDown.Update()
			p.Animations.WalkingUp.Index = 1
		} else {
			p.Animations.WalkingDown.Index = 1
			p.Direction = UP
			p.Animations.WalkingUp.Update()
		}
		p.Velocity.Y = float32(norm.Y) * p.Speed
	} else {
		p.Velocity.Y = 0
		p.Animations.WalkingUp.Index = 1
		p.Animations.WalkingDown.Index = 1
	}
}

func (p *Player) move() {
	p.Position.X += p.Velocity.X * ray.GetFrameTime()
	pRec := p.GetCollision()
	if col, collided := p.collides(); collided {
		p.Position.X -= p.Velocity.X * ray.GetFrameTime()
		if wall, isWall := col.(*WallTile); isWall && wall.I > 0 && wall.J > 0 && wall.I-1 < TileManager.XCount && wall.J < TileManager.YCount-1 {
			rec := col.GetCollision()
			if p.Velocity.Y <= 0 {
				if rec.Y-pRec.Y > rec.Height/3 {
					// p.Velocity.Y = -p.Speed
					p.Position.Y += -p.Speed * ray.GetFrameTime()
					p.Velocity.Y = 0
				}
			}
			if p.Velocity.Y >= 0 {
				if rec.Height+rec.Y-pRec.Y < rec.Height/3 {
					// p.Velocity.Y = p.Speed
					p.Position.Y += p.Speed * ray.GetFrameTime()
					p.Velocity.Y = 0
				}
			}
		}
	}

	p.Position.Y += p.Velocity.Y * ray.GetFrameTime()
	if col, collided := p.collides(); collided {
		p.Position.Y -= p.Velocity.Y * ray.GetFrameTime()
		rec := col.GetCollision()
		if wall, isWall := col.(*WallTile); isWall && wall.I > 0 && wall.J > 0 && wall.I < TileManager.XCount-1 && wall.J < TileManager.YCount-1 {
			// fmt.Println(rec.Width - (pRec.X - rec.X))
			if p.Velocity.X == 0 {
				if pRec.Width-(rec.X-pRec.X) < rec.Width/3 {
					p.Position.X += -p.Speed * ray.GetFrameTime()
				}
			}
			if p.Velocity.X == 0 {
				if rec.Width-(pRec.X-rec.X) < rec.Width/3 {
					p.Position.X += p.Speed * ray.GetFrameTime()
				}
			}
		}
	}
}

func (p *Player) Update() {
	p.Lock()
	defer p.Unlock()
	if !p.IsAlive {
		p.Animations.dying.Update()
		p.Animations.crying1.Update()
		p.Animations.crying2.Update()
		return
	}
	if p.IsControllable {
		p.handleInput()
		if IsKeyPressed(fmt.Sprintf("p%d-placeBomb", p.Number)) {
			i := int(p.Position.X / TILE_LENGTH)
			j := int(p.Position.Y / TILE_LENGTH)
			p.PlaceBomb(i, j)
		}
	}
	p.move()
}

func (p *Player) PlaceBomb(i, j int) {
	if p.upgrades.bombCount <= 0 {
		return
	}
	// fmt.Println(v)
	if i >= 0 && i < TileManager.XCount && j >= 0 && j < TileManager.YCount {
		core.GetSound("place-bomb").Play()
		b := NewBombTile(i, j, p.upgrades.bombRadius, p).AddToTileSet()
		p.LastBombX = i
		p.LastBombY = j
		if b != nil {
			p.upgrades.bombCount--
		}
	}
}

func (p *Player) Draw() {
	p.Lock()
	defer p.Unlock()
	if p.IsAlive {
		switch p.Direction {
		case UP:
			p.Animations.WalkingUp.DrawCenteredAtWithScale(p.Position.X, p.Position.Y, p.Scale)
		case RIGHT:
			p.Animations.WalkingRight.DrawCenteredAtWithScale(p.Position.X, p.Position.Y, p.Scale)
		case DOWN:
			p.Animations.WalkingDown.DrawCenteredAtWithScale(p.Position.X, p.Position.Y, p.Scale)
		case LEFT:
			p.Animations.WalkingLeft.DrawCenteredAtWithScale(p.Position.X, p.Position.Y, p.Scale)
		}
	} else {
		if time.Now().UnixMilli()-p.diedAt < 500 {
			p.Animations.dying.DrawCenteredAtWithScale(p.Position.X, p.Position.Y, p.Scale)
		} else if time.Now().UnixMilli()-p.diedAt < 1500 {
			if !p.Animations.crying1.IsPlaying() {
				p.Animations.crying1.Play()
			}
			p.Animations.crying1.DrawCenteredAtWithScale(p.Position.X, p.Position.Y, p.Scale)
		} else if time.Now().UnixMilli()-p.diedAt < 3200 {
			if !p.Animations.crying2.IsPlaying() {
				p.Animations.crying2.Play()
			}
			p.Animations.crying2.DrawCenteredAtWithScale(p.Position.X, p.Position.Y, p.Scale)
		}
	}
	// p.PlayerAnimations.WalkingDown.DrawAt(ray.NewRectangle(20, 20, 16*3, 25*3))
	// p.PlayerAnimations.WalkingDown.DrawAt(ray.NewRectangle(20, 20, 16*3, 25*3))
	// p.PlayerAnimations.WalkingLeft.DrawCenteredAt(p.position.X, p.position.Y)
	if ShowCollsions {
		ray.DrawRectangleLinesEx(p.GetCollision(), 2, ray.Red)
	}
}

func (p *Player) GetCollision() ray.Rectangle {
	return ray.NewRectangle(p.Position.X-8*p.Scale, p.Position.Y-12*p.Scale, 16*p.Scale, 24*p.Scale)
}

func (p *Player) Enabled() bool {
	return p.HasCollision
}

func (p *Player) collides() (Collision, bool) {
	t := CollisionManager.TryLock()
	if t {
		defer CollisionManager.Unlock()
	}
	if CollisionManager.objects == nil {
		return nil, false
	}
	playerCollisionRec := p.GetCollision()
	for _, collision := range CollisionManager.objects {
		if !collision.Enabled() {
			continue
		}
		switch t := collision.(type) {
		case *FallingTile:
			if ray.CheckCollisionRecs(playerCollisionRec, collision.GetCollision()) {
				return collision, true
			}
		case *WallTile:
			if !t.isDestoryable {
				if ray.CheckCollisionRecs(playerCollisionRec, collision.GetCollision()) {
					return collision, true
				}
			} else if !p.upgrades.passWalls {
				if ray.CheckCollisionRecs(playerCollisionRec, collision.GetCollision()) {
					return collision, true
				}
			}
		case *BombTile:
			if !p.upgrades.passBombs && ray.CheckCollisionRecs(playerCollisionRec, collision.GetCollision()) {
				return collision, true
			}
		case *UpgradeTile:
			if ray.CheckCollisionRecs(playerCollisionRec, collision.GetCollision()) {
				p.pickupUpgrade(t)
			}
		}
	}
	return nil, false
}

func (p *Player) pickupUpgrade(t *UpgradeTile) {
	if t.isDestroyed.Load() {
		return
	}
	core.GetSound("item-get").Play()
	switch t.Type {
	case UPGRADE_ADD_SPEED:
		p.Speed = ray.Clamp(p.Speed+20, p.Speed, maxSpeed)
	case UPGRADE_ADD_FIRE:
		p.bombRadius = int(ray.Clamp(float32(p.bombRadius+1), float32(p.bombRadius), maxFire))
	case UPGRADE_FULL_FIRE:
		p.bombRadius = maxFire
	case UPGRADE_ADD_BOMB:
		p.bombCount = int(ray.Clamp(float32(p.bombCount+1), float32(p.bombCount), maxBomb))
	case UPGRADE_FULL_BOMB:
		p.upgrades.bombCount = maxBomb
	case UPGRADE_ADD_HEALTH:
		p.upgrades.heartsCount += 1
	case UPGRADE_PASS_BOMB:
		p.upgrades.passBombs = true
	case UPGRADE_PASS_WALL:
		p.upgrades.passWalls = true
	case UPGRADE_KICK_BOMB:
		p.upgrades.kickBomb = true
	}
	t.Remove()
}

func (p *Player) Die() {
	if !p.IsAlive {
		return
	}
	p.IsAlive = false
	p.HasCollision = false
	p.diedAt = time.Now().UnixMilli()
	p.Animations.dying.Play()
}
