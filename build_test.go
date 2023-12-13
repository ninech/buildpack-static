package static

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/nginx"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir string
		buffer     *bytes.Buffer
	)

	it.Before(func() {
		buffer = bytes.NewBuffer(nil)
	})

	context("when building", func() {
		it.Before(func() {
			var err error
			workingDir, err = os.MkdirTemp(t.TempDir(), "working-dir-*")
			Expect(err).NotTo(HaveOccurred())
		})

		it("generates an nginx.conf", func() {
			_, err := Build(scribe.NewEmitter(buffer))(packit.BuildContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())

			conf, err := os.ReadFile(filepath.Join(workingDir, nginx.ConfFile))
			Expect(err).NotTo(HaveOccurred())
			Expect(conf).NotTo(BeEmpty())
		})
	})

}
