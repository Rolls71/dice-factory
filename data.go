package main

import (
	_ "image/png"
	"log"
	"math/rand"
)

type Data struct {
	dicePoints uint64
}

func (d *Data) RollDie(min, max int) {
	if min < 0 || max < 0 {
		log.Fatal("Error: Cannot roll a negative number.")
	}
	d.dicePoints += uint64(rand.Intn(max) + min)
}
