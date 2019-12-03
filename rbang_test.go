package rbang

import (
	"math/rand"
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
	t.Run("BenchInsert", func(t *testing.T) {
		geoindex.Tests.TestBenchInsert(t, &RTree{}, 100000)
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
