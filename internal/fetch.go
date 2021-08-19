package internal

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func CancelWorkflow(client Client, ID WorkflowID) (Workflow, error) {
	workflow := Workflow{}

	endpoint := fmt.Sprintf("/api/workflows/%v", ID)
	response, err := client.Delete(endpoint)

	if err != nil {
		return workflow, err
	}

	defer response.Body.Close()

	if err := decodeJSON(response.Body, &workflow); err != nil {
		return workflow, err
	}

	return workflow, nil
}

func FetchWorkflows(client Client) ([]Workflow, error) {
	workflows := []Workflow{}

	response, err := client.Get("/api/workflows?$orderby=CreatedDate%20asc")

	if err != nil {
		return workflows, err
	}

	defer response.Body.Close()

	if err := decodeJSON(response.Body, &workflows); err != nil {
		return workflows, err
	}

	return workflows, nil
}

func FetchWorkflow(client Client, ID WorkflowID) (Workflow, error) {
	workflow := Workflow{}

	endpoint := fmt.Sprintf("/api/workflows/%v", ID)
	response, err := client.Get(endpoint)

	if err != nil {
		return workflow, err
	}

	defer response.Body.Close()

	if err := decodeJSON(response.Body, &workflow); err != nil {
		return workflow, err
	}

	return workflow, nil
}

func SubmitWorkflow(client Client, config SubmitConfig) (Workflow, error) {
	workflow := Workflow{}

	newWorkflow, err := buildSubmitWorkflowPayload(config)

	if err != nil {
		return workflow, err
	}

	response, err := client.Post("/api/workflows", newWorkflow)

	if err != nil {
		return workflow, err
	}

	defer response.Body.Close()

	if err := decodeJSON(response.Body, &workflow); err != nil {
		return workflow, err
	}

	return workflow, nil
}

func decodeJSON(reader io.Reader, value interface{}) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(value)
}

func buildSubmitWorkflowPayload(config SubmitConfig) (NewWorkflow, error) {
	newWorkflow := NewWorkflow{}

	inputBlobServiceClient, err := NewBlobServiceClient(
		config.Input.Storage.AccountName,
		config.Input.Storage.AccountKey,
	)

	if err != nil {
		return newWorkflow, err
	}

	blobName := config.Input.BlobName
	blobSAS, err := inputBlobServiceClient.GenerateBlobSAS(
		config.Input.Storage.ContainerName,
		blobName,
		azblob.BlobSASPermissions{Read: true},
	)

	if err != nil {
		return newWorkflow, err
	}

	blobNameWithSAS := fmt.Sprintf("%s?%s", blobName, blobSAS)

	outputBlobServiceClient, err := NewBlobServiceClient(
		config.Input.Storage.AccountName,
		config.Input.Storage.AccountKey,
	)

	if err != nil {
		return newWorkflow, err
	}

	containerSAS, err := outputBlobServiceClient.GenerateContainerSAS(
		config.Output.Storage.ContainerName,
		azblob.ContainerSASPermissions{Delete: true, Read: true, Write: true},
	)

	if err != nil {
		return newWorkflow, err
	}

	newWorkflow.Process = config.Process.Name
	newWorkflow.ProcessArgs = config.Process.Args
	newWorkflow.Description = config.Description
	newWorkflow.InputStorageType = StorageKindAzureBlockBlob
	newWorkflow.InputArgs = NewWorkflowInputArgs{
		AccountName:      config.Input.Storage.AccountName,
		ContainerName:    config.Input.Storage.ContainerName,
		BlobNames:        blobName,
		BlobNamesWithSAS: blobNameWithSAS,
	}
	newWorkflow.OutputStorageType = StorageKindAzureBlockBlob
	newWorkflow.OutputArgs = NewWorkflowOutputArgs{
		AccountName:           config.Output.Storage.AccountName,
		ContainerName:         config.Output.Storage.ContainerName,
		ContainerSAS:          containerSAS,
		Basename:              config.Output.Basename,
		Overwrite:             config.Output.Overwrite,
		OutputIncludeLogfiles: config.Output.IncludeLog,
	}
	newWorkflow.OptionalArgs = NewWorkflowOptionalArgs{
		GATKEmitRefConfidence: config.OptionalArgs.EmitRefConfidence,
		BgzipOutput:           config.OptionalArgs.BgzipOutput,
	}

	return newWorkflow, nil
}
