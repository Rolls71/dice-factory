package main

import (
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	name  string
	image *ebiten.Image
}

func (g *Game) InitTiles() {
	//replace with images later
	tempGrass := ebiten.NewImage(tileSize, tileSize)
	tempGrass.Fill(color.RGBA{0, 0xFF, 0, 0xFF})

	g.tileSet = []Tile{
		{
			name:  "basicGrass",
			image: tempGrass,
		},
	}

}

func (g *Game) DrawTiles(screen *ebiten.Image) {
	for y := 0; y < stageSizeY; y++ {
		for x := 0; x < stageSizeX; x++ {
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
			screen.DrawImage(g.tileSet[g.tileStage[y][x]].image, options)
		}
	}

}
