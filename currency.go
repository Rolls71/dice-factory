package main

import (
	_ "image/png"
	"math"
)

// SellDie adds a value to DicePoints
func (g *Game) SellDie(item ItemType, value uint64) {
	g.Currencies[item] += value
}

// Cost returns the calculated cost of an ObjectType.
// Defaults to the max uint64 value.
func (g *Game) Cost(object ObjectType) (ItemType, uint64) {
	switch object {
	case ConveyorBelt:
		return Plain, uint64(math.Pow(float64(g.ObjectCount[object])+1, 2))
	case Builder:
		return Plain, uint64(math.Pow(2, float64(g.ObjectCount[object])+1))
	case Upgrader:
		return Plain, uint64(math.Pow(3, float64(g.ObjectCount[object])+1) * 10)
	default:
		return Plain, maxUint64
	}
}

// Pay subtracts given value from DicePoints unless value is less than
// DicePoints. Returns true if the payment was successful
func (g *Game) Pay(item ItemType, value uint64) bool {
	if g.Currencies[item] >= value {
		g.Currencies[item] -= value
		return true
	}
	return false
}

// Buy will attempt to Pay for an object and spawn it if successful
func (g *Game) Buy(object ObjectType, x, y int, objectFacing CardinalDir) {
	if g.Pay(g.Cost(object)) {
		g.SpawnObject(object, x, y, objectFacing)
	}
}
