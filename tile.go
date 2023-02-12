package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TileType int

const (
	BasicGrass = iota
	LongGrass
)

type Tile struct {
	Name  string
	Image *ebiten.Image
}

// NewTile adds a new type of Tiles to the game's tileSet.
func (g *Game) NewTile(
	tile TileType,
	imageName string,
) {
	path := "images/" + imageName
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	g.tileImages[tile] = img
}

// DrawTiles will draw every Tile in the game's list of objects.
// Tiles are drawn on their stored grid coordinate.
func (g *Game) DrawTiles(screen *ebiten.Image) {
	for y := 0; y < stageSizeY; y++ {
		for x := 0; x < stageSizeX; x++ {
			img := g.tileImages[TileType(g.TileStage[y][x])]
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Scale(float64(tileSize)/float64(img.Bounds().Dx()),
				float64(tileSize)/float64(img.Bounds().Dy()))
			options.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
			screen.DrawImage(img, options)
		}
	}

}
