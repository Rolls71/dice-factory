package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

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

type Game struct {
	tileStage  [stageSizeY][stageSizeX]int
	tileSet    []Tile
	objects    []Object
	isDragging bool
}

// NewGame constructs and returns a Game struct.
// Object and Tile initialisation functions are called after construction.
func NewGame() *Game {
	game := Game{
		tileStage:  [stageSizeY][stageSizeX]int{},
		tileSet:    []Tile{},
		objects:    []Object{},
		isDragging: false,
	}

	game.InitTiles()
	game.NewObject("die_1", "d6_6.png", 3, 3)
	game.NewObject("die_2", "d6_6.png", 5, 6)

	return &game
}

func (g *Game) Update() error {
	g.UpdateCursor()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawTiles(screen)
	g.DrawObjects(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (
	_screenWidth, _screenHeight int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
