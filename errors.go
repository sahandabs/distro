package main

import (
	"encoding/csv"
	"errors"
	"fmt"
)

var ErrProcessEaryBail = errors.New("error, processing needed to be finished due to context cancellation")
var ErrInvalidField = errors.New("error, can not parse the field")
var ErrInvalidFieldsCount = errors.New("error, can not parse the row, invalid fields count")

type FieldParseError struct {
	err       error
	fieldName string
}

func (pe FieldParseError) Error() string {
	if pe.fieldName != "" {
		return fmt.Sprintf("error when parsing field %s %v", pe.fieldName, pe.err)
	}
	return pe.err.Error()
}

func (pe FieldParseError) IsRetryable() bool {
	if errors.Is(pe.err, ErrInvalidField) ||
		errors.Is(pe.err, ErrInvalidFieldsCount) ||
		errors.Is(pe.err, csv.ErrFieldCount) {
		return true
	}
	return false
}

func (e FieldParseError) Unwrap() error {
	return e.err
}
