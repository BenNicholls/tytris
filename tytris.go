package main

import (
	"math/rand"
	"slices"

	"github.com/bennicholls/tyumi/engine"
	"github.com/bennicholls/tyumi/engine/platform_sdl"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

var WellDims vec.Dims = vec.Dims{10, 25}
var starting_gravity int = 45
var speed_up_gravity int = 8    //gravity for when the player is holding DOWN
var acceleration_time int = 300 //speed up every 300 ticks (5 seconds)
var gravity_minimum int = 5
var InvalidLines int = 3

func main() {
	log.EnableConsoleOutput()
	engine.InitConsole(48, 27)
	engine.SetPlatform(platform_sdl.New())
	engine.SetupRenderer("res/tytris-glyphs24x24.bmp", "res/font12x24.bmp", "TyTris")

	ui.SetDefaultElementVisuals(gfx.Visuals{
		Mode:    gfx.DRAW_GLYPH,
		Colours: col.Pair{border_colour, background_colour},
	})

	input.SuppressKeyRepeats()

	game := TyTris{}
	game.Init(engine.FIT_CONSOLE, engine.FIT_CONSOLE)
	game.setup()
	engine.SetInitialMainState(&game)

	engine.Run()

	return
}

type TyTris struct {
	engine.StatePrototype

	playField    PlayField
	upcomingArea UpcomingPieceView
	heldArea     GridArea

	current_piece   Piece
	held_piece      Piece
	ghost_position  vec.Coord
	matrix          []Line
	upcoming_pieces []Piece

	last_piece_drop_tick int
	gravity              int
	speed_up             bool // will be true if player is holding down the DOWN key
	swapped_piece        bool // whether or not a swap has taken place for this piece

}

func (t *TyTris) setup() {
	t.SetInputHandler(t.handleInput)

	t.matrix = make([]Line, WellDims.H)
	for i := range t.matrix {
		t.matrix[i].Clear()
	}

	t.setupUI()
	t.shuffle_pieces()
	t.gravity = starting_gravity
	t.held_piece = Piece{pType: NO_PIECE}
	t.spawn_piece(t.get_next_piece())
}

func (t *TyTris) Update() {
	//apply gravity
	current_gravity := t.gravity
	if t.speed_up && t.gravity > speed_up_gravity {
		current_gravity = speed_up_gravity
	}

	if (engine.GetTick()-t.last_piece_drop_tick)%current_gravity == 0 {
		if t.testMove(vec.DIR_DOWN) {
			t.movePiece(vec.DIR_DOWN)
		} else {
			t.lockPiece()
		}
	}
}

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
	ui.GetLabelledElement[*PieceElement](t.Window(), "current piece").UpdatePiece(t.current_piece)
}

func (t *TyTris) testRotate(dir int) (kick vec.Coord, ok bool) {
	test_piece := t.current_piece
	test_piece.Rotate(dir)
	if t.testValidPosition(test_piece) {
		return vec.ZERO_COORD, true
	}

	//try kicks
	for _, test_kick := range test_piece.GetKicks() {
		test_piece.pos.Move(test_kick.X, test_kick.Y)
		if t.testValidPosition(test_piece) {
			return test_kick, true
		}
		test_piece.pos = test_piece.pos.Subtract(test_kick)
	}

	return 
}

func (t *TyTris) movePiece(dir vec.Direction) {
	if !t.testMove(dir) {
		return
	}

	t.current_piece.pos.Move(dir.X, dir.Y)
	t.updateGhost()
	ui.GetLabelledElement[*PieceElement](t.Window(), "current piece").UpdatePiece(t.current_piece)
}

func (t *TyTris) testMove(dir vec.Direction) bool {
	test_piece := t.current_piece
	test_piece.pos = test_piece.pos.Step(dir)
	return t.testValidPosition(test_piece)
}

func (t *TyTris) testValidPosition(piece Piece) bool {
	piece_shape := piece.GetShape()
	for i, block := range piece_shape {
		if block {
			block_pos := piece.pos.Add(vec.IndexToCoord(i, piece.Stride()))
			//not in well
			if !vec.IsInside(block_pos, vec.Rect{vec.ZERO_COORD, WellDims}) {
				return false
			}

			//collide with matrix
			if t.matrix[block_pos.Y].blocks[block_pos.X] != NO_PIECE {
				return false
			}
		}
	}

	return true
}

func (t *TyTris) dropPiece() {
	t.current_piece.pos = t.ghost_position
	t.lockPiece()
}

func (t *TyTris) lockPiece() {
	//write piece in current position to lines buffers
	piece_shape := t.current_piece.GetShape()
	for i, block := range piece_shape {
		if block {
			block_pos := t.current_piece.pos.Add(vec.IndexToCoord(i, t.current_piece.Stride()))
			t.matrix[block_pos.Y].blocks[block_pos.X] = t.current_piece.pType
		}
	}

	t.playField.Updated = true

	//test for full lines
	for i, line := range t.matrix {
		if line.isFull() {
			t.destroyLine(i)
			//hand out points
			//i dunno, do some animations or something??
			log.Info("Nice job!")
		}
	}

	//test for game over
	for i := range InvalidLines {
		if t.matrix[i].hasBlock() {
			log.Info("GAME OVER, YOU STINK LOSER!")
			event.Fire(event.New(engine.EV_QUIT))
			return
		}
	}

	t.last_piece_drop_tick = engine.GetTick()
	t.swapped_piece = false
	t.spawn_piece(t.get_next_piece())
}

func (t *TyTris) destroyLine(line_index int) {
	for i := line_index - 1; i > 0; i-- {
		if t.matrix[i].hasBlock() {
			t.matrix[i+1] = t.matrix[i]
		} else {
			t.matrix[i+1].Clear()
			break
		}
	}
}

func (t *TyTris) updateGhost() {
	test_piece := t.current_piece
	test_piece.pos = test_piece.pos.Step(vec.DIR_DOWN)
	for {
		if t.testValidPosition(test_piece) {
			test_piece.pos = test_piece.pos.Step(vec.DIR_DOWN)
		} else {
			t.ghost_position = test_piece.pos.Step(vec.DIR_UP)
			break
		}
	}

	test_piece.pos = t.ghost_position
	ui.GetLabelledElement[*PieceElement](t.Window(), "ghost").UpdatePiece(test_piece)
}

// adds a shuffled set of the 7 pieces to the upcoming piece list
func (t *TyTris) shuffle_pieces() {
	pieces := []PieceType{I, J, Z, S, T, O, L}
	//pieces := []PieceType{T, T, T, T, T, T, T}
	rand.Shuffle(len(pieces), func(i, j int) {
		pieces[i], pieces[j] = pieces[j], pieces[i]
	})

	for i := range pieces {
		t.upcoming_pieces = append(t.upcoming_pieces, Piece{
			pType: pieces[i],
		})
	}
}

func (t *TyTris) spawn_piece(piece Piece) {
	t.current_piece = piece
	t.current_piece.pos = piece.StartLocation()

	t.updateGhost()
	ui.GetLabelledElement[*PieceElement](t.Window(), "current piece").UpdatePiece(t.current_piece)

	//update gravity if necessary
	t.gravity = util.Clamp(starting_gravity-engine.GetTick()/acceleration_time, gravity_minimum, starting_gravity)
}

func (t *TyTris) get_next_piece() Piece {
	piece := t.upcoming_pieces[0]
	t.upcoming_pieces = slices.Delete(t.upcoming_pieces, 0, 1)
	if len(t.upcoming_pieces) < 6 {
		t.shuffle_pieces()
	}
	t.upcomingArea.UpdatePieces(t.upcoming_pieces[0:6])

	return piece
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

	ui.GetLabelledElement[*PieceElement](t.Window(), "held").UpdatePiece(t.held_piece)
	colour := t.held_piece.Colour()

	flash := gfx.NewFlashAnimation(t.heldArea.DrawableArea(), 0, col.Pair{colour, colour}, 15)
	flash.OneShot = true
	flash.Play()
	t.heldArea.AddAnimation(flash)
}

type Line struct {
	blocks [10]PieceType
}

func (l *Line) Clear() {
	for i := range l.blocks {
		l.blocks[i] = NO_PIECE
	}
}

func (l Line) isFull() bool {
	for _, block := range l.blocks {
		if block == NO_PIECE {
			return false
		}
	}
	return true
}

func (l Line) hasBlock() bool {
	for _, block := range l.blocks {
		if block != NO_PIECE {
			return true
		}
	}

	return false
}
