package main

import (
	"log"
	"sort"
	"sync"
)

type Point struct {
	Row, Column int
}

type PointSet struct {
	set           map[Point]struct{}
	width, height int
	mux           *sync.RWMutex
}

func NewPointSet(width, height int) PointSet {
	return PointSet{
		width:  width,
		height: height,
		set:    make(map[Point]struct{}),
		mux:    &sync.RWMutex{},
	}
}

func (s PointSet) Set(c Point) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if !Between(0, c.Column, s.width) || !Between(0, c.Row, s.height) {
		log.Fatalf(
			"Cannot set point (%d, %d) with w: %d and h: %d",
			c.Row, c.Column, s.width, s.height,
		)
	}
	s.set[c] = struct{}{}
}

func (s PointSet) Has(c Point) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	_, has := s.set[c]
	return has
}

func (s PointSet) Copy() PointSet {
	s.mux.RLock()
	defer s.mux.RUnlock()
	ns := make(map[Point]struct{})
	for k, v := range s.set {
		ns[k] = v
	}
	return PointSet{
		set:    ns,
		width:  s.width,
		height: s.height,
		mux:    &sync.RWMutex{},
	}
}

func (s PointSet) AllPassed(w, h int) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	// Not a safe implementation, it assumes that only Points within the range
	// are added
	return len(s.set) == w*h
}

func (s PointSet) ToSortedSliced() []Point {
	var sl PointSlice
	s.mux.RLock()
	for p := range s.set {
		sl = append(sl, p)
	}
	s.mux.RUnlock()
	sort.Sort(sl)
	return sl
}

type PointSlice []Point

func (s PointSlice) Len() int {
	return len(s)
}

func (s PointSlice) Less(i, j int) bool {
	if s[i].Row == s[j].Row {
		return s[i].Column < s[j].Column
	}
	return s[i].Row < s[j].Row
}

func (s PointSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
