// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package moreland

import (
	"fmt"
	"image/color"
	"math"
	"sort"
)

// Luminance is a color palette that interpolates
// between control colors in a way that ensures a linear relationship
// between the luminance of a color and the value it represents.
type Luminance struct {
	colors  []cieLAB
	scalars []float64

	// Alpha represents the opacity of the returned
	// colors in the range (0,1). It is 1 by default.
	Alpha float64

	// max is the maximum value of the range of scalars that can be
	// mapped to colors using this palette. In a Luminance color map
	// the minimum value is required to be zero so that the luminance
	// of the color is a linear function of the value being represented.
	max float64
}

// NewLuminance creates a new Luminance color scale from the given controlColors.
// If the luminance of the controlColors is not monotonically increasing, an
// error will be returned.
func NewLuminance(controlColors []color.Color) (*Luminance, error) {
	l := Luminance{
		colors:  make([]cieLAB, len(controlColors)),
		scalars: make([]float64, len(controlColors)),
		Alpha:   1,
	}
	max := math.Inf(-1)
	min := math.Inf(1)
	for i, c := range controlColors {
		lab := colorTosRGBA(c).LAB()
		l.colors[i] = lab
		max = math.Max(max, lab.L)
		min = math.Min(min, lab.L)
		if i > 0 && lab.L <= l.colors[i-1].L {
			return nil, fmt.Errorf("moreland: luminance of color %d (%g) is not "+
				"greater than that of color %d (%g)", i, lab.L, i-1, l.colors[i-1].L)
		}
	}
	// Normalize scalar values to the range (0,1).
	rnge := max - min
	for i, c := range l.colors {
		l.scalars[i] = (c.L - min) / rnge
	}
	// Avoid floating point math problems.
	l.scalars[0] = 0
	l.scalars[len(l.scalars)-1] = 1
	return &l, nil
}

// At implements the palette.ColorMap interface for a Luminance object.
func (l *Luminance) At(scalar float64) (color.Color, error) {
	if l.max == 0 {
		return nil, fmt.Errorf("moreland: Luminance color map Max == 0")
	}
	scalar = scalar / l.max
	if scalar < 0 || scalar > 1 {
		return nil, fmt.Errorf("moreland: interpolation value (%g) out of range (0,%g)", scalar, l.max)
	}
	i := sort.SearchFloat64s(l.scalars, scalar)
	if i == 0 {
		return l.colors[i].XYZ().linearRGB().S(l.Alpha), nil
	}
	c1 := l.colors[i-1]
	c2 := l.colors[i]
	frac := (scalar - l.scalars[i-1]) / (l.scalars[i] - l.scalars[i-1])
	o := cieLAB{
		L: frac*(c2.L-c1.L) + c1.L,
		A: frac*(c2.A-c1.A) + c1.A,
		B: frac*(c2.B-c1.B) + c1.B,
	}.XYZ().linearRGB().S(l.Alpha)
	o.fix()
	return o, nil
}

// SetMax implements the palette.ColorMap interface for a Luminance object.
func (l *Luminance) SetMax(v float64) {
	l.max = v
}

// SetMin implements the palette.ColorMap interface for a Luminance object.
// However, it will panic whenever it is called because the
// minimum value must always be zero.
func (l *Luminance) SetMin(v float64) {
	panic("moreland: Luminance minimum value cannot be changed from zero")
}

// Max implements the palette.ColorMap interface for a Luminance object.
func (l *Luminance) Max() float64 {
	return l.max
}

// Min implements the palette.ColorMap interface for a Luminance object.
func (l *Luminance) Min() float64 {
	return 0
}

// Palette fulfils the palette.Palette interface.
type Palette []color.Color

// Colors fulfils the palette.Palette interface.
func (p Palette) Colors() []color.Color {
	return p
}

// Palette returns an object that fulfills the palette.Palette interface,
// where nColors is the number of desired colors.
func (l Luminance) Palette(nColors int) Palette {
	if l.max == 0 {
		l.max = 1
	}
	delta := l.max / float64(nColors-1)
	v := 0.0
	c := make([]color.Color, nColors)
	for i := 0; i < nColors; i++ {
		var err error
		c[i], err = l.At(v)
		if err != nil {
			panic(err)
		}
		v += delta
	}
	return Palette(c)
}

// BlackBody is based on the colors of black body radiation.
// Although the colors are inspired by the wavelengths of light from
// black body radiation, the actual colors used are designed to be
// perceptually uniform. Colors of the desired brightness and hue are chosen,
// and then the colors are adjusted such that the luminance is perceptually
// linear (according to the CIELAB color space).
func BlackBody() *Luminance {
	return &Luminance{
		colors: []cieLAB{
			cieLAB{L: 0, A: 0, B: 0},
			cieLAB{L: 39.112572747719774, A: 55.92470934659227, B: 37.65159714510402},
			cieLAB{L: 58.45705480680232, A: 43.34389690857626, B: 65.95409116544081},
			cieLAB{L: 84.13253643355525, A: -6.459770854468639, B: 82.41994470228775},
			cieLAB{L: 100, A: 0, B: 0}},
		scalars: []float64{0, 0.39112572747719776, 0.5845705480680232, 0.8413253643355525, 1},
		Alpha:   1}
}

// ExtendedBlackBody derives a color map based on the colors of black body radiation
// with some blue and purple hues thrown in at the lower end to add some "color."
// The color map is similar to the default colors used in gnuplot. Colors of
// the desired brightness and hue are chosen, and then the colors are adjusted
// such that the luminance is perceptually linear (according to the CIELAB
// color space).
func ExtendedBlackBody() *Luminance {
	return &Luminance{
		colors: []cieLAB{
			cieLAB{L: 0, A: 0, B: 0},
			cieLAB{L: 21.873483862751876, A: 50.19882295659109, B: -74.66982659778306},
			cieLAB{L: 34.506542513775905, A: 75.41302687474061, B: -88.73807072507786},
			cieLAB{L: 47.02980511087303, A: 70.93217189227919, B: 33.59880053746508},
			cieLAB{L: 65.17482203230537, A: 49.14591409658836, B: 56.86480950937553},
			cieLAB{L: 84.13253643355525, A: -6.459770854468639, B: 82.41994470228775},
			cieLAB{L: 100, A: 0, B: 0},
		},
		scalars: []float64{0, 0.21873483862751875, 0.34506542513775906, 0.4702980511087303,
			0.6517482203230537, 0.8413253643355525, 1},
		Alpha: 1,
	}
}

// Kindlmann uses the colors first proposed in a paper
// by Kindlmann, Reinhard, and Creem. The map is basically the rainbow
// color map with the luminance adjusted such that it monotonically
// changes, making it much more perceptually viable.
//
// Citation:
// Gordon Kindlmann, Erik Reinhard, and Sarah Creem. 2002. Face-based
// luminance matching for perceptual colormap generation. In Proceedings
// of the conference on Visualization '02 (VIS '02). IEEE Computer Society,
// Washington, DC, USA, 299-306.
func Kindlmann() *Luminance {
	return &Luminance{
		colors: []cieLAB{
			cieLAB{L: 0, A: 0, B: 0},
			cieLAB{L: 10.479520542426698, A: 34.05557958902206, B: -34.21934877170809},
			cieLAB{L: 21.03011379005111, A: 52.30473571100955, B: -61.852601228346536},
			cieLAB{L: 31.03098927978494, A: 23.814976212074402, B: -57.73419358300511},
			cieLAB{L: 40.21480513626115, A: -24.858012706049536, B: -7.322176588219942},
			cieLAB{L: 52.73108089333358, A: -19.064976357731634, B: -25.558178073848147},
			cieLAB{L: 60.007326812392634, A: -61.75624590074585, B: 56.43522875191319},
			cieLAB{L: 69.81578343076002, A: -58.33353084882392, B: 68.37457857626646},
			cieLAB{L: 79.55703752324776, A: -22.50477758899383, B: 78.57946686200843},
			cieLAB{L: 89.818961593653, A: 7.586705160677109, B: 15.375961528833981},
			cieLAB{L: 100, A: 0, B: 0},
		},
		scalars: []float64{0, 0.10479520542426699, 0.2103011379005111, 0.3103098927978494,
			0.4021480513626115, 0.5273108089333358, 0.6000732681239264, 0.6981578343076003,
			0.7955703752324775, 0.89818961593653, 1},
		Alpha: 1,
	}
}

// ExtendedKindlmann uses the colors from Kindlmann but also
// adds more hues by doing a more than 360 degree loop around the hues.
// This works because the endpoints have low saturation and very
// different brightness.
func ExtendedKindlmann() *Luminance {
	return &Luminance{
		colors: []cieLAB{
			cieLAB{L: 0, A: 0, B: 0},
			cieLAB{L: 13.371291966477482, A: 40.39368469479174, B: -47.73239449160565},
			cieLAB{L: 25.072421338587574, A: -18.01441053740843, B: -5.313556572210176},
			cieLAB{L: 37.411516363056116, A: -43.058336774976055, B: 39.30203907343062},
			cieLAB{L: 49.75026355291354, A: -15.774050138318895, B: 53.507917567416094},
			cieLAB{L: 61.643756252245225, A: 52.67703578954919, B: 43.82595336046358},
			cieLAB{L: 74.93187540089825, A: 50.92061741619164, B: -30.235411697966242},
			cieLAB{L: 87.64732748562544, A: 14.355163639545697, B: -17.471161313826332},
			cieLAB{L: 100, A: 0, B: 0},
		},
		scalars: []float64{0, 0.13371291966477483, 0.25072421338587575, 0.37411516363056113,
			0.4975026355291354, 0.6164375625224523, 0.7493187540089825, 0.8764732748562544, 1},
		Alpha: 1,
	}
}
