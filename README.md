# msgenctl

[![CI status](https://github.com/stjudecloud/msgenctl/workflows/CI/badge.svg)](https://github.com/stjudecloud/msgenctl/actions/workflows/ci.yml)

**msgenctl** is a CLI to manage [Microsoft Genomics] workflows.

[Microsoft Genomics]: https://azure.microsoft.com/en-us/services/genomics/

## Prerequisites

  * [Go] 1.17+

[Go]: https://golang.org/

## Install

[Precompiled binaries] of msgenctl are built for Linux distributions
(`linux/amd64`), macOS (`darwin/amd64`), and Windows (`windows/amd64`).

Alternatively, it can be compiled from source.

```bash
$ git clone https://github.com/stjudecloud/msgenctl
$ cd msgenctl
$ go build
$ ./msgenctl --help
```

[Precompiled binaries]: https://github.com/stjudecloud/msgenctl/releases

## Usage

```
Query and send commands to Microsoft Genomics

Usage:
  msgenctl [command]

Available Commands:
  cancel      cancels a running workflow
  completion  generate the autocompletion script for the specified shell
  status      prints the status a workflow or all workflows
  submit      submits a new workflow
  wait        polls until the completion of a workflow

Flags:
      --access-key string   Microsoft Genomics API access key
      --base-url string     Microsoft Genomics API base URL
  -h, --help                help for msgenctl
  -v, --version             version for msgenctl

Use "msgenctl [command] --help" for more information about a command.
```

### Examples

#### Submit a workflow

```sh
msgenctl submit \
    --base-url $MSGEN_BASE_URL \
    --access-key $MSGEN_ACCESS_KEY \
    --process-name snapgatk \
    --process-args R=hg38m1x \
    --description sample \
    --input-storage-connection-string "$MSGEN_STORAGE_CONNECTION_STRING" \
    --input-storage-container-name $MSGEN_STORAGE_CONTAINER_NAME \
    --input-blob-name sample.bam \
    --output-storage-connection-string "$MSGEN_STORAGE_CONNECTION_STRING" \
    --output-storage-container-name $MSGEN_STORAGE_CONTAINER_NAME
```

#### Show the status of a workflow

```sh
msgenctl status --base-url $MSGEN_BASE_URL --access-key $MSGEN_ACCESS_KEY <workflow-id>
```

#### Show the statuses of all workflow

```sh
msgenctl status --base-url $MSGEN_BASE_URL --access-key $MSGEN_ACCESS_KEY
```

#### Wait until a workflow completes

```sh
msgenctl wait --base-url $MSGEN_BASE_URL --access-key $MSGEN_ACCESS_KEY <workflow-id>
```

#### Cancel a workflow

```sh
msgenctl cancel --base-url $MSGEN_BASE_URL --access-key $MSGEN_ACCESS_KEY <workflow-id>
```

## Limitations

  * Input (`--input-blob-name`) is expected to be a single BAM blob. A SAS is
    automatically generated for it.

  * The input Azure Storage account must be in the same region as the Microsoft
    Genomics service (`--base-url`).
