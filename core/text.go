package core

import (
	ray "github.com/gen2brain/raylib-go/raylib"
)

type Text struct {
	Text       string
	Font       ray.Font
	Pos        ray.Vector2
	TopLeftPos ray.Vector2
	FontSize   float32
	Spacing    float32
	Color      ray.Color
	Size       ray.Vector2
}

func NewText(text string, Font ray.Font, Pos ray.Vector2, FontSize float32, Spacing float32, Color ray.Color) *Text {
	var txt *Text = &Text{
		Text:       text,
		Font:       Font,
		Pos:        Pos,
		TopLeftPos: Pos,
		FontSize:   FontSize,
		Spacing:    Spacing,
		Color:      Color,
		Size:       ray.NewVector2(0, 0),
	}

	txt.Measure()
	return txt
}

func (t *Text) SetText(s string) {
	t.Text = s
	t.Measure()
}

func (t *Text) Measure() {
	if t.Font.Chars == nil {
		panic("nil font in initialized Text")
	}
	t.Size = ray.MeasureTextEx(t.Font, t.Text, t.FontSize, t.Spacing)
	t.TopLeftPos.X = t.Pos.X - t.Size.X/2
	t.TopLeftPos.Y = t.Pos.Y - t.Size.Y/2
}

func (t *Text) DrawCentered() {
	if t.Size.X == 0.0 {
		t.Measure()
	}
	ray.DrawTextEx(
		t.Font,
		t.Text,
		t.TopLeftPos,
		// ray.Vector2{X: t.Pos.X - t.Size.X/2, Y: t.Pos.Y - t.Size.Y/2},
		t.FontSize,
		t.Spacing,
		t.Color,
	)
}

func (t *Text) Draw() {
	ray.DrawTextEx(t.Font, t.Text, t.Pos, t.FontSize, t.Spacing, t.Color)
}

func (t *Text) DrawAt(p ray.Vector2) {
	ray.DrawTextEx(t.Font, t.Text, p, t.FontSize, t.Spacing, t.Color)
}
