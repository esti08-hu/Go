package main

import "testing"

func TestCalculateGrade(t *testing.T) {
	tests := []struct {
		score float64
		want  string
	}{
		{100, "A"},
		{90, "A"},
		{89.99, "B"},
		{80, "B"},
		{79.99, "C"},
		{70, "C"},
		{69.99, "D"},
		{60, "D"},
		{59.99, "F"},
		{0, "F"},
		{-10, "F"},
	}
	for _, tt := range tests {
		got := calculateGrade(tt.score)
		if got != tt.want {
			t.Errorf("calculateGrade(%v) = %v; want %v", tt.score, got, tt.want)
		}
	}
}

func TestCalculateAvarage(t *testing.T) {
	cases := []struct {
		scores []float64
		want   float64
	}{
		{[]float64{100, 90, 80}, 90.0},
		{[]float64{50, 50, 50, 50}, 50.0},
		{[]float64{}, 0.0},
		{[]float64{100}, 100.0},
		{[]float64{0, 0, 0}, 0.0},
	}
	for _, c := range cases {
		got := calculateAvarage(c.scores)
		if got != c.want {
			t.Errorf("calculateAvarage(%v) = %v; want %v", c.scores, got, c.want)
		}
	}
}
