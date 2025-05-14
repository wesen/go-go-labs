package main

import (
	"time"

	"gonum.org/v1/gonum/stat/distuv"
)

type ParetoBehavior struct {
	Scale float64
	Shape float64
}

func (b *ParetoBehavior) NextBlockSize() int {
	pareto := distuv.Pareto{
		Xm:    b.Scale,
		Alpha: b.Shape,
	}
	return int(pareto.Rand())
}

func (b *ParetoBehavior) Sleep() {
	time.Sleep(100 * time.Millisecond) // Fixed sleep time for simplicity
}
