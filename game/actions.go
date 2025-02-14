package main

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

func (t *TyTris) handleInput(event event.Event) (event_handled bool) {
	if event.ID() == input.EV_KEYBOARD {
		key_event := event.(*input.KeyboardEvent)

		switch key_event.PressType {
		case input.KEY_PRESSED:
			switch key_event.Direction() {
			case vec.DIR_LEFT:
				t.movePiece(vec.DIR_LEFT)
				event_handled = true
			case vec.DIR_RIGHT:
				t.movePiece(vec.DIR_RIGHT)
				event_handled = true
			case vec.DIR_UP:
				t.dropPiece()
				event_handled = true
			case vec.DIR_DOWN:
				t.speed_up = true
				event_handled = true
			}

			switch key_event.Key {
			case input.K_z:
				t.rotatePiece(CCW)
				event_handled = true
			case input.K_c:
				t.rotatePiece(CW)
				event_handled = true
			case input.K_x:
				t.swap_held_piece()
				event_handled = true
			}

		case input.KEY_RELEASED:
			if key_event.Direction() == vec.DIR_DOWN {
				t.speed_up = false
				event_handled = true
			}
		}
	}

	return
}

func (t *TyTris) rotatePiece(dir int) {
	kick, ok := t.testRotate(dir)
	if !ok {
		return
	}

	t.current_piece.Rotate(dir)
	t.current_piece.pos.Move(kick.X, kick.Y)
	t.updateGhost()
	ui.GetLabelled[*PieceElement](t.Window(), "current piece").UpdatePiece(t.current_piece)
}

func (t *TyTris) movePiece(dir vec.Direction) {
	if !t.testMove(dir) {
		return
	}

	t.current_piece.pos.Move(dir.X, dir.Y)
	t.updateGhost()
	ui.GetLabelled[*PieceElement](t.Window(), "current piece").UpdatePiece(t.current_piece)
}

func (t *TyTris) dropPiece() {
	t.current_piece.pos = t.ghost_position
	t.lockPiece()
}

func (t *TyTris) swap_held_piece() {
	if t.held_piece.pType == t.current_piece.pType {
		return
	}

	if t.swapped_piece {
		return
	}

	if t.held_piece.pType == NO_PIECE {
		t.held_piece = Piece{pType: t.current_piece.pType}
		t.spawn_piece(t.get_next_piece())
	} else {
		held := t.held_piece
		t.held_piece = Piece{pType: t.current_piece.pType}
		t.spawn_piece(held)
	}

	t.swapped_piece = true
	t.heldArea.Updated = true

	held_element := ui.GetLabelled[*PieceElement](t.Window(), "held")
	held_element.UpdatePiece(t.held_piece)
	if t.held_piece.pType == O {
		held_element.MoveTo(vec.Coord{2, 1})
	} else {
		held_element.MoveTo(vec.Coord{1, 1})
	}

	colour := t.held_piece.Colour()
	t.held_flash.Colours = col.Pair{colour, colour}
	t.held_flash.Play()
}
