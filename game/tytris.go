package main

import (
	"math/rand"
	"slices"
	"strconv"

	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/platform/sdl"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
	//"github.com/pkg/profile"
)

var well_size vec.Dims = vec.Dims{10, 25}
var starting_gravity int = 45
var speed_up_gravity int = 8    //gravity for when the player is holding DOWN
var acceleration_time int = 300 //speed up every 300 ticks (5 seconds)
var gravity_minimum int = 5
var invalid_lines int = 3

var debug bool

func main() {
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	tyumi.InitConsole(vec.Dims{48, 27})
	tyumi.SetPlatform(sdl.New())
	tyumi.SetupRenderer("res/tytris-glyphs24x24.bmp", "res/font12x24.bmp", "TyTris")

	ui.SetDefaultElementVisuals(gfx.Visuals{
		Mode:    gfx.DRAW_GLYPH,
		Colours: col.Pair{border_colour, background_colour},
	})

	game := TyTris{}
	game.Init(vec.Dims{tyumi.FIT_CONSOLE, tyumi.FIT_CONSOLE})
	game.setup()
	defer game.highScores.WriteToDisk()

	tyumi.SetInitialMainState(&game)
	tyumi.Run()

	return
}

const (
	GAME_START int = iota
	NEW_GAME
	PLAYING
	PAUSED
	GAME_OVER
)

type TyTris struct {
	tyumi.State

	state int // one of the constants above

	//ui elements
	playField     PlayField
	matrixView    MatrixView
	upcomingArea  UpcomingPieceView
	heldArea      GridArea
	highScoreArea HighScoreView

	//animations
	held_flash gfx.FlashAnimation

	current_piece   Piece
	held_piece      Piece
	ghost_position  vec.Coord
	matrix          []Line
	upcoming_pieces []Piece
	highScores      HighScores

	info             GameInfo
	piece_spawn_tick int
	gravity          int
	speed_up         bool // will be true if player is holding down the DOWN key
	swapped_piece    bool // whether or not a swap has taken place for this piece
	spawn_next       bool // true if a new piece needs to be spawned
}

func (t *TyTris) setup() {
	t.Events().AddHandler(t.handle_event)
	t.Events().Listen(EV_CHANGESTATE, EV_HIGHSCORE)

	// do some game and ui setup
	t.matrix = make([]Line, well_size.H)
	for i := range t.matrix {
		t.matrix[i].Clear()
	}

	t.highScores.LoadFromDisk()
	t.setupUI()
}

func (t *TyTris) changeState(new_state int) {
	if new_state == t.state {
		return
	}

	switch new_state {
	case GAME_START:
		t.cleanupUI()
		ui.GetLabelled[*MainMenu](t.Window(), "menu").Activate(GAME_START)
	case GAME_OVER:
		log.Debug("GAME OVER")
		t.SetInputHandler(nil)
		t.info.high_score = t.highScores.IsHighScore(t.info.score)
		ui.GetLabelled[*MainMenu](t.Window(), "menu").Hide()
		ui.GetLabelled[*GameOverScreen](t.Window(), "gameover").Activate(t.info)
	case NEW_GAME:
		t.new_game()
		return
	case PLAYING:
		if t.state == PAUSED {
			log.Debug("UNPAUSING!")
		}

		ui.GetLabelled[*MainMenu](t.Window(), "menu").Hide()
		t.SetInputHandler(t.handleInput_playing)
	case PAUSED:
		if t.state != PLAYING {
			return
		}

		//if previous state was playing, pause game and show pause message, wait for input
		log.Debug("GAME PAUSED")
		ui.GetLabelled[*MainMenu](t.Window(), "menu").Activate(PAUSED)
		t.SetInputHandler(nil)
	default:
		log.Error("Oops, bad state change.")
		return
	}

	t.state = new_state
}

func (t *TyTris) new_game() {
	t.info = GameInfo{}
	t.shuffle_pieces()
	t.gravity = starting_gravity
	t.spawn_next = true
	t.held_piece.pType = NO_PIECE
	fireStateChangeEvent(PLAYING)
}

func (t *TyTris) Update() {
	if t.state != PLAYING {
		return
	}

	if t.spawn_next {
		//test for game over
		for i := range invalid_lines {
			if t.matrix[i].hasBlock() {
				fireStateChangeEvent(GAME_OVER)
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

	if (t.info.time-t.piece_spawn_tick)%current_gravity == 0 {
		if t.testMove(vec.DIR_DOWN) {
			t.movePiece(vec.DIR_DOWN)
		} else {
			t.lockPiece()
		}
	}

	t.info.time += 1
}

func (t *TyTris) updateScore(lines_destroyed int) {
	points := lines_destroyed * 10
	//do more score stuff here????

	t.info.score += points
	score_text := ui.GetLabelled[*ui.Textbox](t.Window(), "score")
	score_text.ChangeText(strconv.Itoa(t.info.score))
	pulse := gfx.NewPulseAnimation(score_text.DrawableArea(), 0, 12, col.Pair{col.YELLOW, col.NONE})
	pulse.OneShot = true
	pulse.Start()
	score_text.AddAnimation(&pulse)
}

func (t *TyTris) cleanupUI() {
	for i := range t.matrix {
		t.matrix[i].Clear()
	}
	t.matrixView.Updated = true

	t.upcomingArea.Reset()

	t.held_piece = Piece{pType: NO_PIECE}
	ui.GetLabelled[*PieceElement](t.Window(), "held").UpdatePiece(t.held_piece)
	t.heldArea.Updated = true

	ui.GetLabelled[*ui.Textbox](t.Window(), "score").ChangeText("0")
	ui.GetLabelled[*ui.Textbox](t.Window(), "time").ChangeText("0")
	ui.GetLabelled[*ui.Textbox](t.Window(), "speed").ChangeText("0")
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
			if !block_pos.IsInside(well_size) {
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
	t.matrixView.Updated = true

	t.current_piece.pType = NO_PIECE
	ui.GetLabelled[*PieceElement](t.Window(), "current piece").UpdatePiece(t.current_piece)

	//test for full lines
	var destroyed_lines int
	for i, line := range t.matrix {
		if line.isFull() {
			lda := NewLineDestroyAnimation(vec.Rect{vec.Coord{0, i}, vec.Dims{well_size.W, 1}})
			t.playField.AddAnimation(&lda)
			destroyed_lines += 1
		}
	}

	if destroyed_lines > 0 {
		t.info.lines_destroyed += destroyed_lines
		switch destroyed_lines {
		case 2:
			t.info.double_kills += 1
		case 3:
			t.info.triple_kills += 1
		case 4:
			t.info.quad_kills += 1
		}
		t.updateScore(destroyed_lines)
	}

	t.info.pieces_dropped += 1
	t.spawn_next = true
}

func (t *TyTris) cleanMatrix() {
	for i, line := range t.matrix {
		if line.isFull() {
			t.destroyLine(i)
		}
	}

	t.playField.Updated = true
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

	t.matrixView.Updated = true
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
	old_gravity := t.gravity
	t.gravity = util.Clamp(starting_gravity-t.info.time/acceleration_time, gravity_minimum, starting_gravity)
	if t.gravity != old_gravity {
		speedup_animation := NewSpeedUpAnimation()
		t.playField.AddAnimation(&speedup_animation)
	}

	t.piece_spawn_tick = t.info.time
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

type GameInfo struct {
	score           int
	time            int
	pieces_dropped  int
	quick_drops     int
	swaps           int
	lines_destroyed int
	double_kills    int
	triple_kills    int
	quad_kills      int

	high_score bool
}
