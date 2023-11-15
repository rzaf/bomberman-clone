package core

import (
	ray "github.com/gen2brain/raylib-go/raylib"
)

type SprtieSheet struct {
	texture      *Texture
	count        int
	rowsCount    int
	columnsCount int
	width        int
	height       int
	seperationX  int
	seperationY  int
	src          ray.Rectangle
	isVertical   bool
	reverse      bool
}

func NewSpriteSheet(texture *Texture, count, rowsCount, columnsCount, spriteWidth, spriteHeight, seperationX, seperationY int, src ray.Rectangle, isVertical, reverse bool) *SprtieSheet {
	return &SprtieSheet{
		texture:      texture,
		count:        count,
		rowsCount:    rowsCount,
		columnsCount: columnsCount,
		width:        spriteWidth,
		height:       spriteHeight,
		seperationX:  seperationX,
		seperationY:  seperationY,
		src:          src,
		isVertical:   isVertical,
		reverse:      reverse,
	}
}
