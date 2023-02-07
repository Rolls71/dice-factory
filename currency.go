package main

import (
	_ "image/png"
	"math"
)

type Currency struct {
	DicePoints uint64
}

func (g *Game) Cost(object ObjectType) uint64 {
	switch object {
	case ConveyorBelt:
		return uint64(math.Pow(float64(g.ObjectCount[object])+1, 2))
	case Builder:
		return uint64(math.Pow(10, float64(g.ObjectCount[object])+1))
	case Upgrader:
		return uint64(math.Pow(10, float64(g.ObjectCount[object])+1) * 30)
	}
	return maxUint64
}

func (g *Game) AddDie(value uint64) {
	g.Balance.DicePoints += (value)
}

func (g *Game) Pay(value uint64) bool {
	if g.Balance.DicePoints >= value {
		g.Balance.DicePoints -= value
		return true
	}
	return false
}

func (g *Game) Buy(object ObjectType, x, y int, objectFacing ObjectFacing) {
	if g.Pay(g.Cost(object)) {
		g.SpawnObject(object, x, y, objectFacing)
	}
}
