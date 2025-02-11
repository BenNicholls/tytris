package main

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/vec"
)

// colours!
var background_colour uint32 = col.MakeOpaque(26, 20, 13)
var border_colour uint32 = col.MakeOpaque(77, 61, 38)
var grid_colour uint32 = col.MakeOpaque(26, 26, 26)
var invalid_line_colour uint32 = col.MakeOpaque(51, 51, 51)

func (t *TyTris) setupUI() {
	//define a custom border style (derived from one of the provided borderstyles) and set it as the default for all
	//borders
	tytris_border := ui.BorderStyles["Thin"]
	tytris_border.Colours = col.Pair{border_colour, background_colour}
	ui.SetDefaultBorderStyle(tytris_border)
	
	logoImage := ui.Image{}
	logoImage.Init(12, 6, vec.Coord{3,1}, 0, "res/logo.xp")
	t.Window().AddChild(&logoImage)

	//initialize the playfield, where the blocks fall and the matrix is drawn.
	t.playField.Init(WellDims.W, WellDims.H, vec.Coord{19, 1}, 0)
	t.playField.EnableBorder()
	t.playField.matrix = &t.matrix

	current_piece := PieceElement{}
	current_piece.Init(3, 2, vec.Coord{0, 0}, 2)
	current_piece.SetLabel("current piece")
	t.playField.AddChild(&current_piece)

	ghost_piece := PieceElement{
		ghost: true,
	}
	ghost_piece.Init(3, 2, vec.Coord{0, 0}, 1)
	ghost_piece.SetLabel("ghost")
	t.playField.AddChild(&ghost_piece)

	t.Window().AddChild(&t.playField)

	t.upcomingArea.Init(18, 4, vec.Coord{30, 3}, 0)
	t.upcomingArea.SetupBorder("Upcoming Pieces", "")
	for range 6 {
		upcoming_piece := PieceElement{}
		upcoming_piece.Init(3, 2, vec.Coord{0,0}, 1)
		t.upcomingArea.AddChild(&upcoming_piece)
	}

	t.heldArea.Init(6, 4, vec.Coord{30, 8}, 0)
	t.heldArea.SetupBorder("", "held piece")
	
	held_piece := PieceElement{}
	held_piece.Init(3,2, vec.Coord{1,1}, 1)
	held_piece.SetLabel("held")
	t.heldArea.AddChild(&held_piece)
	
	t.Window().AddChildren(&t.upcomingArea, &t.heldArea)
}

func drawBlock(canvas *gfx.Canvas, block_pos vec.Coord, glyph int, colour, highlight uint32) {
	canvas.DrawVisuals(block_pos, 1, gfx.NewGlyphVisuals(glyph, col.Pair{highlight, colour}))
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