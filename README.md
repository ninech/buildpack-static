# buildpack-static-confgen

This buildpack builds upon the [Paketo nginx
Buildpack](https://github.com/paketo-buildpacks/nginx) to automate some of the
most common use-cases where the workspace just contains an `index.html` or
`public/index.html`. It will configure nginx accordingly without any required
parameters. It also supports compilation of JavaScript frontend apps when
grouped with the
[web-servers](https://github.com/paketo-buildpacks/web-servers) buildpack.

The `PORT` env variable is required to be set when launching the built
container as nginx will use it to configure the listening port.

As this buildpack writes the nginx config which needs to be present in the
build phase of the nginx buildpack, it needs to be ordered before that. Due to
this ordering it's not possible to require the nginx buildpack. For that we
need to use [`ninech/buildpack-static-require`](https://github.com/ninech/buildpack-static-require)
and order it after the nginx/web-servers buildpack.

To test the build locally, checkout this repository and then build it with:

```bash
pack build static --path ./integration/testdata/default_app/ \
  --builder paketobuildpacks/builder-jammy-base \
  --buildpack . \
  --buildpack paketo-buildpacks/web-servers \
  --buildpack ghcr.io/ninech/buildpack-static-require
```
