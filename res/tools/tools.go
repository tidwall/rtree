package tools

import (
	"fmt"
	"strconv"
	"strings"
)

// RTree interface
type RTree interface {
	Insert(min, max []float64, value interface{})
	Scan(iter func(min, max []float64, value interface{}) bool)
	Search(min, max []float64, iter func(min, max []float64,
		value interface{}) bool)
	Delete(min, max []float64, value interface{})
	Traverse(iter func(min, max []float64, height, level int,
		value interface{}) int)
	Count() int
	TotalOverlapArea() float64
	Nearby(min, max []float64, iter func(min, max []float64,
		item interface{}) bool)
}

func svg(min, max []float64, height int) string {
	var out string
	point := true
	for i := 0; i < 2; i++ {
		if min[i] != max[i] {
			point = false
			break
		}
	}
	if point { // is point
		out += fmt.Sprintf(
			"<rect x=\"%.0f\" y=\"%0.f\" width=\"%0.f\" height=\"%0.f\" "+
				"stroke=\"%s\" fill=\"purple\" "+
				"fill-opacity=\"0\" stroke-opacity=\"1\" "+
				"rx=\"15\" ry=\"15\"/>\n",
			(min[0])*svgScale,
			(min[1])*svgScale,
			(max[0]-min[0]+1/svgScale)*svgScale,
			(max[1]-min[1]+1/svgScale)*svgScale,
			strokes[height%len(strokes)])
	} else { // is rect
		out += fmt.Sprintf(
			"<rect x=\"%.0f\" y=\"%0.f\" width=\"%0.f\" height=\"%0.f\" "+
				"stroke=\"%s\" fill=\"purple\" "+
				"fill-opacity=\"0\" stroke-opacity=\"1\"/>\n",
			(min[0])*svgScale,
			(min[1])*svgScale,
			(max[0]-min[0]+1/svgScale)*svgScale,
			(max[1]-min[1]+1/svgScale)*svgScale,
			strokes[height%len(strokes)])
	}
	return out
}

const (
	// Continue to first child rectangle and/or next sibling.
	Continue = iota
	// Ignore child rectangles but continue to next sibling.
	Ignore
	// Stop iterating
	Stop
)

const svgScale = 4.0

var strokes = [...]string{"black", "#cccc00", "green", "red", "purple"}

// SVG prints 2D rtree in wgs84 coordinate space
func SVG(tr RTree) string {
	var out string
	out += fmt.Sprintf("<svg viewBox=\"%.0f %.0f %.0f %.0f\" "+
		"xmlns =\"http://www.w3.org/2000/svg\">\n",
		-190.0*svgScale, -100.0*svgScale,
		380.0*svgScale, 190.0*svgScale)

	out += fmt.Sprintf("<g transform=\"scale(1,-1)\">\n")
	var outb []byte
	tr.Traverse(func(min, max []float64, height, level int, _ interface{}) int {
		outb = append(outb, svg(min, max, height)...)
		return Continue
	})
	out += string(outb)
	out += fmt.Sprintf("</g>\n")
	out += fmt.Sprintf("</svg>\n")
	return out
}

// Cities returns big list of cities base on json from
// https://github.com/lutangar/cities.json
func Cities(bigJSON string) [][2]float64 {
	var out [][2]float64
	s := bigJSON
	for i := 0; ; i++ {
		idx := strings.Index(s, `"lat": "`)
		if idx == -1 {
			break
		}
		s = s[idx+8:]
		idx = strings.IndexByte(s, '"')
		lat, _ := strconv.ParseFloat(s[:idx], 64)
		idx = strings.Index(s, `"lng": "`)
		s = s[idx+8:]
		idx = strings.IndexByte(s, '"')
		lng, _ := strconv.ParseFloat(s[:idx], 64)
		s = s[idx+1:]
		out = append(out, [2]float64{lng, lat})
	}
	return out
}
