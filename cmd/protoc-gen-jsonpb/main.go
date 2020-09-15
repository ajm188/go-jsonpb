package main

import (
	"io/ioutil"
	"os"

	"github.com/ajm188/go-jsonpb/gen"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func getSrc(filename string) string {
	f, err := os.Open(filename)
	must(err)

	defer f.Close()

	src, err := ioutil.ReadAll(f)
	must(err)

	return string(src)
}

func getDest(filename string) *os.File {
	f, err := os.Create(filename)

	must(err)
	must(f.Chmod(0644))

	return f
}

func main() {
	srcname := os.Args[1]
	destname := os.Args[2]
	src := getSrc(srcname)
	dest := getDest(destname)

	defer dest.Close()

	must(gen.Generate(src, dest))
}
