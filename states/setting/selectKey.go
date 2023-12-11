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
	if currentI < 2 { // keyboard
		k := ray.GetKeyPressed()
		if k > 0 {
			fmt.Println(k, GetKeyName(k))
			if GetKeyName(k) != "-" {
				if currentI == 0 {
					game.AddFirstKey(actions[currentJ], int(k))
				} else {
					game.AddSecondryKey(actions[currentJ], int(k))
				}
				isWaiting = false
				errMessage = ""
				return
			}
			errMessage = "invalid key!"
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
				game.AddGamepadKey(currentGamePadId, actions[currentJ], int(k))
				return
			}
			// errMessage = "invalid key!"
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
