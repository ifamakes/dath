package dath

import (
	"math"
)

// InterpolateType offers an enum to select which color space to interpolate in.
type InterpolateType int64

const (
	None InterpolateType = iota
	UseRGB
	UseCYMK
	UseHSV
	UseHSL
	UseLUV
	UseHCL
	UseLAB
)

// Interpolate returns the interpolation of the given colors at a given ratio.
// vt specifies the ratio and the color space to interpolate within (float64 or InterpolateType).
// i.e. Interpolate(color1, color1, 0.75, UseHSV)
// The default ratio and color space are 0.5 and UseLUV
func Interpolate(c1 *Color, c2 *Color, vt ...interface{}) (c *Color) {
	var v float64
	var t InterpolateType
	if len(vt) > 0 {
		for _, n := range vt {
			switch n := n.(type) {
			case float64:
				v = n
			case InterpolateType:
				t = n
			}
		}
	} else {
		v = 0.5
	}

	switch t {
	case UseRGB:
		c = mixRGB(c1, c2, v)
	case UseCYMK:
	case UseHSV:
		hsv := mixHSV(c1.HSV(), c2.HSV(), v)
		c = NewColor().FromHSV(hsv.H, hsv.S, hsv.V)
	case UseHSL:
		hsl := mixHSL(c1.HSL(), c2.HSL(), v)
		c = NewColor().FromHSL(hsl.H, hsl.S, hsl.L)
	case UseLAB:
		lab := mixLAB(c1.LAB(), c2.LAB(), v)
		c = NewColor().FromLAB(lab.L, lab.A, lab.B)
	case UseHCL:
		fallthrough
	case UseLUV:
		fallthrough
	default:
		luv := mixLUV(c1.LUV(), c2.LUV(), v)
		c = NewColor().FromLUV(luv.L, luv.U, luv.V)
	}
	return
}

func lerp(v1, v2, r float64) float64 {
	if math.IsNaN(v1) {
		v1 = 0.0
	}
	if math.IsNaN(v2) {
		v2 = 0.0
	}
	return (1-r)*v1 + r*v2
}

func hslOrhsv(h1, s1, o1, h2, s2, o2, v float64) (hh, ss, oo float64) {
	if (h2 - h1) > 180 {
		hh = h1 + 360
		hh = math.Mod((1-v)*hh+v*h2, 360.0)
	}
	if (h2 - h1) <= 180 {
		hh = h1 + v*(h2-h1)
	}
	ss = math.Max(0.0, math.Min((1-v)*s1+v*s2, 1.0))
	oo = math.Max(0.0, math.Min((1-v)*o1+v*o2, 1.0))
	return
}

func mixLUV(c1 *LUV, c2 *LUV, v float64) *LUV {
	new := &LUV{}
	new.L = lerp(c1.L, c2.L, v)
	new.U = lerp(c1.U, c2.U, v)
	new.V = lerp(c1.V, c2.V, v)
	return new
}

func mixLAB(c1 *LAB, c2 *LAB, v float64) *LAB {
	new := &LAB{}
	new.L = lerp(c1.L, c2.L, v)
	new.A = lerp(c1.A, c2.A, v)
	new.B = lerp(c1.B, c2.B, v)
	return new
}

func mixRGB(c1 *Color, c2 *Color, v float64) *Color {
	new := &Color{}
	new.r = lerp(c1.r, c2.r, v)
	new.g = lerp(c1.g, c2.g, v)
	new.b = lerp(c1.b, c2.b, v)
	return new
}

func mixHSV(c1 *HSV, c2 *HSV, v float64) *HSV {
	hsv := &HSV{}
	hsv.H, hsv.S, hsv.V = hslOrhsv(c1.H, c1.S, c1.V, c2.H, c2.S, c2.V, v)
	return hsv
}
func mixHSL(c1 *HSL, c2 *HSL, v float64) *HSL {
	hsl := &HSL{}
	hsl.H, hsl.S, hsl.L = hslOrhsv(c1.H, c1.S, c1.L, c2.H, c2.S, c2.L, v)
	return hsl
}

/* func easeInOut(x float64) float64 {
	return -(math.Cos(math.Pi*x) - 1) / 2
} */

/* func Gradient(a *HSL, b *HSL, v float64) *HSL {
	c := &HSL{}
	v2 := easeInOut(v)
	if (b.H - a.H) > 180 {
		c.H = a.H + 360
		c.H = math.Mod((1-v2)*c.H+v2*b.H, 360.0)
	}
	if (b.H - a.H) <= 180 {
		c.H = a.H + v*(b.H-a.H)
	}
	c.S = math.Max(0.0, math.Min((1-v)*a.S+v*b.S, 1.0))
	c.L = math.Max(0.0, math.Min((1-v)*a.L+v*b.L, 1.0))
	return c
}

func Gradient(a *LUV, b *LUV, v float64) *LUV {
	c := &LUV{}
	v = easeInOut(v)

	c.L = a.L + (b.L-a.L)*v
	c.U = a.U + (b.U-a.U)*v
	c.V = a.V + (b.V-a.V)*v
	return c
}
*/
