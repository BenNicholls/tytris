package main

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

type PieceElement struct {
	ui.ElementPrototype

	piece Piece
	ghost bool
}

func (pe *PieceElement) UpdatePiece(p Piece) {
	if p.pType == NO_PIECE {
		pe.SetVisible(false)
		return
	} else {
		pe.SetVisible(true)
	}

	if pe.piece.Dims() != p.Dims() {
		pe.Resize(p.Dims())
		pe.Updated = true
	} else if pe.piece.pType != p.pType {
		pe.Clear()
		pe.Updated = true
	}

	if pe.piece.pos != p.pos {
		pe.MoveTo(p.pos)
	}

	pe.piece = p
}

func (pe *PieceElement) Render() {
	shape := pe.piece.Shape()
	stride := shape.stride
	for i, piece_block := range shape.shape {
		offset := vec.IndexToCoord(i, stride)
		if piece_block {
			glyph := gfx.GLYPH_NONE
			if pe.ghost {
				drawBlock(&pe.Canvas, offset, glyph, col.LIGHTGREY, col.NONE)
			} else {
				if i < stride || !shape.shape[i-stride] {
					glyph = gfx.GLYPH_HALFBLOCK_UP
				}
				drawBlock(&pe.Canvas, offset, glyph, pe.piece.Colour(), pe.piece.Highlight())
			}
		} else {
			pe.DrawNone(offset)
		}
	}
}
