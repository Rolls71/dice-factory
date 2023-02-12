package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// GetCursorCoordinates returns the tile coordinate that the cursor is within.
func GetCursorCoordinates() (int, int) {
	x, y := ebiten.CursorPosition()
	x /= tileSize
	y /= tileSize
	return x, y
}

// IsInGameArea returns true if the coordinate is within the games boundaries
func IsInGameArea(x, y int) bool {
	return (x > 0 &&
		x < screenWidth &&
		y > 0 &&
		y < screenHeight-lowerHUDHeight)
}

// UpdateInput runs all major input functions.
// Keys can be rebound here
func (g *Game) UpdateInput() {
	g.onDragStart(ebiten.MouseButtonLeft)
	g.onDragEnd(ebiten.MouseButtonLeft)
	g.onRotate(ebiten.KeyR)
}

// onDragStart tests if an Object has been selected.
// The Game's isDragging flag and the Object's trackMouse flag is set to true.
func (g *Game) onDragStart(mouseButton ebiten.MouseButton) {
	if inpututil.IsMouseButtonJustPressed(mouseButton) &&
		!g.isDragging {
		x, y := ebiten.CursorPosition()
		if !IsInGameArea(x, y) {
			for _, object := range g.UIObjects {
				if x > object.uiPosition &&
					x < object.uiPosition+tileSize {
					object.isDragged = true
					g.isDragging = true
					return
				}
			}
		} else {
			xTile, yTile := GetCursorCoordinates()
			isObject, object := g.GetObjectAt(xTile, yTile)
			if isObject && object.Object != Collector {
				object.isDragged = true
				g.isDragging = true
			}
		}
	}
}

// onDragEnd tests if a dragged object has been released.
// The Game's isDragging flag and the Object's trackMouse flag is set to false.
func (g *Game) onDragEnd(mouseButton ebiten.MouseButton) {
	if inpututil.IsMouseButtonJustReleased(mouseButton) &&
		g.isDragging {
		x, y := ebiten.CursorPosition()
		isObject, _ := g.GetObjectAt(x, y)
		for _, object := range g.UIObjects {
			if object.isDragged {
				object.isDragged = false
				g.isDragging = false
				if !isObject && IsInGameArea(x, y) {
					g.Buy(object.Object, x/tileSize, y/tileSize, South)
				}
				return
			}
		}
		for _, object := range g.Objects {
			if object.isDragged {
				object.isDragged = false
				g.isDragging = false
				if !isObject && IsInGameArea(x, y) {
					object.X = x / tileSize
					object.Y = y / tileSize
				}
				return
			}
		}
	}
}

// onRotate will rotate an object under the cursor if the right key has been
// pressed. The key is passed as a parameter
func (g *Game) onRotate(key ebiten.Key) {
	if inpututil.IsKeyJustPressed(key) {
		if g.isDragging {
			for index, object := range g.Objects {
				if object.isDragged {
					g.Objects[index].Rotate()
					break
				}
			}
		} else {
			x, y := GetCursorCoordinates()
			isObject, object := g.GetObjectAt(x, y)
			if isObject {
				object.Rotate()
			}
		}
	}
}
