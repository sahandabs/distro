package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type outPutManager struct {
	writer *csv.Writer
}

func newOutPutManager(writePath string, clmns ...string) (*outPutManager, func(), error) {
	f, err := os.Create(writePath)
	if err != nil {
		return &outPutManager{}, func() {}, fmt.Errorf("can not create file error:%w", err)
	}
	// use bufio to buffer the writing to the file
	buf := bufio.NewWriter(f)
	wr := csv.NewWriter(buf)
	if len(clmns) != 0 {
		err = wr.Write(clmns)
		if err != nil {
			return &outPutManager{}, func() {}, fmt.Errorf("can not write columns to the file %w", err)
		}
	}
	return &outPutManager{writer: wr}, func() { wr.Flush(); f.Close() }, nil
}

func (om *outPutManager) Save(data map[int]map[string]map[int]float64) error {
	// TODO(Sahand): call flush occasionally
	for year, yearData := range data {
		for quarter, quarterData := range yearData {
			for houseId, consumption := range quarterData {
				err := om.writer.Write([]string{strconv.Itoa(year), quarter, strconv.Itoa(houseId), fmt.Sprintf("%f", consumption)})
				if err != nil {
					return fmt.Errorf("error when writing data :%w", err)
				}
			}
		}
	}
	return nil
}
