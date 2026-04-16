package anomaly

import (
	"math"
	"testing"
)

func TestMeanAndStddev(t *testing.T) {
	tests := []struct {
		name       string
		values     []float64
		wantMean   float64
		wantStddev float64
	}{
		{
			name:       "empty",
			values:     []float64{},
			wantMean:   0,
			wantStddev: 0,
		},
		{
			name:       "single value",
			values:     []float64{10.0},
			wantMean:   10.0,
			wantStddev: 0,
		},
		{
			name:       "uniform values",
			values:     []float64{5.0, 5.0, 5.0},
			wantMean:   5.0,
			wantStddev: 0,
		},
		{
			name:       "varied values",
			values:     []float64{10, 20, 30, 40, 50},
			wantMean:   30.0,
			wantStddev: math.Sqrt(200),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMean, gotStddev := meanAndStddev(tt.values)
			if math.Abs(gotMean-tt.wantMean) > 0.001 {
				t.Errorf("mean: got %f, want %f", gotMean, tt.wantMean)
			}
			if math.Abs(gotStddev-tt.wantStddev) > 0.001 {
				t.Errorf("stddev: got %f, want %f", gotStddev, tt.wantStddev)
			}
		})
	}
}
