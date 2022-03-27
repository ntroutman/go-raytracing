package rnd

import (
	"math/rand"
	"time"
)

func Init() {
	rand.Seed(time.Now().UnixNano())

}

func Float64() float64 {
	return rand.Float64()
}

func Float64Range(min, max float64) float64 {
	return min + (max-min)*rand.Float64()
}
