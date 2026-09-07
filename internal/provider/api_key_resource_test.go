package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccSearchClient(t *testing.T) *search.APIClient {
	t.Helper()
	c, err := search.NewClient(os.Getenv("ALGOLIA_APP_ID"), os.Getenv("ALGOLIA_API_KEY"))
	if err != nil {
		t.Fatalf("unable to build Algolia test client: %s", err)
	}
	return c
}

func TestAccAPIKeyResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
resource "algolia_api_key" "test" {
  description = "acc test key"
  acl         = ["search", "browse"]
  indexes     = ["acc_test_*"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("algolia_api_key.test", "key"),
					resource.TestCheckResourceAttrSet("algolia_api_key.test", "created_at"),
					resource.TestCheckResourceAttr("algolia_api_key.test", "description", "acc test key"),
					resource.TestCheckResourceAttr("algolia_api_key.test", "acl.#", "2"),
				),
			},
			{
				Config: `
resource "algolia_api_key" "test" {
  description = "acc test key updated"
  acl         = ["search"]
  indexes     = ["acc_test_*"]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("algolia_api_key.test", "description", "acc test key updated"),
					resource.TestCheckResourceAttr("algolia_api_key.test", "acl.#", "1"),
				),
			},
			{
				ResourceName:      "algolia_api_key.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAPIKeyDestroy(s *terraform.State) error {
	c, err := search.NewClient(os.Getenv("ALGOLIA_APP_ID"), os.Getenv("ALGOLIA_API_KEY"))
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "algolia_api_key" {
			continue
		}
		_, err := c.GetApiKey(c.NewApiGetApiKeyRequest(rs.Primary.ID))
		if err == nil {
			return fmt.Errorf("API key %s still exists", rs.Primary.ID)
		}
	}
	return nil
}
