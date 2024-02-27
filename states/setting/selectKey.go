package setting

import (
	"fmt"
	"time"

	ray "github.com/gen2brain/raylib-go/raylib"
	"github.com/rzaf/bomberman-clone/game"
)

const waitingTimout = 5

var (
	errMessage = ""
	enteredT   = time.Now().Second()
)

func updateWaiting() {
	if time.Now().Second()-enteredT >= waitingTimout {
		isWaiting = false
		return
	}
	if currentI < 2 { // keyboard and mouse
		k := ray.GetKeyPressed() //keyboard
		if k > 0 {
			fmt.Println(k, GetKeyName(game.NewKey(game.Keyboard, k)))
			if GetKeyName(game.NewKey(game.Keyboard, k)) != "-" {
				if currentI == 0 {
					game.AddFirstKey(actions[currentJ], game.NewKey(game.Keyboard, int32(k)))
				} else {
					game.AddSecondryKey(actions[currentJ], game.NewKey(game.Keyboard, int32(k)))
				}
				game.SetActionDown(actions[currentJ])
				isWaiting = false
				errMessage = ""
				return
			}
			errMessage = "invalid key!"
		}
		// mouse
		k = -1
		if ray.IsMouseButtonReleased(ray.MouseButtonLeft) {
			k = ray.MouseButtonLeft
		} else if ray.IsMouseButtonReleased(ray.MouseButtonRight) {
			k = ray.MouseButtonRight
		} else if ray.IsMouseButtonReleased(ray.MouseButtonMiddle) {
			k = ray.MouseButtonMiddle
		} else if ray.IsMouseButtonReleased(ray.MouseButtonSide) {
			k = ray.MouseButtonSide
		} else if ray.IsMouseButtonReleased(ray.MouseButtonExtra) {
			k = ray.MouseButtonExtra
		} else if ray.IsMouseButtonReleased(ray.MouseButtonForward) {
			k = ray.MouseButtonForward
		} else if ray.IsMouseButtonReleased(ray.MouseButtonBack) {
			k = ray.MouseButtonBack
		}
		if k >= 0 {
			if currentI == 0 {
				game.AddFirstKey(actions[currentJ], game.NewKey(game.Mouse, int32(k)))
			} else {
				game.AddSecondryKey(actions[currentJ], game.NewKey(game.Mouse, int32(k)))
			}
			game.SetActionDown(actions[currentJ])
			isWaiting = false
			errMessage = ""
			return
		}
	} else { // controller
		if !ray.IsGamepadAvailable(int32(currentGamePadId)) {
			errMessage = fmt.Sprintf("gamepad %d not connected!", currentGamePadId+1)
			return
		} else {
			errMessage = ""
		}
		k := ray.GetGamepadButtonPressed()
		if k > 0 {
			if ray.IsGamepadButtonPressed(int32(currentGamePadId), k) {
				isWaiting = false
				errMessage = ""
				game.AddGamepadKey(currentGamePadId, actions[currentJ], game.NewKey(game.Gamepad, int32(k)))
				game.SetActionDown(actions[currentJ])
				return
			}
			// errMessage = "invalid key!"
		}
		k = -1
		var moveMent float32
		moveMent = ray.GetGamepadAxisMovement(int32(currentGamePadId), ray.GamepadAxisLeftX)
		if moveMent > 0.5 {
			k = game.GamepadAxisLeftXRight
		} else if moveMent < (-0.5) {
			k = game.GamepadAxisLeftXLeft
		}
		moveMent = ray.GetGamepadAxisMovement(int32(currentGamePadId), ray.GamepadAxisLeftY)
		if moveMent > 0.5 {
			k = game.GamepadAxisLeftYDown
		} else if moveMent < (-0.5) {
			k = game.GamepadAxisLeftYUp
		}
		moveMent = ray.GetGamepadAxisMovement(int32(currentGamePadId), ray.GamepadAxisRightX)
		if moveMent > 0.5 {
			k = game.GamepadAxisRightXRight
		} else if moveMent < (-0.5) {
			k = game.GamepadAxisRightXLeft
		}
		moveMent = ray.GetGamepadAxisMovement(int32(currentGamePadId), ray.GamepadAxisRightY)
		if moveMent > 0.5 {
			k = game.GamepadAxisRightYDown
		} else if moveMent < (-0.5) {
			k = game.GamepadAxisRightYUp
		}
		if k != -1 {
			isWaiting = false
			errMessage = ""
			game.AddGamepadKey(currentGamePadId, actions[currentJ], game.NewKey(game.GamepadAxis, int32(k)))
			game.SetActionDown(actions[currentJ])
			return
		}
	}
}

func drawWaiting() {
	ray.DrawRectangle(0, 0, int32(game.Width), int32(game.Height), ray.NewColor(0, 0, 0, 220))
	ray.DrawText("press ...", int32(game.Width)/2-50, int32(game.Height)/2-20, 35, ray.White)
	ray.DrawText(fmt.Sprintf("canceling in %d", (waitingTimout+enteredT)-time.Now().Second()), 25, 20, 35, ray.White)
	if errMessage != "" {
		ray.DrawText(errMessage, int32(game.Width)/2-50, int32(game.Height)/2-50, 25, ray.Red)
	}
}
