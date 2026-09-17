package content

type FontInfo struct {
	BaseFont Name
	Subtype  Name
}

type PathSegment struct {
	Operator string
	Points   []Point
}

func BasicStyle(state GraphicsState) PaintStyle {
	stroke, fill := state.Stroke, state.Fill
	stroke.Components = append([]float64{}, stroke.Components...)
	fill.Components = append([]float64{}, fill.Components...)
	return PaintStyle{
		Stroke:      stroke,
		Fill:        fill,
		LineWidth:   state.LineWidth,
		LineCap:     state.LineCap,
		LineJoin:    state.LineJoin,
		MiterLimit:  state.MiterLimit,
		Dash:        append([]float64{}, state.Dash...),
		DashPhase:   state.DashPhase,
		StrokeAlpha: state.StrokeAlpha,
		FillAlpha:   state.FillAlpha,
		BlendMode:   state.BlendMode,
		Clipped:     len(state.Clip) > 0,
		Complete:    state.Complete,
	}
}
