package client

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	blob "github.com/Azure/azure-storage-blob-go/2016-05-31/azblob"
	"github.com/cloudfoundry/bosh-azureblobcli/config"
)

// AzureBlobBlobstore encapsulates interaction with the Azure Blob Storage blobstore
type AzureBlobBlobstore struct {
	azureBlobContainerURL *blob.ContainerURL
	config                *config.AzureBlobCli
}

//var errorInvalidCredentialsSourceValue = errors.New("the client operates in read only mode. Change 'credentials_source' parameter value ")

// New returns a BlobstoreClient if the configuration file backing configFile is valid
func New(cfg *config.AzureBlobCli) (*AzureBlobBlobstore, error) {
	if cfg == nil {
		return nil, errors.New("expected non-nill config object")
	}

	containerURL := newStorageClient(cfg)

	return &AzureBlobBlobstore{azureBlobContainerURL: containerURL, config: cfg}, nil
}

// Get fetches a blob from an AzureBlob blobstore
// Destination will be overwritten if exists
func (client *AzureBlobBlobstore) Get(ctx context.Context, src string, dest io.Writer) error {
	blobURL := client.azureBlobContainerURL.NewBlockBlobURL(src)

	stream := blob.NewDownloadStream(ctx, blobURL.GetBlob, blob.DownloadStreamOptions{})
	defer stream.Close()

	_, err := io.Copy(dest, stream) // Write to the file by reading from the blob (with intelligent retries).
	if err != nil {
		return err
	}

	return nil
}

// Put uploads a blob to an AzureBlob blobstore
func (client *AzureBlobBlobstore) Put(ctx context.Context, src *os.File, dest string) error {
	blobURL := client.azureBlobContainerURL.NewBlockBlobURL(dest)

	_, err := blob.UploadFileToBlockBlob(ctx, src, blobURL,
		blob.UploadToBlockBlobOptions{
			BlockSize:   4 * 1024 * 1024,
			Parallelism: 16,
		})
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a blob from an AzureBlob blobstore. If the object does
// not exist, Delete does not return an error.
func (client *AzureBlobBlobstore) Delete(ctx context.Context, dest string) error {
	blobURL := client.azureBlobContainerURL.NewBlockBlobURL(dest)

	_, err := blobURL.Delete(ctx, blob.DeleteSnapshotsOptionNone, blob.BlobAccessConditions{})

	if err != nil {
		return err
	}

	return nil
}

// Exists checks if blob exists in an AzureBlob blobstore
func (client *AzureBlobBlobstore) Exists(ctx context.Context, dest string) (bool, error) {
	blobURL := client.azureBlobContainerURL.NewBlockBlobURL(dest)

	_, err := blobURL.GetPropertiesAndMetadata(ctx, blob.BlobAccessConditions{})
	if err == nil {
		return true, nil
	} else if strings.Contains(err.Error(), "404") {
		return false, nil
	}

	return false, err
}
