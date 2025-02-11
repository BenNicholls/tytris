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
	heldArea     ui.ElementPrototype

	current_piece   Piece
	ghost_position  vec.Coord
	matrix          []Line
	upcoming_pieces []Piece

	last_piece_drop_tick int
	gravity              int
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
	t.spawn_piece()
}

func (t *TyTris) Update() {
	//apply gravity
	if (engine.GetTick()-t.last_piece_drop_tick)%t.gravity == 0 {
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
		if dir := key_event.Direction(); dir != vec.DIR_NONE {
			switch dir {
			case vec.DIR_LEFT:
				t.movePiece(vec.DIR_LEFT)
				event_handled = true
			case vec.DIR_RIGHT:
				t.movePiece(vec.DIR_RIGHT)
				event_handled = true
			case vec.DIR_DOWN:
				t.dropPiece()
				event_handled = true
			}
		}

		switch key_event.Key {
		case input.K_z:
			t.rotatePiece(CCW)
			event_handled = true
		case input.K_c:
			t.rotatePiece(CW)
			event_handled = true
		}
	}

	return
}

func (t *TyTris) rotatePiece(dir int) {
	if !t.testRotate(dir) {
		return
	}

	t.current_piece.Rotate(dir)
	t.updateGhost()
	ui.GetLabelledElement[*PieceElement](t.Window(), "current piece").UpdatePiece(t.current_piece)
}

func (t *TyTris) testRotate(dir int) bool {
	test_piece := t.current_piece
	test_piece.Rotate(dir)
	return t.testValidPosition(test_piece)
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
	// test leaving well
	if intersect := vec.FindIntersectionRect(piece.Bounds(), vec.Rect{vec.ZERO_COORD, WellDims}); intersect.Area() != piece.Bounds().Area() {
		return false
	}

	// test collide with matrix
	piece_shape := piece.Shape()
	for i, block := range piece_shape.shape {
		if block {
			block_pos := piece.Bounds().Coord.Add(vec.IndexToCoord(i, piece_shape.stride))
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
	piece_shape := t.current_piece.Shape()
	for i, block := range piece_shape.shape {
		if block {
			block_pos := t.current_piece.pos.Add(vec.IndexToCoord(i, piece_shape.stride))
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
	t.spawn_piece()
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
	rand.Shuffle(len(pieces), func(i, j int) {
		pieces[i], pieces[j] = pieces[j], pieces[i]
	})

	for i := range pieces {
		t.upcoming_pieces = append(t.upcoming_pieces, Piece{
			pType: pieces[i],
			pos:   vec.Coord{3, 0},
		})
	}
}

func (t *TyTris) spawn_piece() {
	t.current_piece = t.upcoming_pieces[0]

	if t.current_piece.pType == I {
		t.current_piece.pos.Y = 1
	}

	t.updateGhost()
	ui.GetLabelledElement[*PieceElement](t.Window(), "current piece").UpdatePiece(t.current_piece)

	t.upcoming_pieces = slices.Delete(t.upcoming_pieces, 0, 1)
	if len(t.upcoming_pieces) < 6 {
		t.shuffle_pieces()
	}
	t.upcomingArea.UpdatePieces(t.upcoming_pieces[0:6])

	//update gravity if necessary
	t.gravity = util.Clamp(starting_gravity-engine.GetTick()/acceleration_time, gravity_minimum, starting_gravity)
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
