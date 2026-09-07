package provider

import (
	"os"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIndexConfigResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCleanupIndex,
		Steps: []resource.TestStep{
			{
				Config: `
resource "algolia_index_config" "test" {
  name                    = "acc_test_index_config"
  searchable_attributes   = ["name", "description"]
  attributes_for_faceting = ["brand"]
  custom_ranking          = ["desc(popularity)"]
  typo_tolerance          = "min"
  hits_per_page           = 25
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("algolia_index_config.test", "id", "acc_test_index_config"),
					resource.TestCheckResourceAttr("algolia_index_config.test", "searchable_attributes.#", "2"),
					resource.TestCheckResourceAttr("algolia_index_config.test", "typo_tolerance", "min"),
					resource.TestCheckResourceAttr("algolia_index_config.test", "hits_per_page", "25"),
				),
			},
			{
				Config: `
resource "algolia_index_config" "test" {
  name                  = "acc_test_index_config"
  searchable_attributes = ["name"]
  typo_tolerance        = "strict"
  hits_per_page         = 50
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("algolia_index_config.test", "searchable_attributes.#", "1"),
					resource.TestCheckResourceAttr("algolia_index_config.test", "typo_tolerance", "strict"),
					resource.TestCheckResourceAttr("algolia_index_config.test", "hits_per_page", "50"),
				),
			},
			{
				ResourceName:      "algolia_index_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccCleanupIndex deletes any test indices left behind. The index_config
// resource intentionally does not delete the index on destroy, so the test
// cleans it up directly.
func testAccCleanupIndex(s *terraform.State) error {
	c, err := search.NewClient(os.Getenv("ALGOLIA_APP_ID"), os.Getenv("ALGOLIA_API_KEY"))
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "algolia_index_config" {
			continue
		}
		if _, err := c.DeleteIndex(c.NewApiDeleteIndexRequest(rs.Primary.ID)); err != nil {
			return err
		}
	}
	return nil
}
