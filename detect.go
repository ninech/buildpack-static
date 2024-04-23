package static

import (
	"errors"
	"os"
	"path/filepath"

	require "github.com/ninech/buildpack-static-require"
	"github.com/paketo-buildpacks/libnodejs"
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

		ok, nodeAppPath, err := isNodeApp(context.WorkingDir)
		if err != nil {
			return packit.DetectResult{}, err
		}
		if ok {
			webRoot = nodeWebRoot(nodeAppPath)
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

// isNodeApp checks if there is a package.json in the "projectPath". By
// default that is the workspace but it could also be overridden using
// BP_NODE_PROJECT_PATH.
func isNodeApp(workingDir string) (bool, string, error) {
	projectPath, err := libnodejs.FindProjectPath(workingDir)
	if err != nil {
		return false, "", err
	}

	if _, err := libnodejs.ParsePackageJSON(projectPath); err == nil {
		relativePath, err := filepath.Rel(workingDir, projectPath)
		return true, relativePath, err
	} else {
		if errors.Is(err, os.ErrNotExist) {
			return false, "", nil
		}

		return false, "", err
	}
}

func nodeWebRoot(nodeAppPath string) string {
	// if explicitly specified, always use the WebRootEnv
	if path, ok := os.LookupEnv(require.WebRootEnv); ok {
		return path
	}

	// build is the default for react apps
	return filepath.Join(nodeAppPath, "build")
}
