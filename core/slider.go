package core

import (
	"math"

	ray "github.com/gen2brain/raylib-go/raylib"
)

type Slider struct {
	Rectangle              ray.Rectangle
	Point                  ray.Vector2
	minimum, maximum       float32
	color1, color2, color3 ray.Color
	isVertical             bool
	isSliderDown           bool
}

func NewSlider(x, y, width, height int32, minimum, maximum, value float32) *Slider {
	isVertical := false
	if height > width {
		isVertical = true
	}
	p := ray.NewVector2(float32(x), float32(y))
	if value < minimum || value > maximum {
		panic("value out of range !!! ")
	}
	v := (value - minimum) / (maximum - minimum)
	if isVertical {
		p.Y += v * float32(height)
	} else {
		p.X += v * float32(width)
	}
	var s *Slider = &Slider{
		Rectangle:    ray.NewRectangle(float32(x), float32(y), float32(width), float32(height)),
		Point:        p,
		minimum:      minimum,
		maximum:      maximum,
		isVertical:   isVertical,
		color1:       ray.Black,
		color2:       ray.NewColor(230, 230, 230, 250),
		color3:       ray.White,
		isSliderDown: false,
	}
	return s
}

func (s *Slider) GetValueInt() int {
	return int(math.Round(float64(s.GetValue())))
}

func (s *Slider) GetValue() float32 {
	r0 := s.maximum - s.minimum
	var percentage float32
	if s.isVertical {
		percentage = (s.Point.Y - s.Rectangle.Y) / s.Rectangle.Height
	} else {
		percentage = (s.Point.X - s.Rectangle.X) / s.Rectangle.Width
	}
	return s.minimum + r0*percentage
}

func (s *Slider) SetValue(value float32) {
	v := (value - s.minimum) / (s.maximum - s.minimum)
	if s.isVertical {
		s.Point.Y = s.Rectangle.Y + v*float32(s.Rectangle.Height)
	} else {
		s.Point.X = s.Rectangle.X + v*float32(s.Rectangle.Width)
	}
}

func (s *Slider) SetPos(x, y int32) {
	v := (s.GetValue() - s.minimum) / (s.maximum - s.minimum)
	s.Point.X = float32(x)
	s.Point.Y = float32(y)
	s.Rectangle.X = float32(x)
	s.Rectangle.Y = float32(y)
	if s.isVertical {
		s.Point.Y += v * float32(s.Rectangle.Height)
	} else {
		s.Point.X += v * float32(s.Rectangle.Width)
	}
}

func (s *Slider) SetColors(c1 ray.Color, c2 ray.Color, c3 ray.Color) {
	s.color1 = c1
	s.color2 = c2
	s.color3 = c3
}

func (s *Slider) Update() {
	if s.isSliderDown {
		if ray.IsMouseButtonDown(ray.MouseLeftButton) {
			if s.isVertical {
				s.Point.Y = ray.Clamp(float32(ray.GetMouseY()), s.Rectangle.Y, s.Rectangle.Y+s.Rectangle.Height)
			} else {
				s.Point.X = ray.Clamp(float32(ray.GetMouseX()), s.Rectangle.X, s.Rectangle.X+s.Rectangle.Width)
			}
		} else {
			s.isSliderDown = false
			ray.SetMouseCursor(ray.MouseCursorDefault)
		}
	} else {
		if ray.IsMouseButtonPressed(ray.MouseLeftButton) {
			if ray.CheckCollisionPointRec(ray.GetMousePosition(), ray.NewRectangle(s.Point.X-2, s.Point.Y-4, 4, s.Rectangle.Height+8)) {
				s.isSliderDown = true
				ray.SetMouseCursor(ray.MouseCursorCrosshair)
			} else if ray.CheckCollisionPointRec(ray.GetMousePosition(), s.Rectangle) {
				s.isSliderDown = true
				ray.SetMouseCursor(ray.MouseCursorCrosshair)
			}
		}
	}
}

func (s *Slider) Draw() {
	ray.DrawRectangle(s.Rectangle.ToInt32().X, s.Rectangle.ToInt32().Y, s.Rectangle.ToInt32().Width, s.Rectangle.ToInt32().Height, s.color1)
	ray.DrawRectangle(int32(s.Point.X)-2, int32(s.Point.Y-4), 4, int32(s.Rectangle.Height+8), s.color2)
	ray.DrawRectangleLines(int32(s.Point.X)-2, int32(s.Point.Y-4), 4, int32(s.Rectangle.Height+8), s.color3)
}
