package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// OBJECT

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

// GAME FUNCTIONS

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

func (g *Game) DrawObjects(screen *ebiten.Image) {
	for _, object := range g.objects {
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(object.x*tileSize), float64(object.y*tileSize))
		screen.DrawImage(object.image, options)
	}
}
