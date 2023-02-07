package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	upperHUDHeight = 32
	lowerHUDHeight = 64
)

// InitHUD adds UIObjects to hotbar.
// Run InitHUD after objectImages are initialised
func (g *Game) InitHUD() {
	for i := 0; i < len(g.objectImages); i++ {
		if i == int(Collector) {
			continue
		}
		g.SpawnUIObject(ObjectType(i))
	}
}

// SpawnUIObject constructs a new object of ObjectType in the UI overlay
func (g *Game) SpawnUIObject(
	object ObjectType,
) {
	g.UIObjects = append(g.UIObjects, &Object{Object: object})
}

func (g *Game) DrawHUD(screen *ebiten.Image) {
	g.DrawUIObjects(screen)

	string := fmt.Sprintf("Dice Points: %d\n", g.Balance.DicePoints)
	string += fmt.Sprintf("Conveyor Belt: %d\n", g.Cost(ConveyorBelt))
	string += fmt.Sprintf("Builder: %d\n", g.Cost(Builder))
	string += fmt.Sprintf("Upgrader: %d\n", g.Cost(Upgrader))
	ebitenutil.DebugPrint(screen, string)
}

func (g Game) DrawUIObjects(screen *ebiten.Image) {
	var onTop *ebiten.Image
	var topOptions *ebiten.DrawImageOptions
	for index, object := range g.UIObjects {
		options := &ebiten.DrawImageOptions{}
		if object.isDragged {
			x, y := ebiten.CursorPosition()
			options.GeoM.Translate(float64(x), float64(y))
			onTop = g.objectImages[object.Object]
			topOptions = options
		} else {
			object.uiPosition = index*tileSize + (screenWidth-len(g.UIObjects)*tileSize)/2
			options.GeoM.Translate(
				float64(object.uiPosition),
				float64(screenHeight-tileSize))
			screen.DrawImage(g.objectImages[object.Object], options)
		}
	}
	if onTop != nil {
		screen.DrawImage(onTop, topOptions)
	}
}
