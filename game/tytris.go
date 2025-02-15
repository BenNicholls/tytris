package main

import (
	"math/rand"
	"slices"
	"strconv"

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

	//ui elements
	playField    PlayField
	upcomingArea UpcomingPieceView
	heldArea     GridArea

	//animations
	held_flash gfx.FlashAnimation

	current_piece   Piece
	held_piece      Piece
	ghost_position  vec.Coord
	matrix          []Line
	upcoming_pieces []Piece

	score            int // HIGHER IS BETTER!!!!!
	piece_spawn_tick int
	gravity          int
	gameTick         int  // ticks since game was started
	speed_up         bool // will be true if player is holding down the DOWN key
	swapped_piece    bool // whether or not a swap has taken place for this piece
	spawn_next       bool // true if a new piece needs to be spawned
}

func (t *TyTris) setup() {
	// set the event handler for game events, engine events, etc. we also tell the event stream which events
	// to listen for. 
	t.SetEventHandler(t.handleEvent)
	t.Events().Listen(gfx.EV_ANIMATION_COMPLETE)

	// set the event handler for input events. these are keypresses, mouse movements, etc. the state object
	// sets the input event stream to listen to input events for us by default
	t.SetInputHandler(t.handleInput)

	// do some game and ui setup
	t.matrix = make([]Line, WellDims.H)
	t.setupUI()

	// setup a new game! with this done, the game will commence once tyumi's gameloop begins running
	t.new_game()
}

func (t *TyTris) new_game() {
	for i := range t.matrix {
		t.matrix[i].Clear()
	}

	t.gameTick = 0
	t.score = 0
	ui.GetLabelled[*ui.Textbox](t.Window(), "score").ChangeText("0")
	t.shuffle_pieces()
	t.gravity = starting_gravity
	t.held_piece = Piece{pType: NO_PIECE}
	t.spawn_next = true
}

func (t *TyTris) Update() {
	if t.spawn_next {
		//test for game over
		for i := range InvalidLines {
			if t.matrix[i].hasBlock() {
				log.Info("GAME OVER, YOU STINK LOSER!")
				event.Fire(event.New(engine.EV_QUIT))
				return
			}
		}

		t.cleanMatrix()

		t.spawn_piece(t.get_next_piece())
		t.spawn_next = false
	}

	if t.current_piece.pType == NO_PIECE {
		return
	}

	//apply gravity
	current_gravity := t.gravity
	if t.speed_up && t.gravity > speed_up_gravity {
		current_gravity = speed_up_gravity
	}

	if (t.gameTick-t.piece_spawn_tick)%current_gravity == 0 {
		if t.testMove(vec.DIR_DOWN) {
			t.movePiece(vec.DIR_DOWN)
		} else {
			t.lockPiece()
		}
	}

	t.gameTick += 1
}

func (t *TyTris) updateScore(lines_destroyed int) {
	points := lines_destroyed * 10
	//do more score stuff here????

	t.score += points
	ui.GetLabelled[*ui.Textbox](t.Window(), "score").ChangeText(strconv.Itoa(t.score))
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

func (t *TyTris) lockPiece() {
	//write piece in current position to lines buffers
	piece_shape := t.current_piece.GetShape()
	for i, block := range piece_shape {
		if block {
			block_pos := t.current_piece.pos.Add(vec.IndexToCoord(i, t.current_piece.Stride()))
			t.matrix[block_pos.Y].blocks[block_pos.X] = t.current_piece.pType
		}
	}

	t.current_piece.pType = NO_PIECE
	ui.GetLabelled[*PieceElement](t.Window(), "current piece").UpdatePiece(t.current_piece)

	t.playField.Updated = true

	//test for full lines
	var destroyed_lines int
	for i, line := range t.matrix {
		if line.isFull() {
			lda := NewLineDestroyAnimation(vec.Rect{vec.Coord{0, i}, vec.Dims{WellDims.W, 1}})
			t.playField.AddAnimation(&lda)
			destroyed_lines += 1
			//i dunno, do some animations or something??
		}
	}

	if destroyed_lines > 0 {
		t.updateScore(destroyed_lines)
	} else {
		t.spawn_next = true
	}
}

func (t *TyTris) cleanMatrix() {
	for i, line := range t.matrix {
		if line.isFull() {
			t.destroyLine(i)
		}
	}
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
	ui.GetLabelled[*PieceElement](t.Window(), "ghost").UpdatePiece(test_piece)
}

// adds a shuffled set of the 7 pieces to the upcoming piece list
func (t *TyTris) shuffle_pieces() {
	pieces := []PieceType{I, J, Z, S, T, O, L}
	//pieces := []PieceType{T, T, T, T, T, T, T}
	//pieces := []PieceType{O, O, O, O, O, O, O}
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
	t.swapped_piece = false

	t.updateGhost()
	ui.GetLabelled[*PieceElement](t.Window(), "current piece").UpdatePiece(t.current_piece)

	//update gravity if necessary
	t.gravity = util.Clamp(starting_gravity-t.gameTick/acceleration_time, gravity_minimum, starting_gravity)
	t.piece_spawn_tick = t.gameTick
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
