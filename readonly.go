// Copyright 2021 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package rtree

import (
	"log"
)

type RRect[T any] struct {
	Value T
	Node  RNode
	Min   [2]float64
	Max   [2]float64
}

type RNode struct {
	Start int
	End   int
}

type ReadOnly[T any] struct {
	Rects []RRect[T]
	Root  RRect[T]
	Count int
}

// contains return struct when b is fully contained inside of n
func (r *RRect[T]) contains(b *RRect[T]) bool {
	if b.Min[0] < r.Min[0] || b.Max[0] > r.Max[0] {
		return false
	}
	if b.Min[1] < r.Min[1] || b.Max[1] > r.Max[1] {
		return false
	}
	return true
}

// contains return struct when b is fully contained inside of n
func (r *RRect[T]) intersects(b *RRect[T]) bool {
	if b.Min[0] > r.Max[0] || b.Max[0] < r.Min[0] {
		return false
	}
	if b.Min[1] > r.Max[1] || b.Max[1] < r.Min[1] {
		return false
	}
	//fmt.Printf("INTERSECTION %v,%v => %v,%v\n", r.Min, r.Max, b.Min, b.Max)
	return true
}

// returning false indicates search is complete
func (tr ReadOnly[T]) search(
	in, target RRect[T], child bool,
	iter func(min, max [2]float64, value T) bool,
) bool {
	n := in.Node

	if !target.intersects(&in) {
		return true
	}
	for i := n.Start; i < n.End; i++ {
		r := tr.Rects[i]
		if target.intersects(&r) {
			if tr.search(r, target, true, iter) {
				return false
			}
		}
	}
	if child {
		return !iter(in.Min, in.Max, in.Value)
	}
	return true
}

func (tr *ReadOnly[T]) Scan(
	iter func(min, max [2]float64, value T) bool,
) {
	tr.scan(tr.Root, iter)
}

func (tr *ReadOnly[T]) scan(
	r RRect[T],
	iter func(min, max [2]float64, value T) bool,
) bool {
	for i := r.Node.Start; i < r.Node.End; i++ {
		r = tr.Rects[i]
		if !iter(r.Min, r.Max, r.Value) {
			return false
		}
	}

	return true
}

// Search ...
func (tr *ReadOnly[T]) Search(
	min, max [2]float64,
	iter func(min, max [2]float64, value T) bool,
) {
	what := RRect[T]{Min: min, Max: max}
	tr.search(tr.Root, what, false, iter)
}

// Len returns the number of items in tree
func (tr *ReadOnly[T]) Len() int {
	return tr.Count
}

// Bounds returns the minimum bounding rect
func (tr *ReadOnly[T]) Bounds() (min, max [2]float64) {
	return tr.Root.Min, tr.Root.Max
}

func (tr *ReadOnly[T]) DupeNode(from *node[T]) RNode {
	start := len(tr.Rects)
	// prefill as we're doing depth first traversal
	for i := 0; i < from.count; i++ {
		tr.Rects = append(tr.Rects, RRect[T]{})
	}
	for i := 0; i < from.count; i++ {
		r := from.rects[i]
		k := start + i
		tr.Rects[k] = tr.DupeRect(r)
	}
	return RNode{
		Start: start,
		End:   start + from.count,
	}
}

func (tr *ReadOnly[T]) DupeRect(from rect[T]) RRect[T] {
	r := RRect[T]{
		Min: from.min,
		Max: from.max,
	}
	switch n := from.data.(type) {
	case T:
		r.Value = n
	case *node[T]:
		r.Node = tr.DupeNode(n)
	default:
		log.Fatalf("invalid type: %T (%+v)", from.data, from.data)
	}
	return r
}

func NewReadOnly[T any](in Generic[T]) ReadOnly[T] {
	rt := ReadOnly[T]{
		Count: in.count,
		Rects: make([]RRect[T], 0, in.count),
	}
	rt.Root = rt.DupeRect(in.root)
	return rt
}
