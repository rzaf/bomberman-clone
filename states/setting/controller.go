package setting

import (
	"fmt"
	"time"

	ray "github.com/gen2brain/raylib-go/raylib"
	"github.com/rzaf/bomberman-clone/game"
)

var (
	currentGamePadId = 0
)

func updateControllerMapping() {
	if isWaiting {
		updateWaiting()
		return
	}
	if game.IsKeyPressed("back") || ray.IsKeyPressed(ray.KeyBackspace) {
		isSelected = false
		Load()
		changeTexts()
		return
	}
	if game.IsKeyPressed("p1-Right") {
		currentGamePadId = (currentGamePadId + 1) % 4
	} else if game.IsKeyPressed("p1-Left") {
		if currentGamePadId == 0 {
			currentGamePadId = 3
		} else {
			currentGamePadId -= 1
		}
	}

	if game.IsKeyPressed("p1-Down") {
		currentJ = (currentJ + 1) % (len(actions) + 1)
	} else if game.IsKeyPressed("p1-Up") {
		if currentJ == 0 {
			currentJ = len(actions)
		} else {
			currentJ--
		}
	}
	if game.IsKeyPressed("accept") {
		if currentJ == len(actions) {
			Save()
			isSelected = false
			changeTexts()
		} else {
			if ray.IsGamepadAvailable(int32(currentGamePadId)) {
				isWaiting = true
				enteredT = time.Now().Second()
			}
		}
	} else if ray.IsKeyPressed(ray.KeySpace) {
		game.AddGamepadKey(currentGamePadId, actions[currentJ], -1)
	} else if ray.IsKeyPressed(ray.KeyR) {
		game.ResetKeys()
	}
}

func drawControllerMapping() {
	ray.DrawText("Actions", tlX+120, 50, 20, ray.White)
	n := ray.GetGamepadName(int32(currentGamePadId))
	if n == "" {
		n = "not found!"
	}
	ray.DrawText(fmt.Sprintf("controller%d : %s", currentGamePadId+1, n), tlX+260, 50, 20, ray.White)
	for i := 0; i < len(actions); i++ {
		ray.DrawText(actions[i], tlX+120, 100+int32(i)*20, 18, ray.White)
	}
	for i := 0; i < len(actions); i++ {
		if currentI == 2 && currentJ == i {
			ray.DrawText(GetGamepadKeyName(game.GetGamepadKey(currentGamePadId, actions[i])), tlX+320, 100+int32(i)*20, 18, ray.Red)
		} else {
			ray.DrawText(GetGamepadKeyName(game.GetGamepadKey(currentGamePadId, actions[i])), tlX+320, 100+int32(i)*20, 18, ray.White)
		}
	}
	if currentJ == len(actions) {
		ray.DrawText("Apply", tlX+320, 110+int32(len(actions))*20, 23, ray.Red)
	} else {
		ray.DrawText("Apply", tlX+320, 110+int32(len(actions))*20, 23, ray.White)
	}

	if isSelected {
		ray.DrawText("press SPACE to reset key", 50, int32(game.Height)-50, 20, ray.White)
	}
	if isWaiting {
		drawWaiting()
	}
}

func GetGamepadKeyName(key int32) string {
	switch key {
	case ray.GamepadButtonRightFaceDown:
		return "A"
	case ray.GamepadButtonRightFaceRight:
		return "B"
	case ray.GamepadButtonRightFaceLeft:
		return "X"
	case ray.GamepadButtonRightFaceUp:
		return "Y"
	case ray.GamepadButtonLeftTrigger1:
		return "Lb"
	case ray.GamepadButtonLeftTrigger2:
		return "Lt"
	case ray.GamepadButtonRightTrigger1:
		return "Rb"
	case ray.GamepadButtonRightTrigger2:
		return "Rt"
	case ray.GamepadButtonLeftFaceUp:
		return "Up"
	case ray.GamepadButtonLeftFaceRight:
		return "Right"
	case ray.GamepadButtonLeftFaceDown:
		return "Down"
	case ray.GamepadButtonLeftFaceLeft:
		return "Left"
	case ray.GamepadButtonMiddleLeft:
		return "Select"
	case ray.GamepadButtonMiddleRight:
		return "Start"
	case ray.GamepadButtonMiddle:
		return "Home"
		// case ray.GamepadXboxAxisLt:
		// 	return "AxisLt"
	}
	return "-"
}
