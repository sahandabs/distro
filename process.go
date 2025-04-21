package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
)

func process(ctx context.Context, inPath, outPath string, wg *sync.WaitGroup) error {
	defer wg.Done()
	readP, inClose, err := newReadProcessorCSV(inPath, true)
	if err != nil {
		return fmt.Errorf("error when initiating read processor: %w", err)
	}
	defer inClose()
	agg, err := newAggregator(QuarterTimePeriod)
	if err != nil {
		return fmt.Errorf("error when initiating aggregator: %w", err)
	}
	outProcessor, outClose, err := newOutPutManager(outPath, "year", "quarter", "houseID", "comsumption")
	if err != nil {
		return fmt.Errorf("error when initiating output manager: %w", err)
	}
	defer outClose()

	var finished bool
	count := 0
	for !finished {
		select {
		case <-ctx.Done():
			// we can also use the ctx err
			return ErrProcessEaryBail
		default:
		}
		row, err := readP.readRow()
		if err != nil {
			var pe FieldParseError
			if errors.As(err, &pe) && pe.IsRetryable() {
				fmt.Printf("ignore parsing error, %s, line: %d \n", pe.Error(), count+1)
				count++
				continue
			}
			if errors.Is(err, io.EOF) {
				fmt.Println("no rows left, so far:", count)
				break
			}
			return fmt.Errorf("error when reading data error: %w, line: %d", err, count+1)
		}
		err = agg.Aggregate(row)
		if err != nil {
			return fmt.Errorf("error when aggregating data error: %w, line: %d", err, count+1)
		}
		count++
		if count%1000 == 0 {
			fmt.Printf("proccessed: %d\n", count)
		}
	}

	err = outProcessor.Save(agg.data)
	if err != nil {
		return fmt.Errorf("error when saving output %w", err)
	}
	fmt.Println(agg.data)

	return nil
}

// for high scale inputs we can use parallel processing,
// the simplest form can be a go routine reading the input and putting the input in a channel
// a number of go routines will fan out and use the same channel as their input, and update a shared map(either sync map or a map with rw lock)
// in higher scale the result can be multiple services a message broker and partiitoning based on the house id?
// and using a key value map
// this is a typical map reduce problem

// func parallelProcess(ctx context.Context, inPath, outPath string, wg *sync.WaitGroup) error {
// read from file in a go routine, put in a channel
// input := make(chan []string, 100)
// go func() {
//     for {
//         rec, err := reader.Read()
//         // parse
//         input <- rec
//     }
//     close(csvChan)
// }()
// jobs := 5
// for i <5 {
// wait group
// go(){
// read from input
// update shared map
// }()
// }
// }
