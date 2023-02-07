package main

import (
	"image"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ItemType int

const (
	PlainD6 ItemType = iota
	GoldD6
)

const (
	d6Min          int    = 1
	d6Max          int    = 6
	goldMultiplier uint64 = 2
)

type Item struct {
	Item               ItemType
	Face               int // value shown on face
	X, Y               float64
	ID                 uint64  // unique generated identifier
	TargetX, TargetY   int     // index of target object
	CatchupX, CatchupY float64 // if item is behind, saves lost distance
}

func (i *Item) SetRealCoordinate(x, y float64) {
	i.X = x
	i.Y = y
}

func (i *Item) SetTargetPosition(x, y int) {
	i.TargetX = x
	i.TargetY = y
}

func (i *Item) SetID(id uint64) {
	i.ID = id
}

func (i *Item) Value() uint64 {
	switch i.Item {
	case PlainD6:
		return uint64(i.Face)
	case GoldD6:
		return uint64(i.Face) * goldMultiplier
	}
	log.Fatal("Error: unknown itemType")
	return 0
}

func (i *Item) Roll() {
	i.Face = rand.Intn(d6Max) + d6Min
}

// Step moves an item conveyorSpeed units per second towards target
// Stores catchup when theres movement left, adds on next movement
func (i *Item) Step() {
	xDelta := ToReal(i.TargetX) - i.X
	if math.Abs(xDelta) < conveyorSpeed*frameDelta {
		if i.X != ToReal(i.TargetX) {
			i.CatchupX += conveyorSpeed*frameDelta - math.Abs(xDelta)
		}
		i.X = ToReal(i.TargetX)
	} else {
		if xDelta > 0 {
			i.X += conveyorSpeed*frameDelta + i.CatchupX
		} else {
			i.X -= conveyorSpeed*frameDelta + i.CatchupX
		}
		i.CatchupX = 0
	}

	yDelta := ToReal(i.TargetY) - i.Y
	if math.Abs(yDelta) < conveyorSpeed*frameDelta {
		if i.Y != ToReal(i.TargetY) {
			i.CatchupY += conveyorSpeed*frameDelta - math.Abs(yDelta)
		}
		i.Y = ToReal(i.TargetY)
	} else {
		if yDelta > 0 {
			i.Y += conveyorSpeed*frameDelta + i.CatchupY
		} else {
			i.Y -= conveyorSpeed*frameDelta + i.CatchupY
		}
		i.CatchupY = 0
	}
}

// UpdateObjects will iterate through each Item and switch,
// depending on their type. Each Item type may have different functionality.
func (g *Game) UpdateItems() {
	for _, copy := range g.Items {
		isObject, _ := g.GetObjectAt(copy.TargetX, copy.TargetY)
		if !isObject {
			delete(g.Items, copy.ID)
			continue
		}
		switch copy.Item {
		case PlainD6:
		}
		// has item reached target position?
		if copy.X == ToReal(copy.TargetX) &&
			copy.Y == ToReal(copy.TargetY) {
			continue
		}
		g.Items[copy.ID].Step()
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
	g.itemImages[itemType] = img
}

// SpawnItem will create an instance of an Item in the set.
// The Item's position and Target position will be set to that of the creator.
func (g *Game) SpawnItem(itemType ItemType, creator *Object) *Item {
	item := &Item{}
	x, y := creator.X, creator.Y
	item.SetID(g.NextID())
	item.SetRealCoordinate(ToReal(x), ToReal(y))
	item.SetTargetPosition(x, y)
	item.Roll()

	g.Items[item.ID] = item
	return item
}

func (g *Game) SetItem(item *Item, itemType ItemType) {
	item.Item = itemType
}

// GetItemTargeting will find an Item targeting a given Object.
// if an item is not found, it will return false and an Empty Object Reference.
func (g *Game) GetItemTargeting(object *Object) (bool, *Item) {
	for _, copy := range g.Items {
		if copy.TargetX == object.X &&
			copy.TargetY == object.Y {
			return true, g.Items[copy.ID]
		}
	}
	return false, &Item{}
}

// DrawItems draws each Item at a pixel coordinate
func (g *Game) DrawItems(screen *ebiten.Image) {
	itemArray := []*Item{}
	for _, copy := range g.Items {
		itemArray = append(itemArray, copy)
	}

	sort.SliceStable(itemArray, func(i, j int) bool {
		return itemArray[i].ID < itemArray[j].ID
	})

	for _, copy := range itemArray {
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(copy.X), float64(copy.Y))

		if copy.Face == 0 {
			log.Fatal("Error: Item has no set face")
		}
		itemIndex := (copy.Face - 1) * tileSize
		screen.DrawImage(g.itemImages[copy.Item].SubImage(image.Rect(
			itemIndex,
			0,
			itemIndex+tileSize,
			tileSize)).(*ebiten.Image), options)
	}

}
