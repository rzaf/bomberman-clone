package core

import (

	// "fmt"

	ray "github.com/gen2brain/raylib-go/raylib"
)

type Animation struct {
	spriteSheet *SprtieSheet
	srcs        []*ray.Rectangle
	Index       int
	Speed       int64 // ms per sprite
	Loop        bool
	elapsedTime int64
}

func (anim *Animation) loadTexturesFromSpriteSheet(sh *SprtieSheet) {
	anim.srcs = make([]*ray.Rectangle, sh.count)
	for i := 0; i < sh.count; i++ {
		src := ray.NewRectangle(0, 0, float32(sh.width), float32(sh.height))
		if !sh.isVertical {
			src.X = sh.src.X + float32(i%sh.rowsCount)*float32(sh.width) + float32(i*sh.seperationX)
			src.Y = sh.src.Y + float32(i/sh.rowsCount)*float32(sh.height) + float32(i*sh.seperationY)
		} else {
			src.X = sh.src.X + float32(i/sh.columnsCount)*float32(sh.width) + float32(i*sh.seperationX)
			src.Y = sh.src.Y + float32(i%sh.columnsCount)*float32(sh.height) + float32(i*sh.seperationY)
		}
		// fmt.Println(src)
		if sh.reverse {
			anim.srcs[sh.count-i-1] = &src
		} else {
			anim.srcs[i] = &src
		}
	}
}

func NewAnimation(spriteSheet *SprtieSheet, speed int64, loop bool) *Animation {
	newAnim := &Animation{
		spriteSheet: spriteSheet,
		srcs:        nil,
		Index:       0,
		Speed:       speed,
		Loop:        loop,
		elapsedTime: 0,
	}
	newAnim.loadTexturesFromSpriteSheet(spriteSheet)
	return newAnim
}

func (a *Animation) Play() {
	a.elapsedTime = 1
	// a.elapsedTime = time.Now().UnixMilli()
}

func (a *Animation) Stop() {
	a.elapsedTime = 0
}

func (a *Animation) IsPlaying() bool {
	return a.elapsedTime >= 1
}

func (a *Animation) Update() {
	if !a.Loop && a.Index == len(a.srcs)-1 {
		a.Stop()
	}
	// fmt.Println(a.elapsedTime, ray.GetFrameTime()*1000, a.Speed)

	if a.elapsedTime > 0 {
		a.elapsedTime += int64(ray.GetFrameTime() * 1000)
	}

	// if a.startTime > 0 && time.Now().UnixMilli()-a.startTime > a.Speed {
	if a.elapsedTime > 0 && a.elapsedTime > a.Speed {
		a.elapsedTime = 1

		// a.elapsedTime = time.Now().UnixMilli()
		a.Index = (a.Index + 1) % a.spriteSheet.count
		// fmt.Println("current sprite:", a.Index, a.srcs[a.Index])
	}
}

func (a *Animation) DrawAt(dst ray.Rectangle) {
	ray.DrawTexturePro(a.spriteSheet.texture.Texture, *a.srcs[a.Index], dst, ray.NewVector2(0, 0), 0, ray.White)
}

// func (a *Animation) DrawAtWithScale(v ray.Vector2, scale float32) {
// ray.DrawTexturePro(a.spriteSheet.texture.Texture, *a.srcs[a.Index], dst, ray.NewVector2(0, 0), 0, ray.White)
// ray.DrawTextureEx(a.spriteSheet.texture.Texture, v, 0, scale, ray.White)
// }

func (a *Animation) DrawCenteredAt(x, y float32) {
	w := float32(a.spriteSheet.width)
	h := float32(a.spriteSheet.height)
	ray.DrawTexturePro(
		a.spriteSheet.texture.Texture,
		*a.srcs[a.Index],
		ray.NewRectangle(x+w/2, y+h/2, w, h),
		ray.NewVector2(0, 0),
		0,
		ray.White,
	)
}

func (a *Animation) DrawCenteredAtWithScale(x, y, s float32) {
	w := float32(a.spriteSheet.width) * s
	h := float32(a.spriteSheet.height) * s
	ray.DrawTexturePro(
		a.spriteSheet.texture.Texture,
		*a.srcs[a.Index],
		ray.NewRectangle(x-w/2, y-h/2, w, h),
		ray.NewVector2(0, 0),
		0,
		ray.White,
	)
	// ray.DrawCircle(int32(x), int32(y), 2.5, ray.Red)
}
