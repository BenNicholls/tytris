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
)

var pieceData [MAX_PIECETYPE]PieceData = [MAX_PIECETYPE]PieceData{
	{ // I
		default_shape: PieceShape{
			shape:  []bool{true, true, true, true},
			stride: 4,
		},
		colour:        col.CYAN,
		max_rotations: 1,
	},

	{ // J
		default_shape: PieceShape{
			shape:  []bool{true, false, false, true, true, true},
			stride: 3,
		},
		colour:        col.NAVY,
		max_rotations: 3,
	},

	{ // L
		default_shape: PieceShape{
			shape:  []bool{false, false, true, true, true, true},
			stride: 3,
		},
		colour:        col.ORANGE,
		max_rotations: 3,
	},

	{ // O
		default_shape: PieceShape{
			shape:  []bool{true, true, true, true},
			stride: 2,
		},
		colour:        col.YELLOW,
		max_rotations: 0,
	},

	{ // S
		default_shape: PieceShape{
			shape:  []bool{false, true, true, true, true, false},
			stride: 3,
		},
		colour:        col.GREEN,
		max_rotations: 1,
	},

	{ // Z
		default_shape: PieceShape{
			shape:  []bool{true, true, false, false, true, true},
			stride: 3,
		},
		colour:        col.RED,
		max_rotations: 1,
	},

	{ // T
		default_shape: PieceShape{
			shape:  []bool{false, true, false, true, true, true},
			stride: 3,
		},
		colour:        col.FUSCHIA,
		max_rotations: 3,
	},
}

type PieceData struct {
	default_shape PieceShape
	colour        uint32
	max_rotations int
}

type Piece struct {
	pType    PieceType
	rotation int
	pos      vec.Coord
	ghost    bool
}

func (p Piece) Colour() uint32 {
	if p.ghost {
		return col.GREY
	} else {
		return pieceData[p.pType].colour
	}
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
			rot_shape.shape[len(def_shape.shape) - 1 - i] = val
		}
		rot_shape.stride = pieceData[p.pType].default_shape.stride	
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
	p.rotation = util.CycleClamp(p.rotation + dir, 0, pieceData[p.pType].max_rotations)
}

type PieceShape struct {
	shape  []bool
	stride int
}

func (ps PieceShape) Dims() vec.Dims {
	return vec.Dims{ps.stride, len(ps.shape) / ps.stride}
}
