package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TruckType int

const (
	BasicTruck = iota
)

const truckArrivalTime float64 = 2

type Truck struct {
	X, Y             float64
	SpawnX, SpawnY   float64
	TargetX, TargetY float64
	ID               uint64 // unique generated identifier
	Truck            TruckType
	Storage          *Storage // associated Storage
	Collectors       []*Object
	Width, Height    int // width along x axis

	PercentComplete float64 // 0 to 1
	IsExiting       bool
}

func (t *Truck) Send() {
	for _, collector := range t.Collectors {
		collector.IsCollecting = false
	}
	t.IsExiting = true
}

// Step moves the truck towards its target with decreasing velocity.
// If truck IsExiting, it moves away from its target with increasing velocity.
// Returns true on the same tick of arrival
func (t *Truck) Step() bool {
	onComplete := false
	if t.IsExiting {
		t.PercentComplete -= frameDelta / truckArrivalTime
		if t.PercentComplete < 0 {
			t.PercentComplete = 0
			onComplete = true
		}
	} else {
		t.PercentComplete += frameDelta / truckArrivalTime
		if t.PercentComplete > 1 {
			t.PercentComplete = 1
			onComplete = true
		}
	}
	totalDistance := t.TargetX - t.SpawnX
	t.X = (-math.Pow(t.PercentComplete-1, 2)+1)*totalDistance + t.SpawnX

	totalDistance = t.TargetY - t.SpawnY
	t.Y = (-math.Pow(t.PercentComplete-1, 2)+1)*totalDistance + t.SpawnY

	return onComplete
}

func (g *Game) UpdateTrucks() {
	for _, truck := range g.Trucks {
		// Is the truck currently being loaded?
		if truck.Collectors[0].IsCollecting {
			continue
		}

		// Is the truck at it's destination?
		if (truck.PercentComplete == 1 && !truck.IsExiting) ||
			(truck.PercentComplete == 0 && truck.IsExiting) {
			continue
		}

		// Step truck, and on the last frame enable collectors if arriving
		if truck.Step() {
			if !truck.IsExiting {
				for _, collector := range truck.Collectors {
					collector.IsCollecting = true
				}
			} else {
				// Spawn a new copy of this truck
				g.SpawnTruck(
					truck.Truck,
					truck.Collectors,
					ToTile(truck.SpawnX),
					ToTile(truck.SpawnY),
					ToTile(truck.TargetX),
					ToTile(truck.TargetY),
					truck.Width,
					truck.Height,
				)
				// Load trucks contents into Warehouse
				g.Warehouse.Load(truck.Storage)
				// Delete old version of truck
				delete(g.Trucks, truck.ID)
			}
		}
	}
}

func (g *Game) SpawnTruck(
	truckType TruckType,
	collectors []*Object,
	spawnX, spawnY int,
	targetX, targetY int,
	width, height int) *Truck {
	storage := g.NewStorage(TruckTrailer, truckCapacity, truckTypeLimit)
	g.Storages[storage.ID] = storage

	truck := &Truck{
		X:          ToReal(spawnX),
		Y:          ToReal(spawnY),
		SpawnX:     ToReal(spawnX),
		SpawnY:     ToReal(spawnY),
		TargetX:    ToReal(targetX),
		TargetY:    ToReal(targetY),
		ID:         g.NextID(),
		Truck:      truckType,
		Storage:    storage,
		Collectors: collectors,
		Width:      width,
		Height:     height,
	}

	if len(truck.Collectors) < 1 {
		log.Fatal("Error: Truck must have at least one collector")
	}

	g.Trucks[truck.ID] = truck

	return truck
}

func (g *Game) NewTruck(truckType TruckType, imageName string) {
	path := "images/" + imageName
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	g.truckImages[truckType] = img
}

func (g *Game) DrawTrucks(screen *ebiten.Image) {
	for _, truck := range g.Trucks {
		img := g.truckImages[truck.Truck]
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Scale(
			float64(tileSize)/(float64(img.Bounds().Dx())/float64(truck.Width)),
			float64(tileSize)/(float64(img.Bounds().Dy())/float64(truck.Height)))
		options.GeoM.Translate(truck.X, truck.Y)
		screen.DrawImage(img, options)
	}
}
