package main

import "github.com/bennicholls/tyumi/event"

var EV_CHANGESTATE int = event.Register("State Change")

type stateChangeEvent struct {
	event.EventPrototype

	new_state int
}

func fireStateChangeEvent(new_state int) {
	sce := stateChangeEvent{
		EventPrototype: *event.New(EV_CHANGESTATE),
		new_state: new_state,
	}
	
	event.Fire(&sce)
}

func (t *TyTris) handle_event(event event.Event) (event_handled bool) {
	switch event.ID() {
	case EV_CHANGESTATE:
		e := event.(*stateChangeEvent)
		t.changeState(e.new_state)
	}
	
	return
}