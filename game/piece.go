package main

import (
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

const (
	CW  = 1
	CCW = -1
)

type PieceType int

const (
	I PieceType = iota
	J
	L
	O
	S
	Z
	T

	MAX_PIECETYPE
	NO_PIECE
)

type PieceData struct {
	default_shape    []bool
	stride           int
	colour           uint32
	highlight_colour uint32
	start_location   vec.Coord
}

var pieceData [MAX_PIECETYPE]PieceData = [MAX_PIECETYPE]PieceData{
	{ // I
		default_shape: []bool{
			false, false, false, false,
			true, true, true, true,
			false, false, false, false,
			false, false, false, false,
		},
		stride:           4,
		colour:           col.MakeOpaque(0, 163, 217),
		highlight_colour: col.MakeOpaque(102, 217, 255),
		start_location:   vec.Coord{3, 0},
	},

	{ // J
		default_shape: []bool{
			true, false, false,
			true, true, true,
			false, false, false,
		},
		stride:           3,
		colour:           col.MakeOpaque(26, 0, 102),
		highlight_colour: col.MakeOpaque(64, 0, 255),
		start_location:   vec.Coord{3, 0},
	},

	{ // L
		default_shape: []bool{
			false, false, true,
			true, true, true,
			false, false, false,
		},
		stride:           3,
		colour:           col.MakeOpaque(255, 128, 0),
		highlight_colour: col.MakeOpaque(255, 178, 102),
		start_location:   vec.Coord{3, 0},
	},

	{ // O
		default_shape: []bool{
			true, true,
			true, true,
		},
		stride:           2,
		colour:           col.MakeOpaque(217, 217, 0),
		highlight_colour: col.MakeOpaque(255, 255, 102),
		start_location:   vec.Coord{4, 0},
	},

	{ // S
		default_shape: []bool{
			false, true, true,
			true, true, false,
			false, false, false,
		},
		stride:           3,
		colour:           col.MakeOpaque(0, 102, 0),
		highlight_colour: col.MakeOpaque(0, 217, 0),
		start_location:   vec.Coord{3, 0},
	},

	{ // Z
		default_shape: []bool{
			true, true, false,
			false, true, true,
			false, false, false,
		},
		stride:           3,
		colour:           col.MakeOpaque(178, 0, 45),
		highlight_colour: col.MakeOpaque(255, 51, 102),
		start_location:   vec.Coord{3, 0},
	},

	{ // T
		default_shape: []bool{
			false, true, false,
			true, true, true,
			false, false, false,
		},
		stride:           3,
		colour:           col.MakeOpaque(140, 0, 140),
		highlight_colour: col.MakeOpaque(255, 51, 255),
		start_location:   vec.Coord{3, 0},
	},
}

type Piece struct {
	pType    PieceType
	rotation int
	pos      vec.Coord
}

func (p Piece) Colour() uint32 {
	return pieceData[p.pType].colour
}

func (p Piece) Highlight() uint32 {
	return pieceData[p.pType].highlight_colour
}

func (p Piece) StartLocation() vec.Coord {
	return pieceData[p.pType].start_location
}

func (p Piece) Stride() int {
	return pieceData[p.pType].stride
}

func (p Piece) GetKicks() []vec.Coord {
	switch p.pType {
	case O:
		return []vec.Coord{}
	case I:
		return []vec.Coord{{-1, 0}, {1, 0}, {-2, 0}, {2, 0}}
	default:
		return []vec.Coord{{-1, 0}, {1, 0}}
	}
}

func (p Piece) GetShape() []bool {
	if p.rotation == 0 {
		return pieceData[p.pType].default_shape
	}

	def_shape := pieceData[p.pType].default_shape
	rot_shape := make([]bool, len(def_shape))
	stride := p.Stride()

	for i, block := range def_shape {
		if block {
			pos := vec.IndexToCoord(i, stride)
			var rot_pos vec.Coord
			switch p.rotation {
			case 1:
				rot_pos = vec.Coord{-pos.Y + stride - 1, pos.X}
			case 2:
				rot_pos = vec.Coord{-pos.X + stride - 1, -pos.Y + stride - 1}
			case 3:
				rot_pos = vec.Coord{pos.Y, -pos.X + stride - 1}
			}
			rot_shape[rot_pos.ToIndex(stride)] = true
		}
	}

	return rot_shape
}

func (p Piece) Dims() vec.Dims {
	return vec.Dims{p.Stride(), p.Stride()}
}

func (p *Piece) Rotate(dir int) {
	if p.pType == O {
		return
	}

	p.rotation = util.CycleClamp(p.rotation+dir, 0, 3)
}
