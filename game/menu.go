package game

import (
	ray "github.com/gen2brain/raylib-go/raylib"
	"github.com/rzaf/bomberman-clone/core"
)

type MenuItem struct {
	text     *core.Text
	onSubmit func()
}

func NewMenuItem(str string, onSubmit func()) *MenuItem {
	return &MenuItem{
		text:     core.NewText(str, ray.GetFontDefault(), ray.NewVector2(0, 0), 15, 4, ray.White),
		onSubmit: onSubmit,
	}
}

type Menu struct {
	items []*MenuItem
	Font  ray.Font

	SelectKey          string
	NextKey, BeforeKey string

	Spacing          float32
	FontSize         float32
	SelectedFontSize float32
	Color            ray.Color
	SelectedColor    ray.Color
	currentIndex     int
	IsCentered       bool
	IsMouseEnabled   bool
	IsClickable      bool

	Pos     ray.Vector2
	Padding int
	Direction
}

func NewMenu(mItems ...*MenuItem) *Menu {
	m := Menu{}
	m.items = mItems
	m.IsMouseEnabled = true
	m.IsClickable = true
	m.currentIndex = 0
	m.IsCentered = true
	m.Color = ray.White
	m.SelectedColor = ray.Red
	m.Padding = 25
	m.Spacing = 4
	m.Font = ray.GetFontDefault()
	m.FontSize = 18
	m.SelectedFontSize = 22
	m.Direction = DOWN
	m.SelectKey = "accept"
	m.NextKey = "p1-Down"
	m.BeforeKey = "p1-Up"
	return &m
}

func (m *Menu) Refresh() {
	for i, item := range m.items {
		item.text.Font = m.Font
		item.text.Spacing = m.Spacing
		switch m.Direction {
		case RIGHT, LEFT:
			item.text.Pos.X = m.Pos.X + float32(i*m.Padding)
			item.text.Pos.Y = m.Pos.Y
		case UP, DOWN:
			item.text.Pos.X = m.Pos.X
			item.text.Pos.Y = m.Pos.Y + float32(i*m.Padding)
		}
		if i == m.currentIndex {
			item.text.Color = m.SelectedColor
			item.text.FontSize = m.SelectedFontSize
		} else {
			item.text.FontSize = m.FontSize
			item.text.Color = m.Color
		}
		if m.IsCentered {
			item.text.Measure()
		}
	}
}

func (m *Menu) getIndexOfMousePos() int {
	mousePos := ray.GetMousePosition()
	rec := ray.NewRectangle(0, 0, 0, 0)
	for i, item := range m.items {
		if m.IsCentered {
			rec.X = item.text.TopLeftPos.X
			rec.Y = item.text.TopLeftPos.Y
		} else {
			rec.X = item.text.Pos.X
			rec.Y = item.text.Pos.Y
		}
		rec.Width = item.text.Size.X
		rec.Height = item.text.Size.Y
		if ray.CheckCollisionPointRec(mousePos, rec) {
			return i
		}
	}
	return -1
}
func (m *Menu) Update() {
	if m.IsMouseEnabled {
		if v := ray.GetMouseDelta(); v.X != 0 || v.Y != 0 { // mouse moved
			i := m.getIndexOfMousePos()
			if i != -1 {
				m.currentIndex = i
				m.Refresh()
			}
		}
	}
	if m.IsClickable {
		if ray.IsMouseButtonReleased(ray.MouseButtonLeft) {
			i := m.getIndexOfMousePos()
			if i != -1 {
				m.currentIndex = i
				m.items[m.currentIndex].onSubmit()
				m.Refresh()
			}
		}
	}
	if IsKeyPressed(m.NextKey) {
		m.currentIndex = (m.currentIndex + 1) % len(m.items)
		m.Refresh()
	} else if IsKeyPressed(m.BeforeKey) {
		if m.currentIndex == 0 {
			m.currentIndex = len(m.items) - 1
		} else {
			m.currentIndex -= 1
		}
		m.Refresh()
	}
	if IsKeyPressed(m.SelectKey) {
		m.items[m.currentIndex].onSubmit()
	}
}

func (m *Menu) Draw() {
	for _, item := range m.items {
		if m.IsCentered {
			item.text.DrawCentered()
		} else {
			item.text.Draw()
		}
	}

	// bounding box
	// for _, item := range m.items {
	// 	rec := ray.NewRectangle(0, 0, 0, 0)
	// 	if m.IsCentered {
	// 		rec.X = item.text.TopLeftPos.X
	// 		rec.Y = item.text.TopLeftPos.Y
	// 	} else {
	// 		rec.X = item.text.Pos.X
	// 		rec.Y = item.text.Pos.Y
	// 	}
	// 	rec.Width = item.text.Size.X
	// 	rec.Height = item.text.Size.Y
	// 	ray.DrawRectangleLinesEx(rec, 2, ray.Blue)
	// }
}

func (m *Menu) Index() int {
	return m.currentIndex
}

func (m *Menu) SetIndex(i int) {
	m.currentIndex = i
	m.Refresh()
}
