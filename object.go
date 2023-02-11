package main

import (
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	buildCycleSeconds = 8 // Seconds per build cycle.
)

const conveyorSpeed float64 = float64(tileSize) / 1.75 // pixels per second

type ObjectType int

const (
	PlainObject  ObjectType = iota
	ConveyorBelt            // Moves items onto facing neighbor.
	Builder                 // Spawns a new item every build cycle and moves.
	Collector               // Deletes items
	Upgrader                // Upgrades items
)

type ObjectFacing int

const (
	South ObjectFacing = iota
	West
	North
	East
)

type Object struct {
	Object ObjectType
	X      int          // tile coord
	Y      int          // tile coord
	ID     uint64       // unique generated identifier
	Facing ObjectFacing // default South

	uiPosition int  // stores position of ui objects
	isDragged  bool // default false
}

func (o *Object) SetID(id uint64) {
	o.ID = id
}

func (o *Object) SetPosition(x, y int) {
	o.X = x
	o.Y = y
}

func (o *Object) SetFacing(dir ObjectFacing) {
	o.Facing = dir
}

func (o *Object) Rotate() {
	o.Facing = (o.Facing + 1) % 4
}

// IsItemOn tests if there is an item targeting the belt, and if it's currently
// on the belt.
// If so, it returns the item
func (g *Game) IsItemOn(object *Object) (bool, *Item) {
	// is there an item targeting the belt?
	isItem, item := g.GetItemTargeting(object)
	if !isItem {
		return false, item
	}

	// is the item currently on the belt?
	if item.X != ToReal(object.X) ||
		item.Y != ToReal(object.Y) {
		return false, item
	}

	return true, item
}

// IsItemMoveable tests if the belt is pointing at an object and if theres an
// item targeting the neighbor.
// If so, it returns the neighbor
func (g *Game) IsItemMoveable(object *Object) (bool, *Object) {
	// is the belt pointing at an object?
	isNeighbor, neighbor := g.GetNeighborOf(object)
	if !isNeighbor {
		return false, neighbor
	}

	// is there an item targeting the neighbor?
	isItem, item := g.GetItemTargeting(neighbor)
	if isItem && item.Item != Truck {
		return false, neighbor
	}

	if object.Object == ConveyorBelt ||
		neighbor.Object == ConveyorBelt {
		return true, neighbor
	}

	return false, neighbor
}

// MoveItemOn tests IsItemOn and IsItemMoveable before setting an item's
// target position to a neighbor.
func (g *Game) MoveItemOn(object *Object) {
	isItemOn, item := g.IsItemOn(object)
	if !isItemOn {
		return
	}

	isItemMoveable, neighbor := g.IsItemMoveable(object)
	if !isItemMoveable {
		return
	}

	// set the item to target that object
	item.SetTargetPosition(neighbor.X, neighbor.Y)
}

// UpdateObjects will iterate through each Object and switch,
// depending on their type. Each Object type may have different functionality
func (g *Game) UpdateObjects() {
	for _, copy := range g.Objects {
		object := g.Objects[copy.ID]
		switch object.Object {
		case ConveyorBelt:
			g.MoveItemOn(object)
		case Builder:
			if g.ticks%uint64(frameRate*buildCycleSeconds) == 0 {
				isItemMoveable, _ := g.IsItemMoveable(object)
				if isItemMoveable {
					g.SpawnItem(Plain, object)
				}
			}
			g.MoveItemOn(object)
		case Collector:
			isItemOn, item := g.IsItemOn(object)
			if isItemOn && item.Item != Truck {
				g.Storages[Warehouse].StoreDie(item.Item, item.Face)
				g.SellDie(item.Item, item.Value())
				delete(g.Items, item.ID)
			}
		case Upgrader:
			isItemOn, item := g.IsItemOn(object)
			if isItemOn && g.ticks%uint64(frameRate*buildCycleSeconds) == 0 {
				g.SetItem(item, Gold)
				g.MoveItemOn(object)
			}
		}
	}
}

// GetNeighborOf looks at the tile the object is facing to check for an Object
// If there is an object, it returns true, and a reference to the Object
// If there is no object, it returns false, and an empty Object
func (g *Game) GetNeighborOf(o *Object) (bool, *Object) {
	switch o.Facing {
	case South:
		isObject, object := g.GetObjectAt(o.X, o.Y+1)
		if isObject {
			return true, object
		}
	case West:
		isObject, object := g.GetObjectAt(o.X-1, o.Y)
		if isObject {
			return true, object
		}
	case North:
		isObject, object := g.GetObjectAt(o.X, o.Y-1)
		if isObject {
			return true, object
		}
	case East:
		isObject, object := g.GetObjectAt(o.X+1, o.Y)
		if isObject {
			return true, object
		}
	}
	return false, &Object{}
}

// GetObjectAt returns true if there is an Object at the given coordinates
// An array of every Object at that coordinate is also returned.
func (g *Game) GetObjectAt(x, y int) (bool, *Object) {
	for _, copy := range g.Objects {
		if copy.X == x &&
			copy.Y == y {
			return true, g.Objects[copy.ID]
		}
	}
	return false, &Object{}
}

// NewObject creates a new type of object.
// New Object is appended to the Game's Object Set
func (g *Game) NewObject(objectType ObjectType, imageName string) {
	path := "images/" + imageName
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	g.objectImages[objectType] = img
}

// SpawnObject constructs a new object of ObjectType
func (g *Game) SpawnObject(
	objectType ObjectType,
	x, y int,
	facing ObjectFacing,
) *Object {
	object := Object{
		Object: objectType,
	}
	object.SetID(g.NextID())
	object.SetPosition(x, y)
	object.SetFacing(facing)

	g.ObjectCount[objectType] += 1
	g.UnlockObject(objectType)

	g.Objects[object.ID] = &object
	return &object
}

// DrawObjects will draw every Tile in the game's list of objects.
// Objects are drawn on their stored grid coordinate.
// An Object flagged with trackMouse will be drawn attached to cursor instead.
func (g Game) DrawObjects(screen *ebiten.Image) {
	var onTop *ebiten.Image
	var topOptions *ebiten.DrawImageOptions
	for _, object := range g.Objects {
		img := g.objectImages[object.Object]
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Scale(float64(tileSize)/float64(img.Bounds().Dx()),
			float64(tileSize)/float64(img.Bounds().Dy()))
		options.GeoM.Rotate(math.Pi / 2 * float64(object.Facing))
		switch object.Facing {
		case West:
			options.GeoM.Translate(float64(tileSize), 0)
		case North:
			options.GeoM.Translate(float64(tileSize), float64(tileSize))
		case East:
			options.GeoM.Translate(0, float64(tileSize))
		}
		if object.isDragged {
			x, y := ebiten.CursorPosition()
			options.GeoM.Translate(float64(x), float64(y))
			onTop = img
			topOptions = options
		} else {
			options.GeoM.Translate(
				float64(object.X*tileSize),
				float64(object.Y*tileSize))
			screen.DrawImage(img, options)
		}
	}
	if onTop != nil {
		screen.DrawImage(onTop, topOptions)
	}
}
