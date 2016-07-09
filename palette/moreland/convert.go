// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package moreland

import (
	"fmt"
	"image/color"
	"math"
)

// linearToS converts from a linear RGB component to an sRGB component.
func linearToS(v float64) float64 {
	if v > 0.0031308 {
		return 1.055*math.Pow(v, 1/2.4) - 0.055
	}
	return 12.92 * v
}

// sToLinear converts from an sRGB component to a linear RGB component.
func sToLinear(v float64) float64 {
	if v > 0.04045 {
		return math.Pow((v+0.055)/1.055, 2.4)
	}
	return v / 12.92
}

// linearRGB repesents a physically linear RGB color.
type linearRGB struct {
	R, G, B float64
}

// XYZ converts a linear RGB color to a CIE XYZ color.
func (c linearRGB) XYZ() cieXYZ {
	return cieXYZ{
		X: 0.4124*c.R + 0.3576*c.G + 0.1805*c.B,
		Y: 0.2126*c.R + 0.7152*c.G + 0.0722*c.B,
		Z: 0.0193*c.R + 0.1192*c.G + 0.9505*c.B,
	}
}

// linearRGB converts a CIEXYZ color to physically linear RGB.
func (c cieXYZ) linearRGB() linearRGB {
	return linearRGB{
		R: c.X*3.2406 + c.Y*-1.5372 + c.Z*-0.4986,
		G: c.X*-0.9689 + c.Y*1.8758 + c.Z*0.0415,
		B: c.X*0.0557 + c.Y*-0.204 + c.Z*1.057,
	}
}

// labTemp is an intermediate step in converting from CIE XYZ to CIE LAB.
func labTemp(v float64) float64 {
	if v > 0.008856 {
		return math.Pow(v, 1.0/3.0)
	}
	return 7.787*v + 16.0/116.0
}

// LAB converts an CIE XYZ color to a CIE LAB color.
func (c cieXYZ) LAB() cieLAB {
	tempX := labTemp(c.X / 0.9505)
	tempY := labTemp(c.Y)
	tempZ := labTemp(c.Z / 1.089)
	return cieLAB{
		L: (116.0 * tempY) - 16.0,
		A: 500.0 * (tempX - tempY),
		B: 200 * (tempY - tempZ),
	}
}

// sRGB represents a color within the sRGB color space, with an alpha channel
// but not premultiplied.
type sRGBA struct {
	R, G, B, A float64
}

func (c sRGBA) linearRGB() linearRGB {
	return linearRGB{
		R: sToLinear(c.R),
		G: sToLinear(c.G),
		B: sToLinear(c.B),
	}
}

// RGBA implements the color.Color interface.
func (c sRGBA) RGBA() (r, g, b, a uint32) {
	return uint32(c.R * c.A * 65535.0), uint32(c.G * c.A * 65535.0), uint32(c.B * c.A * 65535.0), uint32(c.A * 65535.0)
}

// LAB converts a sRGB color to a CIE LAB color.
func (c sRGBA) LAB() cieLAB {
	return c.linearRGB().XYZ().LAB()
}

// sRGB converts a CIE LAB color to an sRGBA color, where alpha is opacity
// between 0 and 1.
func (c cieLAB) sRGB(alpha float64) sRGBA {
	return c.XYZ().linearRGB().S(alpha)
}

// colorTosRGBA converts a color to an sRGBA.
func colorTosRGBA(c color.Color) sRGBA {
	r, g, b, a := c.RGBA()
	alpha := float64(a) / 65535.0
	return sRGBA{
		R: float64(r) / alpha / 65535.0,
		G: float64(g) / alpha / 65535.0,
		B: float64(b) / alpha / 65535.0,
		A: alpha,
	}
}

// ColorToMSH converts a color to MSH space.
func ColorToMSH(c color.Color) MSH {
	return colorTosRGBA(c).LAB().MSH()
}

// S converts a linear RGB color to an sRGBA color, where alpha is the alpha
// channel value (between 0 and 1).
func (c linearRGB) S(alpha float64) sRGBA {
	sc := sRGBA{
		R: linearToS(c.R),
		G: linearToS(c.G),
		B: linearToS(c.B),
		A: alpha,
	}
	return sc
}

// check returns an error if any of the channels in c are out of range.
func (c sRGBA) check() error {
	if c.R > 1 || c.G > 1 || c.B > 1 || c.A > 1 || c.R < 0 || c.G < 0 || c.B < 0 || c.A < 0 {
		return fmt.Errorf("moreland: invalid color r:%g, g:%g, b:%g, a:%g", c.R, c.G, c.B, c.A)
	}
	return nil
}

// fix forces all channels in c to be within the range (0,1).
func (c *sRGBA) fix() {
	if c.R > 1 {
		c.R = 1
	}
	if c.G > 1 {
		c.G = 1
	}
	if c.B > 1 {
		c.B = 1
	}
	if c.A > 1 {
		c.A = 1
	}
	if c.R < 0 {
		c.R = 0
	}
	if c.G < 0 {
		c.G = 0
	}
	if c.B < 0 {
		c.B = 0
	}
	if c.A < 0 {
		c.A = 0
	}
}

// cieLAB represents a color in CIE LAB space.
type cieLAB struct {
	L, A, B float64
}

// cieXYZ represents a color in CIE XYZ space.
type cieXYZ struct {
	X, Y, Z float64
}

// xyzTemp is an intermediate step in converting from CIE LAB to CIE XYZ.
func xyzTemp(v float64) float64 {
	const (
		xlim = 0.008856
		a    = 7.787
		b    = 16. / 116.
		ylim = a*xlim + b
	)
	if v > ylim {
		return v * v * v
	}
	return (v - b) / a
}

// XYZ converts a CIELAB color to a CIEXYZ color.
func (c cieLAB) XYZ() cieXYZ {
	// Reference white-point D65
	const xn, yn, zn = 0.95047, 1.0, 1.08883
	return cieXYZ{
		X: xn * xyzTemp((c.A/500)+(c.L+16)/116),
		Y: yn * xyzTemp((c.L+16.)/116.),
		Z: zn * xyzTemp((c.L+16.)/116.-(c.B/200.)),
	}
}

func (c cieLAB) MSH() MSH {
	m := math.Pow(c.L*c.L+c.A*c.A+c.B*c.B, 0.5)
	return MSH{
		M: m,
		S: math.Acos(c.L / m),
		H: math.Atan2(c.B, c.A),
	}
}

// lab converts a MSH color to a CIELAB color.
func (c MSH) lab() cieLAB {
	return cieLAB{
		L: c.M * math.Cos(c.S),
		A: c.M * math.Sin(c.S) * math.Cos(c.H),
		B: c.M * math.Sin(c.S) * math.Sin(c.H),
	}
}

// MSH defines a color in Magnitude-Saturation-Hue color space.
type MSH struct {
	M, S, H float64
}

func hueTwist(c MSH, convergeM float64) float64 {
	signH := c.H / math.Abs(c.H)
	return signH * c.S * math.Sqrt(convergeM*convergeM-c.M*c.M) / (c.M * math.Sin(c.S))
}
