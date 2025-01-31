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
				offset := vec.IndexToCoord(i, stride)
				pf.drawBlock((*pf.ghost_pos).Add(offset), gfx.GLYPH_FILL_DENSE, col.DARKGREY)
				pf.drawBlock(pf.current_piece.pos.Add(offset), gfx.GLYPH_BLOCK, pf.current_piece.Colour())
			}
		}
	}

	//render matrix
	for y, line := range *pf.matrix {
		for x, block := range line.blocks {
			if block != 0 {
				pf.drawBlock(vec.Coord{x, y}, gfx.GLYPH_BLOCK, block)
			}
		}
	}
}

func (pf *PlayField) drawBlock(block_pos vec.Coord, glyph int, colour uint32) {
	pos := block_pos.Scale(BlockSize)
	pf.DrawRect(vec.Rect{pos, vec.Dims{BlockSize, BlockSize}}, 1, gfx.NewGlyphVisuals(glyph, col.Pair{colour, col.NONE}))
}
