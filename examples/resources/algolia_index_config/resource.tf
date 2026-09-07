resource "algolia_index_config" "products" {
  name = "products"

  searchable_attributes = [
    "name",
    "unordered(description)",
  ]

  attributes_for_faceting = [
    "brand",
    "searchable(categories)",
    "filterOnly(price)",
  ]

  attributes_to_retrieve = ["*"]

  ranking = [
    "typo",
    "geo",
    "words",
    "filters",
    "proximity",
    "attribute",
    "exact",
    "custom",
  ]

  custom_ranking = ["desc(popularity)"]

  max_values_per_facet = 100
  sort_facet_values_by = "count"

  min_word_size_for_1_typo      = 4
  min_word_size_for_2_typos     = 8
  typo_tolerance                = "min"
  allow_typos_on_numeric_tokens = false

  hits_per_page         = 20
  pagination_limited_to = 1000

  attributes_to_highlight = ["name", "description"]
  highlight_pre_tag       = "<mark>"
  highlight_post_tag      = "</mark>"
}
