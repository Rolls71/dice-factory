package main

import (
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ObjectType int

const (
	PlainObject ObjectType = iota
	ConveyorBelt
)

type ObjectFacing int

const (
	South ObjectFacing = iota
	West
	North
	East
)

type Object struct {
	objectType ObjectType
	image      *ebiten.Image
	x          int // tile coord
	y          int // tile coord

	id        uint64       // unique generated identifier
	facing    ObjectFacing // default South
	isDragged bool         // default false
}

func (o *Object) SetID(id uint64) {
	o.id = id
}

func (o *Object) SetPosition(x, y int) {
	o.x = x
	o.y = y
}

func (o *Object) SetFacing(dir ObjectFacing) {
	o.facing = dir
}

func (o *Object) Rotate() {
	o.facing = (o.facing + 1) % 4
}

// MoveItemOn will set the target of an item on the object to a neighbor.
// First it checks if theres an item targeting the belt, if there is it checks
// if it's on the object. If its on the object, it checks if the obhect is
// pointing at another object. If theres somewhere to send the item, it checks
// if the next object is targeted by another item. If its not it will retarget
// the item.
func (g *Game) MoveItemOn(object *Object) {
	// is there an item targeting the belt?
	isItem, item := g.GetItemTargeting(object)
	if !isItem {
		return
	}

	// is the item currently on the belt?
	if item.x != ToReal(object.x) ||
		item.y != ToReal(object.y) {
		return
	}

	// is the belt pointing at an object?
	isNeighbor, neighbor := g.GetNeighborOf(object)
	if !isNeighbor {
		return
	}

	// is there an item targeting the neighbor?
	isItem, _ = g.GetItemTargeting(neighbor)
	if isItem {
		return
	}

	// set the item to target that object
	item.SetTargetPosition(neighbor.x, neighbor.y)
}

// UpdateObjects will iterate through each Object and switch,
// depending on their type. Each Object type may have different functionality
func (g *Game) UpdateObjects() {
	for _, copy := range g.objects {
		object := g.objects[copy.id]
		switch object.objectType {
		case ConveyorBelt:
			g.MoveItemOn(object)
		}
	}
}

// GetNeighborOf looks at the tile the object is facing to check for an Object
// If there is an object, it returns true, and a reference to the Object
// If there is no object, it returns false, and an empty Object
func (g *Game) GetNeighborOf(o *Object) (bool, *Object) {
	switch o.facing {
	case South:
		isObject, object := g.GetObjectAt(o.x, o.y+1)
		if isObject {
			return true, object
		}
	case West:
		isObject, object := g.GetObjectAt(o.x-1, o.y)
		if isObject {
			return true, object
		}
	case North:
		isObject, object := g.GetObjectAt(o.x, o.y-1)
		if isObject {
			return true, object
		}
	case East:
		isObject, object := g.GetObjectAt(o.x+1, o.y)
		if isObject {
			return true, object
		}
	}
	return false, &Object{}
}

// GetObjectAt returns true if there is an Object at the given coordinates
// An array of every Object at that coordinate is also returned.
func (g *Game) GetObjectAt(x, y int) (bool, *Object) {
	for _, copy := range g.objects {
		if copy.x == x &&
			copy.y == y {
			return true, g.objects[copy.id]
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
	g.objectSet[objectType] = &Object{
		objectType: objectType,
		image:      img,
	}
}

// SpawnObject constructs a new object of ObjectType
func (g *Game) SpawnObject(
	objectType ObjectType,
	x, y int,
	facing ObjectFacing,
) *Object {
	object := *g.objectSet[objectType]
	object.SetID(g.NextID())
	object.SetPosition(x, y)
	object.SetFacing(facing)

	g.objects[object.id] = &object
	return &object
}

// DrawObjects will draw every Tile in the game's list of objects.
// Objects are drawn on their stored grid coordinate.
// An Object flagged with trackMouse will be drawn attached to cursor instead.
func (g *Game) DrawObjects(screen *ebiten.Image) {
	var onTop *ebiten.Image
	var topOptions *ebiten.DrawImageOptions
	for _, copy := range g.objects {
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Rotate(math.Pi / 2 * float64(copy.facing))
		switch copy.facing {
		case West:
			options.GeoM.Translate(float64(tileSize), 0)
		case North:
			options.GeoM.Translate(float64(tileSize), float64(tileSize))
		case East:
			options.GeoM.Translate(0, float64(tileSize))
		}
		if copy.isDragged {
			x, y := ebiten.CursorPosition()
			options.GeoM.Translate(float64(x), float64(y))
			onTop = copy.image
			topOptions = options
		} else {
			options.GeoM.Translate(
				float64(copy.x*tileSize),
				float64(copy.y*tileSize))
			screen.DrawImage(copy.image, options)
		}
	}
	if onTop != nil {
		screen.DrawImage(onTop, topOptions)
	}
}
