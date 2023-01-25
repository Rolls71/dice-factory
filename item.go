package main

import (
	_ "image/png"
	"log"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ItemType int

const (
	PlainItem ItemType = iota
)

const itemSpeed float64 = 32 // pixels per second

type Item struct {
	itemType         ItemType
	image            *ebiten.Image
	x, y             float64
	id               uint64 // unique generated identifier
	xTarget, yTarget int    // index of target object
}

func (i *Item) SetRealCoordinate(x, y float64) {
	i.x = x
	i.y = y
}

func (i *Item) SetTargetPosition(x, y int) {
	i.xTarget = x
	i.yTarget = y
}

func (i *Item) SetID(id uint64) {
	i.id = id
}

// Step moves an item itemSpeed units per second towards target
//
// TODO: fix minor delay when target is reached
func (i *Item) Step() {
	xDelta := ToReal(i.xTarget) - i.x
	if math.Pow(xDelta, 2) < 1 {
		// not moving full speed causes building delay
		i.x = ToReal(i.xTarget)
	} else {
		if xDelta > 0 {
			i.x += itemSpeed * frameDelta
		} else {
			i.x -= itemSpeed * frameDelta
		}
	}

	yDelta := ToReal(i.yTarget) - i.y
	if math.Pow(yDelta, 2) < 1 {
		i.y = ToReal(i.yTarget)
	} else {
		if yDelta > 0 {
			i.y += itemSpeed * frameDelta
		} else {
			i.y -= itemSpeed * frameDelta
		}
	}
}

// UpdateObjects will iterate through each Item and switch,
// depending on their type. Each Item type may have different functionality.
func (g *Game) UpdateItems() {
	for _, copy := range g.items {
		isObject, _ := g.GetObjectAt(copy.xTarget, copy.yTarget)
		if !isObject {
			delete(g.items, copy.id)
			continue
		}
		switch copy.itemType {
		case PlainItem:
			// has item reached target position?
			if copy.x == ToReal(copy.xTarget) &&
				copy.y == ToReal(copy.yTarget) {
				continue
			}
			g.items[copy.id].Step()
		}
	}
}

// NewItem will create a new item of given image and type
// Other struct elements will default
func (g *Game) NewItem(itemType ItemType, imageName string) {
	path := "images/" + imageName
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	g.itemSet[itemType] = &Item{
		itemType: itemType,
		image:    img,
	}
}

// SpawnItem will create an instance of an Item in the set.
// The Item's position and Target position will be set to that of the creator.
func (g *Game) SpawnItem(itemType ItemType, creator *Object) *Item {
	item := *g.itemSet[itemType]
	x, y := creator.x, creator.y
	item.SetID(g.NextID())
	item.SetRealCoordinate(ToReal(x), ToReal(y))
	item.SetTargetPosition(x, y)

	g.items[item.id] = &item
	return &item
}

// GetItemTargeting will find an Item targeting a given Object.
// if an item is not found, it will return false and an Empty Object Reference.
func (g *Game) GetItemTargeting(object *Object) (bool, *Item) {
	for _, copy := range g.items {
		if copy.xTarget == object.x &&
			copy.yTarget == object.y {
			return true, g.items[copy.id]
		}
	}
	return false, &Item{}
}

// DrawItems draws each Item at a pixel coordinate
func (g *Game) DrawItems(screen *ebiten.Image) {
	itemArray := []*Item{}
	for _, copy := range g.items {
		itemArray = append(itemArray, copy)
	}

	sort.SliceStable(itemArray, func(i, j int) bool {
		return itemArray[i].id < itemArray[j].id
	})

	for _, copy := range itemArray {
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(copy.x), float64(copy.y))
		screen.DrawImage(copy.image, options)
	}

}
