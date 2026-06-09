package twilio

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccTwilioAccount exercises the hand-written subaccount resource against
// the in-memory mock: create, update, import, and destroy (which closes it).
func TestAccTwilioAccount(t *testing.T) {
	setupMockProvider(t)

	const resourceName = "twilio_api_accounts.tenant"

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "twilio_api_accounts" "tenant" {
					friendly_name = "Tenant A"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "friendly_name", "Tenant A"),
					resource.TestCheckResourceAttrSet(resourceName, "sid"),
					resource.TestCheckResourceAttrSet(resourceName, "auth_token"),
				),
			},
			{
				Config: `resource "twilio_api_accounts" "tenant" {
					friendly_name = "Tenant A renamed"
				}`,
				Check: resource.TestCheckResourceAttr(resourceName, "friendly_name", "Tenant A renamed"),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTwilioAccountKeySecret verifies that the API key secret returned only
// at creation is captured into state by the overridden keys resource.
func TestAccTwilioAccountKeySecret(t *testing.T) {
	setupMockProvider(t)

	const resourceName = "twilio_api_accounts_keys.runtime"

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "twilio_api_accounts_keys" "runtime" {
					friendly_name = "runtime key"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "friendly_name", "runtime key"),
					resource.TestCheckResourceAttrSet(resourceName, "sid"),
					resource.TestCheckResourceAttrSet(resourceName, "secret"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// The secret is only returned at creation and cannot be re-read,
				// so it is expected to be absent from imported state.
				ImportStateVerifyIgnore: []string{"secret"},
			},
		},
	})
}
