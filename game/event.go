package main

import (
	"github.com/bennicholls/tyumi/event"
)

var EV_CHANGESTATE int = event.Register("State Change")
var EV_HIGHSCORE int = event.Register("High Score Recorded!")

type stateChangeEvent struct {
	event.EventPrototype

	new_state int
}

func fireStateChangeEvent(new_state int) {
	sce := stateChangeEvent{
		EventPrototype: *event.New(EV_CHANGESTATE),
		new_state:      new_state,
	}

	event.Fire(&sce)
}

type highScoreEvent struct {
	event.EventPrototype

	name  string
	score int
}

func fireHighScoreEvent(name string, score int) {
	hse := highScoreEvent{
		EventPrototype: *event.New(EV_HIGHSCORE),
		name:           name,
		score:          score,
	}

	event.Fire(&hse)
}

func (t *TyTris) handle_event(event event.Event) (event_handled bool) {
	switch event.ID() {
	case EV_CHANGESTATE:
		e := event.(*stateChangeEvent)
		t.changeState(e.new_state)
	case EV_HIGHSCORE:
		e := event.(*highScoreEvent)
		t.highScores.AddEntry(HighScoreEntry{Name: e.name, Score: e.score})
		t.highScoreArea.UpdateScores(t.highScores)
	}

	return
}
