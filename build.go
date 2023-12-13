package static

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/nginx"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
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
				confGen := nginx.NewDefaultConfigGenerator(logger)
				if err := confGen.Generate(nginx.Configuration{
					NGINXConfLocation: nginxConf,
					WebServerRoot:     webRoot,
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
