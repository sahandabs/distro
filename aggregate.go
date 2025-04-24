package main

import "fmt"

// Note: aggreagator it not thread safe as of now
type aggregator struct {
	// this is a map of selected time periods to map of household id and the aggregated consumption
	data       map[int]map[string]map[int]float64 // this can be sth like map[year_quarter_houseid]float64
	timePeriod TimePeriod
}

type DataAggregator interface {
	Aggregate(row) error
}

var _ DataAggregator = aggregator{}

func newAggregator(tp TimePeriod) (aggregator, error) {
	if !tp.IsValid() {
		return aggregator{}, fmt.Errorf("invalid time period: %s", tp)
	}
	aggData := make(map[int]map[string]map[int]float64)
	return aggregator{data: aggData, timePeriod: tp}, nil
}

func (ag aggregator) Aggregate(r row) error {
	year := r.time.Year()
	if ag.data[year] == nil {
		ag.data[year] = map[string]map[int]float64{}
	}
	p := ag.timePeriod.GetPeriod(r.time)
	if ag.data[year][p] == nil {
		ag.data[year][p] = make(map[int]float64)
	}
	ag.data[year][p][r.houseId] += r.consumption

	// TODO(Sahand): check the 15 min interval,
	// if it passes the current period do liner interpolation and add the rest to the next quarter
	// the 15 min period can be a time duration in the aggregator struct
	return nil
}
