# terraform-provider-graphenedb

## Local dev

```sh
$ export TF_LOG_PROVIDER=DEBUG
$ make install
$ tf init && tf apply
```

## Release and publish on the registry

```sh
$ export GITHUB_TOKEN=...
$ git tag vX.Y.Z
$ git push origin vX.Y.Z
$ goreleaser release --rm-dist --skip-sign
# sign manually the SHA file
$ gpg --detach-sign dist/terraform-provider-graphenedb_X.Y.Z_SHA256SUMS
```
