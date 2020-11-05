package main

// NodeType enum type define for representing node type: normal, village, mine
type NodeType int8

const (
	normalNode NodeType = iota
	villageNode
	mineNode
	startingNode
	endingNode
)

var specialNodeMap = map[string]NodeType{
	"v": villageNode, "m": mineNode, "s": startingNode, "e": endingNode,
}

func (n NodeType) String() string {
	return []string{"Normal", "Village", "Mine", "Start", "End"}[n]
}

// WeatherType enum for different weather
type WeatherType int8

const (
	highTemp WeatherType = iota
	sunny
	sandStorm
)

var weatherMap = map[string]WeatherType{
	"sun": sunny, "sand": sandStorm, "high": highTemp,
}

func (w WeatherType) String() string {
	return []string{"High Tempreture", "Sunny", "SandStorm"}[w]
}

// ResourceType enum for resource type
type ResourceType int8

const (
	resourceWater ResourceType = iota
	resourceFood
)

var resourceMap = map[string]ResourceType{
	"water": resourceWater, "food": resourceFood,
}
