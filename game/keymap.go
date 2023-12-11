package game

import (
	ray "github.com/gen2brain/raylib-go/raylib"
)

type gamePadMap struct {
	keys map[string]int32
	id   int32
}

var (
	firstKeys     map[string]int32 = make(map[string]int32)
	secondoryKeys map[string]int32 = make(map[string]int32)
	gamePadMaps   []gamePadMap     = make([]gamePadMap, 4)
)

func ResetKeys() {
	firstKeys = make(map[string]int32)
	secondoryKeys = make(map[string]int32)
	gamePadMaps[0] = gamePadMap{id: 0, keys: make(map[string]int32)}
	gamePadMaps[1] = gamePadMap{id: 1, keys: make(map[string]int32)}
	gamePadMaps[2] = gamePadMap{id: 2, keys: make(map[string]int32)}
	gamePadMaps[3] = gamePadMap{id: 3, keys: make(map[string]int32)}
}

func LoadDefaultKeys() {
	ResetKeys()

	firstKeys["pause"] = ray.KeyEscape
	firstKeys["accept"] = ray.KeyEnter
	firstKeys["back"] = ray.KeyEscape
	firstKeys["p1-Left"] = ray.KeyA
	firstKeys["p1-Right"] = ray.KeyD
	firstKeys["p1-Up"] = ray.KeyW
	firstKeys["p1-Down"] = ray.KeyS
	firstKeys["p1-placeBomb"] = ray.KeySpace

	secondoryKeys["back"] = ray.KeyBackspace
	secondoryKeys["p1-Left"] = ray.KeyLeft
	secondoryKeys["p1-Right"] = ray.KeyRight
	secondoryKeys["p1-Up"] = ray.KeyUp
	secondoryKeys["p1-Down"] = ray.KeyDown

	gamePadMaps[0].keys["p1-Left"] = ray.GamepadXboxButtonLeft
	gamePadMaps[0].keys["p1-Right"] = ray.GamepadXboxButtonRight
	gamePadMaps[0].keys["p1-Up"] = ray.GamepadXboxButtonUp
	gamePadMaps[0].keys["p1-Down"] = ray.GamepadXboxButtonDown
	gamePadMaps[0].keys["p1-placeBomb"] = ray.GamepadXboxButtonA

	//// p2
	firstKeys["p2-Left"] = ray.KeyKp4  //ray.KeyLeft
	firstKeys["p2-Right"] = ray.KeyKp6 //ray.KeyRight
	firstKeys["p2-Up"] = ray.KeyKp8    //ray.KeyUp
	firstKeys["p2-Down"] = ray.KeyKp5  //ray.KeyDown
	firstKeys["p2-placeBomb"] = ray.KeyL

	gamePadMaps[1].keys["p2-Left"] = ray.GamepadXboxButtonLeft
	gamePadMaps[1].keys["p2-Right"] = ray.GamepadXboxButtonRight
	gamePadMaps[1].keys["p2-Up"] = ray.GamepadXboxButtonUp
	gamePadMaps[1].keys["p2-Down"] = ray.GamepadXboxButtonDown
	gamePadMaps[1].keys["p2-placeBomb"] = ray.GamepadXboxButtonA
}

func AddFirstKey(action string, key int) {
	firstKeys[action] = int32(key)
}

func AddSecondryKey(action string, key int) {
	secondoryKeys[action] = int32(key)
}

func GetFirstKey(action string) int32 {
	return firstKeys[action]
}

func GetSecondryKey(action string) int32 {
	return secondoryKeys[action]
}

func AddGamepadKey(gamePad int, action string, key int) {
	if key == -1 {
		delete(gamePadMaps[gamePad].keys, action)
		return
	}
	gamePadMaps[gamePad].keys[action] = int32(key)
}

func GetGamepadKey(gamepad int, action string) int32 {
	k, found := gamePadMaps[gamepad].keys[action]
	if !found {
		return -1
	}
	return k
}

func IsKeyDown(action string) bool {
	return ray.IsKeyDown(firstKeys[action]) ||
		ray.IsKeyDown(secondoryKeys[action]) ||
		ray.IsGamepadButtonDown(gamePadMaps[0].id, gamePadMaps[0].keys[action]) ||
		ray.IsGamepadButtonDown(gamePadMaps[1].id, gamePadMaps[1].keys[action])
}

// returns true when all maped keys are up
func IsKeyUp(action string) bool {
	return ray.IsKeyUp(firstKeys[action]) &&
		ray.IsKeyUp(secondoryKeys[action]) &&
		ray.IsGamepadButtonUp(gamePadMaps[0].id, gamePadMaps[0].keys[action]) &&
		ray.IsGamepadButtonUp(gamePadMaps[1].id, gamePadMaps[1].keys[action])
}

func IsKeyPressed(action string) bool {
	return ray.IsKeyPressed(firstKeys[action]) ||
		ray.IsKeyPressed(secondoryKeys[action]) ||
		ray.IsGamepadButtonPressed(gamePadMaps[0].id, gamePadMaps[0].keys[action]) ||
		ray.IsGamepadButtonPressed(gamePadMaps[1].id, gamePadMaps[1].keys[action])
}

func IsKeyReleased(action string) bool {
	return ray.IsKeyReleased(firstKeys[action]) &&
		ray.IsKeyReleased(secondoryKeys[action]) &&
		ray.IsGamepadButtonReleased(gamePadMaps[0].id, gamePadMaps[0].keys[action]) &&
		ray.IsGamepadButtonReleased(gamePadMaps[1].id, gamePadMaps[1].keys[action])
}
