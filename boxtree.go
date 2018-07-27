package boxtree

import (
	"github.com/tidwall/boxtree/d2"
	"github.com/tidwall/boxtree/d3"
)

// BoxTree is an rtree by a different name
type BoxTree interface {
	Insert(min, max []float64, value interface{})
	Delete(min, max []float64, value interface{})
	Search(min, max []float64,
		iter func(min, max []float64, value interface{}) bool,
	)
	TotalOverlapArea() float64
	Traverse(iter func(min, max []float64, height, level int,
		value interface{}) int)
	Scan(iter func(min, max []float64, value interface{}) bool)
	Nearby(min, max []float64,
		iter func(min, max []float64, item interface{}) bool,
	)
	Bounds() (min, max []float64)
	Count() int
}

// New returns are new BoxTree, only 2 dims are allows
func New(dims int) BoxTree {
	switch dims {
	default:
		panic("invalid dimensions")
	case 2:
		return new(d2.BoxTree)
	case 3:
		return new(d3.BoxTree)
	}
}
