package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tile struct {
	name  string
	image *ebiten.Image
}

func (g *Game) NewTile(
	name,
	imageName string,
) {
	path := "images/" + imageName
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	g.tileSet = append(g.tileSet, Tile{
		name:  name,
		image: img,
	})
}

// SetTileStage sets the value of every Tile on-screen
// Requires a 2D array of integers
// given array must be of type [stageSizeY][stageSizeX]int
// integer values correspond to TileSet element indices
// e.g. 0 -> the first NewTile(), 1 -> the second NewTile() ...
func (g *Game) SetTileStage(tileArray [stageSizeY][stageSizeX]int) {
	g.tileStage = tileArray
}

// DrawTiles will draw every Tile in the game's list of objects.
// Tiles are drawn on their stored grid coordinate.
func (g *Game) DrawTiles(screen *ebiten.Image) {
	for y := 0; y < stageSizeY; y++ {
		for x := 0; x < stageSizeX; x++ {
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
			screen.DrawImage(g.tileSet[g.tileStage[y][x]].image, options)
		}
	}

}
