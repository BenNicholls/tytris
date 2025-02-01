package main

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

type PlayField struct {
	ui.ElementPrototype

	current_piece *Piece
	ghost_pos     *vec.Coord
	matrix        *[]Line
}

func (pf *PlayField) Render() {
	pf.Clear()

	//render checkerboard background
	for cursor := range vec.EachCoord(pf.Canvas.Bounds()) {
		if (cursor.X/BlockSize+cursor.Y/BlockSize)%2 == 0 {
			pf.DrawColours(cursor, 0, col.Pair{col.NONE, col.LIGHTGREY})
		}
	}

	//render current piece and ghost
	if pf.current_piece != nil {
		stride := pf.current_piece.Shape().stride
		for i, piece_block := range pf.current_piece.Shape().shape {
			if piece_block {
				glyph := gfx.GLYPH_NONE
				if i < stride || !pf.current_piece.Shape().shape[i-stride] {
					glyph = gfx.GLYPH_HALFBLOCK_UP
				}
				offset := vec.IndexToCoord(i, stride)
				pf.drawBlock((*pf.ghost_pos).Add(offset), gfx.GLYPH_BLOCK, col.NONE, col.DARKGREY)
				pf.drawBlock(pf.current_piece.pos.Add(offset), glyph, pf.current_piece.Colour(), pf.current_piece.Highlight())
			}
		}
	}

	//render matrix
	for y, line := range *pf.matrix {
		for x, block := range line.blocks {
			if block != NO_PIECE {
				glyph := gfx.GLYPH_NONE
				if (*pf.matrix)[y-1].blocks[x] == NO_PIECE {
					glyph = gfx.GLYPH_HALFBLOCK_UP
				}
				pf.drawBlock(vec.Coord{x, y}, glyph, pieceData[block].colour, pieceData[block].highlight_colour)
			}
		}
	}
}

func (pf *PlayField) drawBlock(block_pos vec.Coord, glyph int, colour, highlight uint32) {
	pos := block_pos.Scale(BlockSize)
	pf.DrawRect(vec.Rect{pos, vec.Dims{BlockSize, BlockSize}}, 1, gfx.NewGlyphVisuals(glyph, col.Pair{highlight, colour}))
}
