package main

import (
	"slices"

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

var pieceData [MAX_PIECETYPE]PieceData = [MAX_PIECETYPE]PieceData{
	{ // I
		default_shape: PieceShape{
			shape:  []bool{true, true, true, true},
			stride: 4,
		},
		colour:           col.MakeOpaque(0, 163, 217),
		highlight_colour: col.MakeOpaque(102, 217, 255),
		max_rotations:    1,
	},

	{ // J
		default_shape: PieceShape{
			shape:  []bool{true, false, false, true, true, true},
			stride: 3,
		},
		colour:           col.MakeOpaque(26, 0, 102),
		highlight_colour: col.MakeOpaque(64, 0, 255),
		max_rotations:    3,
	},

	{ // L
		default_shape: PieceShape{
			shape:  []bool{false, false, true, true, true, true},
			stride: 3,
		},
		colour:           col.MakeOpaque(255, 128, 0),
		highlight_colour: col.MakeOpaque(255, 178, 102),
		max_rotations:    3,
	},

	{ // O
		default_shape: PieceShape{
			shape:  []bool{true, true, true, true},
			stride: 2,
		},
		colour:           col.MakeOpaque(217, 217, 0),
		highlight_colour: col.MakeOpaque(255, 255, 102),
		max_rotations:    0,
	},

	{ // S
		default_shape: PieceShape{
			shape:  []bool{false, true, true, true, true, false},
			stride: 3,
		},
		colour:           col.MakeOpaque(0, 102, 0),
		highlight_colour: col.MakeOpaque(0, 217, 0),
		max_rotations:    1,
	},

	{ // Z
		default_shape: PieceShape{
			shape:  []bool{true, true, false, false, true, true},
			stride: 3,
		},
		colour:           col.MakeOpaque(178, 0, 45),
		highlight_colour: col.MakeOpaque(255, 51, 102),
		max_rotations:    1,
	},

	{ // T
		default_shape: PieceShape{
			shape:  []bool{false, true, false, true, true, true},
			stride: 3,
		},
		colour:           col.MakeOpaque(140, 0, 140),
		highlight_colour: col.MakeOpaque(255, 51, 255),
		max_rotations:    3,
	},
}

type PieceData struct {
	default_shape    PieceShape
	colour           uint32
	highlight_colour uint32
	max_rotations    int
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

func (p Piece) Shape() PieceShape {
	if p.rotation == 0 {
		return pieceData[p.pType].default_shape
	}

	def_shape := pieceData[p.pType].default_shape
	rot_shape := PieceShape{}
	rot_shape.shape = make([]bool, len(def_shape.shape))

	switch p.rotation {
	case 1:
		if p.pType == I {
			if def_shape.stride == 1 {
				rot_shape.stride = 4
			} else {
				rot_shape.stride = 1
			}
			rot_shape.shape = def_shape.shape
		} else {
			if def_shape.stride == 2 {
				rot_shape.stride = 3
				rot_shape.shape[0] = def_shape.shape[4]
				rot_shape.shape[1] = def_shape.shape[2]
				rot_shape.shape[2] = def_shape.shape[0]
				rot_shape.shape[3] = def_shape.shape[5]
				rot_shape.shape[4] = def_shape.shape[3]
				rot_shape.shape[5] = def_shape.shape[1]
			} else {
				rot_shape.stride = 2
				rot_shape.shape[0] = def_shape.shape[3]
				rot_shape.shape[1] = def_shape.shape[0]
				rot_shape.shape[2] = def_shape.shape[4]
				rot_shape.shape[3] = def_shape.shape[1]
				rot_shape.shape[4] = def_shape.shape[5]
				rot_shape.shape[5] = def_shape.shape[2]
			}
		}
	case 2: //upside-down
		for i, val := range slices.Backward[[]bool](def_shape.shape) {
			rot_shape.shape[len(def_shape.shape)-1-i] = val
		}
		rot_shape.stride = def_shape.stride
	case 3:
		if def_shape.stride == 2 {
			rot_shape.stride = 3
			rot_shape.shape[0] = def_shape.shape[1]
			rot_shape.shape[1] = def_shape.shape[3]
			rot_shape.shape[2] = def_shape.shape[5]
			rot_shape.shape[3] = def_shape.shape[0]
			rot_shape.shape[4] = def_shape.shape[2]
			rot_shape.shape[5] = def_shape.shape[4]
		} else {
			rot_shape.stride = 2
			rot_shape.shape[0] = def_shape.shape[2]
			rot_shape.shape[1] = def_shape.shape[5]
			rot_shape.shape[2] = def_shape.shape[1]
			rot_shape.shape[3] = def_shape.shape[4]
			rot_shape.shape[4] = def_shape.shape[0]
			rot_shape.shape[5] = def_shape.shape[3]
		}
	}

	return rot_shape
}

func (p Piece) Bounds() vec.Rect {
	return vec.Rect{p.pos, p.Dims()}
}

func (p Piece) Dims() vec.Dims {
	return p.Shape().Dims()
}

func (p *Piece) Rotate(dir int) {
	p.rotation = util.CycleClamp(p.rotation+dir, 0, pieceData[p.pType].max_rotations)
}

type PieceShape struct {
	shape  []bool
	stride int
}

func (ps PieceShape) Dims() vec.Dims {
	return vec.Dims{ps.stride, len(ps.shape) / ps.stride}
}
