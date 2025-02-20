package main

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

type GridArea struct {
	ui.ElementPrototype
}

func (ga *GridArea) Render() {
	//render checkerboard background
	for cursor := range vec.EachCoordInArea(ga.Canvas) {
		if (cursor.X+cursor.Y)%2 == 0 {
			ga.DrawColours(cursor, 0, col.Pair{col.NONE, grid_colour})
		} else {
			ga.DrawColours(cursor, 0, col.Pair{col.NONE, background_colour})
		}
	}
}

type PlayField struct {
	GridArea
}

func (pf *PlayField) Render() {
	pf.GridArea.Render()

	//draw invalid line
	invalid_brush := gfx.NewGlyphVisuals(gfx.GLYPH_LOWERCURSOR, col.Pair{invalid_line_colour, col.NONE})
	for i := range 10 {
		pf.DrawVisuals(vec.Coord{i, invalid_lines - 1}, 0, invalid_brush)
	}

}

type MatrixView struct {
	ui.ElementPrototype

	matrix *[]Line
}

func (m *MatrixView) Render() {
	// render matrix
	for y, line := range *m.matrix {
		for x, block := range line.blocks {
			pos := vec.Coord{x, y}
			if block != NO_PIECE {
				glyph := gfx.GLYPH_NONE
				if y != 0 && (*m.matrix)[y-1].blocks[x] == NO_PIECE {
					glyph = gfx.GLYPH_HALFBLOCK_UP
				}
				drawBlock(&m.Canvas, pos, glyph, pieceData[block].colour, pieceData[block].highlight_colour)
			} else {
				m.DrawNone(pos)
			}
		}
	}
}

type UpcomingPieceView struct {
	GridArea
}

func (upv *UpcomingPieceView) Init(size vec.Dims, pos vec.Coord, depth int) {
	upv.GridArea.Init(size, pos, depth)

	upv.SetupBorder("Upcoming Pieces", "")

	for range 6 {
		upcoming_piece := PieceElement{}
		upcoming_piece.Init(vec.Dims{3, 2}, vec.Coord{0, 0}, 1)
		upv.AddChild(&upcoming_piece)
	}
}

func (upv *UpcomingPieceView) UpdatePieces(pieces []Piece) {
	piece_elements := upv.GetChildren()
	for i, piece := range pieces {
		piece_elements[i].(*PieceElement).UpdatePiece(piece)
	}

	x := 1
	for i := range piece_elements {
		piece_elements[i].MoveTo(vec.Coord{x, 1})
		x += piece_elements[i].Size().W + 1
	}
}

func (upv *UpcomingPieceView) Reset() {
	piece_elements := upv.GetChildren()
	for i := range piece_elements {
		piece_elements[i].(*PieceElement).UpdatePiece(Piece{pType: NO_PIECE})
	}
}
