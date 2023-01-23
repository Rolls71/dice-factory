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
	x          int // tile coord
	y          int // tile coord

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

func (g *Game) UpdateObjects() {
	for index, object := range g.objects {
		switch object.objectType {
		case ConveyorBelt:
			isItem, itemIndex := g.GetItemOn(index)
			if isItem {
				isNeighbor, neighborIndex := g.GetNeighborOf(&object)
				if isNeighbor {
					g.items[itemIndex].SetTarget(
						g.objects[neighborIndex].x,
						g.objects[neighborIndex].y,
					)
				}

			}
		}
	}
}

func (g *Game) GetNeighborOf(o *Object) (bool, int) {
	var isObject bool
	var index int
	switch o.facing {
	case South:
		isObject, index = g.GetObjectsAt(o.x, o.y+1)
		if isObject {
			return true, index
		}
	case West:
		isObject, index = g.GetObjectsAt(o.x-1, o.y)
		if isObject {
			return true, index
		}
	case North:
		isObject, index = g.GetObjectsAt(o.x, o.y-1)
		if isObject {
			return true, index
		}
	case East:
		isObject, index = g.GetObjectsAt(o.x+1, o.y)
		if isObject {
			return true, index
		}
	}
	isObject, index = g.GetObjectsAt(o.x, o.y)
	if isObject {
		return true, index
	} else {
		return false, -1
	}
}

// GetObjectsAt returns true if there is an Object at the given coordinates
// An array of every Object at that coordinate is also returned.
func (g *Game) GetObjectsAt(x, y int) (bool, int) {
	for index, object := range g.objects {
		if object.x == x &&
			object.y == y {
			return true, index
		}
	}
	return false, 0
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
	var onTop *ebiten.Image
	var topOptions *ebiten.DrawImageOptions
	for _, object := range g.objects {
		options := &ebiten.DrawImageOptions{}
		if object.trackMouse {
			x, y := ebiten.CursorPosition()
			options.GeoM.Translate(float64(x), float64(y))
			onTop = object.image
			topOptions = options
		} else {
			options.GeoM.Translate(
				float64(object.x*tileSize),
				float64(object.y*tileSize))
			screen.DrawImage(object.image, options)
		}
	}
	if onTop != nil {
		screen.DrawImage(onTop, topOptions)
	}
}
