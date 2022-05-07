package rtree

import (
	"bytes"
	"encoding/gob"
	"testing"
)

var (
	boxen = [][2][2]float64{
		{{310, 320}, {330, 340}},
		{{110, 210}, {120, 220}},
		{{130, 240}, {140, 260}},
		{{150, 260}, {160, 280}},
		{{210, 220}, {330, 440}},
		{{220, 230}, {230, 240}},
		{{1, 1}, {666, 999}},
		{{1, 1}, {999, 666}},
		{{500, 500}, {600, 600}},
		{{600, 600}, {700, 700}},
		{{100, 200}, {300, 400}},
		{{100, 200}, {300, 400}},
	}
)

func exampleGeneric() Generic[int] {
	var tr Generic[int]
	for i, pair := range boxen {
		tr.Insert(pair[0], pair[1], i+1)
	}
	return tr
}

// transcode makes a copy to validate if serialized data is usable
func transcode[T any](in T) T {
	var buf bytes.Buffer
	var dupe T
	if err := gob.NewEncoder(&buf).Encode(in); err != nil {
		panic(err)
	}
	if err := gob.NewDecoder(&buf).Decode(&dupe); err != nil {
		panic(err)
	}
	return dupe
}

func TestReadOnlySearch(t *testing.T) {
	tr := exampleGeneric()
	min := [2]float64{225, 233}
	max := [2]float64{227, 235}
	const expected = 6 // min,max is within 5 of the sample rects
	var count int
	fn := func(min, max [2]float64, value int) bool {
		count++
		t.Logf("check %v -> %v (%d)\n", min, max, value)
		return true
	}

	t.Logf("\n\nsearching (%d) %v,%v\nsearch original generics first\n\n", len(boxen), min, max)
	tr.Search(min, max, fn)
	if count != expected {
		t.Errorf("want %d hits -- have %d", expected, count)
	}

	t.Logf("\n\nread-only version\n\n")
	fake := NewReadOnly(tr)
	count = 0
	fake.Search(min, max, fn)
	if count != expected {
		t.Errorf("want %d hits -- have %d", expected, count)
	}

	t.Logf("\n\nserialized version\n\n")
	dupe := transcode(fake)
	count = 0
	dupe.Search(min, max, fn)
	if count != expected {
		t.Errorf("want %d hits -- have %d", expected, count)
	}
}

func TestReadOnlyScan(t *testing.T) {
	tr := exampleGeneric()
	t.Logf("scanning (%d)\n", len(boxen))
	var count int
	fn := func(min, max [2]float64, value int) bool {
		count++
		t.Logf("check %v -> %v (%d)\n", min, max, value)
		return true
	}
	tr.Scan(fn)
	if count != len(boxen) {
		t.Errorf("want %d rects -- have %d", len(boxen), count)
	}
}
