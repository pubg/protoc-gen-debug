package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

type PluginOptions struct {
	DumpBinary bool
	DumpJson   bool
	DumpText   bool
	FileBinary string
	FileJson   string
	FileText   string
	Parameter  string
}

var version string = "develop"

func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "--version" {
			fmt.Println(version)
		} else if os.Args[1] == "--help" {
			fmt.Println("USAGE:")
			fmt.Println("  protoc-gen-debug --version  : print version")
		}
		return
	}

	var flags flag.FlagSet
	opts := protogen.Options{
		ParamFunc: flags.Set,
	}

	options := PluginOptions{}
	flags.BoolVar(&options.DumpBinary, "dump_binary", true, "Enable dump protobuf request as binary format")
	flags.BoolVar(&options.DumpJson, "dump_json", false, "Enable dump protobuf request as json format")
	flags.BoolVar(&options.DumpText, "dump_text", false, "Enable dump protobuf request as text format")
	flags.StringVar(&options.FileBinary, "file_binary", "request.pb.bin", "binary file path")
	flags.StringVar(&options.FileJson, "file_json", "request.pb.json", "json file path")
	flags.StringVar(&options.FileText, "file_text", "request.pb.txt", "text file path")
	flags.StringVar(&options.Parameter, "parameter", "", "parameter")

	opts.Run(func(plugin *protogen.Plugin) error {
		normalizedParam := strings.Replace(options.Parameter, ":", ",", -1)
		plugin.Request.Parameter = &normalizedParam

		if options.DumpBinary {
			buf, err := proto.Marshal(plugin.Request)
			if err != nil {
				return err
			}
			_, err = plugin.NewGeneratedFile(options.FileBinary, "").Write(buf)
			if err != nil {
				return err
			}

		}

		if options.DumpJson {
			buf, err := protojson.Marshal(plugin.Request)
			if err != nil {
				return err
			}
			_, err = plugin.NewGeneratedFile(options.FileJson, "").Write(buf)
			if err != nil {
				return err
			}
		}

		if options.DumpText {
			buf, err := prototext.Marshal(plugin.Request)
			if err != nil {
				return err
			}
			_, err = plugin.NewGeneratedFile(options.FileText, "").Write(buf)
			if err != nil {
				return err
			}
		}

		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL) | uint64(pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
		return nil
	})
}
