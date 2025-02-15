package main

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

var LDA_Duration int = 18

type PlayField struct {
	GridArea

	matrix *[]Line
}

func (pf *PlayField) Render() {
	pf.GridArea.Render()

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
				if y != 0 && (*pf.matrix)[y-1].blocks[x] == NO_PIECE {
					glyph = gfx.GLYPH_HALFBLOCK_UP
				}
				drawBlock(&pf.Canvas, vec.Coord{x, y}, glyph, pieceData[block].colour, pieceData[block].highlight_colour)
			}
		}
	}
}

type LineDestroyAnimation struct {
	gfx.AnimationChain
}

func NewLineDestroyAnimation(area vec.Rect) (lda LineDestroyAnimation) {
	lda.OneShot = true
	lda.Label = "line destroy"
	lda.Area = area
	lda.Start()

	flash1 := gfx.NewFadeAnimation(area, 1, col.Pair{col.PURPLE, col.PURPLE}, LDA_Duration/2)
	flash2 := gfx.NewFadeAnimation(area, 1, col.Pair{background_colour, background_colour}, LDA_Duration/2)

	lda.Add(&flash1, &flash2)

	return
}
