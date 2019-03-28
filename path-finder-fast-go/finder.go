package main

import (
	"context"
	"sync"
)

func ConsiderExitPoint(
	exit Point,         // the exit that we are heading for
	ch chan<- Point,    // the channel to sent Point to if we found the path
	w, h int,           // width and height of the whole field
	start Point,        // the point we are starting at
	wg *sync.WaitGroup, // the sync.WaitGroup used for wait all goroutine
) {

	// make sure that we have tell the WaitGroup we are done before we exit
	defer wg.Done()

	// prepare a context that is to be cancelled as we found a path for such
	// exit
	ctx, onFound := context.WithCancel(context.TODO())

	// prepare the PointSet for all passed points
	passedPoints := NewPointSet(w, h)

	wg.Add(1)
	go WalkToExit(passedPoints, exit, start, w, h, ctx, onFound, wg, ch)
}

func WalkToExit(
	passedPoints PointSet,      // the set of points that are already passed
	exit Point,                 // the exit we are targeting for
	pos Point,                  // current position
	w, h int,                   // width and height of fields
	ctx context.Context,        // stop here if cancelled
	onFound context.CancelFunc, // used for cancel the context
	wg *sync.WaitGroup,         // use this to wait for every goroutine
	ch chan<- Point,            // use this to signal path found
) {

	// Make sure that WaitGroup will be done before we return
	defer wg.Done()

	DebugLog("at (%d, %d), exit at (%d, %d)", pos.Row, pos.Column,
		exit.Row, exit.Column)

	// Check if ctx is already cancelled
	if ctx.Err() != nil {

		// Already cancelled
		return
	}

	// We assumes that passedPoints is a modifiable copy
	// Add current position
	passedPoints.Set(pos)

	// Check if we arrived target exit
	if pos == exit {
		// Arrived target exit

		// Check if all points are already walked and ctx is not cancelled
		if passedPoints.AllPassed() && ctx.Err() == nil {
			// Nice, this path leads us to the exit

			onFound() // Cancel the context so other goroutines won't emit
			ch <- pos // Emit this exit
			return    // Go back prevent further action
		}

		// Oh, this path stepped on exit without passing all other cells!
		// Or someone (goroutine) already found this exit
		// Anyway we have to return
		return
	}

	// Test we have anyway to walk
	walked := false

	// Consider walking left
	if !passedPoints.Has(pos.Left()) && pos.Left().Valid(w, h) {

		// We can walk, so mark walk as true
		walked = true

		// Tell the WaitGroup to wait for next goroutine
		wg.Add(1)

		// Make a copy of passedPoints and move pos to left, and pass it
		go WalkToExit(passedPoints.Copy(), exit, pos.Left(), w, h, ctx, onFound,
			wg, ch)
	}

	// Consider walking right
	if !passedPoints.Has(pos.Right()) && pos.Right().Valid(w, h) {

		// We can walk, so mark walk as true
		walked = true

		// Tell the WaitGroup to wait for next goroutine
		wg.Add(1)

		// Make a copy of passedPoints and move pos to left, and pass it
		go WalkToExit(passedPoints.Copy(), exit, pos.Right(), w, h, ctx,
			onFound, wg, ch)
	}

	// Consider walking up
	if !passedPoints.Has(pos.Up()) && pos.Up().Valid(w, h) {

		// We can walk, so mark walk as true
		walked = true

		// Tell the WaitGroup to wait for next goroutine
		wg.Add(1)

		// Make a copy of passedPoints and move pos to left, and pass it
		go WalkToExit(passedPoints.Copy(), exit, pos.Up(), w, h, ctx, onFound,
			wg, ch)
	}

	// Consider walking down
	if !passedPoints.Has(pos.Down()) && pos.Down().Valid(w, h) {

		// We can walk, so mark walk as true
		walked = true

		// Tell the WaitGroup to wait for next goroutine
		wg.Add(1)

		// Make a copy of passedPoints and move pos to left, and pass it
		go WalkToExit(passedPoints.Copy(), exit, pos.Down(), w, h, ctx, onFound,
			wg, ch)
	}

	// Actually this is useless since we already checked exit,
	// just to suppress warning
	if !walked {
		return
	}
}
