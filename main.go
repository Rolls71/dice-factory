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

// ToReal converts a tile coordinate to a real coordinate
func ToReal(i int) float64 {
	return float64(i * tileSize)
}

// ToTile converts a real coordinate to a tile coordinate
func ToTile(f float64) int {
	return int(f) / tileSize
}

// Game stores all data relevant to the running game
type Game struct {
	tileSet    []Tile
	tileStage  [stageSizeY][stageSizeX]int
	objectSet  []Object
	objects    []Object
	itemSet    []Item
	items      []Item
	isDragging bool
}

// NewGame constructs and returns a Game struct.
func NewGame() *Game {
	game := Game{
		tileSet:    []Tile{},
		tileStage:  [stageSizeY][stageSizeX]int{},
		objectSet:  []Object{},
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

	game.NewObject(ConveyorBelt, "conveyor_belt.png")
	game.NewObject(PlainObject, "plain_object.png")

	start := game.SpawnObject(ConveyorBelt, 7, 7, South)
	game.SpawnObject(ConveyorBelt, 7, 8, South)
	game.SpawnObject(ConveyorBelt, 7, 9, West)
	game.SpawnObject(ConveyorBelt, 6, 9, West)
	game.SpawnObject(ConveyorBelt, 5, 9, West)
	game.SpawnObject(ConveyorBelt, 5, 8, West)
	game.SpawnObject(ConveyorBelt, 5, 7, West)
	game.SpawnObject(ConveyorBelt, 6, 7, West)

	game.NewItem(PlainItem, "d6_6.png")
	game.SpawnItem(PlainItem, start)

	return &game
}

// Update calls the game's update functions
func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.SpawnItem(PlainItem, &g.objects[0])
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		x, y := g.GetCursorCoordinates()
		g.SpawnObject(ConveyorBelt, x, y, South)
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

// Draw calls the games draf functions and passes the screen
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
	ebiten.SetWindowTitle("Dice Factory")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
