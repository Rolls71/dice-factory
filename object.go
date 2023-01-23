package main

import (
	_ "image/png"
	"log"

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
	x          int
	y          int

	facing     ObjectFacing // default South
	trackMouse bool         // default false
}

func (o *Object) SetPosition(x, y int) {
	o.x = x
	o.y = y
}

func (o *Object) SetFacing(dir ObjectFacing) {
	o.facing = dir
}

// LiftObject pushes an object to the end of the render list
// Dangerous function! Moves object to end of g.objects
func (g *Game) LiftObject(index int) {
	length := len(g.objects)
	if index > length ||
		index < 0 {
		log.Fatal("Error: Index out of range")
	}
	object := g.objects[index]
	slice := g.objects[:index]
	if index != length-1 {
		slice = append(slice, g.objects[index+1:]...)
	}
	g.objects = append(slice, object)
}

// GetObjectsAt returns true if there is an Object at the given coordinates
// An array of every Object at that coordinate is also returned.
func (g *Game) GetObjectsAt(x, y int) (bool, []int) {
	var objectIndices []int
	for index, object := range g.objects {
		if object.x == x &&
			object.y == y {
			objectIndices = append(objectIndices, index)
		}
	}
	return (len(objectIndices) > 0), objectIndices
}

// NewObject constructs a new object of ObjectType
// New Object is appended to the Game's Object Set
func (g *Game) NewObject(
	objectType ObjectType,
	imageName string,
	x, y int,
) {
	path := "images/" + imageName
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	g.objects = append(g.objects, Object{
		objectType: objectType,
		image:      img,
		x:          x,
		y:          y,
	})
}

// DrawObjects will draw every Tile in the game's list of objects.
// Objects are drawn on their stored grid coordinate.
// An Object flagged with trackMouse will be drawn attached to cursor instead.
func (g *Game) DrawObjects(screen *ebiten.Image) {
	for _, object := range g.objects {
		options := &ebiten.DrawImageOptions{}
		if object.trackMouse {
			x, y := ebiten.CursorPosition()
			options.GeoM.Translate(float64(x), float64(y))
		} else {
			options.GeoM.Translate(
				float64(object.x*tileSize),
				float64(object.y*tileSize))
		}
		screen.DrawImage(object.image, options)
	}
}
