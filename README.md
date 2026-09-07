# terraform-provider-algolia

A Terraform provider for [Algolia](https://www.algolia.com/), built on the
official Algolia Go SDK v4. It currently manages two resources:

- `algolia_api_key` — create and manage Algolia API keys.
- `algolia_index_config` — manage the settings (configuration) of an index with
  a fully typed schema.

## Features

- **`algolia_api_key`**: manage ACLs, allowed indexes, referers, query
  parameters, rate limits, and validity. The generated key value is stored as
  the resource ID and exposed as a sensitive `key` attribute.
- **`algolia_index_config`**: a **typed** subset of Algolia index settings —
  searchable attributes, faceting, ranking, typo tolerance, pagination, and
  highlighting/snippeting. This resource manages index *settings* only; it does
  not create or delete the index itself.
- Import support for both resources.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://go.dev/doc/install) >= 1.23 (to build the provider)

## Using the provider

```hcl
terraform {
  required_providers {
    algolia = {
      source = "ECCOShoes/algolia"
    }
  }
}

provider "algolia" {
  app_id  = var.algolia_app_id
  api_key = var.algolia_api_key
}
```

Both settings can also be provided via environment variables: `ALGOLIA_APP_ID`
and `ALGOLIA_API_KEY`. The API key must be an Admin key with permission to
manage index settings and API keys.

### Index configuration lifecycle

`algolia_index_config` manages settings only. On `terraform destroy` the default
behavior (`reset_on_destroy = false`) simply removes the resource from state and
leaves the index and its settings untouched. Set `reset_on_destroy = true` to
reset the managed settings to Algolia's default values on destroy instead.
Creating or deleting the index itself is out of scope for this resource.

See [`examples/`](./examples) and the generated [`docs/`](./docs) for full usage.

## Developing the provider

```sh
# Build
make build

# Format & vet
make fmt
make vet

# Unit tests
make test

# Regenerate documentation (requires Terraform on PATH)
go generate ./...

# Install the provider into the local plugin directory
make install
```

### Acceptance tests

Acceptance tests create and destroy real resources in an Algolia application and
only run when `TF_ACC` is set. Configure the `ALGOLIA_*` environment variables
first:

```sh
export ALGOLIA_APP_ID=... ALGOLIA_API_KEY=...
make testacc
```

## License

This provider is distributed under the [MIT License](./LICENSE).
