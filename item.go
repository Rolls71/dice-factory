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

/*const (
	itemSpeed float64 = 1 // pixels per tick
)*/

type Item struct {
	itemType         ItemType
	image            *ebiten.Image
	x, y             int // pixel coords
	targetX, targetY int // index of target object
}

/*func (i *Item) StepToTarget() {
	xDir := 1
	if i.x > i.targetX {
		xDir = -1
	}
	yDir := 1
	if i.y > i.targetY {
		yDir = -1
	}
	i.x += int(itemSpeed) * xDir
	i.y += int(itemSpeed) * yDir
}*/

func (i *Item) SetPosition(x, y int) {
	i.x = x
	i.y = y
}

func (i *Item) SetTarget(x, y int) {
	i.targetX = x
	i.targetY = y
}

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

func (g *Game) GetItemOn(objectIndex int) (bool, int) {
	for itemIndex, _ := range g.items {
		if g.items[itemIndex].targetX == g.objects[objectIndex].x &&
			g.items[itemIndex].targetY == g.objects[objectIndex].y {
			return true, itemIndex
		}
	}
	return false, -1
}

func (g *Game) SpawnItem(itemType ItemType, creator *Object) int {
	fmt.Printf("spawn\n")
	len := len(g.items)
	g.items = append(g.items, g.itemSet[itemType])
	x, y := creator.x, creator.y
	g.items[len].SetPosition(x*tileSize, y*tileSize)
	isObject, i := g.GetNeighborOf(creator)
	if isObject {
		g.items[len].SetTarget(g.objects[i].x*tileSize, g.objects[i].y*tileSize)
		fmt.Printf("target %d, %d\n", g.objects[i].x, g.objects[i].y)
	} else {
		log.Fatal("GetNeighborOf Failed: No neighbor")
	}
	return len
}

func (g *Game) UpdateItems() {
	for _, item := range g.items {
		switch item.itemType {
		case PlainItem: // doStuff
		}
		/*if item.x != item.targetX ||
			item.y != item.targetY {
			g.items[index].StepToTarget()
		}*/
	}
}

func (g *Game) DrawItems(screen *ebiten.Image) {
	for _, item := range g.items {
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(item.x), float64(item.y))
		screen.DrawImage(item.image, options)
	}
}
