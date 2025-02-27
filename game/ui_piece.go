package main

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

type PieceElement struct {
	ui.Element

	piece Piece
	ghost bool
}

func (pe *PieceElement) Init(size vec.Dims, pos vec.Coord, depth int) {
	pe.Element.Init(size, pos, depth)
	pe.SetDefaultVisuals(gfx.Visuals{
		Mode:    gfx.DRAW_NONE,
		Colours: col.Pair{col.WHITE, col.FUSCHIA},
	})

	pe.piece.pType = NO_PIECE
}

func (pe *PieceElement) UpdatePiece(p Piece) {
	if p.pType == NO_PIECE {
		pe.Hide()
		return
	} else {
		pe.Show()
	}

	if pe.piece.pType == NO_PIECE || pe.piece.Dims() != p.Dims() {
		pe.Resize(p.Dims())
		pe.Updated = true
	} else if pe.piece.pType != p.pType || pe.piece.rotation != p.rotation {
		pe.Clear()
		pe.Updated = true
		pe.GetParent().ForceRedraw()
	}

	if pe.piece.pos != p.pos {
		pe.MoveTo(p.pos)
	}

	pe.piece = p
}

func (pe *PieceElement) Render() {
	if pe.piece.pType == NO_PIECE {
		return
	}

	shape := pe.piece.GetShape()
	stride := pe.piece.Stride()
	for i, piece_block := range shape {
		offset := vec.IndexToCoord(i, stride)
		if piece_block {
			glyph := gfx.GLYPH_NONE
			if pe.ghost {
				drawBlock(&pe.Canvas, offset, glyph, col.DARKGREY, col.NONE)
			} else {
				if i < stride || !shape[i-stride] {
					glyph = gfx.GLYPH_HALFBLOCK_UP
				}
				drawBlock(&pe.Canvas, offset, glyph, pe.piece.Colour(), pe.piece.Highlight())
			}
		}
	}
}
