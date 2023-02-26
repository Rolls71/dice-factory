package dicefactory

import (
	_ "image/png"
	"math"
	"math/rand"
)

type CurrencyType int

const (
	PlainBuck = iota
	GoldBuck
)

const sellRate = 4 // secs per sell

// Cost returns the calculated cost of an ObjectType.
// Defaults to the max uint64 value.
func (g *Game) Cost(object ObjectType) (CurrencyType, uint64) {
	switch object {
	case ConveyorBelt:
		return PlainBuck, uint64(math.Pow(float64(g.ObjectCount[object])+1, 2))
	case Builder:
		return PlainBuck, uint64(math.Pow(2, float64(g.ObjectCount[object])+1))
	case Upgrader:
		return PlainBuck, uint64(math.Pow(3, float64(g.ObjectCount[object])+1) * 10)
	default:
		return PlainBuck, maxUint64
	}
}

// Pay subtracts given value from DicePoints unless value is less than
// DicePoints. Returns true if the payment was successful
func (g *Game) Pay(currencyType CurrencyType, value uint64) bool {
	if g.Currencies[currencyType] >= value {
		g.Currencies[currencyType] -= value
		return true
	}
	return false
}

// Buy will attempt to Pay for an object and spawn it if successful
func (g *Game) Buy(objectType ObjectType, x, y int, objectFacing CardinalDir) {
	if g.Pay(g.Cost(objectType)) {
		g.SpawnObject(objectType, x, y, objectFacing)
	}
}

func (g *Game) UpdateCurrency() {
	if g.ticks%(uint64(frameRate)*sellRate) == 0 {
		g.SellRandom()
	}
}

// Sell adds the face of the die to the correct currency.
// Sell is often best used with RemoveDie
func (g *Game) Sell(itemType ItemType, face int) {
	switch itemType {
	case PlainD6:
		g.Currencies[PlainBuck] += uint64(face)
	case GoldD6:
		g.Currencies[GoldBuck] += uint64(face)
	}

}

// SellRandom sells a random dice in the warehouse
func (g *Game) SellRandom() {
	var item ItemType
	var face int

	// are there any dice to sell
	if g.Warehouse.Count <= 0 {
		return
	}

	// pick a type
	pick := rand.Intn(len(g.Warehouse.Dice))
	for pickItem := range g.Warehouse.Dice {
		if pick == 0 {
			item = pickItem
			break
		}
		pick--
	}

	// pick a face
	pick = rand.Intn(len(g.Warehouse.Dice[item]))
	for pickFace := range g.Warehouse.Dice[item] {
		if pick == 0 {
			face = pickFace
			break
		}
		pick--
	}

	// attempt to remove that Die
	if !g.Warehouse.RemoveDie(item, face) {
		g.SellRandom()
		return
	}

	g.Sell(item, face)
}
