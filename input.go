package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type readProcessor struct {
	reader *csv.Reader
}

func newReadProcessorCSV(path string, shouldReadFirstLine bool) (*readProcessor, func(), error) {
	file, err := os.Open(path)
	if err != nil {
		return &readProcessor{}, func() {}, fmt.Errorf("error when opening input %w", err)
	}
	// use bufio to buffer the reading of the file
	buf := bufio.NewReader(file)
	reader := csv.NewReader(buf)

	if shouldReadFirstLine {
		columns, err := reader.Read()
		if err != nil {
			return &readProcessor{}, func() {}, fmt.Errorf("error when reading the head %w", err)
		}
		fmt.Println("coloumns are", columns)
	}

	return &readProcessor{reader: reader}, func() { file.Close() }, nil
}

func (r *readProcessor) readRow() (row, error) {
	// condiser to set record reuse to true
	records, err := r.reader.Read()
	if err != nil {
		return row{}, err
	}

	// each row has 3 feilds
	if len(records) != 3 {
		return row{}, FieldParseError{err: ErrInvalidFieldsCount}
	}

	houseId, err := strconv.Atoi(records[0])
	if err != nil {
		return row{}, FieldParseError{err: errors.Join(err, ErrInvalidField), fieldName: "house-number"}
	}
	if houseId == 0 {
		return row{}, FieldParseError{err: ErrInvalidField, fieldName: "house-number"}
	}

	consumption, err := strconv.ParseFloat(records[1], 64)
	if err != nil {
		return row{}, FieldParseError{err: errors.Join(err, ErrInvalidField), fieldName: "consumption"}
	}

	seconds, err := strconv.Atoi(records[2])
	if err != nil {
		return row{}, FieldParseError{err: errors.Join(err, ErrInvalidField), fieldName: "time"}
	}
	t := time.Unix(int64(seconds), 0)
	if t.IsZero() {
		return row{}, FieldParseError{err: ErrInvalidField, fieldName: "time"}
	}

	// no need to change the time zone, the standard unix time is in utc and when reading Time library will move it to local time
	// however to make sure this code runs everywhere correctly we can set the time zone to amsterdam time at the begining
	// location, err := time.LoadLocation("Europe/Amsterdam")
	// if err != nil {
	// 	fmt.Println("Error loading location:", err)
	// 	return row{}, FieldParseError{error: ErrInvalidField, fieldName: "time"}
	// }
	// nlTime := t.In(location)

	return row{consumption: consumption, houseId: houseId, time: t}, nil
}
