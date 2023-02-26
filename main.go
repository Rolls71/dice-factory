package main

import (
	"log"
	"os"

	"github.com/Rolls71/dice-factory/dicefactory"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(dicefactory.ScreenWidth, dicefactory.ScreenHeight)
	ebiten.SetWindowTitle("Dice Factory")
	ebiten.SetFullscreen(true)

	var game *dicefactory.Game
	if _, err := os.Stat(dicefactory.SaveFilename); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			game = dicefactory.NewGame()
		} else {
			// other error
			log.Fatal(err)
		}
	} else {
		game = dicefactory.LoadGame(dicefactory.SaveFilename)
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
	game.SaveGame()
}
