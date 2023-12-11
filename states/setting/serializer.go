package setting

import (
	"encoding/json"
	"fmt"
	"os"

	ray "github.com/gen2brain/raylib-go/raylib"
	"github.com/rzaf/bomberman-clone/game"
	"github.com/rzaf/bomberman-clone/states/running"
)

func Load() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
		fmt.Println("loading default config")
		game.LoadDefaultKeys()
		return
	}
	var s any
	err = json.Unmarshal(data, &s)
	if err != nil {
		fmt.Println(err)
		return
	}
	config, ok := s.(map[string]any)
	if !ok {
		fmt.Println("Failed to load config.json")
		return
	}
	audioConfig, _ := config["audio"].(map[string]any)

	game.MasterVolume = float32(audioConfig["masterV"].(float64) / 100)
	game.MusicVolume = float32(audioConfig["musicV"].(float64) / 100)
	game.EffectVolume = float32(audioConfig["sfxV"].(float64) / 100)
	masterV = game.MasterVolume
	musicV = game.MusicVolume
	effectV = game.EffectVolume
	running.SetSoundsVolume()
	ray.SetMasterVolume(masterV)
	fmt.Println(game.MasterVolume, game.MusicVolume, game.EffectVolume)

	game.ResetKeys()
	// keyboard1
	keyboard1Cfg := config["keyboard_1"].(map[string]any)
	if keyboard1Cfg != nil {
		for i := 0; i < len(actions); i++ {
			k, ok := keyboard1Cfg[actions[i]].(float64)
			if ok {
				game.AddFirstKey(actions[i], int(k))
			}
		}
	}

	// keyboard2
	keyboard2Cfg := config["keyboard_2"].(map[string]any)
	if keyboard2Cfg != nil {
		for i := 0; i < len(actions); i++ {
			k, ok := keyboard2Cfg[actions[i]].(float64)
			if ok {
				game.AddSecondryKey(actions[i], int(k))
			}
		}
	}
	// gamepad_1
	gamepad1Cfg := config["gamepad_1"].(map[string]any)
	if gamepad1Cfg != nil {
		for i := 0; i < len(actions); i++ {
			k, ok := gamepad1Cfg[actions[i]].(float64)
			if ok {
				game.AddGamepadKey(0, actions[i], int(k))
			}
		}
	}

	// gamepad_2
	gamepad2Cfg := config["gamepad_2"].(map[string]any)
	if gamepad2Cfg != nil {
		for i := 0; i < len(actions); i++ {
			k, ok := gamepad2Cfg[actions[i]].(float64)
			if ok {
				game.AddGamepadKey(1, actions[i], int(k))
			}
		}
	}

	// gamepad_3
	gamepad3Cfg := config["gamepad_3"].(map[string]any)
	if gamepad3Cfg != nil {
		for i := 0; i < len(actions); i++ {
			k, ok := gamepad3Cfg[actions[i]].(float64)
			if ok {
				game.AddGamepadKey(2, actions[i], int(k))
			}
		}
	}
	// gamepad_4
	gamepad4Cfg := config["gamepad_4"].(map[string]any)
	if gamepad4Cfg != nil {
		for i := 0; i < len(actions); i++ {
			k, ok := gamepad4Cfg[actions[i]].(float64)
			if ok {
				game.AddGamepadKey(3, actions[i], int(k))
			}
		}
	}
}

func Save() {
	config := make(map[string]any)
	config["audio"] = map[string]int{
		"masterV": int(game.MasterVolume * 100),
		"musicV":  int(game.MusicVolume * 100),
		"sfxV":    int(game.EffectVolume * 100),
	}
	config["keyboard_1"] = make(map[string]int32)
	config["keyboard_2"] = make(map[string]int32)
	config["gamepad_1"] = make(map[string]int32)
	config["gamepad_2"] = make(map[string]int32)
	config["gamepad_3"] = make(map[string]int32)
	config["gamepad_4"] = make(map[string]int32)
	keyboard1Cfg := config["keyboard_1"].(map[string]int32)
	keyboard2Cfg := config["keyboard_2"].(map[string]int32)
	gamepad1Cfg := config["gamepad_1"].(map[string]int32)
	gamepad2Cfg := config["gamepad_2"].(map[string]int32)
	gamepad3Cfg := config["gamepad_3"].(map[string]int32)
	gamepad4Cfg := config["gamepad_4"].(map[string]int32)

	for i := 0; i < len(actions); i++ {
		keyboard1Cfg[actions[i]] = game.GetFirstKey(actions[i])
		keyboard2Cfg[actions[i]] = game.GetSecondryKey(actions[i])
		gamepad1Cfg[actions[i]] = game.GetGamepadKey(0, actions[i])
		gamepad2Cfg[actions[i]] = game.GetGamepadKey(1, actions[i])
		gamepad3Cfg[actions[i]] = game.GetGamepadKey(2, actions[i])
		gamepad4Cfg[actions[i]] = game.GetGamepadKey(3, actions[i])
	}
	data, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	os.WriteFile("config.json", data, os.ModePerm)
}
