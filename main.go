package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"flag"
)

func main() {
	fmt.Println("Starting...")
	var inPath, outPath string
	flag.StringVar(&inPath, "inpath", "household_consumption.csv", "path to the csv input file")
	flag.StringVar(&outPath, "outpath", "out.csv", "path to the csv output file")
	flag.Parse()
	fmt.Printf("processing: in %s --- out %s\n", inPath, outPath)

	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		defer wg.Done()
		err := process(ctx, inPath, outPath)
		if err != nil && !errors.Is(err, ErrProcessEaryBail) {
			log.Fatalf("error when calling process, error: %v", err)
		}
		if errors.Is(err, ErrProcessEaryBail) {
			log.Print("process bailed early due to context cancel")
		} else {
			fmt.Printf("process finshed succesfully, exitting...")
			os.Exit(0)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig

	fmt.Println("signal recieved...")
	cancel()

	wg.Wait()
	fmt.Println("Shutdown completed...")
}
