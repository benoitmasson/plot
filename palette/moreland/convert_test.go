// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package moreland

import (
	"image/color"
	"math"
	"testing"
)

func TestLinearToS(t *testing.T) {
	if linearToS(0.015299702) != 0.1298716701086684 {
		t.Errorf("linearToS(0.015299702) should equal 0.1298716701086684 but equals %g", linearToS(0.015299702))
	}
}

func TestLinearRGB(t *testing.T) {
	xyz := cieXYZ{X: 0.128392403, Y: 0.128221351, Z: 0.408477452}
	rgb := xyz.linearRGB()
	wantR, wantG, wantB := 0.015299702837399953, 0.1330700251971, 0.4127549680071
	if rgb.R != wantR {
		t.Errorf("R should equal %g but equals %g", wantR, rgb.R)
	}
	if rgb.G != wantG {
		t.Errorf("G should equal %g but equals %g", wantG, rgb.G)
	}
	if rgb.B != wantB {
		t.Errorf("B should equal %g but equals %g", wantB, rgb.B)
	}
}

func TestLABToXYZ(t *testing.T) {
	lab := cieLAB{L: 42.49401592, A: 4.416911613, B: -43.38526532}
	xyz := lab.XYZ()
	const (
		wantX = 0.12838835051807143
		wantY = 0.12822135121812256
		wantZ = 0.40841368569543157
	)
	if xyz.X != wantX {
		t.Errorf("X should equal %g but equals %g", wantX, xyz.X)
	}
	if xyz.Y != wantY {
		t.Errorf("Y should equal %g but equals %g", wantY, xyz.Y)
	}
	if xyz.Z != wantZ {
		t.Errorf("Z should equal %g but equals %g", wantZ, xyz.Z)
	}
}

func TestMSHToLAB(t *testing.T) {
	c := MSH{M: 80, S: 1.08, H: -1.1}
	lab := c.lab()
	wantL, wantA, wantB := 37.7062691338992, 32.004211237121645, -62.88058310076059
	if lab.L != wantL {
		t.Errorf("L should equal %g but equals %g", wantL, lab.L)
	}
	if lab.A != wantA {
		t.Errorf("A should equal %g but equals %g", wantA, lab.A)
	}
	if lab.B != wantB {
		t.Errorf("B should equal %g but equals %g", wantB, lab.B)
	}
}

func TestLABTosRGB(t *testing.T) {
	lab := cieLAB{}
	rgb := sRGBA{}
	if lab.sRGB(0.0) != rgb {
		t.Errorf("rgb: have %+v, want %+v", lab.sRGB(0.0), rgb)
	}
	lab = cieLAB{L: 43.22418447, A: 59.07682101, B: 32.27381441}
	rgb = sRGBA{R: 0.7596910553350515, G: 0.16292472671190056, B: 0.20600836034382436, A: 1}
	if lab.sRGB(1) != rgb {
		t.Errorf("rgb: have %+v, want %+v", lab.sRGB(1), rgb)
	}
}

func TestHueTwist(t *testing.T) {
	if hueTwist(MSH{M: 80, S: 1.08, H: -1.1}, 88) != -0.5611585624524025 {
		t.Errorf("hueTwist(80, 1.08, -1.1), 88 should equal -0.561158562"+
			" but equals %g", hueTwist(MSH{M: 80, S: 1.08, H: -1.1}, 88))
	}
}

func TestDivergingMSHAt(t *testing.T) {
	// The test tolerance is the precision of a uint8 expressed as a uint32.
	const tolerance = 1.0 / 256.0 * 65535.0

	start := MSH{M: 80, S: 1.08, H: -1.1}
	end := MSH{M: 80, S: 1.08, H: 0.5}
	p := NewDivergingMSH(start, end)
	p.max = 1
	scalar := 0.125
	rgb, err := p.At(scalar)
	if err != nil {
		t.Error(err)
	}
	wantR, wantG, wantB, wantA := color.NRGBA{R: 98, G: 130, B: 234, A: 255}.RGBA()
	r, g, b, a := rgb.RGBA()
	if math.Abs(float64(r)-float64(wantR)) > tolerance {
		t.Errorf("R: want %v but have %v", r, wantR)
	}
	if math.Abs(float64(g)-float64(wantG)) > tolerance {
		t.Errorf("G: want %v but have %v", g, wantG)
	}
	if math.Abs(float64(b)-float64(wantB)) > tolerance {
		t.Errorf("B: want %v but have %v", b, wantB)
	}
	if math.Abs(float64(a)-float64(wantA)) > tolerance {
		t.Errorf("A: want %v but have %v", a, wantA)
	}
}

func TestSToLinear(t *testing.T) {
	s := 0.735356983
	linear := 0.499999999920366
	if sToLinear(s) != linear {
		t.Errorf("linear: have %g, want %g", sToLinear(s), linear)
	}
	s = 0.01292
	linear = 0.001
	if sToLinear(s) != linear {
		t.Errorf("linear: have %g, want %g", sToLinear(s), linear)
	}

	srgb := sRGBA{R: 0.759704028, G: 0.162897038, B: 0.206033415}
	lrgb := linearRGB{R: 0.5377665307661512, G: 0.022698506403451876, B: 0.035015856125996676}
	if srgb.linearRGB() != lrgb {
		t.Errorf("linear: have %+v, want %+v", srgb.linearRGB(), lrgb)
	}
}

func TestXYZToLRGB(t *testing.T) {
	rgb := linearRGB{R: 0.28909265477940005, G: 0.0663313933285, B: 0.0500602839142}
	xyz := cieXYZ{X: 0.151975056, Y: 0.112509738, Z: 0.061066471}
	if xyz.linearRGB() != rgb {
		t.Errorf("rgb: have %+v, want %+v", xyz.linearRGB(), rgb)
	}
	xyz = cieXYZ{X: 0.1519777983318093, Y: 0.11251566341324888, Z: 0.061068490182446714}
	if rgb.XYZ() != xyz {
		t.Errorf("xyz: have %+v, want %+v", rgb.XYZ(), xyz)
	}
}

func TestXYZToLAB(t *testing.T) {
	xyz := cieXYZ{X: 0.151975056, Y: 0.112509738, Z: 0.061066471}
	lab := cieLAB{L: 40.00000000055783, A: 30.000000104296763, B: 19.99999996294335}
	if xyz.LAB() != lab {
		t.Errorf("lab: have %+v, want %+v", xyz.LAB(), lab)
	}
	xyz = cieXYZ{X: 0.15197025931227778, Y: 0.11250973800000005, Z: 0.06105693812573921}
	if lab.XYZ() != xyz {
		t.Errorf("xyz: have %+v, want %+v", lab.XYZ(), xyz)
	}
}

func TestColorToRGB(t *testing.T) {
	c := color.NRGBA{R: 194, G: 42, B: 53, A: 100}
	rgb := sRGBA{R: 0.7607782101167315, G: 0.16466926070038912, B: 0.20782101167315176, A: 0.39215686274509803}
	if colorTosRGBA(c) != rgb {
		t.Errorf("rgb: have %+v, want %+v", colorTosRGBA(c), rgb)
	}
}

func TestLABToMSH(t *testing.T) {
	lab := cieLAB{L: 43.22418447, A: 59.07682101, B: 32.27381441}
	msh := MSH{M: 80.00000000197056, S: 1.0000000000076632, H: 0.5000000000023601}
	if lab.MSH() != msh {
		t.Errorf("msh: have %+v, want %+v", lab.MSH(), msh)
	}
}

func TestColorToMSH(t *testing.T) {
	c := color.NRGBA{B: 255, A: 255}
	msh := ColorToMSH(c)
	wantM, wantS, wantH := 137.64998152940237, 1.333915268336423, -0.9374394027523394
	if msh.M != wantM {
		t.Errorf("M: want %g but have %g", wantM, msh.M)
	}
	if msh.S != wantS {
		t.Errorf("S: want %g but have %g", wantS, msh.S)
	}
	if msh.H != wantH {
		t.Errorf("H: want %g but have %g", wantH, msh.H)
	}
}
