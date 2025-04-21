package main

import "time"

type TimePeriod string

var QuarterTimePeriod TimePeriod = "Quarter"

func (tp TimePeriod) IsValid() bool {
	return tp == QuarterTimePeriod
}

func (tp TimePeriod) GetPeriod(t time.Time) string {
	if t.Before(time.Date(t.Year(), time.April, 1, 0, 0, 0, 0, t.Location())) {
		return "Q1"
	}
	if t.Before(time.Date(t.Year(), time.July, 1, 0, 0, 0, 0, t.Location())) {
		return "Q2"
	}
	if t.Before(time.Date(t.Year(), time.October, 1, 0, 0, 0, 0, t.Location())) {
		return "Q3"
	}
	return "Q4"
}
