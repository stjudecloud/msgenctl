package internal

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
)

const sasLifetime = 72 * time.Hour

type BlobServiceClient struct {
	credential *azblob.SharedKeyCredential
}

func NewBlobServiceClient(accountName string, accountKey string) (BlobServiceClient, error) {
	client := BlobServiceClient{}

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)

	if err != nil {
		return client, err
	}

	client.credential = credential

	return client, nil
}

func (c *BlobServiceClient) GenerateBlobSAS(
	containerName string,
	blobName string,
	permissions sas.BlobPermissions,
) (string, error) {
	now := time.Now().UTC()
	expiryTime := now.Add(sasLifetime)

	values := sas.BlobSignatureValues{
		StartTime:     now,
		ExpiryTime:    expiryTime,
		ContainerName: containerName,
		BlobName:      blobName,
		Permissions:   permissions.String(),
	}

	queryParams, err := values.SignWithSharedKey(c.credential)

	if err != nil {
		return "", err
	}

	sas, err := encodeOrdered(&queryParams)

	if err != nil {
		return "", err
	}

	return sas, nil
}

func (c *BlobServiceClient) GenerateContainerSAS(
	containerName string,
	permissions sas.ContainerPermissions,
) (string, error) {
	now := time.Now().UTC()
	expiryTime := now.Add(sasLifetime)

	values := sas.BlobSignatureValues{
		StartTime:     now,
		ExpiryTime:    expiryTime,
		ContainerName: containerName,
		Permissions:   permissions.String(),
	}

	queryParams, err := values.SignWithSharedKey(c.credential)

	if err != nil {
		return "", err
	}

	sas, err := encodeOrdered(&queryParams)

	if err != nil {
		return "", err
	}

	return sas, nil
}

// encodeOrdered encodes the SAS query parameters with the `signedversion`
// (`sv`) key first.
//
// Microsoft Genomics requires the `signedversion` key to be first in the SAS;
// otherwise, the service mysteriously replies with an HTTP 500.
func encodeOrdered(p *sas.QueryParameters) (string, error) {
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

type ConnectionString struct {
	AccountName string
	AccountKey  string
}

func ParseConnectionString(s string) (ConnectionString, error) {
	const delimiter = ";"
	const componentSeparator = "="
	const maxComponents = 2

	connectionString := ConnectionString{}

	if len(s) == 0 {
		return connectionString, errors.New("invalid connection string: empty input")
	}

	rawFields := strings.Split(strings.TrimRight(s, delimiter), delimiter)

	for _, rawField := range rawFields {
		if len(rawField) == 0 {
			return connectionString, errors.New("invalid connection string: contains an empty field")
		}

		components := strings.SplitN(rawField, componentSeparator, maxComponents)
		key := components[0]

		if len(components) != 2 {
			return connectionString, fmt.Errorf("invalid connection string: %s is missing a value", key)
		}

		value := components[1]

		switch key {
		case "AccountName":
			connectionString.AccountName = value
		case "AccountKey":
			connectionString.AccountKey = value
		default:
			continue
		}
	}

	return connectionString, nil
}
