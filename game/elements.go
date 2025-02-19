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
			ga.DrawColours(cursor, 0, col.Pair{col.NONE, ga.DefaultColours().Back})
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
		pf.DrawVisuals(vec.Coord{i, InvalidLines - 1}, 0, invalid_brush)
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

type MainMenu struct {
	GridArea

	message      ui.Textbox
	instructions ui.Textbox
}

func (mm *MainMenu) Init(size vec.Dims) {
	mm.GridArea.Init(size, vec.ZERO_COORD, ui.BorderDepth)
	mm.SetLabel("menu")
	mm.SetDefaultVisuals(gfx.Visuals{
		Mode:    gfx.DRAW_GLYPH,
		Colours: col.Pair{col.WHITE, col.NONE},
	})

	mm.message.Init(vec.Dims{size.W-2, 6}, vec.Coord{1, 2}, 0, "Do you have what it takes to withstand a TORRENT of COLOURFUL SHAPES???", true)
	mm.message.SetDefaultColours(col.Pair{text_colour, gfx.COL_DEFAULT})
	mm.instructions.Init(vec.Dims{size.W, 2}, vec.Coord{0, 10}, 0, "PRESS ANY KEY TO BEGIN!", true)
	mm.instructions.SetDefaultColours(col.Pair{text_colour, gfx.COL_DEFAULT})

	instructions_pulse := gfx.NewPulseAnimation(mm.instructions.DrawableArea(), 0, 60, col.Pair{col.GREEN, col.NONE})
	instructions_pulse.Repeat = true
	instructions_pulse.Start()
	mm.instructions.AddAnimation(&instructions_pulse)

	mm.AddChildren(&mm.message, &mm.instructions)

	controls := ControlsView{}
	controls.Init(vec.Dims{size.W, 9}, vec.Coord{0, 16}, 0)
	mm.AddChild(&controls)
}

type ControlsView struct {
	ui.ElementPrototype
}

func (cv *ControlsView) Init(size vec.Dims, pos vec.Coord, depth int) {
	cv.ElementPrototype.Init(size, pos, depth)
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
