## Linting

Check the code with:

```
go tool golangci-lint run
```

## Release

Just tag a release and push the tag to build a release. There's a workflow
that picks up the tag, builds all the files, and upload them as a release.

```
git tag -a v0.0.2 -m "Release v0.0.2"
git push origin v0.0.2
```
