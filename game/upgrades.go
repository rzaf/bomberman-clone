package game

import (
	"fmt"
	"github.com/rzaf/bomberman-clone/core"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	ray "github.com/gen2brain/raylib-go/raylib"
)

type UpgradeType uint8

const (
	UPGRADE_ADD_BOMB UpgradeType = iota
	UPGRADE_ADD_FIRE
	UPGRADE_FULL_BOMB
	UPGRADE_FULL_FIRE
	UPGRADE_PASS_BOMB
	UPGRADE_PASS_WALL
	UPGRADE_ADD_SPEED
	UPGRADE_KICK_BOMB
	UPGRADE_BLUEGLOVE
	UPGRADE_ADD_HEALTH
	UPGRADE_NONE
	UPGRADE_COUNT
)

var (
	upgradeChances    []int = nil
	upgradeChancesSum int   = 0

	SendingUpgradeLock sync.Mutex
	SendingUpgrades    []*UpgradeTile = nil
	SaveUpgrades       bool           = false
)

func initUpgradeChances() {
	upgradeChances = make([]int, UPGRADE_COUNT)
	upgradeChances[UPGRADE_ADD_BOMB] = 25
	upgradeChances[UPGRADE_ADD_FIRE] = 25
	upgradeChances[UPGRADE_FULL_BOMB] = 1
	upgradeChances[UPGRADE_FULL_FIRE] = 1
	upgradeChances[UPGRADE_PASS_BOMB] = 1
	upgradeChances[UPGRADE_PASS_WALL] = 1
	upgradeChances[UPGRADE_ADD_SPEED] = 25
	upgradeChances[UPGRADE_KICK_BOMB] = 0
	upgradeChances[UPGRADE_BLUEGLOVE] = 0
	upgradeChances[UPGRADE_ADD_HEALTH] = 0
	upgradeChances[UPGRADE_NONE] = 50
	for i := 0; i < len(upgradeChances); i++ {
		upgradeChancesSum += int(upgradeChances[i])
	}
}

type UpgradeTile struct {
	Tile
	Type UpgradeType
	// isDestroyed bool
	isDestroyed atomic.Bool

	destroyingAnimation *core.Animation
}

func NewUpgradeTile(i, j int, Type UpgradeType) *UpgradeTile {
	t := &UpgradeTile{
		Tile: Tile{i, j, i * TILE_LENGTH, j * TILE_LENGTH, TILE_LENGTH, TILE_LENGTH, nil, true},
		Type: Type,
		// isDestroyed: false,
	}
	t.isDestroyed.Store(false)
	switch Type {
	case UPGRADE_ADD_BOMB:
		t.texture = core.GetTexture("tiles").Crop(ray.NewRectangle(372, 0, 16, 16))
	case UPGRADE_ADD_FIRE:
		t.texture = core.GetTexture("tiles").Crop(ray.NewRectangle(388, 0, 16, 16))
	case UPGRADE_FULL_BOMB:
		t.texture = core.GetTexture("tiles").Crop(ray.NewRectangle(404, 0, 16, 16))
	case UPGRADE_FULL_FIRE:
		t.texture = core.GetTexture("tiles").Crop(ray.NewRectangle(420, 0, 16, 16))
	case UPGRADE_PASS_BOMB:
		t.texture = core.GetTexture("tiles").Crop(ray.NewRectangle(436, 0, 16, 16))
	case UPGRADE_PASS_WALL:
		t.texture = core.GetTexture("tiles").Crop(ray.NewRectangle(372, 16, 16, 16))
	case UPGRADE_ADD_SPEED:
		t.texture = core.GetTexture("tiles").Crop(ray.NewRectangle(388, 16, 16, 16))
	case UPGRADE_KICK_BOMB:
		t.texture = core.GetTexture("tiles").Crop(ray.NewRectangle(404, 16, 16, 16))
	case UPGRADE_BLUEGLOVE:
		t.texture = core.GetTexture("tiles").Crop(ray.NewRectangle(404, 16, 16, 16))
	case UPGRADE_ADD_HEALTH:
		t.texture = core.GetTexture("tiles").Crop(ray.NewRectangle(372, 32, 16, 16))
	}
	t.destroyingAnimation = core.NewAnimation(
		core.NewSpriteSheet(core.GetTexture("tiles"), 5, 1, 5, 16, 16, 0, 1, ray.NewRectangle(272, 51, 16, 85), true, false),
		100,
		false,
	)
	CollisionManager.Add(t)
	return t
}

func GetRandomUpgrade() UpgradeType {
	if upgradeChances == nil {
		initUpgradeChances()
		rand.Seed(time.Now().Unix())
	}
	randN := rand.Intn(upgradeChancesSum)
	for i := range upgradeChances {
		if randN < upgradeChances[i] {
			return UpgradeType(i)
		}
		randN -= upgradeChances[i]
	}
	// r := rand.Intn(100)
	// if r < 50 {
	// 	return UPGRADE_NONE
	// }
	// r = rand.Intn(100)
	// if r < 20 {
	// 	return UPGRADE_ADD_BOMB
	// } else if r < 40 {
	// 	return UPGRADE_ADD_FIRE
	// } else if r < 60 {
	// 	return UPGRADE_ADD_SPEED
	// }
	return UPGRADE_ADD_SPEED
}

func AddUpgradeTile(i, j int, t UpgradeType) {
	fmt.Println("random upgrade:", t)
	if t == UPGRADE_NONE {
		TileManager.Add(NewFloorTile(i, j))
	} else {
		ut := NewUpgradeTile(i, j, t)
		TileManager.Add(ut)
		if SaveUpgrades {
			SendingUpgradeLock.Lock()
			SendingUpgrades = append(SendingUpgrades, ut)
			SendingUpgradeLock.Unlock()
		}
	}
}

func (u *UpgradeTile) Update() {
	if u.isDestroyed.Load() {
		if !u.destroyingAnimation.IsPlaying() {
			u.Remove()
			u.destroyingAnimation.Play()
		} else {
			u.destroyingAnimation.Update()
		}
	}
}

func (u *UpgradeTile) Draw() {
	if u.destroyingAnimation.IsPlaying() {
		u.destroyingAnimation.DrawAt(ray.NewRectangle(float32(u.X), float32(u.Y), float32(u.Width), float32(u.Height)))
	} else {
		if !u.isDestroyed.Load() {
			u.texture.DrawAt(ray.NewRectangle(float32(u.X), float32(u.Y), float32(u.Width), float32(u.Height)))
		}
	}
	if ShowCollisions && u.HasCollision {
		ray.DrawRectangleLinesEx(u.GetCollision(), 2, ray.Blue)
	}
}

func (u *UpgradeTile) Destroy() {
	// u.destroyingAnimation.Play()
	// u.Remove()
	u.isDestroyed.Store(true)
}
func (u *UpgradeTile) Remove() {
	CollisionManager.Remove(u)
	TileManager.Remove(u.I, u.J)
	TileManager.Add(NewFloorTile(u.I, u.J))
}
