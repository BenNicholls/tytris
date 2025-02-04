package main

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

type PlayField struct {
	ui.ElementPrototype

	matrix *[]Line
}

func (pf *PlayField) Render() {
	//draw background grid
	drawGrid(&pf.Canvas, pf.Canvas.Bounds(), grid_colour)

	//draw invalid line
	invalid_brush := gfx.NewGlyphVisuals(gfx.GLYPH_LOWERCURSOR, col.Pair{invalid_line_colour, col.NONE})
	for i := range 10 {
		pf.DrawVisuals(vec.Coord{i, InvalidLines - 1}, 0, invalid_brush)
	}

	//render matrix
	for y, line := range *pf.matrix {
		for x, block := range line.blocks {
			if block != NO_PIECE {
				glyph := gfx.GLYPH_NONE
				if (*pf.matrix)[y-1].blocks[x] == NO_PIECE {
					glyph = gfx.GLYPH_HALFBLOCK_UP
				}
				drawBlock(&pf.Canvas, vec.Coord{x, y}, glyph, pieceData[block].colour, pieceData[block].highlight_colour)
			}
		}
	}
}
