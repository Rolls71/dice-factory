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
	Truck
)

const truckSpeed float64 = float64(tileSize)

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
func (i *Item) Step(speed float64) {
	xDelta := ToReal(i.TargetX) - i.X
	if math.Abs(xDelta) < speed*frameDelta {
		if i.X != ToReal(i.TargetX) {
			i.CatchupX += speed*frameDelta - math.Abs(xDelta)
		}
		i.X = ToReal(i.TargetX)
	} else {
		if xDelta > 0 {
			i.X += speed*frameDelta + i.CatchupX
		} else {
			i.X -= speed*frameDelta + i.CatchupX
		}
		i.CatchupX = 0
	}

	yDelta := ToReal(i.TargetY) - i.Y
	if math.Abs(yDelta) < speed*frameDelta {
		if i.Y != ToReal(i.TargetY) {
			i.CatchupY += speed*frameDelta - math.Abs(yDelta)
		}
		i.Y = ToReal(i.TargetY)
	} else {
		if yDelta > 0 {
			i.Y += speed*frameDelta + i.CatchupY
		} else {
			i.Y -= speed*frameDelta + i.CatchupY
		}
		i.CatchupY = 0
	}
}

// UpdateObjects will iterate through each Item and switch,
// depending on their type. Each Item type may have different functionality.
func (g *Game) UpdateItems() {
	for _, item := range g.Items {
		switch item.Item {
		case Truck:
			// has item reached target position?
			if item.X == ToReal(item.TargetX) &&
				item.Y == ToReal(item.TargetY) {
				continue
			}
			g.Items[item.ID].Step(truckSpeed)
		default:
			isObject, _ := g.GetObjectAt(item.TargetX, item.TargetY)
			if !isObject {
				delete(g.Items, item.ID)
				continue
			}
			// has item reached target position?
			if item.X == ToReal(item.TargetX) &&
				item.Y == ToReal(item.TargetY) {
				continue
			}
			g.Items[item.ID].Step(conveyorSpeed)
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
	g.itemImages[itemType] = img
}

// SpawnItem will create an instance of an Item in the set.
// The Item's position and Target position will be set to that of the creator.
func (g *Game) SpawnItem(itemType ItemType, creator *Object) *Item {
	item := &Item{}
	x, y := creator.X, creator.Y
	item.SetID(g.NextID())
	item.Item = itemType
	switch itemType {
	case Truck:
		item.SetRealCoordinate(ToReal(0), ToReal(y))
	default:
		item.Roll()
		item.SetRealCoordinate(ToReal(x), ToReal(y))

	}
	item.SetTargetPosition(x, y)

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
	for _, item := range g.Items {
		itemArray = append(itemArray, item)
	}

	sort.SliceStable(itemArray, func(i, j int) bool {
		if itemArray[i].Item == Truck && itemArray[j].Item != Truck {
			return false
		} else {
			return itemArray[i].ID > itemArray[j].ID
		}
	})

	for _, item := range itemArray {
		img := g.itemImages[item.Item]
		options := &ebiten.DrawImageOptions{}
		switch item.Item {
		case Truck:
			img = g.itemImages[Truck]
			options.GeoM.Scale(float64(tileSize*4)/float64(img.Bounds().Dx()),
				float64(tileSize*2)/float64(img.Bounds().Dy()))
			options.GeoM.Translate(item.X-float64(tileSize*3), item.Y)
		default:
			itemIndex := (item.Face - 1) * img.Bounds().Dy()
			img = img.SubImage(image.Rect(
				itemIndex,
				0,
				itemIndex+img.Bounds().Dy(),
				img.Bounds().Dy())).(*ebiten.Image)
			options.GeoM.Scale(float64(tileSize)/float64(img.Bounds().Dx()),
				float64(tileSize)/float64(img.Bounds().Dy()))
			options.GeoM.Translate(item.X, item.Y)

			if item.Face == 0 {
				log.Fatal("Error: Item has no set face")
			}
		}
		screen.DrawImage(img, options)
	}

}
