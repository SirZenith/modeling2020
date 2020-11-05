package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

// Node is a single node in graph
type Node struct {
	id         string
	neighbour  map[*Node]struct{}
	nodeType   NodeType
	pathWeight map[string]int
}

func newNode(id string) *Node {
	return &Node{id, map[*Node]struct{}{}, normalNode, map[string]int{}}
}

// Graph graph object
type Graph struct {
	starting *Node
	ending   *Node
	nodes    map[string]*Node
}

func newGraph() *Graph {
	return &Graph{nil, nil, map[string]*Node{}}
}

func (g *Graph) appendAdj(id1 string, id2 string) error {
	node1, ok := g.nodes[id1]
	if !ok {
		node1 = newNode(id1)
		g.nodes[id1] = node1
	}
	node2, ok := g.nodes[id2]
	if !ok {
		node2 = newNode(id2)
		g.nodes[id2] = node2
	}
	node1.neighbour[node2] = struct{}{}
	node2.neighbour[node1] = struct{}{}
	return nil
}

func (g *Graph) String() string {
	buf := bytes.NewBufferString("")
	for _, node := range g.nodes {
		fmt.Fprintf(buf, "Node ID: %s\nNode Type: %s\n", node.id, node.nodeType)
		fmt.Fprintf(buf, "Neighbour: ")
		for neighbour := range node.neighbour {
			fmt.Fprintf(buf, "%s ", neighbour.id)
		}
		fmt.Fprintf(buf, "\n--------------------\n")
	}
	return buf.String()
}

// Stage store graph info in txt file, includeing: node count, stating & ending point,
// special node, adjacent releasion etc.
type Stage struct {
	dayCount          int
	load              int
	baseBudget        int
	baseIncome        int
	resourceWeight    [2]int
	resourceBasePrice [2]int
	resourceBaseCost  [3][2]int
	nodeCount         int
	special           map[string]NodeType
	adjacents         map[string][]string
	weightMap         map[string]map[string]int
	weatherList       []WeatherType
}

func (s *Stage) String() string {
	buf := bytes.NewBufferString("")
	fmt.Fprintln(buf, "Day Count:", s.dayCount)
	fmt.Fprintln(buf, "Load:", s.load)
	fmt.Fprintln(buf, "Budget:", s.baseBudget)
	fmt.Fprintln(buf, "Income:", s.baseBudget)
	fmt.Fprintln(buf, "Weight:", s.resourceWeight)
	fmt.Fprintln(buf, "Price:", s.resourceBasePrice)
	return buf.String()
}

func stageFromFile(filePath string) (*Stage, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var (
		line    string
		section string
		parser  func(*Stage, string) error
		hasIt   bool
	)
	stage := new(Stage)
	stage.special = map[string]NodeType{}
	stage.adjacents = map[string][]string{}
	stage.weightMap = map[string]map[string]int{}
	stage.weatherList = []WeatherType{}
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		} else if strings.HasPrefix(line, "- ") {
			section = line[2:]
			parser, hasIt = ParserMap[section]
			if !hasIt {
				log.Println("Unknown section name:", section)
			}
			continue
		} else if !hasIt {
			continue
		}
		err := parser(stage, line)
		if err != nil {
			log.Printf(
				"Section '%s' parsing error while parsing: %s\n", section, line,
			)
			return nil, err
		}
	}
	fmt.Println("Stage read from file:", filePath)
	return stage, nil
}

func (s *Stage) makeGraph() *Graph {
	g := newGraph()
	for id, neighbours := range s.adjacents {
		g.nodes[id] = newNode(id)
		for _, n := range neighbours {
			g.appendAdj(id, n)
		}
	}
	for id, kind := range s.special {
		g.nodes[id].nodeType = kind
		if kind == startingNode {
			g.starting = g.nodes[id]
		} else if kind == endingNode {
			g.ending = g.nodes[id]
		}
	}
	for id, neightbours := range s.weightMap {
		node := g.nodes[id]
		for idn, weight := range neightbours {
			node.pathWeight[idn] = weight
		}
	}
	return g
}

func (s *Stage) randWeather(pSun float64, pHigh float64, pSand float64) {
	pHigh += pSun
	pSand += pHigh
	pSun /= pSand
	pHigh /= pSand
	pSand = 1
	var (
		value   float64
		weather WeatherType
	)
	for i := 0; i < s.dayCount; i++ {
		value = rand.Float64()
		switch {
		case value < pSun:
			weather = sunny
		case value < pHigh:
			weather = highTemp
		default:
			weather = sandStorm
		}
		s.weatherList = append(s.weatherList, weather)
	}
}
