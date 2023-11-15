package core

import (
	ray "github.com/gen2brain/raylib-go/raylib"
)

var (
	textures map[string]*Texture

	sounds map[string]*Sound
)

func LoadSound(path, tag string) *Sound {
	_, found := sounds[tag]
	if found {
		return sounds[tag]
		// panic("sound with tag '" + tag + "' already exist!")
	}
	if sounds == nil {
		sounds = make(map[string]*Sound)
	}
	sound := &Sound{}
	sound.raySound = ray.LoadSound(path)
	sound.SetVolume(0.5)
	sounds[tag] = sound
	return sound
}

func GetSound(tag string) *Sound {
	return sounds[tag]
}

func UnloadSound(tag string) {
	s, found := sounds[tag]
	if !found {
		return
	}
	ray.UnloadSound(s.raySound)
	delete(sounds, tag)
}

func LoadTexture(path, tag string, src ray.Rectangle) *Texture {
	_, found := textures[tag]
	if found {
		return textures[tag]
		// panic("texture with tag '" + tag + "' already exist!")
	}
	if textures == nil {
		textures = make(map[string]*Texture)
	}
	texture := NewTexture(path, src)
	textures[tag] = texture
	return texture
}

func GetTexture(tag string) *Texture {
	return textures[tag]
}

func UnloadTexture(tag string) {
	t, found := textures[tag]
	if !found {
		return
	}
	t.Unload()
	delete(textures, tag)
}
