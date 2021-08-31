#!/usr/bin/env bash

set -eux

# Builds an archive for release.
#
# The prebuilt binary is packaged along with the project license and README. A
# gzipped tar archive is created for Linux and macOS; whereas, a zip archive is
# created for Windows. The filename for the Windows binary is expected to have
# the `.exe` suffix.
#
# The following environment variables are expected to be set:
#
#   * `GITHUB_REF` (e.g., `refs/tags/v0.1.0`),
#   * `GOOS` (e.g., `linux`),
#   * `GOARCH` (e.g., `amd64`), and
#   * `PROJECT_NAME` (e.g., `msgenctl`).
#
# See `go tool dist list` for a list of supported `GOOS`/`GOARCH` combinations.
function main {
    local version=${GITHUB_REF/refs\/tags\/v/}
    local target=$GOOS-$GOARCH
    local package_name=$PROJECT_NAME-$version-$target

    local staging_prefix
    staging_prefix=$(mktemp -d)

    local working_prefix=$staging_prefix/$package_name

    local dst_prefix
    dst_prefix=$(pwd)

    mkdir "$working_prefix"

    local bin_name=$PROJECT_NAME

    if [[ $GOOS == "windows" ]]; then
        bin_name="$bin_name.exe"
    fi

    cp "$bin_name" "$working_prefix"
    cp LICENSE.txt README.md "$working_prefix"

    pushd "$staging_prefix"

    if [[ $GOOS == "windows" ]]; then
        7z a "$dst_prefix/$package_name.zip" "$package_name"
    else
        tar cfz "$dst_prefix/$package_name.tar.gz" "$package_name"
    fi

    popd

    rm -r "$staging_prefix"
}

main
