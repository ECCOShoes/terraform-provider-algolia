# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial release of the `algolia` provider, built on the official Algolia Go
  SDK v4. Authenticates with an application ID and Admin API key (`app_id` /
  `api_key`, or the `ALGOLIA_APP_ID` / `ALGOLIA_API_KEY` environment variables).
- `algolia_api_key` resource to manage API keys: ACLs, allowed indexes,
  referers, query parameters, rate limits (`max_hits_per_query`,
  `max_queries_per_ip_per_hour`), and `validity`.
- `algolia_index_config` resource to manage a typed subset of index settings:
  searchable attributes, faceting, ranking, typo tolerance, pagination, and
  highlighting/snippeting. Manages settings only; destroy removes the resource
  from state by default, or resets settings to Algolia defaults when
  `reset_on_destroy = true`.
- Resource import support for both resources.

