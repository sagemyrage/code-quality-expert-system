package fuzzy

import (
	"math"
	"testing"
)

const epsilon = 1e-8

func TestTrapezoid(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		b    float64
		c    float64
		d    float64
		x    float64
		want float64
	}{
		{
			name: "x < a",
			a:    0.0,
			b:    2.0,
			c:    4.0,
			d:    6.0,
			x:    -1.0,
			want: 0,
		},
		{
			name: "x = a",
			a:    0.0,
			b:    2.0,
			c:    4.0,
			d:    6.0,
			x:    0.0,
			want: 0,
		},
		{
			name: "a < x < b",
			a:    0.0,
			b:    2.0,
			c:    4.0,
			d:    6.0,
			x:    1.0,
			want: 0.5,
		},
		{
			name: "x = b",
			a:    0.0,
			b:    2.0,
			c:    4.0,
			d:    6.0,
			x:    2.0,
			want: 1,
		},
		{
			name: "b < x < c",
			a:    0.0,
			b:    2.0,
			c:    4.0,
			d:    6.0,
			x:    3.0,
			want: 1,
		},
		{
			name: "x = c",
			a:    0.0,
			b:    2.0,
			c:    4.0,
			d:    6.0,
			x:    4.0,
			want: 1,
		},
		{
			name: "c < x < d",
			a:    0.0,
			b:    2.0,
			c:    4.0,
			d:    6.0,
			x:    5.0,
			want: 0.5,
		},
		{
			name: "x = d",
			a:    0.0,
			b:    2.0,
			c:    4.0,
			d:    6.0,
			x:    6.0,
			want: 0,
		},
		{
			name: "x > d",
			a:    0.0,
			b:    2.0,
			c:    4.0,
			d:    6.0,
			x:    7.0,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trapezoid(tt.x, tt.a, tt.b, tt.c, tt.d)
			if math.Abs(got-tt.want) > epsilon {
				t.Errorf("got = %f, want = %f", got, tt.want)
			}
		})
	}
}
