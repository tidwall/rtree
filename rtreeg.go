// Copyright 2021 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package rtree

import (
	"sync/atomic"
	"unsafe"

	"github.com/tidwall/geoindex/child"
)

// SAFTEY: The unsafe package is used, but with care.
// Using "unsafe" allows for one alloction per node and avoids having to use
// an interface{} type for child nodes; that may either be:
//   - *leafNode[T]
//   - *branchNode[T]
// This library makes it generally safe by guaranteeing that all references to
// nodes are simply to `*node[T]`, which is just the header struct for the leaf
// or branch representation. The difference between a leaf and a branch node is
// that a leaf has an array of item data of generic type T on tail of the
// struct, while a branch has an array of child node pointers on the tail. To
// access the child items `node[T].items()` is called; returning a slice, or
// nil if the node is a branch. To access the child nodes `node[T].children()`
// is called; returning a slice, or nil if the node is a leaf. The `items()`
// and `children()` methods check the `node[T].kind` to determine which kind of
// node it is, which is an enum of `none`, `leaf`, or `branch`. The only valid
// way to create a `*node[T]` is `RTreeG[T].newNode(leaf bool)` which take a
// bool that indicates the new node kind is a `leaf` or `branch`.

const maxEntries = 64
const minEntries = maxEntries * 10 / 100
const orderBranches = true
const orderLeaves = true
const quickChooser = false

// copy-on-write atomic incrementer
var cow uint64

type RTreeG[T any] struct {
	cow      uint64
	count    int
	rect     rect
	root     *node[T]
	reinsert []reinsertItem[T]
	empty    T
}

type rect struct {
	min [2]float64
	max [2]float64
}

func (r *rect) expand(b *rect) {
	if b.min[0] < r.min[0] {
		r.min[0] = b.min[0]
	}
	if b.max[0] > r.max[0] {
		r.max[0] = b.max[0]
	}
	if b.min[1] < r.min[1] {
		r.min[1] = b.min[1]
	}
	if b.max[1] > r.max[1] {
		r.max[1] = b.max[1]
	}
}

type kind int8

const (
	none kind = iota
	leaf
	branch
)

type node[T any] struct {
	cow   uint64
	kind  kind
	count int16
	rects [maxEntries]rect
}

func (n *node[T]) leaf() bool {
	return n.kind == leaf
}

type leafNode[T any] struct {
	node[T]
	items [maxEntries]T
}

type branchNode[T any] struct {
	node[T]
	children [maxEntries]*node[T]
}

func (n *node[T]) children() []*node[T] {
	if n.kind != branch {
		// not a branch
		return nil
	}

	return (*branchNode[T])(unsafe.Pointer(n)).children[:]
}

func (n *node[T]) items() []T {
	if n.kind != leaf {
		// not a leaf
		return nil
	}
	return (*leafNode[T])(unsafe.Pointer(n)).items[:]
}

func (tr *RTreeG[T]) newNode(isleaf bool) *node[T] {
	if isleaf {
		n := &leafNode[T]{node: node[T]{cow: tr.cow, kind: leaf}}
		return (*node[T])(unsafe.Pointer(n))
	} else {
		n := &branchNode[T]{node: node[T]{cow: tr.cow, kind: branch}}
		return (*node[T])(unsafe.Pointer(n))
	}
}

func (n *node[T]) rect() rect {
	rect := n.rects[0]
	for i := 1; i < int(n.count); i++ {
		rect.expand(&n.rects[i])
	}
	return rect
}

// Insert data into tree
func (tr *RTreeG[T]) Insert(min, max [2]float64, data T) {
	ir := rect{min, max}
	if tr.root == nil {
		tr.root = tr.newNode(true)
		tr.rect = ir
	}
	grown := tr.nodeInsert(&tr.rect, &tr.root, &ir, data)
	split := tr.root.count == maxEntries
	if grown {
		tr.rect.expand(&ir)
	}
	if split {
		left := tr.root
		right := tr.splitNode(tr.rect, left)
		tr.root = tr.newNode(false)
		tr.root.rects[0] = left.rect()
		tr.root.rects[1] = right.rect()
		tr.root.children()[0] = left
		tr.root.children()[1] = right
		tr.root.count = 2
	}
	if orderBranches && !tr.root.leaf() && (grown || split) {
		tr.root.sort()
	}
	tr.count++
}

func (tr *RTreeG[T]) splitNode(r rect, left *node[T]) (right *node[T]) {
	return tr.splitNodeLargestAxisEdgeSnap(r, left)
}

func (n *node[T]) orderToRight(idx int) int {
	for idx < int(n.count)-1 && n.rects[idx+1].min[0] < n.rects[idx].min[0] {
		n.swap(idx+1, idx)
		idx++
	}
	return idx
}

func (n *node[T]) orderToLeft(idx int) int {
	for idx > 0 && n.rects[idx].min[0] < n.rects[idx-1].min[0] {
		n.swap(idx, idx-1)
		idx--
	}
	return idx
}

// This operation should not be inlined because it's expensive and rarely
// called outside of heavy copy-on-write situations. Marking it "noinline"
// allows for the parent cowLoad to be inlined.
// go:noinline
func (tr *RTreeG[T]) copy(n *node[T]) *node[T] {
	n2 := tr.newNode(n.leaf())
	*n2 = *n
	if n2.leaf() {
		copy(n2.items()[:n.count], n.items()[:n.count])
	} else {
		copy(n2.children()[:n.count], n.children()[:n.count])
	}
	return n2
}

// cowLoad loads the provided node and, if needed, performs a copy-on-write.
func (tr *RTreeG[T]) cowLoad(cn **node[T]) *node[T] {
	if (*cn).cow != tr.cow {
		*cn = tr.copy(*cn)
	}
	return *cn
}

func (n *node[T]) rsearch(key float64) int {
	rects := n.rects[:n.count]
	for i := 0; i < len(rects); i++ {
		if !(n.rects[i].min[0] < key) {
			return i
		}
	}
	return int(n.count)
}

func (n *node[T]) bsearch(key float64) int {
	low, high := 0, int(n.count)
	for low < high {
		h := int(uint(low+high) >> 1)
		if !(key < n.rects[h].min[0]) {
			low = h + 1
		} else {
			high = h
		}
	}
	return low
}

func (tr *RTreeG[T]) nodeInsert(nr *rect, cn **node[T], ir *rect, data T,
) (grown bool) {
	n := tr.cowLoad(cn)
	if n.leaf() {
		items := n.items()
		index := int(n.count)
		if orderLeaves {
			index = n.rsearch(ir.min[0])
			copy(n.rects[index+1:int(n.count)+1], n.rects[index:int(n.count)])
			copy(items[index+1:int(n.count)+1], items[index:int(n.count)])
		}
		n.rects[index] = *ir
		items[index] = data
		n.count++
		grown = !nr.contains(ir)
		return grown
	}

	// choose a subtree
	rects := n.rects[:n.count]
	index := -1
	narea := 0.0
	// take a quick look for any nodes that contain the rect
	for i := 0; i < len(rects); i++ {
		if rects[i].contains(ir) {
			if quickChooser {
				index = i
				break
			} else {
				area := rects[i].area()
				if index == -1 || area < narea {
					index = i
					narea = area
				}
			}
		}
	}
	if index == -1 {
		index = n.chooseLeastEnlargement(ir)
	}

	children := n.children()
	grown = tr.nodeInsert(&n.rects[index], &children[index], ir, data)
	split := children[index].count == maxEntries
	if grown {
		// The child rectangle must expand to accomadate the new item.
		n.rects[index].expand(ir)
		if orderBranches {
			index = n.orderToLeft(index)
		}
		grown = !nr.contains(ir)
	}
	if split {
		left := children[index]
		right := tr.splitNode(n.rects[index], left)
		n.rects[index] = left.rect()
		if orderBranches {
			copy(n.rects[index+2:int(n.count)+1],
				n.rects[index+1:int(n.count)])
			copy(children[index+2:int(n.count)+1],
				children[index+1:int(n.count)])
			n.rects[index+1] = right.rect()
			children[index+1] = right
			n.count++
			if n.rects[index].min[0] > n.rects[index+1].min[0] {
				n.swap(index+1, index)
			}
			index++
			index = n.orderToRight(index)
		} else {
			n.rects[n.count] = right.rect()
			children[n.count] = right
			n.count++
		}

	}
	return grown
}

func (r *rect) area() float64 {
	return (r.max[0] - r.min[0]) * (r.max[1] - r.min[1])
}

// contains return struct when b is fully contained inside of n
func (r *rect) contains(b *rect) bool {
	if b.min[0] < r.min[0] || b.max[0] > r.max[0] {
		return false
	}
	if b.min[1] < r.min[1] || b.max[1] > r.max[1] {
		return false
	}
	return true
}

// intersects returns true if both rects intersect each other.
func (r *rect) intersects(b *rect) bool {
	if b.min[0] > r.max[0] || b.max[0] < r.min[0] {
		return false
	}
	if b.min[1] > r.max[1] || b.max[1] < r.min[1] {
		return false
	}
	return true
}

func (n *node[T]) chooseLeastEnlargement(ir *rect) (index int) {
	rects := n.rects[:int(n.count)]
	j, jenlargement, jarea := -1, 0.0, 0.0
	for i := 0; i < len(rects); i++ {
		// calculate the enlarged area
		uarea := rects[i].unionedArea(ir)
		area := rects[i].area()
		enlargement := uarea - area
		if j == -1 || enlargement < jenlargement ||
			(!(enlargement > jenlargement) && area < jarea) {
			j, jenlargement, jarea = i, enlargement, area
		}
	}
	return j
}

func fmin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
func fmax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// unionedArea returns the area of two rects expanded
func (r *rect) unionedArea(b *rect) float64 {
	return (fmax(r.max[0], b.max[0]) - fmin(r.min[0], b.min[0])) *
		(fmax(r.max[1], b.max[1]) - fmin(r.min[1], b.min[1]))
}

func (r rect) largestAxis() (axis int) {
	if r.max[1]-r.min[1] > r.max[0]-r.min[0] {
		return 1
	}
	return 0
}

func (tr *RTreeG[T]) splitNodeLargestAxisEdgeSnap(r rect, left *node[T],
) (right *node[T]) {
	axis := r.largestAxis()
	right = tr.newNode(left.leaf())
	for i := 0; i < int(left.count); i++ {
		minDist := left.rects[i].min[axis] - r.min[axis]
		maxDist := r.max[axis] - left.rects[i].max[axis]
		if minDist < maxDist {
			// stay left
		} else {
			// move to right
			tr.moveRectAtIndexInto(left, i, right)
			i--
		}
	}
	// Make sure that both left and right nodes have at least
	// minEntries by moving items into underflowed nodes.
	if left.count < minEntries {
		// reverse sort by min axis
		right.sortByAxis(axis, true, false)
		for left.count < minEntries {
			tr.moveRectAtIndexInto(right, int(right.count)-1, left)
		}
	} else if right.count < minEntries {
		// reverse sort by max axis
		left.sortByAxis(axis, true, true)
		for right.count < minEntries {
			tr.moveRectAtIndexInto(left, int(left.count)-1, right)
		}
	}

	if (orderBranches && !right.leaf()) || (orderLeaves && right.leaf()) {
		right.sort()
		// It's not uncommon that the left node is already ordered
		if !left.issorted() {
			left.sort()
		}
	}
	return right
}

func (tr *RTreeG[T]) moveRectAtIndexInto(from *node[T], index int,
	into *node[T],
) {
	into.rects[into.count] = from.rects[index]
	from.rects[index] = from.rects[from.count-1]
	if from.leaf() {
		into.items()[into.count] = from.items()[index]
		from.items()[index] = from.items()[from.count-1]
		from.items()[from.count-1] = tr.empty
	} else {
		into.children()[into.count] = from.children()[index]
		from.children()[index] = from.children()[from.count-1]
		from.children()[from.count-1] = nil
	}
	from.count--
	into.count++
}

func (n *node[T]) search(target rect,
	iter func(min, max [2]float64, data T) bool,
) bool {
	rects := n.rects[:n.count]
	if n.leaf() {
		items := n.items()
		for i := 0; i < len(rects); i++ {
			if rects[i].intersects(&target) {
				if !iter(rects[i].min, rects[i].max, items[i]) {
					return false
				}
			}
		}
		return true
	}
	children := n.children()
	for i := 0; i < len(rects); i++ {
		if target.intersects(&rects[i]) {
			if !children[i].search(target, iter) {
				return false
			}
		}
	}
	return true
}

// Len returns the number of items in tree
func (tr *RTreeG[T]) Len() int {
	return tr.count
}

// Search for items in tree that intersect the provided rectangle
func (tr *RTreeG[T]) Search(min, max [2]float64,
	iter func(min, max [2]float64, data T) bool,
) {
	target := rect{min, max}
	if tr.root == nil {
		return
	}
	if target.intersects(&tr.rect) {
		tr.root.search(target, iter)
	}
}

// Scane all items in the tree
func (tr *RTreeG[T]) Scan(iter func(min, max [2]float64, data T) bool) {
	if tr.root != nil {
		tr.root.scan(iter)
	}
}

func (n *node[T]) scan(iter func(min, max [2]float64, data T) bool) bool {
	if n.leaf() {
		for i := 0; i < int(n.count); i++ {
			if !iter(n.rects[i].min, n.rects[i].max, n.items()[i]) {
				return false
			}
		}
	} else {
		for i := 0; i < int(n.count); i++ {
			if !n.children()[i].scan(iter) {
				return false
			}
		}
	}
	return true
}

// Copy the tree.
// This is a copy-on-write operation and is very fast because it only performs
// a shadowed copy.
func (tr *RTreeG[T]) Copy() *RTreeG[T] {
	tr2 := new(RTreeG[T])
	*tr2 = *tr
	tr.cow = atomic.AddUint64(&cow, 1)
	tr.reinsert = nil
	tr2.cow = atomic.AddUint64(&cow, 1)
	tr2.reinsert = nil
	return tr2
}

// swap two rectanlges
func (n *node[T]) swap(i, j int) {
	n.rects[i], n.rects[j] = n.rects[j], n.rects[i]
	if n.leaf() {
		n.items()[i], n.items()[j] = n.items()[j], n.items()[i]
	} else {
		n.children()[i], n.children()[j] = n.children()[j], n.children()[i]
	}
}

func (n *node[T]) sortByAxis(axis int, rev, max bool) {
	n.qsort(0, int(n.count), axis, rev, max)
}

func (n *node[T]) sort() {
	n.qsort(0, int(n.count), 0, false, false)
}

func (n *node[T]) issorted() bool {
	rects := n.rects[:n.count]
	for i := 1; i < len(rects); i++ {
		if rects[i].min[0] < rects[i-1].min[0] {
			return false
		}
	}
	return true
}

func (n *node[T]) qsort(s, e int, axis int, rev, max bool) {
	nrects := e - s
	if nrects < 2 {
		return
	}
	left, right := 0, nrects-1
	pivot := nrects / 2 // rand and mod not worth it
	n.swap(s+pivot, s+right)
	rects := n.rects[s:e]
	if !rev {
		if !max {
			for i := 0; i < len(rects); i++ {
				if rects[i].min[axis] < rects[right].min[axis] {
					n.swap(s+i, s+left)
					left++
				}
			}
		} else {
			for i := 0; i < len(rects); i++ {
				if rects[i].max[axis] < rects[right].max[axis] {
					n.swap(s+i, s+left)
					left++
				}
			}
		}
	} else {
		if !max {
			for i := 0; i < len(rects); i++ {
				if rects[right].min[axis] < rects[i].min[axis] {
					n.swap(s+i, s+left)
					left++
				}
			}
		} else {
			for i := 0; i < len(rects); i++ {
				if rects[right].max[axis] < rects[i].max[axis] {
					n.swap(s+i, s+left)
					left++
				}
			}
		}
	}
	n.swap(s+left, s+right)
	n.qsort(s, s+left, axis, rev, max)
	n.qsort(s+left+1, e, axis, rev, max)
}

// Delete data from tree
func (tr *RTreeG[T]) Delete(min, max [2]float64, data T) {
	tr.delete(min, max, data)
}

func (tr *RTreeG[T]) delete(min, max [2]float64, data T) bool {
	ir := rect{min, max}
	if tr.root == nil || !tr.rect.contains(&ir) {
		return false
	}
	removed, _ := tr.nodeDelete(&tr.rect, &tr.root, &ir, data)
	if !removed {
		return false
	}
	tr.count -= len(tr.reinsert) + 1
	if tr.count == 0 {
		tr.root = nil
		tr.rect.min = [2]float64{0, 0}
		tr.rect.max = [2]float64{0, 0}
	} else {
		for !tr.root.leaf() && tr.root.count == 1 {
			tr.root = tr.root.children()[0]
		}
	}
	if len(tr.reinsert) > 0 {
		var def reinsertItem[T]
		for i, item := range tr.reinsert {
			tr.Insert(item.rect.min, item.rect.max, item.data)
			tr.reinsert[i] = def
		}
		tr.reinsert = tr.reinsert[:0]
		if cap(tr.reinsert) > maxEntries {
			tr.reinsert = nil
		}
	}
	return true
}

func compare[T any](a, b T) bool {
	return (interface{})(a) == (interface{})(b)
}

func (tr *RTreeG[T]) nodeDelete(nr *rect, cn **node[T], ir *rect, data T,
) (removed, shrunk bool) {
	n := tr.cowLoad(cn)
	rects := n.rects[:n.count]
	if n.leaf() {
		items := n.items()
		for i := 0; i < len(rects); i++ {
			if ir.contains(&rects[i]) && compare(items[i], data) {
				// found the target item to delete
				if orderLeaves {
					copy(n.rects[i:n.count], n.rects[i+1:n.count])
					copy(items[i:n.count], items[i+1:n.count])
				} else {
					n.rects[i] = n.rects[n.count-1]
					items[i] = items[n.count-1]
				}
				items[len(rects)-1] = tr.empty
				n.count--
				shrunk = ir.onedge(nr)
				if shrunk {
					*nr = n.rect()
				}
				return true, shrunk
			}
		}
		return false, false
	}
	children := n.children()
	for i := 0; i < len(rects); i++ {
		if !rects[i].contains(ir) {
			continue
		}
		crect := rects[i]
		removed, shrunk = tr.nodeDelete(&rects[i], &children[i], ir, data)
		if !removed {
			continue
		}
		if children[i].count < minEntries {
			tr.reinsert = children[i].flatten(tr.reinsert)
			if orderBranches {
				copy(n.rects[i:n.count], n.rects[i+1:n.count])
				copy(children[i:n.count], children[i+1:n.count])
			} else {
				n.rects[i] = n.rects[n.count-1]
				children[i] = children[n.count-1]
			}
			children[n.count-1] = nil
			n.count--
			*nr = n.rect()
			return true, true
		}
		if shrunk {
			shrunk = !rects[i].equals(&crect)
			if shrunk {
				*nr = n.rect()
			}
			if orderBranches {
				i = n.orderToRight(i)
			}
		}
		return true, shrunk
	}
	return false, false
}

func (r *rect) equals(b *rect) bool {
	return !(r.min[0] < b.min[0] || r.min[0] > b.min[0] ||
		r.min[1] < b.min[1] || r.min[1] > b.min[1] ||
		r.max[0] < b.max[0] || r.max[0] > b.max[0] ||
		r.max[1] < b.max[1] || r.max[1] > b.max[1])
}

type reinsertItem[T any] struct {
	rect rect
	data T
}

func (n *node[T]) flatten(reinsert []reinsertItem[T]) []reinsertItem[T] {
	if n.leaf() {
		for i := 0; i < int(n.count); i++ {
			reinsert = append(reinsert, reinsertItem[T]{
				rect: n.rects[i],
				data: n.items()[i],
			})
		}
	} else {
		for i := 0; i < int(n.count); i++ {
			reinsert = n.children()[i].flatten(reinsert)
		}
	}
	return reinsert
}

// onedge returns true when r is on the edge of b
func (r *rect) onedge(b *rect) bool {
	return !(r.min[0] > b.min[0] && r.min[1] > b.min[1] &&
		r.max[0] < b.max[0] && r.max[1] < b.max[1])
}

// Replace an item.
// If the old item does not exist then the new item is not inserted.
func (tr *RTreeG[T]) Replace(
	oldMin, oldMax [2]float64, oldData T,
	newMin, newMax [2]float64, newData T,
) {
	if tr.delete(oldMin, oldMax, oldData) {
		tr.Insert(newMin, newMax, newData)
	}
}

// Bounds returns the minimum bounding rect
func (tr *RTreeG[T]) Bounds() (min, max [2]float64) {
	return tr.rect.min, tr.rect.max
}

// Children is a utility function that returns all children for parent node.
// If parent node is nil then the root nodes should be returned. The min, max,
// data, and items slices all must have the same lengths. And, each element
// from all slices must be associated. Returns true for `items` when the the
// item at the leaf level. The reuse buffers are empty length slices that can
// optionally be used to avoid extra allocations.
func (tr *RTreeG[T]) Children(parent interface{}, reuse []child.Child,
) (children []child.Child) {
	children = reuse
	if parent == nil {
		if tr.Len() > 0 {
			// fill with the root
			children = append(children, child.Child{
				Min:  tr.rect.min,
				Max:  tr.rect.max,
				Data: tr.root,
				Item: false,
			})
		}
	} else {
		// fill with child items
		n := parent.(*node[T])
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

// Generic RTree
// Deprecated: use RTreeG
type Generic[T any] struct {
	RTreeG[T]
}

func (tr *Generic[T]) Copy() *Generic[T] {
	return &Generic[T]{*tr.RTreeG.Copy()}
}
