package internal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/pflag"
)

func TestSubmitConfigFromFlags(t *testing.T) {
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	baseURL := flags.String("base-url", "", "")
	accessKey := flags.String("access-key", "", "")
	processName := flags.String("process-name", "", "")
	processArgs := flags.String("process-args", "", "")
	flags.String("input-storage-connection-string", "", "")
	inputStorageContainerName := flags.String("input-storage-container-name", "", "")
	inputBlobName := flags.String("input-blob-name", "", "")
	description := flags.String("description", "", "")
	flags.String("output-storage-connection-string", "", "")
	outputStorageContainerName := flags.String("output-storage-container-name", "", "")
	outputBasename := flags.String("output-basename", "", "")
	overwrite := flags.Bool("output-overwrite", false, "")
	includeLog := flags.Bool("output-include-log", true, "")
	flags.String("emit-ref-confidence", "", "")
	bgzipOutput := flags.Bool("bgzip-output", false, "")

	args := []string{
		"--base-url", "https://example.com",
		"--access-key", "secret",
		"--process-name", "snapgatk-20190409_1",
		"--process-args", "R=hg38m1x",
		"--input-storage-connection-string", "AccountName=input;AccountKey=input-secret;",
		"--input-storage-container-name", "data",
		"--input-blob-name", "sample.bam",
		"--description", "sample run",
		"--output-storage-connection-string", "AccountName=output;AccountKey=output-secret;",
		"--output-storage-container-name", "results",
		"--output-basename", "sample",
		"--output-overwrite",
		"--output-include-log",
		"--emit-ref-confidence", "GVCF",
		"--bgzip-output",
	}

	if err := flags.Parse(args); err != nil {
		t.Fatal(err)
	}

	actual, err := SubmitConfigFromFlags(flags)

	if err != nil {
		t.Fatal(err)
	}

	expected := SubmitConfig{
		Service: ServiceConfig{
			BaseURL:   *baseURL,
			AccessKey: *accessKey,
		},
		Input: InputConfig{
			Storage: StorageConfig{
				AccountName:   "input",
				AccountKey:    "input-secret",
				ContainerName: *inputStorageContainerName,
			},
			BlobName: *inputBlobName,
		},
		Process: ProcessConfig{
			Name: *processName,
			Args: *processArgs,
		},
		Description: *description,
		Output: OutputConfig{
			Storage: StorageConfig{
				AccountName:   "output",
				AccountKey:    "output-secret",
				ContainerName: *outputStorageContainerName,
			},
			Basename:   *outputBasename,
			Overwrite:  *overwrite,
			IncludeLog: *includeLog,
		},
		OptionalArgs: OptionalArgsConfig{
			EmitRefConfidence: ReferenceConfidenceModeGVCF,
			BgzipOutput:       *bgzipOutput,
		},
	}

	if diff := cmp.Diff(actual, expected); len(diff) != 0 {
		t.Errorf("config mismatch (-actual, +expected):\n%s", diff)
	}
}

func TestServiceConfigFromFlags(t *testing.T) {
	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	baseURL := flags.String("base-url", "", "")
	accessKey := flags.String("access-key", "", "")

	args := []string{
		"--base-url", "https://example.com",
		"--access-key", "secret",
	}

	if err := flags.Parse(args); err != nil {
		t.Fatal(err)
	}

	actual, err := ServiceConfigFromFlags(flags)

	if err != nil {
		t.Fatal(err)
	}

	expected := ServiceConfig{
		BaseURL:   *baseURL,
		AccessKey: *accessKey,
	}

	if diff := cmp.Diff(actual, expected); len(diff) != 0 {
		t.Errorf("config mismatch (-actual, +expected):\n%s", diff)
	}
}
