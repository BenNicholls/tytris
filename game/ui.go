package main

import (
	"strconv"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

// colours!
var background_colour uint32 = col.MakeOpaque(26, 20, 13)
var border_colour uint32 = col.MakeOpaque(128, 102, 64)
var text_colour uint32 = col.MakeOpaque(222, 211, 195)
var grid_colour uint32 = col.MakeOpaque(26, 26, 26)
var invalid_line_colour uint32 = col.MakeOpaque(51, 51, 51)

func (t *TyTris) setupUI() {
	//define a custom border style (derived from one of the provided borderstyles) and set it as the default for all
	//borders
	tytris_border := ui.BorderStyles["Thin"]
	tytris_border.Colours = col.Pair{border_colour, background_colour}
	ui.SetDefaultBorderStyle(tytris_border)

	//logo and subtitle
	logoImage := ui.Image{}
	logoImage.Init(vec.Coord{3, 2}, 0, "res/logo.xp")
	subtitle := ui.NewTextbox(vec.Dims{10, ui.FIT_TEXT}, vec.Coord{4, 9}, 0, "The Fun Game That No One Stole At All", true)
	subtitle.SetDefaultColours(col.Pair{text_colour, gfx.COL_DEFAULT})
	t.Window().AddChildren(&logoImage, subtitle)

	//initialize the playfield, where the blocks fall and the matrix is drawn.
	t.playField.Init(well_size, vec.Coord{19, 1}, 0)
	t.playField.EnableBorder()

	t.matrixView.Init(well_size, vec.ZERO_COORD, 2)
	t.matrixView.matrix = &t.matrix

	t.playField.AddChild(&t.matrixView)

	current_piece := PieceElement{}
	current_piece.Init(vec.Dims{3, 2}, vec.Coord{0, 0}, 2)
	current_piece.SetLabel("current piece")
	t.playField.AddChild(&current_piece)

	ghost_piece := PieceElement{
		ghost: true,
	}
	ghost_piece.Init(vec.Dims{3, 2}, vec.Coord{0, 0}, 1)
	ghost_piece.SetLabel("ghost")
	t.playField.AddChild(&ghost_piece)

	//main menu. this will be a child of the playarea, blocking the view of the matrix and everything else when
	//visibility is toggled on. we'll also use this as the pause menu, with a change in some text
	mainMenu := MainMenu{}
	mainMenu.Init(t.playField.DrawableArea().Dims)
	t.playField.AddChild(&mainMenu)

	t.Window().AddChild(&t.playField)

	// upcoming pieces view, where up to the next 6 pieces are displayed in order from left to right
	t.upcomingArea.Init(vec.Dims{18, 4}, vec.Coord{30, 3}, 0)

	// held piece area, where we show which piece the player is holding (if there is one)
	t.heldArea.Init(vec.Dims{6, 4}, vec.Coord{30, 8}, 0)
	t.heldArea.SetupBorder("", "held piece")
	t.held_flash = gfx.NewFlashAnimation(t.heldArea.DrawableArea(), 1, col.Pair{col.NONE, col.NONE}, 15)
	t.heldArea.AddAnimation(&t.held_flash)

	held_piece := PieceElement{}
	held_piece.Init(vec.Dims{3, 2}, vec.Coord{1, 1}, 1)
	held_piece.SetLabel("held")
	t.heldArea.AddChild(&held_piece)

	t.Window().AddChildren(&t.upcomingArea, &t.heldArea)

	// info area, where we show the player's score, as well as the timer and current game speed
	infoArea := ui.Element{}
	infoArea.Init(vec.Dims{14, 6}, vec.Coord{32, 19}, 1)
	infoArea.EnableBorder()

	scoreLabel := ui.NewTextbox(vec.Dims{14, 1}, vec.Coord{0, 0}, 1, "S C O R E", true)
	scoreLabel.SetDefaultColours(col.Pair{text_colour, col.NONE})
	score := ui.NewTextbox(vec.Dims{14, 1}, vec.Coord{0, 1}, 1, "0", true)
	score.SetLabel("score")
	score.SetDefaultColours(col.Pair{text_colour, col.NONE})
	infoArea.AddChildren(scoreLabel, score)

	timeLabel := ui.NewTextbox(vec.Dims{14, 1}, vec.Coord{0, 2}, 1, "T I M E R", true)
	timeLabel.SetDefaultColours(col.Pair{text_colour, border_colour})
	time := ui.NewTextbox(vec.Dims{14, 1}, vec.Coord{0, 3}, 1, "0", true)
	time.SetLabel("time")
	time.SetDefaultColours(col.Pair{text_colour, border_colour})
	infoArea.AddChildren(timeLabel, time)

	speedLabel := ui.NewTextbox(vec.Dims{14, 1}, vec.Coord{0, 4}, 1, "S P E E D", true)
	speedLabel.SetDefaultColours(col.Pair{text_colour, col.NONE})
	speed := ui.NewTextbox(vec.Dims{14, 1}, vec.Coord{0, 5}, 1, "0", true)
	speed.SetLabel("speed")
	speed.SetDefaultColours(col.Pair{text_colour, col.NONE})
	infoArea.AddChildren(speedLabel, speed)

	t.Window().AddChild(&infoArea)

	// highscore area
	t.highScoreArea.Init(vec.Dims{12, 13}, vec.Coord{3, 13}, 1)
	t.highScoreArea.EnableBorder()
	t.highScoreArea.UpdateScores(t.highScores)
	t.Window().AddChild(&t.highScoreArea)

	gameover := GameOverScreen{}
	gameover.Init(t.Window().DrawableArea().Dims.Shrink(14, 14), vec.Coord{7, 7}, 10)
	t.Window().AddChild(&gameover)
}

func drawBlock(canvas *gfx.Canvas, block_pos vec.Coord, glyph gfx.Glyph, colour, highlight uint32) {
	canvas.DrawVisuals(block_pos, 1, gfx.NewGlyphVisuals(glyph, col.Pair{highlight, colour}))
}

func (t *TyTris) UpdateUI() {
	if t.state != PLAYING {
		return
	}

	timer := ui.GetLabelled[*ui.Textbox](t.Window(), "time")
	timer.ChangeText(strconv.Itoa(t.info.time / 60))

	speed := ui.GetLabelled[*ui.Textbox](t.Window(), "speed")
	if t.gravity != gravity_minimum {
		speed.ChangeText(strconv.Itoa(starting_gravity - t.gravity))
	} else {
		speed.ChangeText("MAXIMUM SPEED!!")
	}
}
