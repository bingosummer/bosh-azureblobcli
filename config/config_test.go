package config_test

import (
	"bytes"
	"errors"

	"github.com/cloudfoundry/bosh-azureblobcli/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BlobstoreClient configuration", func() {
	Describe("building a configuration", func() {
		Describe("when container is not specified", func() {
			emptyJSONBytes := []byte(`{"storage_account_name": "name", "storage_account_access_key": "key"}`)
			emptyJSONReader := bytes.NewReader(emptyJSONBytes)

			It("returns an error", func() {
				_, err := config.NewFromReader(emptyJSONReader)
				Expect(err).To(MatchError("container_name must be set"))
			})
		})

		Describe("when container is specified", func() {
			emptyJSONBytes := []byte(`{"storage_account_name": "name", "storage_account_access_key": "key", "container_name": "some-container"}`)
			emptyJSONReader := bytes.NewReader(emptyJSONBytes)

			It("uses the given bucket", func() {
				c, err := config.NewFromReader(emptyJSONReader)
				Expect(err).ToNot(HaveOccurred())
				Expect(c.ContainerName).To(Equal("some-container"))
			})
		})

		Context("when the configuration file cannot be read", func() {
			It("returns an error", func() {
				f := explodingReader{}

				_, err := config.NewFromReader(f)
				Expect(err).To(MatchError("explosion"))
			})
		})

		Context("when the configuration file is invalid JSON", func() {
			It("returns an error", func() {
				invalidJSONBytes := []byte(`invalid-json`)
				invalidJSONReader := bytes.NewReader(invalidJSONBytes)

				_, err := config.NewFromReader(invalidJSONReader)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("validating credentials", func() {
		Describe("when credentials source is not specified", func() {
			Context("when a storage account name and access key are provided", func() {
				It("defaults to static credentials", func() {
					dummyJSONBytes := []byte(`{"storage_account_name": "name", "storage_account_access_key": "key", "container_name": "some-container"}`)
					dummyJSONReader := bytes.NewReader(dummyJSONBytes)

					c, err := config.NewFromReader(dummyJSONReader)
					Expect(err).ToNot(HaveOccurred())
					Expect(c.CredentialsSource).To(Equal("static"))
				})
			})

			Context("when either the storage account name and access key are missing", func() {
				It("raises an error", func() {
					dummyJSONBytes := []byte(`{"storage_account_access_key": "key", "container_name": "some-container"}`)
					dummyJSONReader := bytes.NewReader(dummyJSONBytes)

					_, err := config.NewFromReader(dummyJSONReader)
					Expect(err).To(MatchError("storage_account_name and storage_account_access_key must be provided"))
				})
			})

			Context("when neither storage account name and access key are provided", func() {
				It("defaults credentials source to anonymous", func() {
					dummyJSONBytes := []byte(`{"container_name": "some-container"}`)
					dummyJSONReader := bytes.NewReader(dummyJSONBytes)

					c, err := config.NewFromReader(dummyJSONReader)
					Expect(err).ToNot(HaveOccurred())
					Expect(c.CredentialsSource).To(Equal("none"))
				})
			})

			Describe("when credentials source is invalid", func() {
				It("returns an error", func() {
					dummyJSONBytes := []byte(`{"container_name": "some-container", "credentials_source": "magical_unicorns"}`)
					dummyJSONReader := bytes.NewReader(dummyJSONBytes)

					_, err := config.NewFromReader(dummyJSONReader)
					Expect(err).To(MatchError("Invalid credentials_source: magical_unicorns"))
				})
			})

		})

		Context("when credential source is `static`", func() {
			It("validates that storage account name and access key are set", func() {
				dummyJSONBytes := []byte(`{"container_name": "some-container", "storage_account_name": "some_id"}`)
				dummyJSONReader := bytes.NewReader(dummyJSONBytes)
				_, err := config.NewFromReader(dummyJSONReader)
				Expect(err).To(MatchError("storage_account_name and storage_account_access_key must be provided"))

				dummyJSONBytes = []byte(`{"container_name": "some-container", "storage_account_name": "some_id", "storage_account_access_key": "some_secret"}`)
				dummyJSONReader = bytes.NewReader(dummyJSONBytes)
				_, err = config.NewFromReader(dummyJSONReader)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when credentials source is `env_or_profile`", func() {
			It("validates that storage account name and access key are not set", func() {
				dummyJSONBytes := []byte(`{"container_name": "some-container", "credentials_source": "env_or_profile"}`)
				dummyJSONReader := bytes.NewReader(dummyJSONBytes)

				_, err := config.NewFromReader(dummyJSONReader)
				Expect(err).ToNot(HaveOccurred())

				dummyJSONBytes = []byte(`{"container_name": "some-container", "credentials_source": "env_or_profile", "storage_account_name": "some_id"}`)
				dummyJSONReader = bytes.NewReader(dummyJSONBytes)
				_, err = config.NewFromReader(dummyJSONReader)
				Expect(err).To(MatchError("can't use storage_account_name and storage_account_access_key with env_or_profile credentials_source"))

				dummyJSONBytes = []byte(`{"container_name": "some-container", "credentials_source": "env_or_profile", "storage_account_name": "some_id", "storage_account_access_key": "some_secret"}`)
				dummyJSONReader = bytes.NewReader(dummyJSONBytes)
				_, err = config.NewFromReader(dummyJSONReader)
				Expect(err).To(MatchError("can't use storage_account_name and storage_account_access_key with env_or_profile credentials_source"))
			})
		})

		Context("when the credentials source is `none`", func() {
			It("validates that storage account name and access key are not set", func() {
				dummyJSONBytes := []byte(`{"container_name": "some-container", "credentials_source": "none", "storage_account_name": "some_id"}`)
				dummyJSONReader := bytes.NewReader(dummyJSONBytes)
				_, err := config.NewFromReader(dummyJSONReader)
				Expect(err).To(MatchError("can't use storage_account_name and storage_account_access_key with none credentials_source"))
			})
		})
	})
})

type explodingReader struct{}

func (e explodingReader) Read([]byte) (int, error) {
	return 0, errors.New("explosion")
}
