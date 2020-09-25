package main

import (
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
	)

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		if *unmarshal {
			log.Print("told to unmarshal as well")
		}
		for _, f := range plugin.Files {
			filename := fmt.Sprintf("%s_json.pb.go", f.GeneratedFilenamePrefix)
			g := plugin.NewGeneratedFile(filename, f.GoImportPath)

			gen.GenerateFile(g, plugin.Request, f)
		}
		return nil
	})
}
