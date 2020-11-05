package main

import (
	"fmt"
	"strings"
)

// StateRecorder is a history state for traveling process
type StateRecorder struct {
	currPos  int // current position in record stack
	commands []string
	states   []Traveler
}

func newRecorder(t *Traveler) *StateRecorder {
	return &StateRecorder{0, []string{}, []Traveler{*t}}
}

func (r *StateRecorder) appendRecord(t *Traveler) {
	r.currPos++
	if r.currPos < len(r.states) {
		r.states[r.currPos] = *t
	} else {
		r.states = append(r.states, *t)
	}
}

func (r *StateRecorder) appendCommand(command string) {
	skipPrefix := []string{"redo", "undo", "history", "log"}
	for _, prefix := range skipPrefix {
		if strings.HasPrefix(command, prefix) {
			return
		}
	}
	if r.currPos > 0 && r.currPos-1 < len(r.commands) {
		r.commands[r.currPos-1] = command
	} else {
		r.commands = append(r.commands, command)
	}
}

func (r *StateRecorder) undo(t *Traveler) error {
	if r.currPos == 0 {
		return fmt.Errorf("No older states")
	}
	r.currPos--
	r.readState(t)
	return nil
}

func (r *StateRecorder) redo(t *Traveler) error {
	if r.currPos == len(r.states)-1 {
		return fmt.Errorf("No newer states")
	}
	r.currPos++
	r.readState(t)
	return nil
}

func (r *StateRecorder) undoToDate(t *Traveler, date int) error {
	if r.currPos == 0 {
		return fmt.Errorf("No older states")
	}
	pos := r.currPos
	for pos > -1 && r.states[pos].date != date {
		pos--
	}
	if pos >= 0 {
		r.currPos = pos
	} else {
		return fmt.Errorf("Invalid date")
	}
	r.readState(t)
	return nil
}

func (r *StateRecorder) redoToDate(t *Traveler, date int) error {
	if r.currPos == len(r.states)-1 {
		return fmt.Errorf("No newer state")
	}
	pos := r.currPos
	for pos > -1 && r.states[pos].date != date {
		pos++
	}
	if pos >= 0 {
		r.currPos = pos
	} else {
		return fmt.Errorf("Invalid date")
	}
	r.readState(t)
	return nil
}

func (r *StateRecorder) readState(t *Traveler) {
	old := r.states[r.currPos]
	t.date = old.date
	t.position = old.position
	t.loadSpace = old.loadSpace
	t.money = old.money
	t.water = old.water
	t.food = old.food
	t.ok = old.ok
	t.firstBuy = old.firstBuy
}
