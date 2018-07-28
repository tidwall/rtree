package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	b, err := ioutil.ReadFile("d2/boxtree.go")
	if err != nil {
		log.Fatal(err)
	}
	s := string(b)
	s = strings.Replace(s, "package d2", "package d3", 1)
	s = strings.Replace(s, "const dims = 2", "const dims = 3", 1)
	if err := os.MkdirAll("d3", 0777); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("d3/boxtree.go", []byte(s), 0666); err != nil {
		log.Fatal(err)
	}
	b, err = ioutil.ReadFile("d2/boxtree_test.go")
	if err != nil {
		log.Fatal(err)
	}
	b = []byte(strings.Replace(string(b), "package d2", "package d3", 1))
	if err := ioutil.WriteFile("d3/boxtree_test.go", b, 0666); err != nil {
		log.Fatal(err)
	}

}
