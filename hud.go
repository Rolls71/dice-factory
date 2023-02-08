package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	upperHUDHeight = 32
	lowerHUDHeight = tileSize + hotbarSpacing*2
	hotbarSpacing  = 5
)

var darkGrey color.RGBA = color.RGBA{0x55, 0x55, 0x55, 0xff}

// UnlockObject attempts to add an object to a hotbar if it has a specific count
func (g *Game) UnlockObject(objectType ObjectType) {
	switch objectType {
	case Builder:
		if g.ObjectCount[objectType] >= 3 {
			g.SpawnUIObject(Upgrader)
		}
	}
}

// InitHUD adds UIObjects to hotbar.
// Run InitHUD after objectImages are initialised
func (g *Game) InitHUD() {
	g.SpawnUIObject(ConveyorBelt)
	g.SpawnUIObject(Builder)
}

// SpawnUIObject constructs a new object of ObjectType in the UI overlay
func (g *Game) SpawnUIObject(
	object ObjectType,
) {
	g.UIObjects = append(g.UIObjects, &Object{Object: object})
}

// DrawHUD calls HUD-related draw functions
func (g *Game) DrawHUD(screen *ebiten.Image) {
	g.DrawHotbar(screen)

	string := fmt.Sprintf("Dice Points: %d\n", g.Balance.DicePoints)
	string += fmt.Sprintf("Conveyor Belt: %d\n", g.Cost(ConveyorBelt))
	string += fmt.Sprintf("Builder: %d\n", g.Cost(Builder))
	string += fmt.Sprintf("Upgrader: %d\n", g.Cost(Upgrader))
	ebitenutil.DebugPrint(screen, string)
}

// DrawHotbar draws objects in the UIObjects array
func (g Game) DrawHotbar(screen *ebiten.Image) {
	hotbar := ebiten.NewImage(screenWidth, lowerHUDHeight)
	hotbar.Fill(darkGrey)
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(0, float64(screenHeight-lowerHUDHeight))
	screen.DrawImage(hotbar, options)

	var onTop *ebiten.Image
	var topOptions *ebiten.DrawImageOptions
	for index, object := range g.UIObjects {
		options = &ebiten.DrawImageOptions{}
		if object.isDragged {
			x, y := ebiten.CursorPosition()
			options.GeoM.Translate(float64(x), float64(y))
			onTop = g.objectImages[object.Object]
			topOptions = options
		} else {
			object.uiPosition = index*(tileSize+hotbarSpacing) + (screenWidth-len(g.UIObjects)*(tileSize+hotbarSpacing))/2
			options.GeoM.Translate(
				float64(object.uiPosition),
				float64(screenHeight-tileSize-hotbarSpacing))
			screen.DrawImage(g.objectImages[object.Object], options)
		}
	}
	if onTop != nil {
		screen.DrawImage(onTop, topOptions)
	}
}
