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

// InitObjects initialises any Objects that should start on-screen.
// Will be converted to take array of objects after Object.New() is created.
func (g *Game) InitObjects() {
	var img *ebiten.Image
	var err error
	img, _, err = ebitenutil.NewImageFromFile("images/d6_6.png")
	if err != nil {
		log.Fatal(err)
	}

	g.objects = append(g.objects, Object{
		name:       "d6",
		image:      img,
		x:          (screenWidth / 2) / tileSize,
		y:          (screenHeight / 2) / tileSize,
		trackMouse: false,
	})

	g.objects = append(g.objects, Object{
		name:       "d6too",
		image:      img,
		x:          (screenWidth / 3) / tileSize,
		y:          (screenHeight / 3) / tileSize,
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
