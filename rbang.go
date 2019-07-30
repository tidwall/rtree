package rbang

import (
	"github.com/tidwall/rbang-go/d2"
	"github.com/tidwall/rbang-go/d3"
)

// RTree is an rtree
type RTree interface {
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
func New(dims int) RTree {
	switch dims {
	default:
		panic("invalid dimensions")
	case 2:
		return new(d2.RTree)
	case 3:
		return new(d3.RTree)
	}
}
