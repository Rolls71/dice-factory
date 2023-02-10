package main

import (
	_ "image/png"
	"math"
)

type Currency struct {
	DicePoints uint64
}

// Cost returns the calculated cost of an ObjectType.
// Defaults to the max uint64 value.
func (g *Game) Cost(object ObjectType) uint64 {
	switch object {
	case ConveyorBelt:
		return uint64(math.Pow(float64(g.ObjectCount[object])+1, 2))
	case Builder:
		return uint64(math.Pow(2, float64(g.ObjectCount[object])+1))
	case Upgrader:
		return uint64(math.Pow(3, float64(g.ObjectCount[object])+1) * 10)
	}
	return maxUint64
}

// AddDie adds a value to DicePoints
func (g *Game) AddDie(value uint64) {
	g.Balance.DicePoints += (value)
}

// Pay subtracts given value from DicePoints unless value is less than
// DicePoints. Returns true if the payment was successful
func (g *Game) Pay(value uint64) bool {
	if g.Balance.DicePoints >= value {
		g.Balance.DicePoints -= value
		return true
	}
	return false
}

// Buy will attempt to Pay for an object and spawn it if successful
func (g *Game) Buy(object ObjectType, x, y int, objectFacing ObjectFacing) {
	if g.Pay(g.Cost(object)) {
		g.SpawnObject(object, x, y, objectFacing)
	}
}
