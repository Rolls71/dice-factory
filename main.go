package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  int = 640
	screenHeight int = 480
)

const (
	tileSize int = 32
)

const (
	frameDelta float64 = 1.0 / 60.0
)

const (
	stageSizeX int = screenWidth / tileSize
	stageSizeY int = screenHeight / tileSize
)

type Game struct {
	tileStage  [stageSizeY][stageSizeX]int
	tileSet    []Tile
	objects    []Object
	itemSet    []Item
	items      []Item
	isDragging bool
}

// NewGame constructs and returns a Game struct.
// Tiles to be used in the game are initialised here with Game.NewTile()
// The Tile stage is initialised here with Game.SetTileArray()
// Objects to be added at the start are initialised here with Game.NewObject()
func NewGame() *Game {
	game := Game{
		tileStage:  [stageSizeY][stageSizeX]int{},
		tileSet:    []Tile{},
		objects:    []Object{},
		itemSet:    []Item{},
		items:      []Item{},
		isDragging: false,
	}

	game.NewTile("basic_grass", "basic_grass.png") // ID = 0
	game.NewTile("long_grass", "long_grass.png")   // ID = 1

	game.SetTileStage([stageSizeY][stageSizeX]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	})

	game.SpawnObject(ConveyorBelt, "conveyor_belt.png", 7, 7).SetFacing(0)
	game.SpawnObject(ConveyorBelt, "conveyor_belt.png", 7, 8).SetFacing(0)
	game.SpawnObject(ConveyorBelt, "conveyor_belt.png", 7, 9).SetFacing(1)
	game.SpawnObject(ConveyorBelt, "conveyor_belt.png", 6, 9).SetFacing(1)
	game.SpawnObject(ConveyorBelt, "conveyor_belt.png", 5, 9).SetFacing(2)
	game.SpawnObject(ConveyorBelt, "conveyor_belt.png", 5, 8).SetFacing(2)
	game.SpawnObject(ConveyorBelt, "conveyor_belt.png", 5, 7).SetFacing(3)
	game.SpawnObject(ConveyorBelt, "conveyor_belt.png", 6, 7).SetFacing(3)

	game.NewItem(PlainItem, "d6_6.png")

	return &game
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.SpawnItem(PlainItem, &g.objects[0])
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		x, y := g.GetCursorCoordinates()
		g.SpawnObject(0, "conveyor_belt.png", x, y)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		fmt.Printf("objects: ")
		fmt.Println(g.objects)
		fmt.Print("items: ")
		fmt.Println(g.items)
	}
	g.UpdateInput()
	g.UpdateObjects()
	g.UpdateItems()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawTiles(screen)
	g.DrawObjects(screen)
	g.DrawItems(screen)
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
