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
