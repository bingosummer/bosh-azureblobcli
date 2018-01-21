# bosh-azureblobcli

A Golang CLI for uploading, fetching and deleting content to/from [Azure Blob Storage](https://azure.microsoft.com/en-us/services/storage/blobs/).
TODO: This tool exists to work with the [bosh-cli](https://github.com/cloudfoundry/bosh-cli) and [director](https://github.com/cloudfoundry/bosh).

## Installation

```bash
git clone https://github.com/bingosummer/bosh-azureblobcli $GOPATH/src/github.com/cloudfoundry/bosh-azureblobcli
```

## Commands

### Usage
```bash
bosh-azureblobcli --help
```
### Upload an object
```bash
bosh-azureblobcli -c config.json put <path/to/file> <remote-blob>
```
### Fetch an object
```bash
bosh-azureblobcli -c config.json get <remote-blob> <path/to/file>
```
### Delete an object
```bash
bosh-azureblobcli -c config.json delete <remote-blob>
```
### Check if an object exists
```bash
bosh-azureblobcli -c config.json exists <remote-blob>```
```

## Configuration
The command line tool expects a JSON configuration file. Run `bosh-azureblobcli --help` for details.

### Authentication Methods (`credentials_source`)
* `static`: `storage_account_name` and `storage_account_access_key` will be provided.
* `none`: No credentials are provided. The client is reading from a public bucket.
* `env_or_profile`: The environment `storage_account_name` and `storage_account_access_key` will be provided.

## Running Integration Tests

## Development

* A Makefile is provided that automates integration testing. Try `make help` to get started.
* [gvt](https://godoc.org/github.com/FiloSottile/gvt) is used for vendoring.

## Contributing

For details on how to contribute to this project - including filing bug reports and contributing code changes - please see [CONTRIBUTING.md](./CONTRIBUTING.md).

## License

This tool is licensed under Apache 2.0. Full license text is available in [LICENSE](LICENSE).
