package game

import (
	ray "github.com/gen2brain/raylib-go/raylib"
	"github.com/rzaf/bomberman-clone/core"
)

type PlayerAnimations struct {
	WalkingLeft  *core.Animation
	WalkingUp    *core.Animation
	WalkingRight *core.Animation
	WalkingDown  *core.Animation
	dying        *core.Animation
	crying1      *core.Animation
	crying2      *core.Animation
	// winning      *core.Animation
}

func (pa *PlayerAnimations) load(playerNumber int) {
	t1 := core.GetTexture("anims")

	var offset float32 = 0
	if playerNumber == 2 {
		offset += 103
	}
	pa.WalkingLeft = core.NewAnimation(
		core.NewSpriteSheet(t1, 3, 3, 1, 16, 25, 0, 0, ray.NewRectangle(offset, 0, 48, 25), false, false),
		200,
		true,
	)
	pa.WalkingRight = core.NewAnimation(
		core.NewSpriteSheet(t1, 3, 3, 1, 16, 25, 0, 0, ray.NewRectangle(offset, 26, 48, 25), false, false),
		200,
		true,
	)
	pa.WalkingUp = core.NewAnimation(
		core.NewSpriteSheet(t1, 3, 3, 1, 16, 25, 0, 0, ray.NewRectangle(48+offset, 26, 48, 25), false, false),
		200,
		true,
	)
	pa.WalkingDown = core.NewAnimation(
		core.NewSpriteSheet(t1, 3, 3, 1, 16, 25, 0, 0, ray.NewRectangle(48+offset, 0, 48, 25), false, false),
		200,
		true,
	)
	pa.dying = core.NewAnimation(
		core.NewSpriteSheet(t1, 4, 4, 1, 16, 25, 0, 0, ray.NewRectangle(offset, 51, 63, 26), false, false),
		60,
		true,
	)
	pa.crying1 = core.NewAnimation(
		core.NewSpriteSheet(t1, 3, 3, 1, 16, 24, 0, 0, ray.NewRectangle(offset, 78, 48, 24), false, false),
		200,
		false,
	)
	pa.crying2 = core.NewAnimation(
		core.NewSpriteSheet(t1, 2, 2, 1, 16, 24, 0, 0, ray.NewRectangle(48+offset, 78, 31, 24), false, false),
		200,
		true,
	)
}
