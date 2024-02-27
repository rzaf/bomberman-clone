package setting

import (
	"time"

	ray "github.com/gen2brain/raylib-go/raylib"
	"github.com/rzaf/bomberman-clone/game"
)

var (
	actions = []string{
		"pause",
		"accept",
		"back",
		"p1-Up",
		"p1-Right",
		"p1-Down",
		"p1-Left",
		"p1-placeBomb",
		"p2-Up",
		"p2-Right",
		"p2-Down",
		"p2-Left",
		"p2-placeBomb",
	}
	currentI, currentJ = 0, 0
	isWaiting          = false
)

func updateKeyboardMapping() {
	if isWaiting {
		updateWaiting()
		return
	}
	if game.IsKeyPressed("back") {
		isSelected = false
		changeTexts()
		Load()
		return
	}
	if game.IsKeyPressed("p1-Right") || game.IsKeyPressed("p1-Left") {
		if currentI == 1 {
			currentI = 0
		} else {
			currentI = 1
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
			if currentI == 1 {
				Save()
			} else {
				Load()
			}
			isSelected = false
			changeTexts()
		} else {
			isWaiting = true
			errMessage = ""
			enteredT = time.Now().Second()
		}
	} else if ray.IsKeyPressed(ray.KeyF1) {
		if currentJ < len(actions) {
			if currentI == 0 {
				game.AddFirstKey(actions[currentJ], game.NewKey(game.Keyboard, 0))
			} else {
				game.AddSecondryKey(actions[currentJ], game.NewKey(game.Keyboard, 0))
			}
		}
	} else if ray.IsKeyPressed(ray.KeyF2) {
		game.LoadDefaultKeys()
	}
}

func drawKeyboardMapping() {
	ray.DrawText("Actions", tlX+120, 50, 20, ray.White)
	ray.DrawText("keyboard", tlX+260, 50, 20, ray.White)
	ray.DrawText("secondry", tlX+390, 50, 20, ray.White)
	for i := 0; i < len(actions); i++ {
		ray.DrawText(actions[i], tlX+120, 100+int32(i)*20, 18, ray.White)
	}
	for i := 0; i < len(actions); i++ {
		if isSelected && currentI == 0 && currentJ == i {
			ray.DrawText(GetKeyName(game.GetFirstKey(actions[i])), tlX+270, 100+int32(i)*20, 18, ray.Red)
		} else {
			ray.DrawText(GetKeyName(game.GetFirstKey(actions[i])), tlX+270, 100+int32(i)*20, 18, ray.White)
		}
	}
	for i := 0; i < len(actions); i++ {
		if isSelected && currentI == 1 && currentJ == i {
			ray.DrawText(GetKeyName(game.GetSecondryKey(actions[i])), tlX+400, 100+int32(i)*20, 18, ray.Red)
		} else {
			ray.DrawText(GetKeyName(game.GetSecondryKey(actions[i])), tlX+400, 100+int32(i)*20, 18, ray.White)
		}
	}

	if isSelected && currentI == 0 && currentJ == len(actions) {
		ray.DrawText("Cancel", tlX+270, 110+int32(len(actions))*20, 23, ray.Red)
	} else {
		ray.DrawText("Cancel", tlX+270, 110+int32(len(actions))*20, 23, ray.White)
	}
	if isSelected && currentI == 1 && currentJ == len(actions) {
		ray.DrawText("Apply", tlX+400, 110+int32(len(actions))*20, 23, ray.Red)
	} else {
		ray.DrawText("Apply", tlX+400, 110+int32(len(actions))*20, 23, ray.White)
	}

	if isSelected {
		ray.DrawText("press F1 to reset current key", 50, int32(game.Height)-55, 20, ray.White)
		ray.DrawText("press F2 to reset to default", 50, int32(game.Height)-25, 20, ray.White)
	}
	if isWaiting {
		drawWaiting()
	}
}

func GetKeyName(gKey game.Key) string {
	key := gKey.Keyid
	if gKey.KeyType == game.Keyboard {
		switch key {
		case ray.KeyZero:
			return "Zero"
		case ray.KeyOne:
			return "One"
		case ray.KeyTwo:
			return "Two"
		case ray.KeyThree:
			return "Three"
		case ray.KeyFour:
			return "Four"
		case ray.KeyFive:
			return "Five"
		case ray.KeySix:
			return "Six"
		case ray.KeySeven:
			return "Seven"
		case ray.KeyEight:
			return "Eight"
		case ray.KeyNine:
			return "Nine"
		case ray.KeySemicolon:
			return "Semicolon"
		case ray.KeyEqual:
			return "="
		case ray.KeyA:
			return "A"
		case ray.KeyB:
			return "B"
		case ray.KeyC:
			return "C"
		case ray.KeyD:
			return "D"
		case ray.KeyE:
			return "E"
		case ray.KeyF:
			return "F"
		case ray.KeyG:
			return "G"
		case ray.KeyH:
			return "H"
		case ray.KeyI:
			return "I"
		case ray.KeyJ:
			return "J"
		case ray.KeyK:
			return "K"
		case ray.KeyL:
			return "L"
		case ray.KeyM:
			return "M"
		case ray.KeyN:
			return "N"
		case ray.KeyO:
			return "O"
		case ray.KeyP:
			return "P"
		case ray.KeyQ:
			return "Q"
		case ray.KeyR:
			return "R"
		case ray.KeyS:
			return "S"
		case ray.KeyT:
			return "T"
		case ray.KeyU:
			return "U"
		case ray.KeyV:
			return "V"
		case ray.KeyW:
			return "W"
		case ray.KeyX:
			return "X"
		case ray.KeyY:
			return "Y"
		case ray.KeyZ:
			return "Z"
		///numpad
		case ray.KeyKp0:
			return "Num 0"
		case ray.KeyKp1:
			return "Num 1"
		case ray.KeyKp2:
			return "Num 2"
		case ray.KeyKp3:
			return "Num 3"
		case ray.KeyKp4:
			return "Num 4"
		case ray.KeyKp5:
			return "Num 5"
		case ray.KeyKp6:
			return "Num 6"
		case ray.KeyKp7:
			return "Num 7"
		case ray.KeyKp8:
			return "Num 8"
		case ray.KeyKp9:
			return "Num 9"
		case ray.KeyKpDecimal:
			return "Num ."
		case ray.KeyKpDivide:
			return "Num /"
		case ray.KeyKpMultiply:
			return "Num *"
		case ray.KeyKpSubtract:
			return "Num -"
		case ray.KeyKpAdd:
			return "Num +"
		case ray.KeyKpEnter:
			return "Num Enter"
		case ray.KeyKpEqual:
			return "Num ="
		case ray.KeyBackspace:
			return "Backspace"
		//// invalid keys
		case ray.KeyCapsLock, ray.KeyScrollLock, ray.KeyNumLock, ray.KeyPrintScreen, ray.KeyPause, ray.KeyF1, ray.KeyF2, ray.KeyF3, ray.KeyF4, ray.KeyF5, ray.KeyF6, ray.KeyF7, ray.KeyF8, ray.KeyF9, ray.KeyF10, ray.KeyF11, ray.KeyF12:
			return "-"
		case ray.KeySpace:
			return "Space"
		case ray.KeyEscape:
			return "Escape"
		case ray.KeyEnter:
			return "Enter"
		case ray.KeyTab:
			return "Tab"
		case ray.KeyInsert:
			return "Insert"
		case ray.KeyDelete:
			return "Delete"
		case ray.KeyRight:
			return "Right"
		case ray.KeyLeft:
			return "Left"
		case ray.KeyDown:
			return "Down"
		case ray.KeyUp:
			return "Up"
		case ray.KeyPageUp:
			return "PageUp"
		case ray.KeyPageDown:
			return "PageDown"
		case ray.KeyHome:
			return "Home"
		case ray.KeyEnd:
			return "End"
		case ray.KeyLeftShift:
			return "LeftShift"
		case ray.KeyLeftControl:
			return "LeftControl"
		case ray.KeyLeftAlt:
			return "LeftAlt"
		case ray.KeyLeftSuper:
			return "LeftSuper"
		case ray.KeyRightShift:
			return "RightShift"
		case ray.KeyRightControl:
			return "RightControl"
		case ray.KeyRightAlt:
			return "RightAlt"
		case ray.KeyRightSuper:
			return "RightSuper"
		case ray.KeyKbMenu:
			return "KbMenu"
		case ray.KeyLeftBracket:
			return "LeftBracket"
		case ray.KeyBackSlash:
			return "BackSlash"
		case ray.KeyRightBracket:
			return "RightBracket"
		case ray.KeyGrave:
			return "Grave"
		}
	} else if gKey.KeyType == game.Mouse {
		switch key {
		case ray.MouseButtonLeft:
			return "LeftMouseButton"
		case ray.MouseButtonRight:
			return "RightMouseButton"
		case ray.MouseButtonMiddle:
			return "MiddleMouseButton"
		case ray.MouseButtonSide:
			return "SideMouseButton"
		case ray.MouseButtonExtra:
			return "ExtraMouseButton"
		case ray.MouseButtonForward:
			return "ForwardMouseButton"
		case ray.MouseButtonBack:
			return "BackMouseButton"
		}
	}
	return "-"
}
