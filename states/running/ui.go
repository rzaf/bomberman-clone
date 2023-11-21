package running

import (
	"fmt"
	"github.com/rzaf/bomberman-clone/core"

	ray "github.com/gen2brain/raylib-go/raylib"
)

var (
	p1Texture *core.Texture
	p2Texture *core.Texture
	timeStr   string
)

func loadUi() {
	p1Texture = core.GetTexture("anims").Crop(ray.NewRectangle(64, 0, 17, 14))
	p2Texture = core.GetTexture("anims").Crop(ray.NewRectangle(167, 0, 17, 14))
}

func updateUi() {
	if roundTimeSeconds-int(elapsedTime) >= 0 {
		t := roundTimeSeconds - int(elapsedTime)
		timeStr = fmt.Sprintf("%d:%02d", t/60, t%60)
	}
}

func drawUi() {
	p1Texture.DrawAt(ray.NewRectangle(0, 20, 25, 30))
	p2Texture.DrawAt(ray.NewRectangle(0, 70, 25, 30))

	p1.Lock()
	p2.Lock()
	ray.DrawText(fmt.Sprintf("%d", p1.Wins), 28, 26, 25, ray.White)
	ray.DrawText(fmt.Sprintf("%d", p2.Wins), 28, 76, 25, ray.Black)
	p1.Unlock()
	p2.Unlock()
	ray.DrawText(timeStr, 30, 115, 25, ray.White)
}
