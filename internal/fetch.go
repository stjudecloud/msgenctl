package internal

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
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

	blobNameWithSAS, err := generateInputBlobSAS(config.Input)

	if err != nil {
		return newWorkflow, err
	}

	containerSAS, err := generateOutputContainerSAS(config.Output)

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
		BlobNames:        config.Input.BlobName,
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
	newWorkflow.IgnoreAzureRegion = config.IgnoreAzureRegion

	return newWorkflow, nil
}

func generateInputBlobSAS(config InputConfig) (string, error) {
	inputBlobServiceClient, err := NewBlobServiceClient(
		config.Storage.AccountName,
		config.Storage.AccountKey,
	)

	if err != nil {
		return "", err
	}

	blobName := config.BlobName
	blobSAS, err := inputBlobServiceClient.GenerateBlobSAS(
		config.Storage.ContainerName,
		blobName,
		sas.BlobPermissions{Read: true},
	)

	if err != nil {
		return "", err
	}

	blobNameWithSAS := fmt.Sprintf("%s?%s", blobName, blobSAS)

	return blobNameWithSAS, nil
}

func generateOutputContainerSAS(config OutputConfig) (string, error) {
	outputBlobServiceClient, err := NewBlobServiceClient(
		config.Storage.AccountName,
		config.Storage.AccountKey,
	)

	if err != nil {
		return "", err
	}

	return outputBlobServiceClient.GenerateContainerSAS(
		config.Storage.ContainerName,
		sas.ContainerPermissions{Delete: true, Read: true, Write: true},
	)
}
