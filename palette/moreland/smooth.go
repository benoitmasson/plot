// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package moreland provides color palettes for psuedocoloring scalar fields.
// The color palettes are described at http://www.kennethmoreland.com/color-advice/
// and in the following publications:
//
// "Why We Use Bad Color Maps and What You Can Do About It." Kenneth Moreland.
// In Proceedings of Human Vision and Electronic Imaging (HVEI), 2016. (To appear)
//
// "Diverging Color Maps for Scientific Visualization." Kenneth Moreland.
// In Proceedings of the 5th International Symposium on Visual Computing,
// December 2009. DOI 10.1007/978-3-642-10520-3_9.
package moreland

import (
	"fmt"
	"image/color"
)

// DivergingMSH is a smooth diverging color palette as described in
// "Diverging Color Maps for Scientific Visualization." by Kenneth Moreland,
// in Proceedings of the 5th International Symposium on Visual Computing,
// December 2009. DOI 10.1007/978-3-642-10520-3_9.
type DivergingMSH struct {
	// start and end are the beginning and ending colors
	start, end MSH

	// convergePoint is a number between 0 and
	// 1 where the colors should converge. It is 0.5 by default.
	ConvergePoint float64

	// ConvergeM is the MSH magnitude of the convergence point.
	// It is 88 by default.
	ConvergeM float64

	// Alpha represents the opacity of the returned
	// colors in the range (0,1). It is 1 by default.
	Alpha float64

	// Min and Max are the minimum and maximum values of the range of
	// scalars that can be mapped to colors using this palette.
	min, max float64
}

// NewDivergingMSH creates a new diverging color map where start and end
// are the start and end point colors in MSH space,
func NewDivergingMSH(start, end MSH) *DivergingMSH {
	return &DivergingMSH{
		start:         start,
		end:           end,
		ConvergeM:     88,
		ConvergePoint: 0.5,
		Alpha:         1,
	}
}

// At implements the palette.ColorMap interface for a DivergingMSH object.
func (p *DivergingMSH) At(scalar float64) (color.Color, error) {
	if p.min == p.max {
		return nil, fmt.Errorf("moreland: DivergingMSH color map Max == Min")
	}
	scalar = (scalar - p.min) / p.max
	o := p.interpolateMSHDiverging(scalar).lab().XYZ().linearRGB().S(p.Alpha)
	return o, o.check()
}

// SetMax implements the palette.ColorMap interface for a Luminance object.
func (p *DivergingMSH) SetMax(v float64) {
	p.max = v
}

// SetMin implements the palette.ColorMap interface for a DivergingMSH object.
func (p *DivergingMSH) SetMin(v float64) {
	p.min = v
}

// Max implements the palette.ColorMap interface for a DivergingMSH object.
func (p *DivergingMSH) Max() float64 {
	return p.max
}

// Min implements the palette.ColorMap interface for a DivergingMSH object.
func (p *DivergingMSH) Min() float64 {
	return p.min
}

// interpolateMSHDiverging performs a color interpolation through MSH space,
// where start and end are the beginning and ending colors, *HTwist are hue
// twists at the start and end, convergeM is M at the point where the
// two color scales converge, scalar is a number between 0 and 1 that the
// color should be evaluated at, and convergePoint is a number between 0 and
// 1 (typically 0.5) where the colors should converge.
func (p *DivergingMSH) interpolateMSHDiverging(scalar float64) MSH {
	startHTwist := hueTwist(p.start, p.ConvergeM)
	endHTwist := hueTwist(p.end, p.ConvergeM)
	if scalar < p.ConvergePoint {
		// interpolation factor
		interp := scalar / p.ConvergePoint
		return MSH{
			M: (p.ConvergeM-p.start.M)*interp + p.start.M,
			S: p.start.S * (1 - interp),
			H: p.start.H + startHTwist*interp,
		}
	}
	// interpolation factors
	interp1 := (scalar - 1) / (p.ConvergePoint - 1)
	interp2 := (scalar/p.ConvergePoint - 1)
	var H float64
	if scalar > p.ConvergePoint {
		H = p.end.H + endHTwist*interp1
	}
	return MSH{
		M: (p.ConvergeM-p.end.M)*interp1 + p.end.M,
		S: p.end.S * interp2,
		H: H,
	}
}

// Palette returns an object that fulfills the palette.Palette interface,
// where nColors is the number of desired colors.
func (p DivergingMSH) Palette(nColors int) Palette {
	if p.max == 0 && p.min == 0 {
		p.min = 0
		p.max = 1
	}
	delta := (p.max - p.min) / float64(nColors-1)
	v := p.min
	c := make([]color.Color, nColors)
	for i := 0; i < nColors; i++ {
		var err error
		c[i], err = p.At(v)
		if err != nil {
			panic(err)
		}
		v += delta
	}
	return Palette(c)
}

// SmoothBlueRed is a smooth diverging color palette ranging from blue to red.
func SmoothBlueRed() *DivergingMSH {
	start := MSH{
		M: 80,
		S: 1.08,
		H: -1.1,
	}
	end := MSH{
		M: 80,
		S: 1.08,
		H: 0.5,
	}
	return NewDivergingMSH(start, end)
}

// SmoothPurpleOrange is a smooth diverging color palette ranging from purple to orange.
func SmoothPurpleOrange() *DivergingMSH {
	start := MSH{
		M: 64.97539711,
		S: 0.899434815,
		H: -0.899431964,
	}
	end := MSH{
		M: 85.00850996,
		S: 0.949730284,
		H: 0.950636521,
	}
	return NewDivergingMSH(start, end)
}

// SmoothGreenPurple is a smooth diverging color palette ranging from green to purple.
func SmoothGreenPurple() *DivergingMSH {
	start := MSH{
		M: 78.04105346,
		S: 0.885011982,
		H: 2.499491379,
	}
	end := MSH{
		M: 64.97539711,
		S: 0.899434815,
		H: -0.899431964,
	}
	return NewDivergingMSH(start, end)
}

// SmoothBlueTan is a smooth diverging color palette ranging from blue to tan.
func SmoothBlueTan() *DivergingMSH {
	start := MSH{
		M: 79.94788321,
		S: 0.798754784,
		H: -1.401313221,
	}
	end := MSH{
		M: 80.07193125,
		S: 0.799798811,
		H: 1.401089787,
	}
	return NewDivergingMSH(start, end)
}

// SmoothGreenRed is a smooth diverging color palette ranging from green to red.
func SmoothGreenRed() *DivergingMSH {
	start := MSH{
		M: 78.04105346,
		S: 0.885011982,
		H: 2.499491379,
	}
	end := MSH{
		M: 76.96722122,
		S: 0.949483656,
		H: 0.499492043,
	}
	return NewDivergingMSH(start, end)
}
