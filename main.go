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
	tileSet    []Tile                      // Stores different types of Tiles.
	tileStage  [stageSizeY][stageSizeX]int // Stores Tile instances to be drawn.
	objectSet  map[ObjectType]*Object      // Stores different types of Objects.
	objects    map[uint64]*Object          // Stores Object instances to be drawn.
	objectID   uint64                      // Stores id of last object spawned.
	itemSet    map[ItemType]*Item          // Stores different types of Items.
	items      map[uint64]*Item            // Stores Item instances to be drawn.
	id         uint64                      // Stores id of last item/object made.
	isDragging bool                        // Is an Object being dragged
}

func (g *Game) NextID() uint64 {
	g.id += 1
	return g.id
}

// NewGame constructs and returns a Game struct.
func NewGame() *Game {
	game := Game{
		tileSet:    []Tile{},
		tileStage:  [stageSizeY][stageSizeX]int{},
		objectSet:  map[ObjectType]*Object{},
		objects:    map[uint64]*Object{},
		itemSet:    map[ItemType]*Item{},
		items:      map[uint64]*Item{},
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
	game.SpawnObject(ConveyorBelt, 5, 9, North)
	game.SpawnObject(ConveyorBelt, 5, 8, North)
	game.SpawnObject(ConveyorBelt, 5, 7, East)
	game.SpawnObject(ConveyorBelt, 6, 7, East)

	game.NewItem(PlainItem, "d6_6.png")
	game.SpawnItem(PlainItem, start)

	return &game
}

// Update calls the game's update functions
func (g *Game) Update() error {
	x, y := g.GetCursorCoordinates()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		isObject, object := g.GetObjectAt(x, y)
		if isObject {
			g.SpawnItem(PlainItem, object)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		isObject, object := g.GetObjectAt(x, y)
		if isObject {
			delete(g.objects, object.id)
		} else {
			g.SpawnObject(ConveyorBelt, x, y, South)
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		fmt.Printf("objects: ")
		fmt.Println(g.objects)
		fmt.Print("items: ")
		fmt.Println(*g.items[1])
	}
	g.UpdateInput()
	g.UpdateObjects()
	g.UpdateItems()
	return nil
}

// Draw calls the games drag functions and passes the screen
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
