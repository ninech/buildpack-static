package static

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitBuildpackStatic(t *testing.T) {
	suite := spec.New("buildpack-static-confgen", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Detect", testDetect)
	suite("Build", testBuild)
	suite.Run(t)
}
