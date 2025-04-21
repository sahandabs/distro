package main

import (
	"testing"
	"time"
)

func TestAggregator_Aggregate(t *testing.T) {
	agg, err := newAggregator(QuarterTimePeriod)
	if err != nil {
		t.Fatalf("could not create aggregator: %v", err)
	}

	// January 1st, 2025 → Q1
	loc, _ := time.LoadLocation("Europe/Amsterdam")
	ts := time.Date(2025, 1, 1, 10, 0, 0, 0, loc)
	r := row{time: ts, houseId: 1, consumption: 3.5}

	err = agg.Aggregate(r)
	if err != nil {
		t.Fatalf("aggregate failed: %v", err)
	}

	q := agg.data[2025]["Q1"][1]
	if q != 3.5 {
		t.Errorf("expected 3.5, got %f", q)
	}

	r.time = r.time.Add(24 * time.Hour)
	r.consumption = 1.5

	err = agg.Aggregate(r)
	if err != nil {
		t.Fatalf("aggregate failed: %v", err)
	}

	q = agg.data[2025]["Q1"][1]
	if q != 5 {
		t.Errorf("expected 5, got %f", q)
	}
}
