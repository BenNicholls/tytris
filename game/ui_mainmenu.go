package main

import (
	"github.com/bennicholls/tyumi"
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

type MainMenu struct {
	GridArea

	new_game_menu ui.List
	pause_menu    ui.List
	pause_message ui.Textbox

	about_text ui.Textbox
}

func (mm *MainMenu) Init(size vec.Dims) {
	mm.GridArea.Init(size, vec.ZERO_COORD, ui.BorderDepth)
	mm.SetLabel("menu")
	mm.SetDefaultVisuals(gfx.Visuals{
		Mode:    gfx.DRAW_GLYPH,
		Colours: col.Pair{col.WHITE, col.NONE},
	})

	mm.pause_message.Init(vec.Dims{size.W, 2}, vec.Coord{0, 1}, 0, "Game/nPaused", true)
	mm.pause_message.SetDefaultColours(col.Pair{text_colour, gfx.COL_DEFAULT})
	mm.pause_message.Hide()
	pulse := gfx.NewPulseAnimation(mm.pause_message.DrawableArea(), 0, 120, col.Pair{col.GREEN, col.NONE})
	pulse.Repeat = true
	pulse.Start()
	mm.pause_message.AddAnimation(&pulse)
	mm.AddChild(&mm.pause_message)

	mm.new_game_menu.Init(vec.Dims{6, 5}, vec.Coord{2, 5}, 1)
	mm.new_game_menu.ToggleHighlight()
	mm.new_game_menu.SetPadding(1)
	mm.new_game_menu.EnableBorder()
	mm.new_game_menu.AddChildren(
		ui.NewTextbox(vec.Dims{6, 1}, vec.ZERO_COORD, 1, "New Game", true),
		ui.NewTextbox(vec.Dims{6, 1}, vec.ZERO_COORD, 1, "About", true),
		ui.NewTextbox(vec.Dims{6, 1}, vec.ZERO_COORD, 1, "Quit", true),
	)
	mm.new_game_menu.OnChangeSelection = func() {
		sounds.Play("move")
	}

	mm.pause_menu.Init(vec.Dims{6, 5}, vec.Coord{2, 5}, 1)
	mm.pause_menu.ToggleHighlight()
	mm.pause_menu.SetPadding(1)
	mm.pause_menu.EnableBorder()
	mm.pause_menu.AddChildren(
		ui.NewTextbox(vec.Dims{6, 1}, vec.ZERO_COORD, 1, "Resume", true),
		ui.NewTextbox(vec.Dims{6, 1}, vec.ZERO_COORD, 1, "Give Up", true),
		ui.NewTextbox(vec.Dims{6, 1}, vec.ZERO_COORD, 1, "Quit", true),
	)
	mm.pause_menu.OnChangeSelection = func() {
		sounds.Play("move")
	}

	mm.AddChildren(&mm.new_game_menu, &mm.pause_menu)

	mm.about_text.Init(vec.Dims{size.W - 2, ui.FIT_TEXT}, vec.Coord{1, 5}, 2, "TYTRIS/n/nCreated by/nBEN NICHOLLS/n/nPlease do not sue me for this thanks", true)
	mm.about_text.SetDefaultColours(col.Pair{text_colour, background_colour})
	mm.about_text.SetupBorder("About", "[Enter]")
	mm.about_text.Hide()
	mm.AddChild(&mm.about_text)

	controls := ControlsView{}
	controls.Init(vec.Dims{size.W, 9}, vec.Coord{0, 16}, 0)
	mm.AddChild(&controls)

	mm.Activate(GAME_START)
}

func (mm *MainMenu) Activate(state int) {
	switch state {
	case GAME_START:
		mm.pause_message.Hide()
		mm.pause_menu.Hide()
		mm.new_game_menu.Show()
	case PAUSED:
		mm.pause_message.Show()
		mm.pause_menu.Show()
		mm.new_game_menu.Hide()
	default:
		return
	}

	mm.Show()
}

func (mm *MainMenu) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	if key_event.PressType == input.KEY_RELEASED {
		return
	}

	switch key_event.Key {
	case input.K_RETURN:
		if mm.about_text.IsVisible() {
			mm.about_text.Hide()
			event_handled = true
		} else if mm.new_game_menu.IsVisible() {
			switch mm.new_game_menu.GetSelectionIndex() {
			case 0: // New Game
				fireStateChangeEvent(NEW_GAME)
				event_handled = true
			case 1: // About
				mm.about_text.Show()
				event_handled = true
			case 2: //quit
				event.Fire(event.New(tyumi.EV_QUIT))
				event_handled = true
			}
		} else if mm.pause_menu.IsVisible() {
			switch mm.pause_menu.GetSelectionIndex() {
			case 0: // Resume
				fireStateChangeEvent(PLAYING)
				event_handled = true
			case 1: // Give Up
				fireStateChangeEvent(GAME_OVER)
				event_handled = true
			case 2: //quit
				event.Fire(event.New(tyumi.EV_QUIT))
				event_handled = true
			}
		}
	}

	return
}

type ControlsView struct {
	ui.Element
}

func (cv *ControlsView) Init(size vec.Dims, pos vec.Coord, depth int) {
	cv.Element.Init(size, pos, depth)
	cv.SetupBorder("Controls", "")
	cv.SetDefaultVisuals(gfx.Visuals{
		Mode:    gfx.DRAW_NONE,
		Colours: col.Pair{text_colour, background_colour},
	})
}

func (cv *ControlsView) Render() {
	cv.DrawText(vec.Coord{0, 1}, 1, "Move Piece", col.Pair{gfx.COL_DEFAULT, col.NONE}, gfx.DRAW_TEXT_LEFT)
	cv.DrawGlyph(vec.Coord{cv.Size().W - 2, 1}, 1, gfx.GLYPH_ARROW_LEFT)
	cv.DrawGlyph(vec.Coord{cv.Size().W - 1, 1}, 1, gfx.GLYPH_ARROW_RIGHT)

	cv.DrawText(vec.Coord{0, 2}, 1, "Fast Drop (hold) ", col.Pair{gfx.COL_DEFAULT, col.NONE}, gfx.DRAW_TEXT_LEFT)
	cv.DrawGlyph(vec.Coord{cv.Size().W - 1, 2}, 1, gfx.GLYPH_ARROW_DOWN)

	cv.DrawText(vec.Coord{0, 3}, 1, "Instant Drop", col.Pair{gfx.COL_DEFAULT, col.NONE}, gfx.DRAW_TEXT_LEFT)
	cv.DrawGlyph(vec.Coord{cv.Size().W - 1, 3}, 1, gfx.GLYPH_ARROW_UP)

	cv.DrawText(vec.Coord{0, 4}, 1, "Rotate CW", col.Pair{gfx.COL_DEFAULT, col.NONE}, gfx.DRAW_TEXT_LEFT)
	cv.DrawGlyph(vec.Coord{cv.Size().W - 1, 4}, 1, gfx.GLYPH_C)

	cv.DrawText(vec.Coord{0, 5}, 1, "Rotate CCW", col.Pair{gfx.COL_DEFAULT, col.NONE}, gfx.DRAW_TEXT_LEFT)
	cv.DrawGlyph(vec.Coord{cv.Size().W - 1, 5}, 1, gfx.GLYPH_Z)

	cv.DrawText(vec.Coord{0, 6}, 1, "Hold/Swap Piece", col.Pair{gfx.COL_DEFAULT, col.NONE}, gfx.DRAW_TEXT_LEFT)
	cv.DrawGlyph(vec.Coord{cv.Size().W - 1, 6}, 1, gfx.GLYPH_X)

	cv.DrawText(vec.Coord{0, 7}, 1, "Pause", col.Pair{gfx.COL_DEFAULT, col.NONE}, gfx.DRAW_TEXT_LEFT)
	cv.DrawGlyph(vec.Coord{cv.Size().W - 3, 7}, 1, gfx.GLYPH_E)
	cv.DrawGlyph(vec.Coord{cv.Size().W - 2, 7}, 1, gfx.GLYPH_S)
	cv.DrawGlyph(vec.Coord{cv.Size().W - 1, 7}, 1, gfx.GLYPH_C)
}
