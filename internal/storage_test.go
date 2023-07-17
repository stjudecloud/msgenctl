package internal

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/google/go-cmp/cmp"
)

func TestBlobServiceClientGenerateBlobSAS(t *testing.T) {
	const accountName = "msgenctl"
	const accountKey = "bXNnZW5jdGw="

	blobServiceClient, err := NewBlobServiceClient(accountName, accountKey)

	if err != nil {
		t.Fatal(err)
	}

	permissions := sas.BlobPermissions{Read: true}
	rawSAS, err := blobServiceClient.GenerateBlobSAS("test", "in.bam", permissions)

	if err != nil {
		t.Fatal(err)
	}

	sas, err := url.ParseQuery(rawSAS)

	if err != nil {
		t.Fatal(err)
	}

	if actual, ok := sas["sr"]; ok {
		expected := []string{"b"}

		if !cmp.Equal(actual, expected) {
			t.Errorf("expected sr=%v, got sr=%s", expected, actual)
		}
	} else {
		t.Error("missing sr entry")
	}
}

func TestBlobServiceClientGenerateContainerSAS(t *testing.T) {
	const accountName = "msgenctl"
	const accountKey = "bXNnZW5jdGw="

	blobServiceClient, err := NewBlobServiceClient(accountName, accountKey)

	if err != nil {
		t.Fatal(err)
	}

	permissions := sas.ContainerPermissions{Delete: true, Read: true, Write: true}
	rawSAS, err := blobServiceClient.GenerateContainerSAS("test", permissions)

	if err != nil {
		t.Fatal(err)
	}

	sas, err := url.ParseQuery(rawSAS)

	if err != nil {
		t.Fatal(err)
	}

	if actual, ok := sas["sr"]; ok {
		expected := []string{"c"}

		if !cmp.Equal(actual, expected) {
			t.Errorf("expected sr=%v, got sr=%s", expected, actual)
		}
	} else {
		t.Error("missing sr entry")
	}
}

func TestEncodeOrdered(t *testing.T) {
	const accountName = "msgenctl"
	const accountKey = "bXNnZW5jdGw="

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)

	if err != nil {
		t.Fatal(err)
	}

	expiryTime := time.Now().UTC()
	permissions := sas.BlobPermissions{Read: true}

	values := sas.BlobSignatureValues{
		ExpiryTime:  expiryTime,
		Permissions: permissions.String(),
	}

	queryParams, err := values.SignWithSharedKey(credential)

	if err != nil {
		t.Fatal(err)
	}

	sas, err := encodeOrdered(&queryParams)

	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(sas, "sv=") {
		t.Errorf("query missing starting `sv` key: %s", sas)
	}
}

func TestParseConnectionString(t *testing.T) {
	s := "AccountName=msgenctl;AccountKey=secret;DefaultEndpointsProtocol=https;"
	actual, err := ParseConnectionString(s)

	if err != nil {
		t.Fatal(err)
	}

	expected := ConnectionString{
		AccountName: "msgenctl",
		AccountKey:  "secret",
	}

	if diff := cmp.Diff(actual, expected); len(diff) != 0 {
		t.Errorf("connection string mismatch (-actual, +expected):\n%s", diff)
	}
}
