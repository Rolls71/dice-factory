package main

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	//"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// CONSTANTS
const (
	screenWidth  int = 640
	screenHeight int = 480
)

const (
	tileSize int = 32
)

const (
	stageSizeX int = screenWidth / tileSize
	stageSizeY int = screenHeight / tileSize
)

// VARIABLES

// TILES
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

// Objects

type Object struct {
	name  string
	image *ebiten.Image
	x     int
	y     int
}

func (g *Game) InitObjects() {
	var img *ebiten.Image
	var err error
	img, _, err = ebitenutil.NewImageFromFile("images/d6_6.png")
	if err != nil {
		log.Fatal(err)
	}

	g.objects = append(g.objects, Object{
		name:  "d6",
		image: img,
		x:     (screenWidth / 2) / tileSize,
		y:     (screenHeight / 2) / tileSize,
	})
}

func (g *Game) DrawObjects(screen *ebiten.Image) {
	for _, object := range g.objects {
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(object.x*tileSize), float64(object.y*tileSize))
		screen.DrawImage(object.image, options)
	}
}

// GAME
type Game struct {
	tileStage [stageSizeY][stageSizeX]int
	tileSet   []Tile
	objects   []Object
}

func NewGame() *Game {
	game := Game{
		// return a 2D array of zeroes
		tileStage: [stageSizeY][stageSizeX]int{},
		tileSet:   []Tile{},
		objects:   []Object{},
	}

	game.InitTiles()
	game.InitObjects()

	return &game
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawTiles(screen)
	g.DrawObjects(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (_screenWidth, _screenHeight int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
