package internal

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"go.uber.org/zap"
)

const SASLifetime = 72 * time.Hour

type BlobServiceClient struct {
	credential azblob.SharedKeyCredential
	serviceURL azblob.ServiceURL
}

func NewBlobServiceClient(accountName string, accountKey string) (BlobServiceClient, error) {
	client := BlobServiceClient{}

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)

	if err != nil {
		return client, err
	}

	serviceURL, err := newServiceURL(credential)

	if err != nil {
		return client, err
	}

	client.credential = *credential
	client.serviceURL = *serviceURL

	return client, nil
}

func newServiceURL(credential *azblob.SharedKeyCredential) (*azblob.ServiceURL, error) {
	accountName := credential.AccountName()
	baseURL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))

	if err != nil {
		return nil, err
	}

	options := azblob.PipelineOptions{}
	pipeline := azblob.NewPipeline(credential, options)

	serviceURL := azblob.NewServiceURL(*baseURL, pipeline)

	return &serviceURL, nil
}

func (c *BlobServiceClient) GenerateBlobSAS(
	containerName string,
	blobName string,
	permissions azblob.BlobSASPermissions,
) (string, error) {
	now := time.Now().UTC()
	expiryTime := now.Add(SASLifetime)

	values := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,
		StartTime:     now,
		ExpiryTime:    expiryTime,
		ContainerName: containerName,
		BlobName:      blobName,
		Permissions:   permissions.String(),
	}

	queryParams, err := values.NewSASQueryParameters(c.credential)

	if err != nil {
		return "", err
	}

	sas, err := encodeOrdered(&queryParams)

	if err != nil {
		return "", err
	}

	blobURL := fmt.Sprintf(
		"https://%s.blob.core.windows.net/%s/%s?%s",
		c.credential.AccountName(), containerName, blobName, sas,
	)
	zap.S().Debugw("generated blob URL", "url", blobURL)

	return sas, nil
}

func (c *BlobServiceClient) GenerateContainerSAS(
	containerName string,
	permissions azblob.ContainerSASPermissions,
) (string, error) {
	now := time.Now().UTC()
	expiryTime := now.Add(SASLifetime)

	values := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,
		StartTime:     now,
		ExpiryTime:    expiryTime,
		ContainerName: containerName,
		Permissions:   permissions.String(),
	}

	queryParams, err := values.NewSASQueryParameters(c.credential)

	if err != nil {
		return "", err
	}

	sas, err := encodeOrdered(&queryParams)

	if err != nil {
		return "", err
	}

	blobURL := fmt.Sprintf(
		"https://%s.blob.core.windows.net/%s?%s",
		c.credential.AccountName(), containerName, sas,
	)
	zap.S().Debugw("generated container URL", "url", blobURL)

	return sas, nil
}

// encodeOrdered encodes the SAS query parameters with the `signedversion`
// (`sv`) key first.
//
// Microsoft Genomics requires the `signedversion` key to be first in the SAS;
// otherwise, the service mysteriously replies with an HTTP 500.
func encodeOrdered(p *azblob.SASQueryParameters) (string, error) {
	const signedVersionKey = "sv"

	values, err := url.ParseQuery(p.Encode())

	if err != nil {
		return "", err
	}

	sv := values.Get(signedVersionKey)

	if len(sv) == 0 {
		return "", fmt.Errorf("missing signedversion (%s) field", signedVersionKey)
	}

	values.Del(signedVersionKey)

	var buf strings.Builder

	buf.WriteString(signedVersionKey)
	buf.WriteByte('=')
	buf.WriteString(url.QueryEscape(sv))
	buf.WriteByte('&')
	buf.WriteString(values.Encode())

	return buf.String(), nil
}
