package game

import (
	ray "github.com/gen2brain/raylib-go/raylib"
)

var (
	keys map[string]int32 = make(map[string]int32)
)

func ResetKeys() {
	keys = make(map[string]int32)
}

func LoadDefaultKeys() {
	keys["pause"] = ray.KeyEscape

	keys["p1-moveLeft"] = ray.KeyA
	keys["p1-moveRight"] = ray.KeyD
	keys["p1-moveUp"] = ray.KeyW
	keys["p1-moveDown"] = ray.KeyS
	keys["p1-placeBomb"] = ray.KeySpace

	keys["p2-moveLeft"] = ray.KeyKp4  //ray.KeyLeft
	keys["p2-moveRight"] = ray.KeyKp6 //ray.KeyRight
	keys["p2-moveUp"] = ray.KeyKp8    //ray.KeyUp
	keys["p2-moveDown"] = ray.KeyKp5  //ray.KeyDown
	keys["p2-placeBomb"] = ray.KeyL
}

func AddKey(action string, key int) {
	keys[action] = int32(key)
}

func GetKey(action string) int32 {
	return keys[action]
}

func IsKeyDown(action string) bool {
	return ray.IsKeyDown(keys[action])
}

func IsKeyUp(action string) bool {
	return ray.IsKeyUp(keys[action])
}

func IsKeyPressed(action string) bool {
	return ray.IsKeyPressed(keys[action])
}

func IsKeyReleased(action string) bool {
	return ray.IsKeyReleased(keys[action])
}
