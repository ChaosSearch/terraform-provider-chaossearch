package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"

	// goreleaser can also pass the specific commit if you want
	// commit  string = ""
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return Provider()
		},
	}

	if debugMode {
		// TODO: update this string with the full name of your provider as used in your configs
		err := plugin.Debug(context.Background(), "registry.terraform.io/hashicorp/scaffolding", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)

	//plugin.Serve(&plugin.ServeOpts{
	//	ProviderFunc: func() *schema.Provider {
	//		return Provider()
	//	},
	//})
}

func init() {
	fileName, ok := os.LookupEnv("TF_LOG_PATH")
	if ok {

		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		Formatter := new(log.TextFormatter)

		Formatter.TimestampFormat = "02-01-2006 15:04:05"
		Formatter.FullTimestamp = true
		log.SetFormatter(Formatter)
		if err != nil {
			fmt.Println(err)
		} else {
			log.SetOutput(f)
		}
		lvl, ok := os.LookupEnv("TF_LOG")

		if !ok {
			lvl = "warn"
		}
		level := lvl
		ll, err := log.ParseLevel(level)
		if err != nil {
			ll = log.DebugLevel
		}
		log.SetLevel(ll)
	}

}
