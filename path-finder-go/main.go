package main

import (
	"flag"
	"fmt"
	"sync"
)

func main() {

	// Get the parameters from argument
	widthPtr := flag.Int("width", 2, "width of fields")
	heightPtr := flag.Int("height", 2, "height of fields")
	flag.Parse()

	// Prepare for the first Walk
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Prepare a fresh, empty set of passedPoints
	passedPoints := NewPointSet(*widthPtr, *heightPtr)

	// Prepare the channel
	ch := make(chan Point)

	// Fire our recursive walker!
	go Walk(passedPoints, 0, 0, *widthPtr, *heightPtr, &wg, ch)

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

func Walk(
	passedPoints PointSet,
	posR, posC int,
	w, h int,
	wg *sync.WaitGroup,
	ch chan Point,
) {
	// We assumes that passedPoints is a modifiable copy

	// Add current position
	passedPoints.Set(Point{Row: posR, Column: posC})

	// Test if we have any way to walk
	walked := false

	// Consider walking left
	if !passedPoints.Has(Point{posR, posC - 1}) &&
		Between(0, posC-1, w) {

		// We can walk, so mark walk as true
		walked = true

		// Tell the wait group to wait for next go routine
		wg.Add(1)

		// Make a copy of passedPoints and change posC to left, and then pass
		// it down
		go Walk(passedPoints.Copy(), posR, posC-1, w, h, wg, ch)
	}

	// Consider walking right
	if !passedPoints.Has(Point{posR, posC + 1}) &&
		Between(0, posC+1, w) {
		walked = true
		wg.Add(1)
		go Walk(passedPoints.Copy(), posR, posC+1, w, h, wg, ch)
	}

	// Consider walking up
	if !passedPoints.Has(Point{posR + 1, posC}) &&
		Between(0, posR+1, h) {
		walked = true
		wg.Add(1)
		go Walk(passedPoints.Copy(), posR+1, posC, w, h, wg, ch)
	}

	// Consider walking down
	if !passedPoints.Has(Point{posR - 1, posC}) &&
		Between(0, posR-1, h) {
		walked = true
		wg.Add(1)
		go Walk(passedPoints.Copy(), posR-1, posC, w, h, wg, ch)
	}

	if !walked {
		// There are no way to walk any more

		// Let's check if this is an exit
		if (posR == 0 || // Bottom
			posR == h-1 || // Top
			posC == 0 || // Left
			posC == w-1) && // Right
			passedPoints.AllPassed(w, h) { // Check all fields are passed
			ch <- Point{posR, posC}
		}
	}

	// tell the WaitGroup that we are done
	wg.Done()

}
