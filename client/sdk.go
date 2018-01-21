package client

import (
	"fmt"
	"net/url"

	"github.com/Azure/azure-pipeline-go/pipeline"
	blob "github.com/Azure/azure-storage-blob-go/2016-05-31/azblob"
	"github.com/cloudfoundry/bosh-azureblobcli/config"
)

const (
	uaString         = "bosh-azureblobcli"
	blobFormatString = `https://%s.blob.core.windows.net`
)

func newStorageClient(cfg *config.AzureBlobCli) *blob.ContainerURL {
	accountName := cfg.StorageAccountName
	key := cfg.StorageAccountAccessKey
	var p pipeline.Pipeline

	if cfg.CredentialsSource == config.StaticCredentialsSource {
		c := blob.NewSharedKeyCredential(accountName, key)
		p = blob.NewPipeline(c, blob.PipelineOptions{
			Telemetry: blob.TelemetryOptions{Value: uaString},
		})
	}

	if cfg.CredentialsSource == config.NoneCredentialsSource {
		ac := blob.NewAnonymousCredential()
		p = blob.NewPipeline(ac, blob.PipelineOptions{
			Telemetry: blob.TelemetryOptions{Value: uaString},
		})
	}

	u, _ := url.Parse(fmt.Sprintf(blobFormatString, accountName))
	service := blob.NewServiceURL(*u, p)
	container := service.NewContainerURL(cfg.ContainerName)
	return &container
}
