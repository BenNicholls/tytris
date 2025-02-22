package main

import (
	"fmt"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/vec"
)

type GameOverScreen struct {
	ui.Element

	statsBox   ui.Textbox
	message    ui.Textbox
	name_input ui.InputBox

	info GameInfo
}

func (gos *GameOverScreen) Init(size vec.Dims, pos vec.Coord, depth int) {
	gos.Element.Init(size, pos, depth)
	gos.EnableBorder()

	gameOverImage := ui.Image{}
	gameOverImage.Init(vec.Coord{2, 1}, 0, "res/gameover.xp")
	gos.AddChild(&gameOverImage)

	gos.message.Init(vec.Dims{size.W - 11, 3}, vec.Coord{0, 9}, 0, "", true)
	gos.name_input.Init(vec.Dims{3, 1}, vec.Coord{10, 11}, 1, 5)
	gos.name_input.SetDefaultColours(col.Pair{background_colour, border_colour})
	gos.AddChildren(&gos.message, &gos.name_input)

	gos.statsBox.Init(vec.Dims{10, size.H}, vec.Coord{size.W - 10, 0}, ui.BorderDepth, "", false)
	gos.statsBox.SetupBorder("S T A T S", "")
	gos.statsBox.SetDefaultColours(col.Pair{text_colour, background_colour})
	gos.AddChild(&gos.statsBox)

	gos.SetLabel("gameover")
	gos.Hide()
}

func (gos *GameOverScreen) Activate(info GameInfo) {
	flash := gfx.NewFlashAnimation(gos.Canvas.Bounds(), ui.BorderDepth+1, col.Pair{col.FUSCHIA, col.FUSCHIA}, 30)
	flash.OneShot = true
	flash.Blocking = true
	flash.Start()
	gos.AddAnimation(&flash)

	stats := fmt.Sprintf(`
		Score         %6d/n
		Total Time    %6d/n
		Pieces        %6d/n
		Quick Drops   %6d/n
		Swaps         %6d/n/n
		Lines Cleared %6d/n
		Double Kills  %6d/n
		Triple Kills  %6d/n
		QUAD Kills    %6d/n`, info.score, info.time/60, info.pieces_dropped, info.quick_drops, info.swaps, info.lines_destroyed, info.double_kills, info.triple_kills, info.quad_kills)

	gos.statsBox.ChangeText(stats)

	if info.high_score {
		gos.message.ChangeText("Huzzah, you got a highscore! Enter your name and be remembered for eternity!")
		gos.name_input.Show()
	} else {
		gos.message.ChangeText("The endless torrent of menacing colours and shapes has defeated you. The poets will sing of your demise.")
		gos.name_input.Hide()
	}

	gos.info = info
	gos.Show()
}

func (gos *GameOverScreen) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	if key_event.PressType == input.KEY_RELEASED {
		return
	}

	if gos.name_input.IsVisible() {
		if key_event.Key == input.K_RETURN {
			fireHighScoreEvent(gos.name_input.InputtedText(), gos.info.score)
			gos.Hide()
			fireStateChangeEvent(GAME_START)
			event_handled = true
		}
	} else { // no high score being entered, back to menu on any keypress
		gos.Hide()
		fireStateChangeEvent(GAME_START)
		event_handled = true
	}
	return
}
