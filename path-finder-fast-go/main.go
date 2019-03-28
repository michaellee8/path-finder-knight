package main

import (
	"flag"
	"fmt"
	"sync"
)

func main() {

	// Improvement against the previously slow and resource hungry version
	//

	// Get the parameters from argument
	widthPtr := flag.Int("width", 2, "width of fields")
	heightPtr := flag.Int("height", 2, "height of fields")
	flag.Parse()

	// Prepare the WaitGroup
	wg := &sync.WaitGroup{}

	// Prepare the channel
	ch := make(chan Point)

	// Consider all exit on vertical sides
	for row := 0; row < *heightPtr; row++ {
		wg.Add(1)
		go ConsiderExitPoint(Point{Row: row, Column: 0}, ch, *widthPtr,
			*heightPtr, Point{Row: 0, Column: 0}, wg)
		wg.Add(1)
		go ConsiderExitPoint(Point{Row: row, Column: *widthPtr - 1}, ch,
			*widthPtr, *heightPtr, Point{Row: 0, Column: 0}, wg)
	}

	// Consider all exit on horizontal sides
	for col := 0; col < *widthPtr; col++ {
		wg.Add(1)
		go ConsiderExitPoint(Point{Row: 0, Column: col}, ch, *widthPtr,
			*heightPtr, Point{Row: 0, Column: 0}, wg)
		wg.Add(1)
		go ConsiderExitPoint(Point{Row: *heightPtr - 1, Column: col}, ch,
			*widthPtr, *heightPtr, Point{Row: 0, Column: 0}, wg)
	}

	// Use a PointSet to collect our exit points
	ps := NewPointSet(*widthPtr, *heightPtr)

	waitCollector := make(chan struct{})

	// Fire another go routine to collect our exit points
	go func() {
		for p := range ch {
			//fmt.Printf("%d,%d\n", p.Row, p.Column)
			ps.Set(p)
		}
		waitCollector <- struct{}{}
	}()

	// Wait until everything is finished
	wg.Wait()

	// Close the channel so we can do range on it
	close(ch)

	// Wait for the collector to collect all the points
	<-waitCollector

	psl := ps.ToSortedSliced()

	for _, p := range psl {
		fmt.Printf("%d,%d\n", p.Row, p.Column)
	}
}
