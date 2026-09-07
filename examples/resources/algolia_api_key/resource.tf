resource "algolia_api_key" "search_only" {
  description = "Search-only key for the storefront"
  acl         = ["search", "browse"]
  indexes     = ["products", "products_*"]

  max_hits_per_query          = 100
  max_queries_per_ip_per_hour = 10000
}
