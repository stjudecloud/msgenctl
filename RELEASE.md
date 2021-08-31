# Release

  * [ ] Update version in `internal/version.go`.
  * [ ] Update `CHANGELOG.md` with version and publication date.
  * [ ] Run tests: `go test ./...`.
  * [ ] Stage changes: `git add internal/version.go CHANGELOG.md`.
  * [ ] Create git commit: `git commit -m "Bump version to $VERSION"`.
  * [ ] Create git tag: `git tag -m "" -a v$VERSION`.
  * [ ] Push release: `git push --follow-tags`.
