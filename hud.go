package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) DrawHUD(screen *ebiten.Image) {
	dicePointString := fmt.Sprintf("Dice Points: %d\n", g.Balance.DicePoints)
	ebitenutil.DebugPrint(screen, dicePointString)
}
