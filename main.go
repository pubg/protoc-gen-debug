package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	var flags flag.FlagSet
	opts := protogen.Options{
		ParamFunc: flags.Set,
	}

	dumpFile := flags.String("dump_file", "request.pb.bin", `dump file path`)
	parameter := flags.String("parameter", "", `parameter`)

	opts.Run(func(plugin *protogen.Plugin) error {
		plugin.Request.Parameter = parameter
		buf, err := proto.Marshal(plugin.Request)
		if err != nil {
			return err
		}
		_, err = plugin.NewGeneratedFile(*dumpFile, "").Write(buf)
		if err != nil {
			return err
		}

		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL) | uint64(pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
		return nil
	})
}
