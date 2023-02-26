package dicefactory

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

var opaqueGrey color.RGBA = color.RGBA{0x55, 0x55, 0x55, 0x99}

// UnlockObject attempts to add an object to a hotbar if it has a specific count
func (g *Game) UnlockObject(objectType ObjectType) {
	switch objectType {
	case Builder:
		if g.ObjectCount[objectType] == 4 {
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

	printString := ""

	if g.Currencies[PlainBuck] > 0 {
		printString += fmt.Sprintf("PlainBucks: %d\n", g.Currencies[PlainBuck])
	}
	if g.Currencies[GoldBuck] > 0 {
		printString += fmt.Sprintf("GoldBucks: %d\n", g.Currencies[GoldBuck])
	}
	if g.Currencies[PlainBuck] > 0 || g.Currencies[GoldBuck] > 0 {
		printString += "\n"
	}
	_, val := g.Cost(ConveyorBelt)
	printString += fmt.Sprintf("Conveyor Belt: %d PlainBucks\n", val)
	_, val = g.Cost(Builder)
	printString += fmt.Sprintf("Builder: %d PlainBucks\n", val)

	if len(g.UIObjects) > 2 {
		_, val = g.Cost(Upgrader)
		printString += fmt.Sprintf("Upgrader: %d PlainBucks\n", val)
	}

	if g.Warehouse.Count > 0 {
		printString += "\n"
		printString += fmt.Sprintf("Warehouse Dice: %d/%d\n", g.Warehouse.Count, g.Warehouse.Capacity)
		printString += fmt.Sprintf("Dice Sell Rate: 1 Dice/%d secs\n", sellRate)
	}

	val = 0
	for id := range g.Trucks {
		if id > val {
			val = id
		}
	}
	printString += "\n"
	printString += fmt.Sprintf("Dice Stored in Truck: %d/%d\n", g.Trucks[val].Storage.Count, g.Trucks[val].Storage.Capacity)
	for itemType := range g.Trucks[val].Storage.Dice {
		printString += fmt.Sprintf("Type Stored in Truck: %s\n", itemType.String())
	}
	if g.Trucks[val].Storage.Count >= g.Trucks[val].Storage.Capacity {
		printString += "Click truck to deliver dice to warehouse"
	}

	ebitenutil.DebugPrint(screen, printString)
}

// DrawHotbar draws objects in the UIObjects array
func (g Game) DrawHotbar(screen *ebiten.Image) {
	hotbar := ebiten.NewImage(ScreenWidth, lowerHUDHeight)
	hotbar.Fill(opaqueGrey)
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(0, float64(ScreenHeight-lowerHUDHeight))
	screen.DrawImage(hotbar, options)

	var onTop *ebiten.Image
	var topOptions *ebiten.DrawImageOptions
	for index, object := range g.UIObjects {
		img := g.objectImages[object.Object]
		options = &ebiten.DrawImageOptions{}
		options.GeoM.Scale(float64(tileSize)/float64(img.Bounds().Dx()),
			float64(tileSize)/float64(img.Bounds().Dy()))
		if object.isDragged {
			x, y := ebiten.CursorPosition()
			options.GeoM.Translate(float64(x), float64(y))
			onTop = img
			topOptions = options
		} else {
			object.uiPosition = index*(tileSize+hotbarSpacing) + (ScreenWidth-len(g.UIObjects)*(tileSize+hotbarSpacing))/2
			options.GeoM.Translate(
				float64(object.uiPosition),
				float64(ScreenHeight-tileSize-hotbarSpacing))
			screen.DrawImage(img, options)
		}
	}
	if onTop != nil {
		screen.DrawImage(onTop, topOptions)
	}
}
