package setting

import (
	"fmt"
	ray "github.com/gen2brain/raylib-go/raylib"
	"github.com/rzaf/bomberman-clone/game"
	"github.com/rzaf/bomberman-clone/states/running"
)

var (
	currentAudioMenuIndex int
	masterV               float32
	musicV                float32
	effectV               float32
)

func updateAudio() {
	if game.IsKeyPressed("back") || ray.IsKeyPressed(ray.KeyBackspace) {
		isSelected = false
		changeTexts()
		return
	}
	if game.IsKeyPressed("p1-Down") {
		currentAudioMenuIndex = (currentAudioMenuIndex + 1) % 4
	} else if game.IsKeyPressed("p1-Up") {
		if currentAudioMenuIndex == 0 {
			currentAudioMenuIndex = 3
		} else {
			currentAudioMenuIndex -= 1
		}
	}

	var rl float32 = 0
	if game.IsKeyPressed("p1-Right") {
		rl = 1
	} else if game.IsKeyPressed("p1-Left") {
		rl = -1
	}

	switch currentAudioMenuIndex {
	case 0:
		masterV = ray.Clamp(masterV+rl*0.1, 0, 1)
	case 1:
		musicV = ray.Clamp(musicV+rl*0.1, 0, 1)
	case 2:
		effectV = ray.Clamp(effectV+rl*0.1, 0, 1)
	case 3:
		if game.IsKeyPressed("accept") {
			isSelected = false
			changeTexts()
			if masterV != game.MasterVolume {
				ray.SetMasterVolume(masterV)
				game.MasterVolume = masterV
			}
			if musicV != game.MusicVolume || effectV != game.EffectVolume {
				game.MusicVolume = musicV
				game.EffectVolume = effectV
				running.SetSoundsVolume()
			}
			Save()
		}
	}
}

func drawAudio() {
	if isSelected {
		switch currentAudioMenuIndex {
		case 0:
			ray.DrawText("master volume:", tlX+80, 100, 35, ray.Red)
			ray.DrawText("music volume:", tlX+100, 150, 30, ray.White)
			ray.DrawText("effects volume:", tlX+100, 200, 30, ray.White)
			ray.DrawText("APPLY", tlX+150, 250, 30, ray.White)
		case 1:
			ray.DrawText("master volume:", tlX+100, 100, 30, ray.White)
			ray.DrawText("music volume:", tlX+80, 150, 35, ray.Red)
			ray.DrawText("effects volume:", tlX+100, 200, 30, ray.White)
			ray.DrawText("APPLY", tlX+150, 250, 30, ray.White)
		case 2:
			ray.DrawText("master volume:", tlX+100, 100, 30, ray.White)
			ray.DrawText("music volume:", tlX+100, 150, 30, ray.White)
			ray.DrawText("effects volume:", tlX+80, 200, 35, ray.Red)
			ray.DrawText("APPLY", tlX+150, 250, 30, ray.White)
		case 3:
			ray.DrawText("master volume:", tlX+100, 100, 30, ray.White)
			ray.DrawText("music volume:", tlX+100, 150, 30, ray.White)
			ray.DrawText("effects volume:", tlX+100, 200, 30, ray.White)
			ray.DrawText("APPLY", tlX+150, 250, 35, ray.Red)
		}
	} else {
		ray.DrawText("master volume:", tlX+100, 100, 30, ray.White)
		ray.DrawText("music volume:", tlX+100, 150, 30, ray.White)
		ray.DrawText("effects volume:", tlX+100, 200, 30, ray.White)
	}
	ray.DrawText(fmt.Sprintf("%.2f", masterV), tlX+360, 100, 30, ray.White)
	ray.DrawText(fmt.Sprintf("%.2f", musicV), tlX+360, 150, 30, ray.White)
	ray.DrawText(fmt.Sprintf("%.2f", effectV), tlX+360, 200, 30, ray.White)
}
