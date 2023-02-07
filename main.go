package main

import (
	"encoding/json"
	"fmt"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  int = 640
	screenHeight int = 640 // 480
)

const (
	tileSize int = 32
)

const (
	frameRate  int     = 60
	frameDelta float64 = 1.0 / float64(frameRate)
)

const (
	stageSizeX int = screenWidth / tileSize
	stageSizeY int = screenHeight / tileSize
)

const (
	saveFilename string = "save.json"
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
// TODO: Reorganise tileSet to store only images
type Game struct {
	tileSet      []Tile                       // Stores different types of Tiles.
	objectImages map[ObjectType]*ebiten.Image // Stores different object images.
	itemImages   map[ItemType]*ebiten.Image   // Stores different item images.

	TileStage   [stageSizeY][stageSizeX]int // Stores Tile instances to be drawn.
	Objects     map[uint64]*Object          // Stores Object instances to be drawn.
	ObjectCount map[ObjectType]uint64       // Tracks the number of Objects
	Items       map[uint64]*Item            // Stores Item instances to be drawn.
	Balance     Currency                    // Stores currency data
	ID          uint64                      // Stores id of last item/object made.

	ticks      uint64 // Stores tick count
	isDragging bool   // Is an Object being dragged

}

// NextID increments the stored id and returns it
func (g *Game) NextID() uint64 {
	g.ID += 1
	return g.ID
}

// InitImages will initialise all images
func (g *Game) InitImages() {
	// initialise tiles
	g.NewTile("basic_grass", "basic_grass.png") // ID = 0
	g.NewTile("long_grass", "long_grass.png")   // ID = 1

	// initialise objects
	g.NewObject(PlainObject, "plain_object.png")
	g.NewObject(ConveyorBelt, "conveyor_belt.png")
	g.NewObject(Builder, "plain_object.png")
	g.NewObject(Collector, "plain_object.png")
	g.NewObject(Upgrader, "plain_object.png")

	//initialise items
	g.NewItem(PlainD6, "d6.png")
	g.NewItem(GoldD6, "gold_d6.png")
}

// NewGame constructs and returns a Game struct.
func NewGame() *Game {

	game := Game{
		tileSet:      []Tile{},
		objectImages: map[ObjectType]*ebiten.Image{},
		itemImages:   map[ItemType]*ebiten.Image{},

		TileStage:   [stageSizeY][stageSizeX]int{},
		Objects:     map[uint64]*Object{},
		ObjectCount: map[ObjectType]uint64{},
		Items:       map[uint64]*Item{},
		Balance:     Currency{},
	}

	// set up tile stage
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

	game.InitImages()

	// builders
	builder1 := game.SpawnObject(Builder, 2, 2, South)
	builder2 := game.SpawnObject(Builder, 4, 2, South)
	builder3 := game.SpawnObject(Builder, 6, 2, South)

	// builder extensions
	game.SpawnObject(ConveyorBelt, 2, 3, South)
	game.SpawnObject(ConveyorBelt, 4, 3, South)
	game.SpawnObject(ConveyorBelt, 6, 3, South)

	// collecting east belt
	game.SpawnObject(Upgrader, 2, 4, East)
	game.SpawnObject(ConveyorBelt, 3, 4, East)
	game.SpawnObject(ConveyorBelt, 4, 4, East)
	game.SpawnObject(ConveyorBelt, 5, 4, East)
	game.SpawnObject(ConveyorBelt, 6, 4, East)

	// connecting south belt
	game.SpawnObject(ConveyorBelt, 7, 4, South)
	game.SpawnObject(ConveyorBelt, 7, 5, South)
	game.SpawnObject(ConveyorBelt, 7, 6, South)

	// cube
	game.SpawnObject(ConveyorBelt, 7, 7, South)
	game.SpawnObject(ConveyorBelt, 7, 8, South)
	game.SpawnObject(ConveyorBelt, 7, 9, West)
	game.SpawnObject(ConveyorBelt, 6, 9, West)
	game.SpawnObject(ConveyorBelt, 5, 9, North)
	game.SpawnObject(ConveyorBelt, 5, 8, North)
	game.SpawnObject(ConveyorBelt, 5, 7, East)
	game.SpawnObject(Collector, 6, 7, East)

	game.SpawnItem(PlainD6, builder1)
	game.SpawnItem(PlainD6, builder2)
	game.SpawnItem(PlainD6, builder3)

	return &game
}

func (g *Game) SaveGame() {
	bytes, err := json.Marshal(g)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(saveFilename, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadGame(filePath string) *Game {
	f, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var game Game
	json.Unmarshal(f, &game)

	game.tileSet = []Tile{}
	game.objectImages = map[ObjectType]*ebiten.Image{}
	game.itemImages = map[ItemType]*ebiten.Image{}
	game.InitImages()

	return &game
}

// Update calls the game's update functions and iterates the games tick count
func (g *Game) Update() error {
	g.ticks += 1

	// Temporary inputs before system is put in place
	x, y := g.GetCursorCoordinates()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		isObject, object := g.GetObjectAt(x, y)
		if isObject {
			g.SpawnItem(PlainD6, object)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		isObject, object := g.GetObjectAt(x, y)
		if isObject {
			g.ObjectCount[object.ObjectType] -= 1
			delete(g.Objects, object.ID)
		} else {
			g.Buy(ConveyorBelt, x, y, South)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		isObject, object := g.GetObjectAt(x, y)
		if isObject {
			g.ObjectCount[object.ObjectType] -= 1
			delete(g.Objects, object.ID)
		} else {
			g.SpawnObject(Builder, x, y, South)
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		fmt.Printf("objects: ")
		fmt.Println(g.Objects)
		fmt.Print("items: ")
		fmt.Println(*g.Items[1])
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
	g.DrawHUD(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (
	_screenWidth, _screenHeight int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Dice Factory")

	var game *Game
	if _, err := os.Stat(saveFilename); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			game = NewGame()
		} else {
			// other error
			log.Fatal(err)
		}
	} else {
		game = LoadGame(saveFilename)
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
	game.SaveGame()
}
