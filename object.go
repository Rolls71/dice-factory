package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Object struct {
	name       string
	image      *ebiten.Image
	x          int
	y          int
	trackMouse bool //default false for constructor
}

func (o *Object) MoveTo(x, y int) {
	o.x = x
	o.y = y
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

// NewObject constructs a new Object with given parameters.
// New Object is appended to the Game's Object list.
func (g *Game) NewObject(
	objectName,
	imageName string,
	x, y int,
) {
	path := "images/" + imageName
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	g.objects = append(g.objects, Object{
		name:       objectName,
		image:      img,
		x:          x,
		y:          y,
		trackMouse: false,
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
