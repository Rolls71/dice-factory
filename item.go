package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ItemType int

const (
	PlainItem ItemType = iota
)

type Item struct {
	itemType         ItemType
	image            *ebiten.Image
	x, y             int // pixel coords
	targetX, targetY int // index of target object
}

func (i *Item) SetPosition(x, y int) {
	i.x = x
	i.y = y
}

func (i *Item) SetTarget(x, y int) {
	i.targetX = x
	i.targetY = y
}

// NewItem will create a new item of given image and type
// Other struct elements will default
func (g *Game) NewItem(item ItemType, imageName string) {
	path := "images/" + imageName
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	g.itemSet = append(g.itemSet, Item{
		itemType: item,
		image:    img,
	})
}

// GetItemOn will find an Item targetting a given Object.
// if an item is not found, it will return false and an Empty Object Reference.
func (g *Game) GetItemOn(object *Object) (bool, *Item) {
	for index, copy := range g.items {
		if copy.targetX == object.x &&
			copy.targetY == object.y {
			return true, &g.items[index]
		}
	}
	return false, &Item{}
}

// SpawnItem will create an instance of an Item in the set.
// The Item's position and Target position will be set to that of the creator.
func (g *Game) SpawnItem(itemType ItemType, creator *Object) *Item {
	fmt.Printf("spawn\n")
	len := len(g.items)
	g.items = append(g.items, g.itemSet[itemType])
	item := &g.items[len]
	x, y := creator.x, creator.y
	item.SetPosition(x*tileSize, y*tileSize)
	item.SetTarget(x*tileSize, y*tileSize)
	return item
}

// UpdateObjects will iterate through each Item and switch,
// depending on their type. Each Item type may have different functionality.
func (g *Game) UpdateItems() {
	for _, item := range g.items {
		switch item.itemType {
		case PlainItem:
		}
	}
}

// DrawItems draws each Item at a pixel coordinate
func (g *Game) DrawItems(screen *ebiten.Image) {
	for _, copy := range g.items {
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(copy.x), float64(copy.y))
		screen.DrawImage(copy.image, options)
	}
}
