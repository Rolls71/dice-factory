package main

import (
	_ "image/png"
)

type Data struct {
	dicePoints uint64
}

func (d *Data) AddDie(value uint64) {
	d.dicePoints += (value)
}
