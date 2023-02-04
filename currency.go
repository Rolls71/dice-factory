package main

import (
	_ "image/png"
)

type Currency struct {
	DicePoints uint64
}

func (d *Currency) AddDie(value uint64) {
	d.DicePoints += (value)
}
