package game

import (
	ray "github.com/gen2brain/raylib-go/raylib"
)

type KeyType uint8

const (
	Keyboard KeyType = iota
	Mouse
	Gamepad
	GamepadAxis
)

type Key struct {
	KeyType
	Keyid int32
}

func NewKey(t KeyType, i int32) Key {
	return Key{t, i}
}

const (
	GamepadAxisLeftX = iota
	GamepadAxisLeftXRight
	GamepadAxisLeftXLeft
	GamepadAxisLeftY
	GamepadAxisLeftYUp
	GamepadAxisLeftYDown
	GamepadAxisRightX
	GamepadAxisRightXRight
	GamepadAxisRightXLeft
	GamepadAxisRightY
	GamepadAxisRightYUp
	GamepadAxisRightYDown
	GamepadAxisLeftTrigger
	GamepadAxisRightTrigger
)

type gamePadMap struct {
	keys map[string]Key
	id   int32
}

var (
	/// maps
	firstKeyMap    map[string]Key = make(map[string]Key)
	secondKeyMap   map[string]Key = make(map[string]Key)
	gamePadKeyMaps []gamePadMap   = make([]gamePadMap, 4)
	/// states
	lastKeyStates = make(map[string]bool)
	keyStates     = make(map[string]bool)
)

func ResetKeys() {
	firstKeyMap = make(map[string]Key)
	secondKeyMap = make(map[string]Key)
	gamePadKeyMaps[0] = gamePadMap{id: 0, keys: make(map[string]Key)}
	gamePadKeyMaps[1] = gamePadMap{id: 1, keys: make(map[string]Key)}
	gamePadKeyMaps[2] = gamePadMap{id: 2, keys: make(map[string]Key)}
	gamePadKeyMaps[3] = gamePadMap{id: 3, keys: make(map[string]Key)}
}

func LoadDefaultKeys() {
	ResetKeys()
	firstKeyMap["pause"] = Key{Keyboard, ray.KeyEscape}
	firstKeyMap["accept"] = Key{Keyboard, ray.KeyEnter}
	firstKeyMap["back"] = Key{Keyboard, ray.KeyEscape}
	firstKeyMap["p1-Left"] = Key{Keyboard, ray.KeyA}
	firstKeyMap["p1-Right"] = Key{Keyboard, ray.KeyD}
	firstKeyMap["p1-Up"] = Key{Keyboard, ray.KeyW}
	firstKeyMap["p1-Down"] = Key{Keyboard, ray.KeyS}
	firstKeyMap["p1-placeBomb"] = Key{Keyboard, ray.KeySpace}

	secondKeyMap["back"] = Key{Keyboard, ray.KeyBackspace}
	secondKeyMap["p1-Left"] = Key{Keyboard, ray.KeyLeft}
	secondKeyMap["p1-Right"] = Key{Keyboard, ray.KeyRight}
	secondKeyMap["p1-Up"] = Key{Keyboard, ray.KeyUp}
	secondKeyMap["p1-Down"] = Key{Keyboard, ray.KeyDown}

	gamePadKeyMaps[0].keys["accept"] = Key{Gamepad, ray.GamepadButtonRightFaceDown}
	gamePadKeyMaps[0].keys["back"] = Key{Gamepad, ray.GamepadButtonRightFaceRight}
	gamePadKeyMaps[0].keys["pause"] = Key{Gamepad, ray.GamepadButtonMiddleRight}
	gamePadKeyMaps[0].keys["p1-Left"] = Key{Gamepad, ray.GamepadButtonLeftFaceLeft}
	gamePadKeyMaps[0].keys["p1-Right"] = Key{Gamepad, ray.GamepadButtonLeftFaceRight}
	gamePadKeyMaps[0].keys["p1-Up"] = Key{Gamepad, ray.GamepadButtonLeftFaceUp}
	gamePadKeyMaps[0].keys["p1-Down"] = Key{Gamepad, ray.GamepadButtonLeftFaceDown}
	gamePadKeyMaps[0].keys["p1-placeBomb"] = Key{Gamepad, ray.GamepadButtonRightFaceDown}

	//// p2
	firstKeyMap["p2-Left"] = Key{Keyboard, ray.KeyKp4}  //ray.KeyLeft
	firstKeyMap["p2-Right"] = Key{Keyboard, ray.KeyKp6} //ray.KeyRight
	firstKeyMap["p2-Up"] = Key{Keyboard, ray.KeyKp8}    //ray.KeyUp
	firstKeyMap["p2-Down"] = Key{Keyboard, ray.KeyKp5}  //ray.KeyDown
	firstKeyMap["p2-placeBomb"] = Key{Keyboard, ray.KeyL}

	gamePadKeyMaps[1].keys["p2-Left"] = Key{Gamepad, ray.GamepadButtonLeftFaceLeft}
	gamePadKeyMaps[1].keys["p2-Right"] = Key{Gamepad, ray.GamepadButtonLeftFaceRight}
	gamePadKeyMaps[1].keys["p2-Up"] = Key{Gamepad, ray.GamepadButtonLeftFaceUp}
	gamePadKeyMaps[1].keys["p2-Down"] = Key{Gamepad, ray.GamepadButtonLeftFaceDown}
	gamePadKeyMaps[1].keys["p2-placeBomb"] = Key{Gamepad, ray.GamepadButtonRightFaceDown}
}

func AddFirstKey(action string, key Key) {
	firstKeyMap[action] = key
}

func AddSecondryKey(action string, key Key) {
	secondKeyMap[action] = key
}

func GetFirstKey(action string) Key {
	return firstKeyMap[action]
}

func GetSecondryKey(action string) Key {
	return secondKeyMap[action]
}

func AddGamepadKey(gamePadId int, action string, key Key) {
	if key.Keyid == -1 {
		delete(gamePadKeyMaps[gamePadId].keys, action)
		return
	}
	gamePadKeyMaps[gamePadId].keys[action] = key
}

func GetGamepadKey(gamePadId int, action string) Key {
	k, found := gamePadKeyMaps[gamePadId].keys[action]
	if !found {
		return Key{Gamepad, -1}
	}
	return k
}

func IsKeyDown(action string) bool {
	return keyStates[action]
}

// returns true when all maped keys are up
func IsKeyUp(action string) bool {
	return !keyStates[action]
}

func IsKeyPressed(action string) bool {
	return keyStates[action] && !lastKeyStates[action]
}

func IsKeyReleased(action string) bool {
	return !keyStates[action] && lastKeyStates[action]
}

func SetActionUp(action string) {
	keyStates[action] = false
	lastKeyStates[action] = false
}

func SetActionDown(action string) {
	keyStates[action] = true
	lastKeyStates[action] = true
}

func UpdateKeys() {
	for k := range keyStates {
		lastKeyStates[k] = keyStates[k]
	}
	for k, v := range firstKeyMap {
		if v.KeyType == Mouse {
			keyStates[k] = ray.IsMouseButtonDown(v.Keyid)
		} else {
			keyStates[k] = ray.IsKeyDown(v.Keyid)
		}
	}
	for k, v := range secondKeyMap {
		if keyStates[k] {
			continue
		}
		if v.KeyType == Mouse {
			keyStates[k] = ray.IsMouseButtonDown(v.Keyid)
		} else {
			keyStates[k] = ray.IsKeyDown(v.Keyid)
		}
	}
	for gamePadId := int32(0); gamePadId < 4; gamePadId++ {
		if ray.IsGamepadAvailable(gamePadId) {
			for k, v := range gamePadKeyMaps[gamePadId].keys {
				if keyStates[k] {
					continue
				}
				if v.KeyType == Gamepad {
					keyStates[k] = ray.IsGamepadButtonDown(gamePadKeyMaps[gamePadId].id, v.Keyid)
				} else if v.KeyType == GamepadAxis {
					var movement float32
					switch v.Keyid {
					//Lx
					case GamepadAxisLeftXLeft:
						movement = ray.GetGamepadAxisMovement(gamePadId, ray.GamepadAxisLeftX)
						if movement > (-0.5) {
							movement = 0
						}
					case GamepadAxisLeftXRight:
						movement = ray.GetGamepadAxisMovement(gamePadId, ray.GamepadAxisLeftX)
						if movement < (0.5) {
							movement = 0
						}
					//Ly
					case GamepadAxisLeftYDown:
						movement = ray.GetGamepadAxisMovement(gamePadId, ray.GamepadAxisLeftY)
						if movement < (0.5) {
							movement = 0
						}
					case GamepadAxisLeftYUp:
						movement = ray.GetGamepadAxisMovement(gamePadId, ray.GamepadAxisLeftY)
						if movement > (-0.5) {
							movement = 0
						}
					//L3
					case GamepadAxisLeftTrigger:
						movement = ray.GetGamepadAxisMovement(gamePadId, ray.GamepadAxisLeftTrigger)
					//Rx
					case GamepadAxisRightXLeft:
						movement = ray.GetGamepadAxisMovement(gamePadId, ray.GamepadAxisRightX)
						if movement > (-0.5) {
							movement = 0
						}
					case GamepadAxisRightXRight:
						movement = ray.GetGamepadAxisMovement(gamePadId, ray.GamepadAxisRightX)
						if movement < (0.5) {
							movement = 0
						}
					//Ry
					case GamepadAxisRightYUp:
						movement = ray.GetGamepadAxisMovement(gamePadId, ray.GamepadAxisRightY)
						if movement < (0.5) {
							movement = 0
						}
					case GamepadAxisRightYDown:
						movement = ray.GetGamepadAxisMovement(gamePadId, ray.GamepadAxisRightY)
						if movement > (-0.5) {
							movement = 0
						}
					//R3
					case GamepadAxisRightTrigger:
						movement = ray.GetGamepadAxisMovement(gamePadId, ray.GamepadAxisRightTrigger)
					}
					if movement != 0 {
						keyStates[k] = true
					}
				}
			}
		}
	}
}
