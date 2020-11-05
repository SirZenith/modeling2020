package main

import (
	"flag"
	"log"
)

const (
	colorNone   string = "\033[0m"
	colorRed           = "\033[1;31m"
	colorGreen         = "\033[32m"
	colorYellow        = "\033[33m"
)

var stageFile = flag.String("stage", "stage.txt", "stage file to read from")
var initDate = flag.Int("date", 0, "Intial date of process")
var initPosition = flag.String("pos", "st", "Initial position of traveler")
var initMoney = flag.Int("money", -1, "Initial budget for travler")
var initWater = flag.Int("water", 0, "Initial water for travler")
var initFood = flag.Int("food", 0, "Initial food for travler")
var notFirstBuy = flag.Bool("first", false, "Initial state of first state")
var isInverse = flag.Bool("inverse", false, "Inverse traveling process, cumulate resource from current place")

func main() {
	flag.Parse()
	t := newTraveler(*stageFile)
	if t == nil {
		log.Println("Traveler initialize failed")
		return
	}
	travelerInit(t)
	states := newRecorder(t)
	shell(t, states)
}

func travelerInit(t *Traveler) {
	t.date = *initDate
	t.position = *initPosition
	if *initMoney >= 0 {
		t.money = *initMoney
	}
	t.food = *initFood
	t.water = *initWater
	t.loadSpace = t.loadSpace - t.food*t.resourceWeight[resourceFood] + t.water*t.resourceWeight[resourceWater]
	t.firstBuy = !*notFirstBuy
}
