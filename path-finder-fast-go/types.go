package main

import (
	"log"
	"sort"
	//"sync"
)

type Point struct {
	Row, Column int
}

func (p Point) Left() Point {
	return Point{Row: p.Row, Column: p.Column - 1}
}

func (p Point) Right() Point {
	return Point{Row: p.Row, Column: p.Column + 1}
}

func (p Point) Up() Point {
	return Point{Row: p.Row + 1, Column: p.Column}
}

func (p Point) Down() Point {
	return Point{Row: p.Row - 1, Column: p.Column}
}

func (p Point) Valid(w, h int) bool {
	return 0 <= p.Row && p.Row < h && 0 <= p.Column && p.Column < w
}

type PointSet struct {
	setMap        map[Point]struct{}
	width, height int
	//mux           *sync.RWMutex
}

func NewPointSet(width, height int) PointSet {
	return PointSet{
		width:  width,
		height: height,
		setMap: make(map[Point]struct{}),
		//mux:    &sync.RWMutex{},
	}
}

func (s PointSet) Set(c Point) {
	//s.mux.Lock()
	//defer s.mux.Unlock()
	if !Between(0, c.Column, s.width) || !Between(0, c.Row, s.height) {
		log.Fatalf(
			"Cannot set point (%d, %d) with w: %d and h: %d",
			c.Row, c.Column, s.width, s.height,
		)
	}
	s.setMap[c] = struct{}{}
}

func (s PointSet) Has(c Point) bool {
	//s.mux.RLock()
	//defer s.mux.RUnlock()
	_, has := s.setMap[c]
	return has
}

func (s PointSet) Copy() PointSet {
	//s.mux.RLock()
	//defer s.mux.RUnlock()
	ns := make(map[Point]struct{})
	for k, v := range s.setMap {
		ns[k] = v
	}
	return PointSet{
		setMap: ns,
		width:  s.width,
		height: s.height,
		//mux:    &sync.RWMutex{},
	}
}

func (s PointSet) AllPassed() bool {
	//s.mux.RLock()
	//defer s.mux.RUnlock()
	// Not a safe implementation, it assumes that only Points within the range
	// are added
	return len(s.setMap) == s.width*s.height
}

func (s PointSet) ToSortedSliced() []Point {
	var sl PointSlice
	//s.mux.RLock()
	for p := range s.setMap {
		sl = append(sl, p)
	}
	//s.mux.RUnlock()
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
