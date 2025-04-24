package main

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadRow_ValidRow(t *testing.T) {
	input := "1,4.5,1735689600\n" // 2025-01-01 00:00:00 UTC
	f := strings.NewReader("houseId,consumption,timestamp\n" + input)
	reader := csv.NewReader(f)
	_, _ = reader.Read() // skip header

	rp := &readProcessor{reader: reader}
	row, err := rp.ReadRow()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, 1, row.houseId)
	assert.Equal(t, 4.5, row.consumption)
	assert.Equal(t, row.time.Year(), 2025)
	assert.Equal(t, row.time.Month(), time.Month(1))
	assert.Equal(t, row.time.Day(), 1)
}

func TestReadRow_InvalidNumber(t *testing.T) {
	input := "notanumber,4.5,1735689600\n"
	f := strings.NewReader("houseId,consumption,timestamp\n" + input)
	reader := csv.NewReader(f)
	_, _ = reader.Read()

	rp := &readProcessor{reader: reader}
	_, err := rp.ReadRow()
	assert.ErrorIs(t, err, ErrInvalidField)
}
