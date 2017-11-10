package msg

import "math"

func (v *Vector2) IsZero() bool {
	return math.Abs(float64(v.X)) < 0.00000001 && math.Abs(float64(v.Y)) < 0.00000001
}
