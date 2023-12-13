package main

import (
	"os"

	confgen "github.com/ninech/buildpack-static-confgen"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

type BuildPlanMetadata struct {
	Version       string `toml:"version,omitempty"`
	VersionSource string `toml:"version-source,omitempty"`
	Launch        bool   `toml:"launch"`
	PublicDir     string `toml:"public-dir,omitempty"`
}

func main() {
	logger := scribe.NewEmitter(os.Stdout).WithLevel(os.Getenv("BP_LOG_LEVEL"))
	packit.Run(confgen.Detect(logger), confgen.Build(logger))
}
