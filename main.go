package main

import (
	"encoding/json"
	"fmt"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  int = 1366
	screenHeight int = 768
)

const tileSize int = 64

const (
	frameRate  int     = 60
	frameDelta float64 = 1.0 / float64(frameRate)
)

const (
	stageSizeX int = screenWidth/tileSize + 1
	stageSizeY int = screenHeight / tileSize
)

const saveFilename string = "save.json"

const maxUint64 = ^uint64(0)

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
	tileImages   map[TileType]*ebiten.Image   // Stores different types of Tiles.
	objectImages map[ObjectType]*ebiten.Image // Stores different object images.
	itemImages   map[ItemType]*ebiten.Image   // Stores different item images.
	truckImages  map[TruckType]*ebiten.Image

	TileStage   [stageSizeY][stageSizeX]int // Stores Tile instances to be drawn.
	Objects     map[uint64]*Object          // Stores Object instances to be drawn.
	UIObjects   []*Object                   // Stores Objects in the UI Overlay
	ObjectCount map[ObjectType]uint64       // Tracks the number of Objects
	Items       map[uint64]*Item            // Stores Item instances to be drawn.
	Currencies  map[CurrencyType]uint64     // Stores different currencies
	Storages    map[uint64]*Storage         // Stores a list of trucks and warehouses
	Trucks      map[uint64]*Truck
	Warehouse   *Storage // Stores the main storage stuct
	ID          uint64   // Stores id of last item/object made.

	ticks      uint64 // Stores tick count
	isDragging bool   // Is an Object being dragged

}

// NextID increments the stored id and returns it
func (g *Game) NextID() uint64 {
	g.ID++
	return g.ID
}

// InitImages will initialise all images
func (g *Game) InitImages() {
	g.NewTile(BasicGrass, "basic_grass.png")
	g.NewTile(LongGrass, "long_grass.png")

	g.NewObject(PlainObject, "plain_object.png")
	g.NewObject(ConveyorBelt, "conveyor_belt.png")
	g.NewObject(Builder, "builder.png")
	g.NewObject(Collector, "plain_object.png")
	g.NewObject(Upgrader, "builder.png")

	g.NewItem(PlainD6, "d6.png")
	g.NewItem(GoldD6, "gold_d6.png")

	g.NewTruck(BasicTruck, "truck.png")
}

// NewGame constructs and returns a Game struct.
func NewGame() *Game {

	game := Game{
		tileImages:   map[TileType]*ebiten.Image{},
		objectImages: map[ObjectType]*ebiten.Image{},
		itemImages:   map[ItemType]*ebiten.Image{},
		truckImages:  map[TruckType]*ebiten.Image{},

		TileStage:   [stageSizeY][stageSizeX]int{},
		Objects:     map[uint64]*Object{},
		UIObjects:   []*Object{},
		ObjectCount: map[ObjectType]uint64{},
		Items:       map[uint64]*Item{},
		Currencies:  map[CurrencyType]uint64{},
		Storages:    map[uint64]*Storage{},
		Trucks:      map[uint64]*Truck{},
	}

	game.Warehouse = game.NewStorage(Warehouse, warehouseCapacity, 0)

	// set up tile stage
	game.TileStage = [stageSizeY][stageSizeX]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	game.InitImages()
	game.InitHUD()

	builder := game.SpawnObject(Builder, 6, 4, South)
	game.SpawnObject(ConveyorBelt, 6, 5, West)

	collector1 := game.SpawnObject(Collector, 5, 5, South)
	collector2 := game.SpawnObject(Collector, 5, 6, South)

	game.SpawnTruck(BasicTruck, []*Object{collector1, collector2},
		-5, 5, 2, 5, 4, 2)

	game.SpawnItem(PlainD6, builder)

	return &game
}

// SaveGame stores the game struct in a JSON file
func (g *Game) SaveGame() {
	bytes, err := json.Marshal(g)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(saveFilename, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// LoadGame returns the game struct stored in given JSON file.
// Images are reinitialised
func LoadGame(filePath string) *Game {
	f, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var game Game
	json.Unmarshal(f, &game)

	game.tileImages = map[TileType]*ebiten.Image{}
	game.objectImages = map[ObjectType]*ebiten.Image{}
	game.itemImages = map[ItemType]*ebiten.Image{}
	game.truckImages = map[TruckType]*ebiten.Image{}
	game.InitImages()

	return &game
}

// Update calls the game's update functions and iterates the games tick count
func (g *Game) Update() error {
	g.ticks += 1

	// Temporary inputs before system is put in place
	x, y := GetCursorCoordinates()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		isObject, object := g.GetObjectAt(x, y)
		if isObject {
			g.SpawnItem(PlainD6, object)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		isObject, object := g.GetObjectAt(x, y)
		if isObject {
			g.ObjectCount[object.Object] -= 1
			delete(g.Objects, object.ID)
		} else {
			g.Buy(ConveyorBelt, x, y, South)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		isObject, object := g.GetObjectAt(x, y)
		if isObject {
			g.ObjectCount[object.Object] -= 1
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
	g.UpdateTrucks()
	g.UpdateCurrency()
	return nil
}

// Draw calls the games drag functions and passes the screen
func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawTiles(screen)
	g.DrawObjects(screen)
	g.DrawItems(screen)
	g.DrawHUD(screen)
	g.DrawTrucks(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (
	_screenWidth, _screenHeight int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Dice Factory")
	ebiten.SetFullscreen(true)

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
	game.ticks = 60 * 7
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
	game.SaveGame()
}
