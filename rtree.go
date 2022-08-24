// Copyright 2021 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package rtree

import "github.com/tidwall/geoindex/child"

type RTree struct {
	base RTreeG[any]
}

// Insert an item into the structure
func (tr *RTree) Insert(min, max [2]float64, data interface{}) {
	tr.base.Insert(min, max, data)
}

// Delete an item from the structure
func (tr *RTree) Delete(min, max [2]float64, data interface{}) {
	tr.base.Delete(min, max, data)
}

// Replace an item in the structure. This is effectively just a Delete
// followed by an Insert. But for some structures it may be possible to
// optimize the operation to avoid multiple passes
func (tr *RTree) Replace(
	oldMin, oldMax [2]float64, oldData interface{},
	newMin, newMax [2]float64, newData interface{},
) {
	tr.base.Replace(
		oldMin, oldMax, oldData,
		newMin, newMax, newData,
	)
}

// Search the structure for items that intersects the rect param
func (tr *RTree) Search(
	min, max [2]float64,
	iter func(min, max [2]float64, data interface{}) bool,
) {
	tr.base.Search(min, max, iter)
}

// Scan iterates through all data in tree in no specified order.
func (tr *RTree) Scan(iter func(min, max [2]float64, data interface{}) bool) {
	tr.base.Scan(iter)
}

// Len returns the number of items in tree
func (tr *RTree) Len() int {
	return tr.base.Len()
}

// Bounds returns the minimum bounding box
func (tr *RTree) Bounds() (min, max [2]float64) {
	return tr.base.Bounds()
}

// Children returns all children for parent node. If parent node is nil
// then the root nodes should be returned.
// The reuse buffer is an empty length slice that can optionally be used
// to avoid extra allocations.
func (tr *RTree) Children(parent interface{}, reuse []child.Child) (children []child.Child) {
	return tr.base.Children(parent, reuse)
}

// Copy the tree.
// This is a copy-on-write operation and is very fast because it only performs
// a shadowed copy.
func (tr *RTree) Copy() *RTree {
	return &RTree{base: *tr.base.Copy()}
}
