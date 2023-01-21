package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) GetCursorCoordinates() (int, int) {
	x, y := ebiten.CursorPosition()
	x /= tileSize
	y /= tileSize
	return x, y
}

func (g *Game) UpdateCursor() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) &&
		!g.isDragging {
		x, y := g.GetCursorCoordinates()
		isObject, objectIndices := g.GetObjectsAt(x, y)
		if isObject {
			g.objects[objectIndices[0]].trackMouse = true
			g.isDragging = true
		}
	}

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
