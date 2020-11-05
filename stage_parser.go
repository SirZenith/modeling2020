package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// ParserMap is a map recording parser for each section
var ParserMap = map[string]func(*Stage, string) error{}

func init() {
	ParserMap["day count"] = parseDayCount
	ParserMap["load"] = parseLoad
	ParserMap["base budget"] = parseBudget
	ParserMap["base income"] = parseBaseIncome
	ParserMap["weight & base price"] = parseWeightAndBasePrice
	ParserMap["base cost"] = parseBaseCost
	ParserMap["node count"] = parseNodeCount
	ParserMap["special node"] = parseSpecial
	ParserMap["adjacent releation"] = parseAdj
	ParserMap["weather"] = parseWeather
	ParserMap["path weight"] = parsePathWeight
}

func parseDayCount(s *Stage, line string) error {
	count, err := strconv.Atoi(line)
	if err != nil {
		return err
	}
	s.dayCount = count
	return nil
}

func parseLoad(s *Stage, line string) error {
	count, err := strconv.Atoi(line)
	if err != nil {
		return err
	}
	s.load = count
	return nil
}

func parseBudget(s *Stage, line string) error {
	count, err := strconv.Atoi(line)
	if err != nil {
		return err
	}
	s.baseBudget = count
	return nil
}

func parseBaseIncome(s *Stage, line string) error {
	count, err := strconv.Atoi(line)
	if err != nil {
		return err
	}
	s.baseIncome = count
	return nil
}

func parseNodeCount(s *Stage, line string) error {
	count, err := strconv.Atoi(line)
	if err != nil {
		return err
	}
	s.nodeCount = count
	return nil
}

func parseWeightAndBasePrice(s *Stage, line string) error {
	parts := strings.Split(line, ":")
	if len(parts) != 3 {
		return errors.New("Wrong Sperator Usage")
	}
	kind, ok := resourceMap[parts[0]]
	if !ok {
		return errors.New("Unknown weather type")
	}
	weight, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	s.resourceWeight[kind] = weight
	price, err := strconv.Atoi(parts[2])
	if err != nil {
		return err
	}
	s.resourceBasePrice[kind] = price
	return nil
}

func parseBaseCost(s *Stage, line string) error {
	parts := strings.Split(line, ":")
	if len(parts) != 3 {
		return errors.New("Wrong Sperator Usage")
	}
	weatherKind, ok := weatherMap[parts[0]]
	if !ok {
		return fmt.Errorf("Unknown weather type %s", parts[0])
	}
	resourceKind, ok := resourceMap[parts[1]]
	if !ok {
		return fmt.Errorf("Unknown resource type %s", parts[1])
	}
	cost, err := strconv.Atoi(parts[2])
	if err != nil {
		return err
	}
	s.resourceBaseCost[weatherKind][resourceKind] = cost
	return nil
}

func parseSpecial(s *Stage, line string) error {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return errors.New("Wrong Sperator Usage")
	}
	id, kind := parts[0], parts[1]
	specialKind, ok := specialNodeMap[kind]
	if !ok {
		return errors.New("Error node type")
	}
	s.special[id] = specialKind
	return nil
}

func parseAdj(s *Stage, line string) error {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return errors.New("Wrong Sperator Usage")
	}
	id, neighbourStr := parts[0], parts[1]
	adjList, ok := s.adjacents[id]
	if !ok {
		adjList = []string{}
	}
	for _, neighbour := range strings.Split(neighbourStr, ",") {
		adjList = append(adjList, neighbour)
	}
	s.adjacents[id] = adjList
	return nil
}

func parseWeather(s *Stage, line string) error {
	weathers := strings.Split(line, ",")
	for _, weather := range weathers {
		weatherKind, ok := weatherMap[weather]
		if !ok {
			log.Println("Unknown weather type:", weather)
			return errors.New("Unknown weather type")
		}
		s.weatherList = append(s.weatherList, weatherKind)
	}
	return nil
}

func parsePathWeight(s *Stage, line string) error {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return errors.New("Wrong sperator usage")
	}
	weight, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	nodeIDs := strings.Split(parts[0], ",")
	if len(parts) != 2 {
		return errors.New("Wrong sperator usage")
	}
	for i := 0; i < 2; i++ {
		j := 1 - i
		weightMap, ok := s.weightMap[nodeIDs[i]]
		if !ok {
			weightMap = map[string]int{}
		}
		weightMap[nodeIDs[j]] = weight
		s.weightMap[nodeIDs[i]] = weightMap
	}
	return nil
}
