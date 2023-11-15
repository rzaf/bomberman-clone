package core

import (
	ray "github.com/gen2brain/raylib-go/raylib"
)

type Sound struct {
	raySound ray.Sound
	volume   float32
}

func (s *Sound) SetVolume(v float32) {
	s.volume = v
	ray.SetSoundVolume(s.raySound, v)
}

func (s *Sound) GetVolume() float32 {
	return s.volume
}

func (s *Sound) Play() {
	ray.PlaySound(s.raySound)
}

func (s *Sound) Stop() {
	if s.IsPlaying() {
		ray.StopSound(s.raySound)
	}
}

func (s *Sound) Pause() {
	ray.PauseSound(s.raySound)
}

func (s *Sound) Resume() {
	ray.ResumeSound(s.raySound)
}

func (s *Sound) IsPlaying() bool {
	return ray.IsSoundPlaying(s.raySound)
}
