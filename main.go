package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

// INPUT HANDLING

// generate ghost of object before moving

func (g *Game) ClickedObject() (bool, int) {
	x, y := ebiten.CursorPosition()
	x /= tileSize
	y /= tileSize
	for i, object := range g.objects {
		if object.x == x &&
			object.y == y {
			return true, i
		}
	}
	return false, 0
}

func (g *Game) UpdateCursor() {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) &&
		!g.isDragging {
		isObject, i := g.ClickedObject()
		if isObject {
			g.objects[i].trackMouse = true
			g.isDragging = true
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) &&
		g.isDragging {
		x, y := ebiten.CursorPosition()
		for i, object := range g.objects {
			if object.trackMouse {
				g.objects[i].trackMouse = false
				g.objects[i].MoveTo(x/tileSize, y/tileSize)
				g.isDragging = false
				break
			}
		}
	}
}

// GAME

type Game struct {
	tileStage  [stageSizeY][stageSizeX]int
	tileSet    []Tile
	objects    []Object
	isDragging bool
}

func NewGame() *Game {
	game := Game{
		tileStage:  [stageSizeY][stageSizeX]int{},
		tileSet:    []Tile{},
		objects:    []Object{},
		isDragging: false,
	}

	game.InitTiles()
	game.InitObjects()

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
