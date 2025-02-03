package main

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

var BlockSize int = 1

// colours!
var background_colour uint32 = col.MakeOpaque(26, 20, 13)
var border_colour uint32 = col.MakeOpaque(77, 61, 38)
var grid_colour uint32 = col.MakeOpaque(26, 26, 26)
var invalid_line_colour uint32 = col.MakeOpaque(51, 51, 51)

func (t *TyTris) setupUI() {
	//define a custon border style (derived from one of the provided borderstyles) and set it as the default for all
	//borders
	tytris_border := ui.BorderStyles["Thick"]
	tytris_border.Colours = col.Pair{border_colour, background_colour}
	ui.SetDefaultBorderStyle(tytris_border)

	t.Window().SetDefaultColours(col.Pair{border_colour, background_colour})

	//initialize the playfield, where the blocks fall and the matrix is drawn.
	t.playField.Init(WellDims.W*BlockSize, WellDims.H*BlockSize, vec.Coord{19, 1}, 0)
	t.playField.SetDefaultColours(col.Pair{col.LIME, background_colour})
	t.playField.EnableBorder()
	t.playField.matrix = &t.matrix
	t.playField.ghost_pos = &t.ghost_position
	t.Window().AddChild(&t.playField)
}

func drawGrid(canvas *gfx.Canvas, area vec.Rect, grid_colour uint32) {
	//render checkerboard background
	for cursor := range vec.EachCoord(area) {
		if (cursor.X/BlockSize+cursor.Y/BlockSize)%2 == 0 {
			canvas.DrawColours(cursor, 0, col.Pair{col.NONE, grid_colour})
		}
	}
}
