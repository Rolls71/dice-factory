package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
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

// GAME
type Game struct {
	tileStage [stageSizeY][stageSizeX]int
	tileSet   []Tile
	objects   []Object
}

func NewGame() *Game {
	game := Game{
		// return a 2D array of zeroes
		tileStage: [stageSizeY][stageSizeX]int{},
		tileSet:   []Tile{},
		objects:   []Object{},
	}

	game.InitTiles()
	game.InitObjects()

	return &game
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawTiles(screen)
	g.DrawObjects(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (_screenWidth, _screenHeight int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
