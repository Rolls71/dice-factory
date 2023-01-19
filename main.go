package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// CONSTANTS
const (
	screenWidth  int = 640
	screenHeight int = 480
)

const (
	tileSize int = 32
)

const (
	stageSizeX int = screenWidth / tileSize
	stageSizeY int = screenHeight / tileSize
)

// VARIABLES

// OBJECTS

// TILES
type Tile struct {
	name  string
	image *ebiten.Image
}

const tileCount = 1

var tiles []Tile

func initTiles() {
	//replace with images later
	tempGrass := ebiten.NewImage(tileSize, tileSize)
	tempGrass.Fill(color.RGBA{0, 0xFF, 0, 0xFF})

	basicGrass := Tile{
		name:  "basicGrass",
		image: tempGrass,
	}

	tiles = []Tile{basicGrass}
}

func GetTile(index int) *ebiten.Image {
	return tiles[index].image
}

func (g Game) DrawTiles(screen *ebiten.Image) {
	for y := 0; y < stageSizeY; y++ {
		for x := 0; x < stageSizeX; x++ {
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
			screen.DrawImage(GetTile(g.tileStage[y][x]), options)
		}
	}

}

// GAME
type Game struct {
	tileStage   [stageSizeY][stageSizeX]int
	objectStage [stageSizeY][stageSizeX]int
}

func NewGame() *Game {
	return &Game{
		// return a 2D array of zeroes
		tileStage:   [stageSizeY][stageSizeX]int{},
		objectStage: [stageSizeY][stageSizeX]int{},
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
	initTiles()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
