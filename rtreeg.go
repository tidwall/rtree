// Copyright 2021 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package rtree

import (
	"github.com/tidwall/geoindex/child"
)

type RTreeG[T any] struct {
	base RTreeGN[float64, T]
}

// Insert data into tree
func (tr *RTreeG[T]) Insert(min, max [2]float64, data T) {
	tr.base.Insert(min, max, data)
}

// Len returns the number of items in tree
func (tr *RTreeG[T]) Len() int {
	return tr.base.Len()
}

// Search for items in tree that intersect the provided rectangle
func (tr *RTreeG[T]) Search(min, max [2]float64,
	iter func(min, max [2]float64, data T) bool,
) {
	tr.base.Search(min, max, iter)
}

// Scan all items in the tree
func (tr *RTreeG[T]) Scan(iter func(min, max [2]float64, data T) bool) {
	tr.base.Scan(iter)
}

// Copy the tree.
// This is a copy-on-write operation and is very fast because it only performs
// a shadowed copy.
func (tr *RTreeG[T]) Copy() *RTreeG[T] {
	return &RTreeG[T]{*tr.base.Copy()}
}

// Delete data from tree
func (tr *RTreeG[T]) Delete(min, max [2]float64, data T) {
	tr.base.Delete(min, max, data)
}

// Replace an item.
// If the old item does not exist then the new item is not inserted.
func (tr *RTreeG[T]) Replace(
	oldMin, oldMax [2]float64, oldData T,
	newMin, newMax [2]float64, newData T,
) {
	tr.base.Replace(
		oldMin, oldMax, oldData,
		newMin, newMax, newData,
	)
}

// Bounds returns the minimum bounding rect
func (tr *RTreeG[T]) Bounds() (min, max [2]float64) {
	return tr.base.Bounds()
}

// children is a utility function that returns all children for parent node.
// If parent node is nil then the root nodes should be returned. The min, max,
// data, and items slices all must have the same lengths. And, each element
// from all slices must be associated. Returns true for `items` when the the
// item at the leaf level. The reuse buffers are empty length slices that can
// optionally be used to avoid extra allocations.
func (tr *RTreeG[T]) children(parent interface{}, reuse []child.Child,
) (children []child.Child) {
	children = reuse
	if parent == nil {
		if tr.Len() > 0 {
			// fill with the root
			children = append(children, child.Child{
				Min:  tr.base.rect.min,
				Max:  tr.base.rect.max,
				Data: tr.base.root,
				Item: false,
			})
		}
	} else {
		// fill with child items
		n := parent.(*node[float64, T])
		for i := 0; i < int(n.count); i++ {
			c := child.Child{
				Min: n.rects[i].min, Max: n.rects[i].max, Item: n.leaf(),
			}
			if c.Item {
				c.Data = n.items()[i]
			} else {
				c.Data = n.children()[i]
			}
			children = append(children, c)
		}
	}
	return children
}

// Nearby performs a kNN-type operation on the index.
// It's expected that the caller provides its own the `dist` function, which
// is used to calculate a distance to rectangles and data.
// The `iter` function will return all items from the smallest distance to the
// largest distance.
//
// BoxDist is included with this package for simple box-distance
// calculations. For example, say you want to return the closest items to
// Point(10 20):
//
//	tr.Nearby(
//		rtree.BoxDist([2]float64{10, 20}, [2]float64{10, 20}, nil),
//		func(min, max [2]float64, data int, dist float64) bool {
//			return true
//		},
//	)
func (tr *RTreeG[T]) Nearby(
	dist func(min, max [2]float64, data T, item bool) float64,
	iter func(min, max [2]float64, data T, dist float64) bool,
) {
	tr.base.Nearby(dist, iter)
}

// Clear will delete all items.
func (tr *RTreeG[T]) Clear() {
	tr.base.Clear()
}

// Generic RTree
// Deprecated: use RTreeG
type Generic[T any] struct {
	RTreeG[T]
}

func (tr *Generic[T]) Copy() *Generic[T] {
	return &Generic[T]{*tr.RTreeG.Copy()}
}
