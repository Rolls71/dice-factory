package main

import (
	//"image"
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// CONSTANTS
const (
	screenWidth  = 640
	screenHeight = 480
)

const (
	tileSize = 32
)

const (
	stageSizeX = 5
	stageSizeY = 5
)

// VARIABLES
var basicStage [stageSizeY][stageSizeX]int

var basicGrassTile *ebiten.Image

// INITIALISATION
func init() {

}

// OBJECTS

// TILES
func (g Game) DrawTiles(screen *ebiten.Image) {
	fmt.Println(len(g.stage[0]))
	for y := 0; y < len(g.stage); y++ {
		for x := 0; x < len(g.stage[y]); x++ {
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
			screen.DrawImage(basicGrassTile, options)
		}
	}

}

// GAME
type Game struct {
	stage [stageSizeY][stageSizeX]int
}

func NewGame() *Game {
	return &Game{
		stage: basicStage,
	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawTiles(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (_screenWidth, _screenHeight int) {
	return screenWidth, screenHeight
}

func main() {
	basicStage = [stageSizeY][stageSizeX]int{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}

	basicGrassTile = ebiten.NewImage(tileSize, tileSize)
	basicGrassTile.Fill(color.RGBA{0, 0xFF, 0, 0xFF})

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
