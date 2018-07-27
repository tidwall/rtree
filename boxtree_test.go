package boxtree

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/tidwall/lotsa"
)

func TestBoxTree(t *testing.T) {
	New(2)
	New(3)
	defer func() {
		s := recover().(string)
		if s != "invalid dimensions" {
			t.Fatalf("expected '%s', got '%s'", "invalid dimensions", s)
		}
	}()
	New(4)
	// there are more test in the d2/d3 directories
}
func TestBenchInsert2D(t *testing.T) {
	testBenchInsert(t, 100000, 2)
}

func TestBenchInsert3D(t *testing.T) {
	testBenchInsert(t, 100000, 3)
}

func testBenchInsert(t *testing.T, N, D int) {
	rand.Seed(time.Now().UnixNano())
	points := make([]float64, N*D)
	for i := 0; i < N; i++ {
		for j := 0; j < D; j++ {
			points[i*D+j] = rand.Float64()*100 - 50
		}
	}
	tr := New(D)
	lotsa.Output = os.Stdout
	fmt.Printf("Insert(%dD): ", D)
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Insert(points[i*D+0:i*D+D], nil, i)
	})
	fmt.Printf("Search(%dD): ", D)
	var count int
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Search(points[i*D+0:i*D+D], points[i*D+0:i*D+D],
			func(min, max []float64, value interface{}) bool {
				count++
				return true
			},
		)
	})
	if count != N {
		t.Fatalf("expected %d, got %d", N, count)
	}
	fmt.Printf("Delete(%dD): ", D)
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Delete(points[i*D+0:i*D+D], points[i*D+0:i*D+D], i)
	})
	if tr.Count() != 0 {
		t.Fatalf("expected %d, got %d", N, tr.Count())
	}
}

type tItem2 struct {
	point [2]float64
}

func (item *tItem2) Point() (x, y float64) {
	return item.point[0], item.point[1]
}
func (item *tItem2) Rect() (minX, minY, maxX, maxY float64) {
	return item.point[0], item.point[1], item.point[0], item.point[1]
}

///////////////////////////////////////////////
// Old Tile38 Index < July 27, 2018
///////////////////////////////////////////////
// func TestBenchInsert2D_Old(t *testing.T) {
// // import "github.com/tidwall/tile38/pkg/index"
// 	N := 100000
// 	D := 2
// 	rand.Seed(time.Now().UnixNano())
// 	items := make([]*tItem2, N*D)
// 	for i := 0; i < N; i++ {
// 		items[i] = new(tItem2)
// 		for j := 0; j < D; j++ {
// 			items[i].point[j] = rand.Float64()*100 - 50
// 		}
// 	}

// 	tr := index.New()
// 	lotsa.Output = os.Stdout
// 	fmt.Printf("Insert(%dD): ", D)
// 	lotsa.Ops(N, 1, func(i, _ int) {
// 		tr.Insert(items[i])
// 	})
// 	fmt.Printf("Search(%dD): ", D)
// 	var count int
// 	lotsa.Ops(N, 1, func(i, _ int) {
// 		tr.Search(
// 			items[i].point[0], items[i].point[1],
// 			items[i].point[0], items[i].point[1],
// 			func(_ interface{}) bool {
// 				count++
// 				return true
// 			},
// 		)
// 	})
// 	if count != N {
// 		t.Fatalf("expected %d, got %d", N, count)
// 	}
// 	fmt.Printf("Delete(%dD): ", D)
// 	lotsa.Ops(N, 1, func(i, _ int) {
// 		tr.Remove(items[i])
// 	})
// 	if tr.Count() != 0 {
// 		t.Fatalf("expected %d, got %d", N, tr.Count())
// 	}

// }
