// generated; DO NOT EDIT!
package rtree

import (
	"math"

	d1 "github.com/tidwall/rtree/dims/d1"
	d2 "github.com/tidwall/rtree/dims/d2"
	d3 "github.com/tidwall/rtree/dims/d3"
	d4 "github.com/tidwall/rtree/dims/d4"
	d5 "github.com/tidwall/rtree/dims/d5"
	d6 "github.com/tidwall/rtree/dims/d6"
	d7 "github.com/tidwall/rtree/dims/d7"
	d8 "github.com/tidwall/rtree/dims/d8"
	d9 "github.com/tidwall/rtree/dims/d9"
	d10 "github.com/tidwall/rtree/dims/d10"
	d11 "github.com/tidwall/rtree/dims/d11"
	d12 "github.com/tidwall/rtree/dims/d12"
	d13 "github.com/tidwall/rtree/dims/d13"
	d14 "github.com/tidwall/rtree/dims/d14"
	d15 "github.com/tidwall/rtree/dims/d15"
	d16 "github.com/tidwall/rtree/dims/d16"
	d17 "github.com/tidwall/rtree/dims/d17"
	d18 "github.com/tidwall/rtree/dims/d18"
	d19 "github.com/tidwall/rtree/dims/d19"
	d20 "github.com/tidwall/rtree/dims/d20"
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
	tr5 *d5.RTree
	tr6 *d6.RTree
	tr7 *d7.RTree
	tr8 *d8.RTree
	tr9 *d9.RTree
	tr10 *d10.RTree
	tr11 *d11.RTree
	tr12 *d12.RTree
	tr13 *d13.RTree
	tr14 *d14.RTree
	tr15 *d15.RTree
	tr16 *d16.RTree
	tr17 *d17.RTree
	tr18 *d18.RTree
	tr19 *d19.RTree
	tr20 *d20.RTree
}

func New(ctx interface{}) *RTree {
	return &RTree{
		ctx: ctx,
		tr1: d1.NewRTree(),
		tr2: d2.NewRTree(),
		tr3: d3.NewRTree(),
		tr4: d4.NewRTree(),
		tr5: d5.NewRTree(),
		tr6: d6.NewRTree(),
		tr7: d7.NewRTree(),
		tr8: d8.NewRTree(),
		tr9: d9.NewRTree(),
		tr10: d10.NewRTree(),
		tr11: d11.NewRTree(),
		tr12: d12.NewRTree(),
		tr13: d13.NewRTree(),
		tr14: d14.NewRTree(),
		tr15: d15.NewRTree(),
		tr16: d16.NewRTree(),
		tr17: d17.NewRTree(),
		tr18: d18.NewRTree(),
		tr19: d19.NewRTree(),
		tr20: d20.NewRTree(),
	}
}

func (tr *RTree) Insert(item Item) {
	if item == nil {
		panic("nil item being added to RTree")
	}
	min, max := item.Rect(tr.ctx)
	if len(min) != len(max) {
		return // just return
		panic("invalid item rectangle")
	}
	switch len(min) {
	default:
		return // just return
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
	case 5:
		var amin, amax [5]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr5.Insert(amin, amax, item)
	case 6:
		var amin, amax [6]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr6.Insert(amin, amax, item)
	case 7:
		var amin, amax [7]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr7.Insert(amin, amax, item)
	case 8:
		var amin, amax [8]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr8.Insert(amin, amax, item)
	case 9:
		var amin, amax [9]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr9.Insert(amin, amax, item)
	case 10:
		var amin, amax [10]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr10.Insert(amin, amax, item)
	case 11:
		var amin, amax [11]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr11.Insert(amin, amax, item)
	case 12:
		var amin, amax [12]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr12.Insert(amin, amax, item)
	case 13:
		var amin, amax [13]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr13.Insert(amin, amax, item)
	case 14:
		var amin, amax [14]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr14.Insert(amin, amax, item)
	case 15:
		var amin, amax [15]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr15.Insert(amin, amax, item)
	case 16:
		var amin, amax [16]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr16.Insert(amin, amax, item)
	case 17:
		var amin, amax [17]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr17.Insert(amin, amax, item)
	case 18:
		var amin, amax [18]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr18.Insert(amin, amax, item)
	case 19:
		var amin, amax [19]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr19.Insert(amin, amax, item)
	case 20:
		var amin, amax [20]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr20.Insert(amin, amax, item)
	}
}

func (tr *RTree) Remove(item Item) {
	if item == nil {
		panic("nil item being added to RTree")
	}
	min, max := item.Rect(tr.ctx)
	if len(min) != len(max) {
		return // just return
		panic("invalid item rectangle")
	}
	switch len(min) {
	default:
		return // just return
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
	case 5:
		var amin, amax [5]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr5.Remove(amin, amax, item)
	case 6:
		var amin, amax [6]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr6.Remove(amin, amax, item)
	case 7:
		var amin, amax [7]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr7.Remove(amin, amax, item)
	case 8:
		var amin, amax [8]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr8.Remove(amin, amax, item)
	case 9:
		var amin, amax [9]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr9.Remove(amin, amax, item)
	case 10:
		var amin, amax [10]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr10.Remove(amin, amax, item)
	case 11:
		var amin, amax [11]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr11.Remove(amin, amax, item)
	case 12:
		var amin, amax [12]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr12.Remove(amin, amax, item)
	case 13:
		var amin, amax [13]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr13.Remove(amin, amax, item)
	case 14:
		var amin, amax [14]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr14.Remove(amin, amax, item)
	case 15:
		var amin, amax [15]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr15.Remove(amin, amax, item)
	case 16:
		var amin, amax [16]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr16.Remove(amin, amax, item)
	case 17:
		var amin, amax [17]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr17.Remove(amin, amax, item)
	case 18:
		var amin, amax [18]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr18.Remove(amin, amax, item)
	case 19:
		var amin, amax [19]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr19.Remove(amin, amax, item)
	case 20:
		var amin, amax [20]float64
		for i := 0; i < len(min); i++ {
			amin[i], amax[i] = min[i], max[i]
		}
		tr.tr20.Remove(amin, amax, item)
	}
}
func (tr *RTree) Reset() {
	tr.tr1 = d1.NewRTree()
	tr.tr2 = d2.NewRTree()
	tr.tr3 = d3.NewRTree()
	tr.tr4 = d4.NewRTree()
	tr.tr5 = d5.NewRTree()
	tr.tr6 = d6.NewRTree()
	tr.tr7 = d7.NewRTree()
	tr.tr8 = d8.NewRTree()
	tr.tr9 = d9.NewRTree()
	tr.tr10 = d10.NewRTree()
	tr.tr11 = d11.NewRTree()
	tr.tr12 = d12.NewRTree()
	tr.tr13 = d13.NewRTree()
	tr.tr14 = d14.NewRTree()
	tr.tr15 = d15.NewRTree()
	tr.tr16 = d16.NewRTree()
	tr.tr17 = d17.NewRTree()
	tr.tr18 = d18.NewRTree()
	tr.tr19 = d19.NewRTree()
	tr.tr20 = d20.NewRTree()
}
func (tr *RTree) Count() int {
	count := 0
	count += tr.tr1.Count()
	count += tr.tr2.Count()
	count += tr.tr3.Count()
	count += tr.tr4.Count()
	count += tr.tr5.Count()
	count += tr.tr6.Count()
	count += tr.tr7.Count()
	count += tr.tr8.Count()
	count += tr.tr9.Count()
	count += tr.tr10.Count()
	count += tr.tr11.Count()
	count += tr.tr12.Count()
	count += tr.tr13.Count()
	count += tr.tr14.Count()
	count += tr.tr15.Count()
	count += tr.tr16.Count()
	count += tr.tr17.Count()
	count += tr.tr18.Count()
	count += tr.tr19.Count()
	count += tr.tr20.Count()
	return count
}
func (tr *RTree) Search(bounds Item, iter Iterator) {
	if bounds == nil {
		panic("nil bounds being used for search")
	}
	min, max := bounds.Rect(tr.ctx)
	if len(min) != len(max) {
		return // just return
		panic("invalid item rectangle")
	}
	switch len(min) {
	default:
		return // just return
		panic("invalid dimension")
	case 1:
	case 2:
	case 3:
	case 4:
	case 5:
	case 6:
	case 7:
	case 8:
	case 9:
	case 10:
	case 11:
	case 12:
	case 13:
	case 14:
	case 15:
	case 16:
	case 17:
	case 18:
	case 19:
	case 20:
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
	if !tr.search5(min, max, iter) {
		return
	}
	if !tr.search6(min, max, iter) {
		return
	}
	if !tr.search7(min, max, iter) {
		return
	}
	if !tr.search8(min, max, iter) {
		return
	}
	if !tr.search9(min, max, iter) {
		return
	}
	if !tr.search10(min, max, iter) {
		return
	}
	if !tr.search11(min, max, iter) {
		return
	}
	if !tr.search12(min, max, iter) {
		return
	}
	if !tr.search13(min, max, iter) {
		return
	}
	if !tr.search14(min, max, iter) {
		return
	}
	if !tr.search15(min, max, iter) {
		return
	}
	if !tr.search16(min, max, iter) {
		return
	}
	if !tr.search17(min, max, iter) {
		return
	}
	if !tr.search18(min, max, iter) {
		return
	}
	if !tr.search19(min, max, iter) {
		return
	}
	if !tr.search20(min, max, iter) {
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

func (tr *RTree) search5(min, max []float64, iter Iterator) bool {
	var amin, amax [5]float64
	for i := 0; i < 5; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr5.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}

func (tr *RTree) search6(min, max []float64, iter Iterator) bool {
	var amin, amax [6]float64
	for i := 0; i < 6; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr6.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}

func (tr *RTree) search7(min, max []float64, iter Iterator) bool {
	var amin, amax [7]float64
	for i := 0; i < 7; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr7.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}

func (tr *RTree) search8(min, max []float64, iter Iterator) bool {
	var amin, amax [8]float64
	for i := 0; i < 8; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr8.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}

func (tr *RTree) search9(min, max []float64, iter Iterator) bool {
	var amin, amax [9]float64
	for i := 0; i < 9; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr9.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}

func (tr *RTree) search10(min, max []float64, iter Iterator) bool {
	var amin, amax [10]float64
	for i := 0; i < 10; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr10.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}

func (tr *RTree) search11(min, max []float64, iter Iterator) bool {
	var amin, amax [11]float64
	for i := 0; i < 11; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr11.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}

func (tr *RTree) search12(min, max []float64, iter Iterator) bool {
	var amin, amax [12]float64
	for i := 0; i < 12; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr12.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}

func (tr *RTree) search13(min, max []float64, iter Iterator) bool {
	var amin, amax [13]float64
	for i := 0; i < 13; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr13.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}

func (tr *RTree) search14(min, max []float64, iter Iterator) bool {
	var amin, amax [14]float64
	for i := 0; i < 14; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr14.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}

func (tr *RTree) search15(min, max []float64, iter Iterator) bool {
	var amin, amax [15]float64
	for i := 0; i < 15; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr15.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}

func (tr *RTree) search16(min, max []float64, iter Iterator) bool {
	var amin, amax [16]float64
	for i := 0; i < 16; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr16.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}

func (tr *RTree) search17(min, max []float64, iter Iterator) bool {
	var amin, amax [17]float64
	for i := 0; i < 17; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr17.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}

func (tr *RTree) search18(min, max []float64, iter Iterator) bool {
	var amin, amax [18]float64
	for i := 0; i < 18; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr18.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}

func (tr *RTree) search19(min, max []float64, iter Iterator) bool {
	var amin, amax [19]float64
	for i := 0; i < 19; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr19.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}

func (tr *RTree) search20(min, max []float64, iter Iterator) bool {
	var amin, amax [20]float64
	for i := 0; i < 20; i++ {
		if i < len(min) {
			amin[i] = min[i]
			amax[i] = max[i]
		} else {
			amin[i] = math.Inf(-1)
			amax[i] = math.Inf(+1)
		}
	}
	ended := false
	tr.tr20.Search(amin, amax, func(dataID interface{}) bool {
		if !iter(dataID.(Item)) {
			ended = true
			return false
		}
		return true
	})
	return !ended
}


