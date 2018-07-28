package main

import (
	"io/ioutil"

	"github.com/tidwall/boxtree"
	"github.com/tidwall/boxtree/res/tools"
	"github.com/tidwall/cities"
)

func main() {
	tr := boxtree.New(2)
	for _, city := range cities.Cities {
		tr.Insert([]float64{city.Longitude, city.Latitude}, nil, &city)
	}
	ioutil.WriteFile("cities.svg", []byte(tools.SVG(tr)), 0600)
}
