package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/ajm188/go-jsonpb/gen"
	"github.com/ajm188/go-jsonpb/template"
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

const modeHelp = `mode to run in. options are ("inplace", "template").
In inplace mode, protoc-gen-jsonpb will modify the AST directly to inject json.Marshaler
implementations. In this mode, you can specify the same source and dest, and the
file will be rewritten to include those functions in addition to the rest of the file.

In template mode, protoc-gen-jsonpb will only use the AST to figure out which types need
json.Marshaler implementations. It then uses text/template to generate a completely
separate file (in the same package as src) containing those implementations. This mode
is **not safe** to run where src == dest, since it would result in the original
protoc-generated code would be overwritten. The tool will refuse to operate in this
situation.
`

func main() {
	srcfile := flag.String("src", "", "go file to use as source for transformation")
	destfile := flag.String("dest", "", "destination to write transformed file. blank to write to stdout")
	mode := flag.String("mode", "inplace", modeHelp)
	debug := flag.Bool("debug", false, "include debug logs")

	flag.Parse()

	setLogger(*debug)

	if *srcfile == "" {
		log.SetOutput(os.Stderr)
		log.Fatal("must pass -src")
	}

	if *mode != "template" && *mode != "inplace" {
		log.SetOutput(os.Stderr)
		log.Fatalf(`invalid option for -mode: %s. must be one of "inplace", "template"`, *mode)
	}

	if *mode == "template" && *destfile != "" && *srcfile == *destfile {
		log.SetOutput(os.Stderr)
		log.Fatal("in template mode, -src cannot be the same as -dest. See the help text of -mode for details")
	}

	src := getSrc(*srcfile)
	dest := getDest(*destfile)

	defer dest.Close()

	switch *mode {
	case "inplace":
		must(gen.Generate(src, dest))
	case "template":
		must(template.Generate(src, dest))
	}
}
