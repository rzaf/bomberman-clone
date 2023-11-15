package core

import (
	ray "github.com/gen2brain/raylib-go/raylib"
)

type Button struct {
	Texture   *Texture
	Text      *Text
	Boundary  ray.Rectangle
	OnClick   func()
	IsEnabled bool
}

func NewImageButton(texture *Texture, boundary ray.Rectangle) *Button {
	return &Button{texture, nil, boundary, nil, true}
}

func NewTextButton(text *Text, boundary ray.Rectangle) *Button {
	return &Button{nil, text, boundary, nil, true}
}

func (b *Button) SetTexture(texture *Texture) {
	b.Texture = texture
}

func (b *Button) Update() {
	if !b.IsEnabled {
		return
	}
	if b.OnClick != nil && ray.IsMouseButtonPressed(ray.MouseLeftButton) {
		if ray.CheckCollisionPointRec(ray.GetMousePosition(), b.Boundary) {
			b.OnClick()
		}
	}

}

func (b *Button) Draw() {
	if b.Text != nil {
		b.Text.Draw()
	}
	if b.Texture != nil {
		b.Texture.Draw()
	}
}

func (b *Button) DrawAt(dst ray.Rectangle) {
	if b.Text != nil {
		b.Text.DrawAt(ray.NewVector2(dst.X, dst.Y))
	}
	if b.Texture != nil {
		b.Texture.DrawAt(dst)
	}
	// b.Texture.DrawAt(dst)
}
