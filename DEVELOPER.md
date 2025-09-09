# README for developers

## New features flow

The `master` branch contains the latest and most current code. But it may contain bugs.  
That's why all releases go through tags.

# README for maintainers

Create a git tag from master
```shell
git tag -a v0.0.1-rc.1 -m "Release candidate v0.0.1-rc.1
```

Publish this release
```shell
make release
```