package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var commandMap = map[string]func([]string, *Traveler, *StateRecorder) error{}

func init() {
	commandMap["repeat"] = commandRepeat
	commandMap["undo"] = commandUndo
	commandMap["redo"] = commandRedo
	commandMap["undoto"] = commandUndoUntil
	commandMap["redoto"] = commandRedoUntil
	commandMap["history"] = commandHistory
	commandMap["go"] = commandGoto
	commandMap["log"] = commandLogState
	commandMap["graph-info"] = commandGraphInfo
	commandMap["stage-info"] = commandStageInfo
	commandMap["weather"] = commandWeather
	commandMap["mine"] = commandMining
	commandMap["buy"] = commandBuy
	commandMap["stay"] = commandStay
	commandMap["random-run"] = randomRun
}

func shell(t *Traveler, states *StateRecorder) {
	command := ""
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for command != "quite" {
		fmt.Print(colorYellow, "Traveling> ", colorNone)
		scanner.Scan()
		input := scanner.Text()
		commands := strings.Split(input, ";")
		for _, command := range commands {
			command = strings.TrimSpace(command)
			if len(command) == 0 {
				continue
			}
			parts := strings.Split(command, " ")
			err := singleCommand(command, parts, t, states)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func singleCommand(command string, parts []string, t *Traveler, states *StateRecorder) error {
	handler, ok := commandMap[parts[0]]
	if !ok {
		return fmt.Errorf("Unknown command `%s`", parts[0])
	}
	fmt.Println()
	err := handler(parts[1:], t, states)
	if err != nil {
		return err
	}
	if parts[0] != "repeat" {
		states.appendCommand(command)
	}
	fmt.Println()
	return nil
}

func commandRepeat(args []string, t *Traveler, states *StateRecorder) error {
	if len(args) < 2 {
		return fmt.Errorf("Not enough argument")
	}
	times, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	for i := 0; i < times; i++ {
		err = singleCommand(strings.Join(args[1:], ""), args[1:], t, states)
		if err != nil {
			return err
		}
	}
	return nil
}

func commandGoto(args []string, t *Traveler, states *StateRecorder) error {
	if !t.ok {
		return fmt.Errorf("Traveler not in normal state")
	}
	if len(args) == 0 {
		return fmt.Errorf("Not enough argument for command")
	}
	for _, arg := range args {
		ok := t.moveTo(arg)
		if !ok {
			return fmt.Errorf("No such node with id '%s'", arg)
		}
	}
	t.checkState()
	states.appendRecord(t)
	commandLogState(args, t, states)
	return nil
}

func commandUndo(args []string, t *Traveler, states *StateRecorder) error {
	if len(args) != 0 {
		return fmt.Errorf("Too much argument")
	}
	err := states.undo(t)
	if err != nil {
		return err
	}
	commandLogState(args, t, states)
	return nil
}

func commandRedo(args []string, t *Traveler, states *StateRecorder) error {
	if len(args) != 0 {
		return fmt.Errorf("Too much argument for command 'back'")
	}
	err := states.redo(t)
	if err != nil {
		return err
	}
	commandLogState(args, t, states)
	return nil
}

func commandUndoUntil(args []string, t *Traveler, states *StateRecorder) error {
	if len(args) != 1 {
		return fmt.Errorf("Wrong number of argument")
	}
	date, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	err = states.undoToDate(t, date)
	if err != nil {
		return err
	}
	commandLogState(args, t, states)
	return nil
}

func commandRedoUntil(args []string, t *Traveler, states *StateRecorder) error {
	if len(args) != 1 {
		return fmt.Errorf("Wrong number of argument")
	}
	date, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	err = states.redoToDate(t, date)
	if err != nil {
		return err
	}
	commandLogState(args, t, states)
	return nil
}

func commandLogState(_ []string, t *Traveler, _ *StateRecorder) error {
	if t == nil {
		return fmt.Errorf("Invalide traveler")
	}
	fmt.Println(t.String())
	return nil
}

func commandHistory(args []string, t *Traveler, states *StateRecorder) error {
	for i, command := range states.commands[:states.currPos] {
		fmt.Printf("%d: %s\n", i+1, command)
	}
	return nil
}

func commandMining(args []string, t *Traveler, states *StateRecorder) error {
	if !t.ok {
		return fmt.Errorf("Traveler not in normal state")
	}
	ok := t.mining()
	if !ok {
		return fmt.Errorf("You have to go to mine to do this")
	}
	t.checkState()
	states.appendRecord(t)
	commandLogState(args, t, states)
	return nil
}

func commandBuy(args []string, t *Traveler, states *StateRecorder) error {
	if !t.ok {
		return fmt.Errorf("Traveler not in normal state")
	} else if len(args) > 2 {
		return fmt.Errorf("too much arguments")
	}
	var (
		foodAmount  int
		waterAmount int
	)
	for _, arg := range args {
		parts := strings.Split(arg, ":")
		if len(parts) != 2 {
			return fmt.Errorf("Wrong sperator usage in argument '%s'", arg)
		}
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}
		if parts[0] == "food" {
			foodAmount = value
		} else if parts[0] == "water" {
			waterAmount = value
		}
	}
	ok := t.buyResource(foodAmount, waterAmount)
	if !ok {
		return fmt.Errorf("You have to go to village to do this or buy resource for the first time at starting point")
	}
	t.checkState()
	states.appendRecord(t)
	commandLogState(args, t, states)
	return nil
}

func commandGraphInfo(_ []string, t *Traveler, _ *StateRecorder) error {
	if t == nil || t.Graph == nil {
		return fmt.Errorf("Invalid traveler")
	}
	fmt.Println(t.Graph.String())
	return nil
}

func commandStageInfo(_ []string, t *Traveler, _ *StateRecorder) error {
	if t == nil || t.Stage == nil {
		return fmt.Errorf("Invalid traveler")
	}
	fmt.Println(t.Stage.String())
	return nil
}

func commandWeather(_ []string, t *Traveler, _ *StateRecorder) error {
	if t == nil || t.Stage == nil {
		return fmt.Errorf("Invalid traveler")
	}
	for i, weather := range t.Stage.weatherList {
		fmt.Printf("%d: %s\n", i, weather)
	}
	return nil
}

func commandStay(args []string, t *Traveler, states *StateRecorder) error {
	t.stay()
	t.checkState()
	states.appendRecord(t)
	commandLogState(args, t, states)
	return nil
}
