package main

import (
	"testing"
	"time"
)

func TestGetPeriod(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/Amsterdam")
	tp := QuarterTimePeriod

	tests := []struct {
		date     time.Time
		expected string
	}{
		{time.Date(2025, 1, 10, 0, 0, 0, 0, loc), "Q1"},
		{time.Date(2025, 4, 2, 0, 0, 0, 0, loc), "Q2"},
		{time.Date(2025, 8, 1, 0, 0, 0, 0, loc), "Q3"},
		{time.Date(2025, 11, 1, 0, 0, 0, 0, loc), "Q4"},
	}

	for _, test := range tests {
		got := tp.GetPeriod(test.date)
		if got != test.expected {
			t.Errorf("expected %s, got %s", test.expected, got)
		}
	}
}
