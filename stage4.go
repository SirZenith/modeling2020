package main

import (
	"fmt"
	"math/rand"
)

var foodExpect = map[string]map[string]float64{}
var waterExpect = map[string]map[string]float64{}

func randWeather(pSun float64, pHigh float64, pSand float64, dayCount int) []WeatherType {
	pSand += pHigh
	pSun /= pSand
	pHigh /= pSand
	pSand = 1
	weatherList := []WeatherType{}
	var (
		value   float64
		weather WeatherType
	)
	for i := 0; i < dayCount; i++ {
		value = rand.Float64()
		switch {
		case value < pSun:
			weather = sunny
		case value < pHigh:
			weather = highTemp
		default:
			weather = sandStorm
		}
		weatherList = append(weatherList, weather)
	}
	return weatherList
}

func randomRun(_ []string, t *Traveler, state *StateRecorder) error {
	pSun, pHigh, pSand := 0.5, 0.4, 0.1
	loopCount := 30
	for i := 0; i < loopCount; i++ {
		t.Stage.weatherList = randWeather(pSun, pHigh, pSand, 8)
		for id := range t.Graph.nodes {
			for _, pos := range []string{"v", "m", "ed"} {
				resetState(t, id)
				t.moveTo(pos)
				addToExpectation(&foodExpect, id, pos, float64(t.food))
				addToExpectation(&waterExpect, id, pos, float64(t.water))

			}
		}
	}
	fmt.Println("Food Expectation")
	for id, special := range foodExpect {
		fmt.Println(id)
		for pos, value := range special {
			value /= float64(loopCount)
			fmt.Printf("\t%s: %.2f\n", pos, value)
			foodExpect[id][pos] = value
		}
	}
	fmt.Println("Water Expectation")
	for id, special := range waterExpect {
		fmt.Println(id)
		for pos, value := range special {
			value /= float64(loopCount)
			fmt.Printf("\t%s: %.2f\n", pos, value)
			waterExpect[id][pos] = value
		}
	}
	return nil
}

func resetState(t *Traveler, pos string) {
	t.date = 0
	t.position = pos
	t.food = 0
	t.water = 0
	t.loadSpace = t.Stage.load
	t.money = t.Stage.baseBudget
}

func addToExpectation(expectation *map[string]map[string]float64, id string, pos string, value float64) {
	expMap, ok := (*expectation)[id]
	if !ok {
		expMap = map[string]float64{}
	}
	expMap[pos] -= value
	(*expectation)[id] = expMap
}
