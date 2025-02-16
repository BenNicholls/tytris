package main

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

var LDA_Duration int = 20

type LineDestroyAnimation struct {
	gfx.AnimationChain
}

func NewLineDestroyAnimation(area vec.Rect) (lda LineDestroyAnimation) {
	lda.OneShot = true
	lda.Blocking = true
	lda.Area = area

	flash1 := gfx.NewFadeAnimation(area, 2, LDA_Duration/2, col.Pair{col.PURPLE, col.PURPLE})
	flash2 := gfx.NewFadeAnimation(area, 2, LDA_Duration/2, col.Pair{background_colour, background_colour}, col.Pair{col.PURPLE, col.PURPLE})

	lda.Add(&flash1, &flash2)
	lda.Start()

	return
}

func (lda *LineDestroyAnimation) Render(canvas *gfx.Canvas) {
	lda.AnimationChain.Render(canvas)

	x := int(float64(lda.GetTicks()) / float64(lda.GetDuration()/2) * float64(lda.Area.W/2))
	mid := lda.Area.Coord
	mid.X += lda.Area.W / 2
	left := mid.StepN(vec.DIR_LEFT, x+1)
	right := mid.StepN(vec.DIR_RIGHT, x)
	for cursor := range vec.EachCoordInArea(lda.Area) {
		if cursor.X < left.X || cursor.X > right.X {
			continue
		}

		if cursor.X >= mid.X {
			canvas.DrawGlyph(cursor, 2, gfx.GLYPH_ARROW_RIGHT)
		} else if cursor.X < mid.X {
			canvas.DrawGlyph(cursor, 2, gfx.GLYPH_ARROW_LEFT)
		}
	}
}
