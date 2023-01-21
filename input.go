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

// UpdateCursor runs updateOnMouseDown and updateOnMouseUp
func (g *Game) UpdateCursor() {
	g.updateOnMouseDown()
	g.updateOnMouseUp()
}

// updateOnMouseDown tests if an Object has been selected.
// The Game's isDragging flag and the Object's trackMouse flag is set to true.
func (g *Game) updateOnMouseDown() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) &&
		!g.isDragging {
		x, y := g.GetCursorCoordinates()
		isObject, objectIndices := g.GetObjectsAt(x, y)
		if isObject {
			g.objects[objectIndices[0]].trackMouse = true
			g.isDragging = true
			g.LiftObject(objectIndices[0])
		}
	}
}

// updateOnMouseUp tests if a dragged object has been released.
// The Game's isDragging flag and the Object's trackMouse flag is set to false.
func (g *Game) updateOnMouseUp() {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
		g.isDragging {
		x, y := g.GetCursorCoordinates()
		isObject, _ := g.GetObjectsAt(x, y)
		for i, object := range g.objects {
			if object.trackMouse {
				g.objects[i].trackMouse = false
				g.isDragging = false
				if !isObject {
					g.objects[i].MoveTo(x, y)
				}
				break
			}
		}
	}
}
