package running

import (
	"fmt"
	"github.com/rzaf/bomberman-clone/core"

	ray "github.com/gen2brain/raylib-go/raylib"
)

var (
	p1Texture    *core.Texture
	p2Texture    *core.Texture
	clockTexture *core.Texture
	timeStr      string
)

func loadUi() {
	p1Texture = core.GetTexture("tiles").Crop(ray.NewRectangle(291, 0, 15, 14))
	p2Texture = core.GetTexture("tiles").Crop(ray.NewRectangle(306, 0, 15, 14))
	clockTexture = core.GetTexture("tiles").Crop(ray.NewRectangle(321, 0, 20, 20))
}

func updateUi() {
	if roundTimeSeconds-int(elapsedTime) >= 0 {
		t := roundTimeSeconds - int(elapsedTime)
		timeStr = fmt.Sprintf("%d:%02d", t/60, t%60)
	}
}

func drawUi() {
	p1Texture.DrawAt(ray.NewRectangle(4, 20, 30, 28))
	p2Texture.DrawAt(ray.NewRectangle(4, 70, 30, 28))
	clockTexture.DrawAt(ray.NewRectangle(4, 110, 30, 30))
	p1.Lock()
	p2.Lock()
	ray.DrawText(fmt.Sprintf("%d", p1.Wins), 38, 26, 25, ray.White)
	ray.DrawText(fmt.Sprintf("%d", p2.Wins), 38, 76, 25, ray.Black)
	p1.Unlock()
	p2.Unlock()
	ray.DrawText(timeStr, 38, 115, 25, ray.White)
}
