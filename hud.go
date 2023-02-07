package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) DrawHUD(screen *ebiten.Image) {
	string := fmt.Sprintf("Dice Points: %d\n", g.Balance.DicePoints)
	string += fmt.Sprintf("Conveyor Belt: %d\n", g.Cost(ConveyorBelt))
	string += fmt.Sprintf("Builder: %d\n", g.Cost(Builder))
	string += fmt.Sprintf("Upgrader: %d\n", g.Cost(Upgrader))
	ebitenutil.DebugPrint(screen, string)

}
