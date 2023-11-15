package running

import (
	"bomberman/core"
	"bomberman/game"
	"strings"

	ray "github.com/gen2brain/raylib-go/raylib"
)

var (
	isHost         bool = true
	isGuestWaiting bool = false
	isHostWaiting  bool = false

	hostText     *core.Text
	joinText     *core.Text
	waitingText1 *core.Text
	waitingText2 *core.Text

	addrs []string = []string{"127", "0", "0", "1"}

	addrIndex int = 0
)

type OnlineMenu struct{}

func (OnlineMenu) OnEnter() {
	w := float32(ray.GetScreenWidth())
	h := float32(ray.GetScreenHeight())

	isHost = true
	hostText = core.NewText("HOST", ray.GetFontDefault(), ray.NewVector2(w/2, h/2-50), 45, 8, ray.Red)
	joinText = core.NewText("JOIN", ray.GetFontDefault(), ray.NewVector2(w/2, h/2+50), 45, 8, ray.White)
	waitingText1 = core.NewText("WAITING FOR HOST", ray.GetFontDefault(), ray.NewVector2(w/2, 100), 50, 8, ray.White)
	waitingText2 = core.NewText("WAITING FOR GUEST", ray.GetFontDefault(), ray.NewVector2(w/2, 100), 50, 8, ray.White)
}

func (OnlineMenu) OnExit() {
	isGuestWaiting = false
	isHostWaiting = false

	if game.State.Get() != game.ONLINE_BATTLE {
		if isHost {
			DisconnectServer()
		} else {
			CloseClient()
		}
	}
}

func (OnlineMenu) OnWindowResized() {
	w := float32(ray.GetScreenWidth())
	h := float32(ray.GetScreenHeight())
	hostText.Pos = ray.NewVector2(w/2, h/2-50)
	joinText.Pos = ray.NewVector2(w/2, h/2+50)
	waitingText1.Pos = ray.NewVector2(w/2, 100)
}

func (OnlineMenu) Update() {
	if ray.IsKeyPressed(ray.KeyEscape) {
		game.State.Change(game.MENU)
	}

	if isHostWaiting {
		return
	}
	if isGuestWaiting {
		return
	}

	if ray.IsKeyPressed(ray.KeyUp) || ray.IsKeyPressed(ray.KeyDown) {
		if isHost {
			isHost = false
			hostText.Color = ray.White
			joinText.Color = ray.Red
		} else {
			isHost = true
			hostText.Color = ray.Red
			joinText.Color = ray.White
		}
	}
	if ray.IsKeyPressed(ray.KeyEscape) {
		game.State.Change(game.MENU)
		// quitOnlineMenu()
	}
	if ray.IsKeyPressed(ray.KeyEnter) {
		if isHost {
			game.State.Change(game.BATTLE_MENU)
			go host()
		} else {
			isGuestWaiting = true
			connectToServer()
			go sendInfoReq()
		}
	}
	if !isHost {
		k := ray.GetKeyPressed()
		ok := true
		switch k {
		case ray.KeyZero, ray.KeyOne, ray.KeyTwo, ray.KeyThree, ray.KeyFour, ray.KeyFive, ray.KeySix, ray.KeySeven, ray.KeyEight, ray.KeyNine:
		case ray.KeyBackspace, ray.KeyTab:
		case ray.KeyKp0:
			k = ray.KeyZero
		case ray.KeyKp1:
			k = ray.KeyOne
		case ray.KeyKp2:
			k = ray.KeyTwo
		case ray.KeyKp3:
			k = ray.KeyThree
		case ray.KeyKp4:
			k = ray.KeyFour
		case ray.KeyKp5:
			k = ray.KeyFive
		case ray.KeyKp6:
			k = ray.KeySix
		case ray.KeyKp7:
			k = ray.KeySeven
		case ray.KeyKp8:
			k = ray.KeyEight
		case ray.KeyKp9:
			k = ray.KeyNine
		default:
			ok = false
		}
		if ok {
			if k == ray.KeyTab {
				addrIndex = (addrIndex + 1) % 4
				return
			}

			if k == ray.KeyBackspace {
				if len(addrs[addrIndex]) > 0 {
					addrs[addrIndex] = addrs[addrIndex][:len(addrs[addrIndex])-1]
				}
			} else {
				if string(k) == "0" {
					if len(addrs[addrIndex]) == 1 && addrs[addrIndex][0] == '0' {
						return
					}
				}

				if len(addrs[addrIndex]) < 3 {
					addrs[addrIndex] += string(k)
				} else {
					addrs[addrIndex] = addrs[addrIndex][1:] + string(k)
				}
				if len(addrs[addrIndex]) > 1 {
					addrs[addrIndex] = strings.TrimLeft(addrs[addrIndex], "0")
				}
			}
		}
	}
}

func (OnlineMenu) Draw() {
	if isGuestWaiting {
		waitingText1.DrawCentered()
		return
	}
	if isHostWaiting {
		waitingText2.DrawCentered()
		return
	}
	hostText.DrawCentered()
	joinText.DrawCentered()
	if !isHost {
		w := ray.GetScreenWidth()
		h := ray.GetScreenHeight()

		ray.DrawRectangle(int32(w/2-200), int32(h-40), 400, 38, ray.Black)
		ray.DrawText(addrs[0], int32(w/2-150-25), int32(h-35), 30, ray.White)

		ray.DrawText(".", int32(w/2-100), int32(h-35), 30, ray.White)
		ray.DrawText(addrs[1], int32(w/2-50-25), int32(h-35), 30, ray.White)

		ray.DrawText(".", int32(w/2), int32(h-35), 30, ray.White)
		ray.DrawText(addrs[2], int32(w/2+50-25), int32(h-35), 30, ray.White)

		ray.DrawText(".", int32(w/2+100), int32(h-35), 30, ray.White)
		ray.DrawText(addrs[3], int32(w/2+150-25), int32(h-35), 30, ray.White)

		ray.DrawLineEx(ray.NewVector2(float32(w/2-180+addrIndex*100), float32(h-4)), ray.NewVector2(float32(w/2-130+addrIndex*100), float32(h-4)), 3, ray.Red)
	}
}
