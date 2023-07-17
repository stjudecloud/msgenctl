package internal

import (
	"fmt"
	"time"
)

type StorageKind string

const StorageKindAzureBlockBlob = "AZURE_BLOCK_BLOB"

type ReferenceConfidenceMode string

// / https://gatk.broadinstitute.org/hc/en-us/articles/4404604697243-HaplotypeCaller#--emit-ref-confidence
const (
	ReferenceConfidenceModeNone         = "NONE"
	ReferenceConfidenceModeBPResolution = "BP_RESOLUTION"
	ReferenceConfidenceModeGVCF         = "GVCF"
)

func ParseReferenceConfidenceMode(s string) (ReferenceConfidenceMode, error) {
	switch s {
	case "NONE":
		return ReferenceConfidenceModeNone, nil
	case "BP_RESOLUTION":
		return ReferenceConfidenceModeBPResolution, nil
	case "GVCF":
		return ReferenceConfidenceModeGVCF, nil
	default:
		return "", fmt.Errorf("invalid reference confidence mode: %q", s)
	}
}

type WorkflowID int

type Workflow struct {
	ID             WorkflowID `json:"Id"`
	TenantID       int        `json:"TenantId"`
	Status         Status
	CreatedDate    time.Time
	EndDate        *time.Time
	FailureCode    int
	Message        string
	Description    string
	Process        string
	BasesProcessed uint64
}

func (w *Workflow) Duration() time.Duration {
	endDate := w.EndDate

	if endDate == nil {
		now := time.Now().UTC()
		endDate = &now
	}

	return endDate.Sub(w.CreatedDate)
}

type NewWorkflowInputArgs struct {
	AccountName      string `json:"ACCOUNT"`
	ContainerName    string `json:"CONTAINER"`
	BlobNames        string `json:"BLOBNAMES"`
	BlobNamesWithSAS string `json:"BLOBNAMES_WITH_SAS"`
}

type NewWorkflowOutputArgs struct {
	AccountName           string `json:"ACCOUNT"`
	ContainerName         string `json:"CONTAINER"`
	ContainerSAS          string `json:"CONTAINER_SAS"`
	Basename              string `json:"OUTPUT_FILENAME_BASE"`
	Overwrite             bool   `json:"OVERWRITE"`
	OutputIncludeLogfiles bool   `json:"OUTPUT_INCLUDE_LOGFILES"`
}

type NewWorkflowOptionalArgs struct {
	GATKEmitRefConfidence ReferenceConfidenceMode `json:"GatkEmitRefConfidence"`
	BgzipOutput           bool
}

type NewWorkflow struct {
	WorkflowClass     string
	Process           string
	ProcessArgs       string
	Description       string
	InputStorageType  StorageKind
	InputArgs         NewWorkflowInputArgs
	OutputStorageType StorageKind
	OutputArgs        NewWorkflowOutputArgs
	OptionalArgs      NewWorkflowOptionalArgs
	IgnoreAzureRegion bool
}
