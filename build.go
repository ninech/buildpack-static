package static

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/paketo-buildpacks/nginx"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/paketo-buildpacks/packit/v2/servicebindings"
)

func Build(logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		var webRoot string
		for _, entry := range context.Plan.Entries {
			if v, ok := entry.Metadata[webRootKey]; ok {
				webRoot, ok = v.(string)
				if !ok {
					return packit.BuildResult{}, fmt.Errorf("%s in metadata is not a string", webRootKey)
				}
				break
			}
		}
		webRoot = filepath.Join(context.WorkingDir, webRoot)
		logger.Process("%s is set to %s", webRootKey, webRoot)

		nginxConf := filepath.Join(context.WorkingDir, nginx.ConfFile)
		if _, err := os.Stat(nginxConf); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// we re-use LoadConfiguration from the nginx buildpack to support some of the env variables.
				cfg, err := nginx.LoadConfiguration(os.Environ(), servicebindings.NewResolver(), os.Getenv("CNB_PLATFORM_DIR"))
				if err != nil {
					return packit.BuildResult{}, packit.Fail.WithMessage("unable to load config: %s", err)
				}
				cfg.NGINXConfLocation = nginxConf
				cfg.WebServerRoot = webRoot

				if err := NewDefaultConfigGenerator(logger).Generate(Configuration{
					// we set the last-modified header to the current time
					// during build. This works around the issue described in:
					// https://github.com/paketo-buildpacks/nginx/issues/447
					LastModifiedValue: time.Now().UTC().Format(http.TimeFormat),
					ETag:              false,
					// allow from Pod CIDR
					SetRealIPFrom: "10.42.0.0/16",
					Configuration: cfg,
				}); err != nil {
					return packit.BuildResult{}, packit.Fail.WithMessage("unable to create nginx.conf: %s", err)
				}
				logger.Process("created default nginx config %s", nginxConf)
				return packit.BuildResult{}, nil
			}
			return packit.BuildResult{}, err
		}
		return packit.BuildResult{}, nil
	}
}
