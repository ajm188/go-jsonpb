package main

import (
	"flag"
	"io/ioutil"
	"log"
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
	if filename == "" {
		log.Printf("-dest left blank, defaulting to stdout")
		return os.Stdout
	}

	f, err := os.Create(filename)

	must(err)
	must(f.Chmod(0644))

	return f
}

func setLogger(debug bool) {
	if debug {
		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stderr)

		return
	}

	log.SetOutput(ioutil.Discard)
}

func main() {
	srcfile := flag.String("src", "", "go file to use as source for transformation")
	destfile := flag.String("dest", "", "destination to write transformed file. blank to write to stdout")
	debug := flag.Bool("debug", false, "include debug logs")

	flag.Parse()

	setLogger(*debug)

	if *srcfile == "" {
		log.SetOutput(os.Stderr)
		log.Fatal("must pass -src")
	}

	src := getSrc(*srcfile)
	dest := getDest(*destfile)

	defer dest.Close()

	must(gen.Generate(src, dest))
}
