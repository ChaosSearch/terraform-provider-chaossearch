package main

import (
	cs "cs-tf-provider/provider"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	log "github.com/sirupsen/logrus"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return cs.Provider()
		},
	})
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
