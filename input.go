package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// GetCursorCoordinates returns the tile coordinate that the cursor is within.
func (g *Game) GetCursorCoordinates() (int, int) {
	x, y := ebiten.CursorPosition()
	x /= tileSize
	y /= tileSize
	return x, y
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
		x, y := g.GetCursorCoordinates()
		isObject, object := g.GetObjectAt(x, y)
		if isObject {
			object.isDragged = true
			g.isDragging = true
		}
	}
}

// onDragEnd tests if a dragged object has been released.
// The Game's isDragging flag and the Object's trackMouse flag is set to false.
func (g *Game) onDragEnd(mouseButton ebiten.MouseButton) {
	if inpututil.IsMouseButtonJustReleased(mouseButton) &&
		g.isDragging {
		x, y := g.GetCursorCoordinates()
		isObject, _ := g.GetObjectAt(x, y)
		for i, copy := range g.Objects {
			if copy.isDragged {
				g.Objects[i].isDragged = false
				g.isDragging = false
				if !isObject {
					g.Objects[i].SetPosition(x, y)
				}
				break
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
			x, y := g.GetCursorCoordinates()
			isObject, object := g.GetObjectAt(x, y)
			if isObject {
				object.Rotate()
			}
		}
	}
}
