# msgenctl

**msgenctl** is a CLI to manage [Microsoft Genomics] workflows.

[Microsoft Genomics]: https://azure.microsoft.com/en-us/services/genomics/

## Prerequisites

  * [Go] ~1.16

[Go]: https://golang.org/

## Quickstart

```bash
$ git clone https://github.com/stjudecloud/msgenctl
$ cd msgenctl
$ go build
$ ./msgenctl --help
```

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

Flags:
      --access-key string   Microsoft Genomics API access key
      --base-url string     Microsoft Genomics API base URL
  -h, --help                help for msgenctl
  -v, --version             version for msgenctl

Use "msgenctl [command] --help" for more information about a command.
```

### Examples

Submit a workflow

```sh
msgenctl submit \
    --base-url $MSGEN_BASE_URL \
    --access-key $MSGEN_ACCESS_KEY \
    --process-name snapgatk \
    --process-args R=hg38m1x \
    --description sample \
    --input-storage-account-name $MSGEN_STORAGE_ACCOUNT_NAME \
    --input-storage-account-key $MSGEN_STORAGE_ACCOUNT_KEY \
    --input-storage-container-name $MSGEN_STORAGE_CONTAINER_NAME \
    --input-blob-name sample.bam \
    --output-storage-account-name $MSGEN_STORAGE_ACCOUNT_NAME \
    --output-storage-account-key $MSGEN_STORAGE_ACCOUNT_KEY \
    --output-storage-container-name $MSGEN_STORAGE_CONTAINER_NAME
```

#### Show the status of a workflow

```sh
msgenctl status --base-url $MSGEN_BASE_URL --access-key $MSGEN_ACCESS_KEY 10000
```

#### Show the statuses of all workflow

```sh
msgenctl status --base-url $MSGEN_BASE_URL --access-key $MSGEN_ACCESS_KEY
```

#### Cancel a workflow

```sh
msgenctl cancel --base-url $MSGEN_BASE_URL --access-key $MSGEN_ACCESS_KEY 10000
```

## Limitations

  * Input (`--input-blob-name`) is expected to be a single BAM blob. A SAS is
    automatically generated for it.

  * The input Azure Storage account (`--input-storage-account-name`) must be in
    the same region as the Microsoft Genomics service (`--base-url`).
