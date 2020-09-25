package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/ajm188/go-jsonpb/gen"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var (
		flags flag.FlagSet

		unmarshal = flags.Bool("unmarshal", false, "Include unmarshaler implementations.")
		protobuf  = flags.String("protobuf", "google.golang.org", "which protobuf package to use, valid choices are google.golang.org or github.com")
	)

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		if *unmarshal {
			log.Print("(UNSUPPORTED) told to unmarshal as well. This will be implemented later.")
		}

		opts := gen.Options{}

		switch *protobuf {
		case "google", "google.golang.org", "google.golang.org/protobuf":
			opts.GithubProtobuf = false
		case "github", "github.com", "github.com/golang", "github.com/golang/protobuf":
			opts.GithubProtobuf = true
		default:
			return errors.New("invalid choice for -protobuf")
		}

		for _, f := range plugin.Files {
			filename := fmt.Sprintf("%s_json.pb.go", f.GeneratedFilenamePrefix)
			g := plugin.NewGeneratedFile(filename, f.GoImportPath)

			gen.GenerateFile(g, plugin.Request, f, opts)
		}
		return nil
	})
}
