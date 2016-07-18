package rtree

import (
	d1 "d1"
	d2 "d2"
	d3 "d3"
	d4 "d4"
	"math"
)

type Iterator func(item Item) bool
type Item interface {
	Rect(ctx interface{}) (min []float64, max []float64)
}

type RTree struct {
	ctx interface{}
	tr1 *d1.RTree
	tr2 *d2.RTree
	tr3 *d3.RTree
	tr4 *d4.RTree
}

func New(ctx interface{}) *RTree {
	return &RTree{
		ctx: ctx,
		tr1: d1.NewRTree(),
		tr2: d2.NewRTree(),
		tr3: d3.NewRTree(),
		tr4: d4.NewRTree(),
	}
}

func (tr *RTree) Insert(item Item) {
	if item == nil {
		panic("nil item being added to RTree")
	}
	min, max := item.Rect(tr.ctx)
	if len(min) != len(max) {
		panic("invalid item rectangle")
	}
	switch len(min) {
	default:
		panic("invalid dimension")
	case 1:
		var amin, amax [1]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr1.Insert(amin, amax, item)
	case 2:
		var amin, amax [2]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr2.Insert(amin, amax, item)
	case 3:
		var amin, amax [3]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr3.Insert(amin, amax, item)
	case 4:
		var amin, amax [4]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr4.Insert(amin, amax, item)
	}
}

func (tr *RTree) Remove(item Item) {
	if item == nil {
		panic("nil item being added to RTree")
	}
	min, max := item.Rect(tr.ctx)
	if len(min) != len(max) {
		panic("invalid item rectangle")
	}
	switch len(min) {
	default:
		panic("invalid dimension")
	case 1:
		var amin, amax [1]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr1.Remove(amin, amax, item)
	case 2:
		var amin, amax [2]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr2.Remove(amin, amax, item)
	case 3:
		var amin, amax [3]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr3.Remove(amin, amax, item)
	case 4:
		var amin, amax [4]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr4.Remove(amin, amax, item)
	}
}
func (tr *RTree) Reset() {
	tr.tr1 = d1.NewRTree()
	tr.tr2 = d2.NewRTree()
	tr.tr3 = d3.NewRTree()
	tr.tr4 = d4.NewRTree()
}
func (tr *RTree) Count() int {
	return tr.tr1.Count() + tr.tr2.Count() + tr.tr3.Count() + tr.tr4.Count()
}
func (tr *RTree) Search(bounds Item, iter Iterator) {
	if bounds == nil {
		panic("nil item being added to RTree")
	}
	min, max := bounds.Rect(tr.ctx)
	if len(min) != len(max) {
		panic("invalid item rectangle")
	}
	switch len(min) {
	default:
		panic("invalid dimension")
	case 1, 2, 3, 4:
	}
	if !tr.search1(min, max, iter) {
		return
	}
	if !tr.search2(min, max, iter) {
		return
	}
	if !tr.search3(min, max, iter) {
		return
	}
	if !tr.search4(min, max, iter) {
		return
	}
}

func (tr *RTree) search1(min, max []float64, iter Iterator) bool {
	var amin, amax [1]float64
	for i := 0; i < 1; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr1.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}
func (tr *RTree) search2(min, max []float64, iter Iterator) bool {
	var amin, amax [2]float64
	for i := 0; i < 2; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr2.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}
func (tr *RTree) search3(min, max []float64, iter Iterator) bool {
	var amin, amax [3]float64
	for i := 0; i < 3; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr3.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}
func (tr *RTree) search4(min, max []float64, iter Iterator) bool {
	var amin, amax [4]float64
	for i := 0; i < 4; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr4.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}
