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

# Both settings may instead be supplied via environment variables:
#   ALGOLIA_APP_ID, ALGOLIA_API_KEY
