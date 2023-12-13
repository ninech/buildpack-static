package static

import (
	require "github.com/ninech/buildpack-static-require"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

const (
	indexFile  = "index.html"
	webRootKey = "web-root"
	name       = "static-confgen"
)

type BuildPlanMetadata struct {
	Version       string `toml:"version,omitempty"`
	VersionSource string `toml:"version-source,omitempty"`
	Launch        bool   `toml:"launch"`
	WebRoot       string `toml:"web-root,omitempty"`
}

func Detect(logger scribe.Emitter) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		webRoot, err := require.WebRoot(context.WorkingDir)
		if err != nil {
			return packit.DetectResult{}, err
		}

		result := packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{{Name: name}},
				Requires: []packit.BuildPlanRequirement{
					{
						Name:     name,
						Metadata: BuildPlanMetadata{WebRoot: webRoot},
					},
				},
			},
		}

		return result, nil
	}
}
