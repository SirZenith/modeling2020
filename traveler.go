package main

import (
	"bytes"
	"fmt"
	"log"
)

// Traveler is object that record traveling state, like position, date, resource
type Traveler struct {
	*Stage
	*Graph
	date          int
	position      string
	loadSpace     int // load space left for resource
	money         int
	water         int
	food          int
	ok            bool // wheather traveler is in a normal state
	firstBuy      bool
	heuristicFunc func(*Traveler) string
}

func newTraveler(stageFile string) *Traveler {
	stage, err := stageFromFile(stageFile)
	if err != nil {
		log.Println(err)
		return nil
	}
	graph := stage.makeGraph()
	return &Traveler{
		Stage:     stage,
		Graph:     graph,
		date:      0,
		position:  graph.starting.id,
		loadSpace: stage.load,
		money:     stage.baseBudget,
		water:     0,
		food:      0,
		ok:        true,
		firstBuy:  true,
	}
}

func (t *Traveler) String() string {
	buf := bytes.NewBufferString("")
	var state string
	if t.ok {
		state = colorGreen + "Ok" + colorNone
	} else {
		state = colorRed + "Error" + colorNone
	}
	fmt.Fprintf(buf, "State: %s\n", state)
	if t.date < 30 {
		fmt.Fprintf(buf, "Weather Tommorrow: %s\n", t.Stage.weatherList[t.date])
	}
	fmt.Fprint(buf, "| Date | Position | Load Space |  Money  | Food | Water |\n")
	fmt.Fprintf(buf, "|%6d|%10s|%12d|%9d|%6d|%7d|", t.date, t.position, t.loadSpace, t.money, t.food, t.water)
	return buf.String()
}

func (t *Traveler) checkState() bool {
	inTime := (t.date < t.dayCount) || (t.date == t.dayCount && t.position == t.ending.id)
	t.ok = inTime && t.loadSpace >= 0 // && t.water >= 0 && t.food >= 0
	return t.ok
}

func (t *Traveler) stay() bool {
	t.consumeResource(1)
	t.date++
	return true
}

func (t *Traveler) moveTo(id string) bool {
	_, ok := t.nodes[id]
	if !ok {
		return false
	}
	distance, ok := t.nodes[t.position].pathWeight[id]
	if !ok {
		distance = 1
	}
	for distance > 0 {
		if t.weatherList[t.date] == sandStorm {
			t.consumeResource(1)
		} else {
			t.consumeResource(2)
			distance--
		}
		t.date++
		// fmt.Println(t.date, distance)
	}
	t.position = id
	return true
}

func (t *Traveler) buyResource(foodAmount int, waterAmount int) bool {
	// if t.starting.id == t.position && t.firstBuy || t.nodes[t.position].nodeType == villageNode {
	// 	t.firstBuy = false
	// } else {
	// 	// return false
	// }
	cost := foodAmount*t.resourceBasePrice[resourceFood] + waterAmount*t.resourceBasePrice[resourceWater]
	if !t.firstBuy {
		cost *= 2
	}
	weight := foodAmount*t.resourceWeight[resourceFood] + waterAmount*t.resourceWeight[resourceWater]
	t.money -= cost
	t.loadSpace -= weight
	t.food += foodAmount
	t.water += waterAmount
	t.firstBuy = false
	return true
}

func (t *Traveler) mining() bool {
	if t.nodes[t.position].nodeType != mineNode {
		return false
	}
	t.consumeResource(3)
	t.money += t.baseIncome
	t.date++
	return true
}

func (t *Traveler) consumeResource(multiplier int) {
	var weather WeatherType
	if *isInverse {
		weather = t.weatherList[t.date-1]
	} else {
		weather = t.weatherList[t.date]
	}
	foodCost := t.resourceBaseCost[weather][resourceFood] * multiplier
	waterCost := t.resourceBaseCost[weather][resourceWater] * multiplier
	t.food -= foodCost
	t.water -= waterCost
	if t.loadSpace < t.Stage.load {
		t.loadSpace += foodCost*t.resourceWeight[resourceFood] + waterCost*t.resourceWeight[resourceWater]
	}
}

func (t *Traveler) moveWithHeuristic() {
	if t.heuristicFunc == nil {
		return
	}
	t.moveTo(t.heuristicFunc(t))
}
