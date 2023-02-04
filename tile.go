package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tile struct {
	Name  string
	Image *ebiten.Image
}

// NewTile adds a new type of Tiles to the game's tileSet.
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
		Name:  name,
		Image: img,
	})
}

// SetTileStage sets the value of every Tile on-screen
// Requires a 2D array of Tile ID integers
// Given array must be of type [stageSizeY][stageSizeX]int
// Tile IDs correspond to TileSet element indices
// e.g. 0 -> the first NewTile(), 1 -> the second NewTile() ...
func (g *Game) SetTileStage(tileArray [stageSizeY][stageSizeX]int) {
	g.TileStage = tileArray
}

// DrawTiles will draw every Tile in the game's list of objects.
// Tiles are drawn on their stored grid coordinate.
func (g *Game) DrawTiles(screen *ebiten.Image) {
	for y := 0; y < stageSizeY; y++ {
		for x := 0; x < stageSizeX; x++ {
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
			screen.DrawImage(g.tileSet[g.TileStage[y][x]].Image, options)
		}
	}

}
