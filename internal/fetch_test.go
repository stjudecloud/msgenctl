package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCancelWorkflow(t *testing.T) {
	createdDate := time.Now()
	endDate := createdDate.Add(2 * time.Hour)

	expected := Workflow{
		ID:             1597,
		TenantID:       144,
		Status:         StatusCancelled,
		CreatedDate:    createdDate,
		EndDate:        &endDate,
		FailureCode:    0,
		Message:        "",
		Description:    "",
		Process:        "snapgatk-20190409_1",
		BasesProcessed: 0,
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		method := "DELETE"
		path := fmt.Sprintf("/api/workflows/%v", expected.ID)

		if r.Method != method {
			t.Errorf("expected request method %s, got %s", method, r.Method)
		}

		if r.URL.Path != path {
			t.Errorf("expected URL path %q, got %q", path, r.URL.Path)
		}

		payload, err := json.Marshal(expected)

		if err != nil {
			t.Fatal(err)
		}

		rw.Write(payload)
	}))

	defer server.Close()

	client := NewClient(server.URL, "secret")
	actual, err := CancelWorkflow(client, expected.ID)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(actual, expected); len(diff) != 0 {
		t.Errorf("workflow mismatch (-expected, +got):\n%s", diff)
	}
}

func TestFetchWorkflow(t *testing.T) {
	createdDate := time.Now()
	endDate := createdDate.Add(2 * time.Hour)

	expected := Workflow{
		ID:             1597,
		TenantID:       144,
		Status:         StatusSuccess,
		CreatedDate:    createdDate,
		EndDate:        &endDate,
		FailureCode:    0,
		Message:        "",
		Description:    "",
		Process:        "snapgatk-20190409_1",
		BasesProcessed: 21,
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		method := "GET"
		path := fmt.Sprintf("/api/workflows/%v", expected.ID)

		if r.Method != method {
			t.Errorf("expected request method %s, got %s", method, r.Method)
		}

		if r.URL.Path != path {
			t.Errorf("expected URL path %q, got %q", path, r.URL.Path)
		}

		payload, err := json.Marshal(expected)

		if err != nil {
			t.Fatal(err)
		}

		rw.Write(payload)
	}))

	defer server.Close()

	client := NewClient(server.URL, "secret")
	actual, err := FetchWorkflow(client, expected.ID)

	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(actual, expected); len(diff) != 0 {
		t.Errorf("workflow mismatch (-expected, +got):\n%s", diff)
	}
}

func TestBuildSubmitWorkflowPayload(t *testing.T) {
	config := SubmitConfig{
		Service: ServiceConfig{
			BaseURL:   "http://example.com",
			AccessKey: "secret",
		},
		Input: InputConfig{
			Storage: StorageConfig{
				AccountName:   "input",
				AccountKey:    "bXNnZW5jdGw=",
				ContainerName: "data",
			},
			BlobName: "sample.bam",
		},
		Process: ProcessConfig{
			Name: "snapgatk-20190409_1",
			Args: "R=hg38m1x",
		},
		Description: "msgenctl/test/sample",
		Output: OutputConfig{
			Storage: StorageConfig{
				AccountName:   "output",
				AccountKey:    "bXNnZW5jdGw=",
				ContainerName: "results",
			},
			Basename:   "sample",
			Overwrite:  true,
			IncludeLog: true,
		},
		OptionalArgs: OptionalArgsConfig{
			EmitRefConfidence: ReferenceConfidenceModeGVCF,
			BgzipOutput:       true,
		},
		IgnoreAzureRegion: true,
	}

	actual, err := buildSubmitWorkflowPayload(config)

	if err != nil {
		t.Fatal(err)
	}

	actual.InputArgs.BlobNamesWithSAS = ""
	actual.OutputArgs.ContainerSAS = ""

	expected := NewWorkflow{
		WorkflowClass:    "",
		Process:          config.Process.Name,
		ProcessArgs:      config.Process.Args,
		Description:      config.Description,
		InputStorageType: StorageKindAzureBlockBlob,
		InputArgs: NewWorkflowInputArgs{
			AccountName:   config.Input.Storage.AccountName,
			ContainerName: config.Input.Storage.ContainerName,
			BlobNames:     config.Input.BlobName,
		},
		OutputStorageType: StorageKindAzureBlockBlob,
		OutputArgs: NewWorkflowOutputArgs{
			AccountName:           config.Output.Storage.AccountName,
			ContainerName:         config.Output.Storage.ContainerName,
			Basename:              config.Output.Basename,
			Overwrite:             config.Output.Overwrite,
			OutputIncludeLogfiles: config.Output.IncludeLog,
		},
		OptionalArgs: NewWorkflowOptionalArgs{
			GATKEmitRefConfidence: config.OptionalArgs.EmitRefConfidence,
			BgzipOutput:           config.OptionalArgs.BgzipOutput,
		},
		IgnoreAzureRegion: config.IgnoreAzureRegion,
	}

	if diff := cmp.Diff(actual, expected); len(diff) != 0 {
		t.Errorf("workflow mismatch (-expected, +got):\n%s", diff)
	}
}
