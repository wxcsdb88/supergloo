package main

import (
	"fmt"
	"os"
	"time"

	"github.com/solo-io/supergloo/cli/pkg/cmd"
	"github.com/solo-io/supergloo/cli/pkg/util"
)

var Version = "0.0.1"

func main() {
	start := time.Now()
	defer util.Telemetry(Version, start)

	app := cmd.App(Version)
	if err := app.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
