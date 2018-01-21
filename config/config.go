package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

// The AzureBlobCli represents configuration for the azureBlobcli
type AzureBlobCli struct {
	StorageAccountName      string `json:"storage_account_name"`
	StorageAccountAccessKey string `json:"storage_account_access_key"`
	ContainerName           string `json:"container_name"`
	CredentialsSource       string `json:"credentials_source"`
	BlockSize               uint64 `json:"block_size"`
	Parallelism             uint16 `json:"parallelism"`
}

// StaticCredentialsSource specifies that credentials will be supplied using storage_account_name and storage_account_access_key
const StaticCredentialsSource = "static"

// NoneCredentialsSource specifies that credentials will be empty. The blobstore client operates in read only mode.
const NoneCredentialsSource = "none"

const credentialsSourceEnvOrProfile = "env_or_profile"

// Nothing was provided in configuration
const noCredentialsSourceProvided = ""

var errorStaticCredentialsMissing = errors.New("storage_account_name and storage_account_access_key must be provided")

type errorStaticCredentialsPresent struct {
	credentialsSource string
}

func (e errorStaticCredentialsPresent) Error() string {
	return fmt.Sprintf("can't use storage_account_name and storage_account_access_key with %s credentials_source", e.credentialsSource)
}

func newStaticCredentialsPresentError(desiredSource string) error {
	return errorStaticCredentialsPresent{credentialsSource: desiredSource}
}

// NewFromReader returns a new azureBlobcli configuration struct from the contents of reader.
// reader.Read() is expected to return valid JSON
func NewFromReader(reader io.Reader) (AzureBlobCli, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return AzureBlobCli{}, err
	}

	c := AzureBlobCli{}

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return AzureBlobCli{}, err
	}

	if c.ContainerName == "" {
		return AzureBlobCli{}, errors.New("container_name must be set")
	}

	switch c.CredentialsSource {
	case StaticCredentialsSource:
		if c.StorageAccountName == "" || c.StorageAccountAccessKey == "" {
			return AzureBlobCli{}, errorStaticCredentialsMissing
		}
	case credentialsSourceEnvOrProfile:
		if c.StorageAccountName != "" || c.StorageAccountAccessKey != "" {
			return AzureBlobCli{}, newStaticCredentialsPresentError(credentialsSourceEnvOrProfile)
		}
	case NoneCredentialsSource:
		if c.StorageAccountName != "" || c.StorageAccountAccessKey != "" {
			return AzureBlobCli{}, newStaticCredentialsPresentError(NoneCredentialsSource)
		}

	case noCredentialsSourceProvided:
		if c.StorageAccountAccessKey != "" && c.StorageAccountName != "" {
			c.CredentialsSource = StaticCredentialsSource
		} else if c.StorageAccountAccessKey == "" && c.StorageAccountName == "" {
			c.CredentialsSource = NoneCredentialsSource
		} else {
			return AzureBlobCli{}, errorStaticCredentialsMissing
		}
	default:
		return AzureBlobCli{}, fmt.Errorf("Invalid credentials_source: %s", c.CredentialsSource)
	}

	return c, nil
}
