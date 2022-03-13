// Copyright 2021 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package rtree

import (
	"errors"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/tidwall/geoindex"
)

func init() {
	seed := time.Now().UnixNano()
	println("seed:", seed)
	rand.Seed(seed)
}

func TestGeoIndex(t *testing.T) {
	t.Run("BenchVarious", func(t *testing.T) {
		geoindex.Tests.TestBenchVarious(t, &RTree{}, 1000000)
	})
	t.Run("RandomRects", func(t *testing.T) {
		geoindex.Tests.TestRandomRects(t, &RTree{}, 10000)
	})
	t.Run("RandomPoints", func(t *testing.T) {
		geoindex.Tests.TestRandomPoints(t, &RTree{}, 10000)
	})
	t.Run("ZeroPoints", func(t *testing.T) {
		geoindex.Tests.TestZeroPoints(t, &RTree{})
	})
	t.Run("CitiesSVG", func(t *testing.T) {
		geoindex.Tests.TestCitiesSVG(t, &RTree{})
	})
}

func BenchmarkRandomInsert(b *testing.B) {
	geoindex.Tests.BenchmarkRandomInsert(b, &RTree{})
}

func TestSane(t *testing.T) {
	if os.Getenv("SANETEST") != "1" {
		println("Use SANETEST=1 to run sane tester")
		return
	}
	println("Running Sane Test... (Press Ctrl-C to cancel)")
	rng := rand.New(rand.NewSource(0))
	N := 100_000
	points := make([][2]float64, N)
	start := time.Now()
	for time.Since(start) < time.Second*1000 {
		seed := time.Now().UnixNano()
		// Add the error seed below
		// seed = 1624721877993099000
		rng.Seed(seed)

		n := rng.Intn(N)
		if n%2 == 1 {
			n++
		}
		points[0][0] = 360*rng.Float64() - 180
		points[0][1] = 180*rng.Float64() - 90
		for i := 1; i < n; i++ {
			points[i][0] = points[i-1][0] + rng.Float64() - 0.5
			points[i][1] = points[i-1][1] + rng.Float64() - 0.5
		}
		var tr Generic[any]
		for i := 0; i < n-1; i += 2 {
			minx := points[i+0][0]
			miny := points[i+0][1]
			maxx := points[i+1][0]
			maxy := points[i+1][1]
			if minx > maxx {
				minx, maxx = maxx, minx
			}
			if miny > maxy {
				miny, maxy = maxy, miny
			}
			tr.Insert([2]float64{minx, miny}, [2]float64{maxx, maxy}, i)
		}
		if err := rSane(&tr); err != nil {
			t.Fatalf("rtree not sane (phase 1): %s (seed: %d)", err, seed)
		}

		// Delete half the items
		for i := 2; i < n-1; i += 4 {
			minx := points[i+0][0]
			miny := points[i+0][1]
			maxx := points[i+1][0]
			maxy := points[i+1][1]
			if minx > maxx {
				minx, maxx = maxx, minx
			}
			if miny > maxy {
				miny, maxy = maxy, miny
			}
			tr.Delete([2]float64{minx, miny}, [2]float64{maxx, maxy}, i)
		}
		if err := rSane(&tr); err != nil {
			t.Fatalf("rtree not sane (phase 2): %s (seed: %d)", err, seed)
		}

	}
}

func rSane(tr *Generic[any]) error {
	height := tr.height
	if height > 0 && tr.root.data == nil {
		return errors.New("not nil root")
	} else if tr.count > 0 && tr.root.data == nil {
		return errors.New("invalid count")
	} else if tr.count == 0 && tr.root.data != nil {
		return errors.New("invalid count")
	} else if len(tr.reinsert) > 0 {
		return errors.New("items in reinsert")
	}
	if tr.root.data == nil {
		return nil
	}
	return rSaneNode(tr, &tr.root, height)
}

func rSaneRect[T any](r rect[T]) error {
	if r.min[0] > r.max[0] || r.min[1] > r.max[1] {
		return errors.New("invalid rect")
	}
	return nil
}

func rSaneNode[T any](tr *Generic[any], r *rect[T], height int) error {
	n := r.data.(*node[T])
	if n.count >= maxEntries {
		return errors.New("invalid count")
	}
	for i := 0; i < n.count; i++ {
		if err := rSaneRect(n.rects[i]); err != nil {
			return err
		}
		if !r.contains(&n.rects[i]) {
			return errors.New("not contains item")
		}
		if n.rects[i].data == nil {
			return errors.New("nil data")
		}
	}
	for i := n.count; i < len(n.rects); i++ {
		if n.rects[i].data != nil {
			return errors.New("not nil data")
		}
	}
	if height > 0 {
		for i := 0; i < n.count; i++ {
			if err := rSaneNode(tr, &n.rects[i], height-1); err != nil {
				return err
			}
		}
	}
	return nil
}
