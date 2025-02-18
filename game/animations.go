package main

import (
	"math/rand"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/vec"
)

var LDA_Duration int = 20

var SUA_Sweep_Duration int = (WellDims.H)
var SUA_Particle_Decay int = 15

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

type SpeedUpAnimation struct {
	gfx.Animation

	particles    []SweepParticle
	particle_vis gfx.Visuals
}

func NewSpeedUpAnimation() (sa SpeedUpAnimation) {
	sa.OneShot = true
	sa.Duration = SUA_Particle_Decay + SUA_Sweep_Duration
	sa.AlwaysUpdates = true

	sa.particle_vis = gfx.NewGlyphVisuals(gfx.GLYPH_ARROW_UP, col.Pair{col.LIGHTGREY, col.DARKGREY})

	sa.Start()

	return
}

func (sa *SpeedUpAnimation) Render(canvas *gfx.Canvas) {
	area := canvas.Bounds()
	progress := float64(sa.GetTicks()) / float64(SUA_Sweep_Duration)
	y := area.H - 1 - int(progress*float64(area.H))

	//draw active particles
	for i, particle := range sa.particles {
		if particle.ticks_remaining > 0 {
			cell := canvas.GetCell(particle.pos)
			vis := sa.particle_vis
			vis.Colours = vis.Colours.Lerp(cell.Colours, SUA_Particle_Decay-particle.ticks_remaining, SUA_Particle_Decay)
			canvas.DrawVisuals(particle.pos, 1, vis)
			sa.particles[i].ticks_remaining -= 1
		}
	}

	if sa.GetTicks() > SUA_Sweep_Duration {
		return
	}

	//render sweeping line
	line := vec.Line{
		Start: vec.Coord{0, y},
		End:   vec.Coord{area.W - 1, y},
	}

	brush := gfx.NewGlyphVisuals(gfx.GLYPH_NONE, col.Pair{col.MAROON, col.DARKGREY})
	canvas.DrawLine(line, 1, brush)

	if y < area.H-2 {
		//add new particles maybe
		for i := range area.W {
			if rand.Intn(8) == 0 {
				if sa.particles == nil {
					sa.particles = make([]SweepParticle, 0, 10)
				}

				sa.particles = append(sa.particles,
					SweepParticle{
						pos:             vec.Coord{i, y},
						ticks_remaining: SUA_Particle_Decay,
					},
				)
			}
		}
	}
}

type SweepParticle struct {
	pos             vec.Coord
	ticks_remaining int
}
