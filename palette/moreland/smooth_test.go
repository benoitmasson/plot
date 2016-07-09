// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package moreland

import (
	"image/color"
	"math"
	"testing"
)

func TestInterpolateMSHDiverging(t *testing.T) {
	type test struct {
		start, end                       MSH
		convergeM, scalar, convergePoint float64
		result                           MSH
	}
	tests := []test{
		test{
			start:         MSH{M: 80, S: 1.08, H: -1.1},
			end:           MSH{M: 80, S: 1.08, H: 0.5},
			convergeM:     88,
			convergePoint: 0.5,
			scalar:        0.125,
			result:        MSH{M: 82, S: 0.81, H: -1.2402896406131008},
		},
		test{
			start:         MSH{M: 80, S: 1.08, H: -1.1},
			end:           MSH{M: 80, S: 1.08, H: 0.5},
			convergeM:     88,
			convergePoint: 0.5,
			scalar:        0.5,
			result:        MSH{M: 88, S: 0, H: 0},
		},
		test{
			start:         MSH{M: 80, S: 1.08, H: -1.1},
			end:           MSH{M: 80, S: 1.08, H: 0.5},
			convergeM:     88,
			convergePoint: 0.5,
			scalar:        0.75,
			result:        MSH{M: 84, S: 0.54, H: 0.7805792812262012},
		},
		test{
			start:         MSH{M: 80, S: 1.08, H: -1.1},
			end:           MSH{M: 80, S: 1.08, H: 0.5},
			convergeM:     88,
			convergePoint: 0.75,
			scalar:        0.7499999999999999,
			result:        MSH{M: 88, S: 1.1990408665951691e-16, H: -1.6611585624524023},
		},
		test{
			start:         MSH{M: 80, S: 1.08, H: -1.1},
			end:           MSH{M: 80, S: 1.08, H: 0.5},
			convergeM:     88,
			convergePoint: 0.75,
			scalar:        0.75,
			result:        MSH{M: 88, S: 0, H: 0},
		},
		test{
			start:         MSH{M: 80, S: 1.08, H: -1.1},
			end:           MSH{M: 80, S: 1.08, H: 0.5},
			convergeM:     88,
			convergePoint: 0.75,
			scalar:        0.7500000000000001,
			result:        MSH{M: 88, S: 2.3980817331903383e-16, H: 1.0611585624524023},
		},
	}
	for i, test := range tests {
		p := NewDivergingMSH(test.start, test.end)
		p.ConvergeM = test.convergeM
		p.ConvergePoint = test.convergePoint
		result := p.interpolateMSHDiverging(test.scalar)
		if result != test.result {
			t.Errorf("test %d: expected %v; got %v", i, test.result, result)
		}
	}
}

func TestSmoothBlueRed(t *testing.T) {
	p := SmoothBlueRed()
	wantP := []color.NRGBA{
		{59, 76, 192, 255},
		{68, 90, 204, 255},
		{77, 104, 215, 255},
		{87, 117, 225, 255},
		{98, 130, 234, 255},
		{108, 142, 241, 255},
		{119, 154, 247, 255},
		{130, 165, 251, 255},
		{141, 176, 254, 255},
		{152, 185, 255, 255},
		{163, 194, 255, 255},
		{174, 201, 253, 255},
		{184, 208, 249, 255},
		{194, 213, 244, 255},
		{204, 217, 238, 255},
		{213, 219, 230, 255},
		{221, 221, 221, 255},
		{229, 216, 209, 255},
		{236, 211, 197, 255},
		{241, 204, 185, 255},
		{245, 196, 173, 255},
		{247, 187, 160, 255},
		{247, 177, 148, 255},
		{247, 166, 135, 255},
		{244, 154, 123, 255},
		{241, 141, 111, 255},
		{236, 127, 99, 255},
		{229, 112, 88, 255},
		{222, 96, 77, 255},
		{213, 80, 66, 255},
		{203, 62, 56, 255},
		{192, 40, 47, 255},
		{180, 4, 38, 255},
	}
	c := p.Palette(33)
	if len(c.Colors()) != len(wantP) {
		t.Errorf("length doesn't match: %d != %d", len(c.Colors()), len(wantP))
	}

	// The test tolerance is the precision of a uint8 expressed as a uint32.
	const tolerance = 1.0 / 256.0 * 65535.0

	for i, c := range c.Colors() {
		w := wantP[i]
		wr, wg, wb, wa := w.RGBA()
		r, g, b, a := c.RGBA()
		if math.Abs(float64(r)-float64(wr)) > tolerance {
			t.Errorf("%d R: want %d but have %d", i, wr, r)
		}
		if math.Abs(float64(g)-float64(wg)) > tolerance {
			t.Errorf("%d G: want %d but have %d", i, wg, g)
		}
		if math.Abs(float64(b)-float64(wb)) > tolerance {
			t.Errorf("%d B: want %d but have %d", i, wb, b)
		}
		if a != wa {
			t.Errorf("%d A: want %d but have %d", i, wa, a)
		}
	}
}

func TestSmoothCoolWarm(t *testing.T) {
	type test struct {
		start [3]float64
		f     func(int) Palette
		end   [3]float64
	}
	tests := []test{
		{[3]float64{0.230, 0.299, 0.754}, SmoothBlueRed().Palette, [3]float64{0.706, 0.016, 0.150}},
		{[3]float64{0.436, 0.308, 0.631}, SmoothPurpleOrange().Palette, [3]float64{0.759, 0.334, 0.046}},
		{[3]float64{0.085, 0.532, 0.201}, SmoothGreenPurple().Palette, [3]float64{0.436, 0.308, 0.631}},
		{[3]float64{0.217, 0.525, 0.910}, SmoothBlueTan().Palette, [3]float64{0.677, 0.492, 0.093}},
		{[3]float64{0.085, 0.532, 0.201}, SmoothGreenRed().Palette, [3]float64{0.758, 0.214, 0.233}},
	}
	midPoint := [3]float64{0.865, 0.865, 0.865}

	for i, test := range tests {
		c := test.f(3).Colors()
		testRGB(t, i, "start", c[0], test.start)
		testRGB(t, i, "mid", c[1], midPoint)
		testRGB(t, i, "end", c[2], test.end)
	}
}

func fracToByte(v float64) uint8 {
	return uint8(v*255 + 0.5)
}

func testRGB(t *testing.T, i int, label string, c1 color.Color, c2 [3]float64) {
	// The test tolerance is the precision of a uint8 expressed as a uint32.
	const tolerance = 1.0 / 256.0 * 65535.0

	r, g, b, a := c1.RGBA()
	c3 := color.NRGBA{
		R: fracToByte(c2[0]),
		G: fracToByte(c2[1]),
		B: fracToByte(c2[2]),
		A: 255,
	}
	wr, wg, wb, wa := c3.RGBA()
	if math.Abs(float64(r)-float64(wr)) > tolerance {
		t.Errorf("%d %s R: want %d but have %d", i, label, wr, r)
	}
	if math.Abs(float64(g)-float64(wg)) > tolerance {
		t.Errorf("%d %s G: want %d but have %d", i, label, wg, g)
	}
	if math.Abs(float64(b)-float64(wb)) > tolerance {
		t.Errorf("%d %s B: want %d but have %d", i, label, wb, b)
	}
	if a != wa {
		t.Errorf("%d %s A: want %d but have %d", i, label, wa, a)
	}
}
