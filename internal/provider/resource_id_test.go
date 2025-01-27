package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type idLens struct {
	b64UrlLen int
	b64StdLen int
	hexLen    int
}

func TestAccResourceID(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIDConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccResourceIDCheck("random_id.foo", &idLens{
						b64UrlLen: 6,
						b64StdLen: 8,
						hexLen:    8,
					}),
				),
			},
			{
				ResourceName:      "random_id.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceID_importWithPrefix(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIDConfigWithPrefix,
				Check: resource.ComposeTestCheckFunc(
					testAccResourceIDCheck("random_id.bar", &idLens{
						b64UrlLen: 12,
						b64StdLen: 14,
						hexLen:    14,
					}),
				),
			},
			{
				ResourceName:        "random_id.bar",
				ImportState:         true,
				ImportStateIdPrefix: "cloud-,",
				ImportStateVerify:   true,
			},
		},
	})
}

func testAccResourceIDCheck(id string, want *idLens) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		b64UrlStr := rs.Primary.Attributes["b64_url"]
		b64StdStr := rs.Primary.Attributes["b64_std"]
		hexStr := rs.Primary.Attributes["hex"]
		decStr := rs.Primary.Attributes["dec"]

		if got, want := len(b64UrlStr), want.b64UrlLen; got != want {
			return fmt.Errorf("base64 URL string length is %d; want %d", got, want)
		}
		if got, want := len(b64StdStr), want.b64StdLen; got != want {
			return fmt.Errorf("base64 STD string length is %d; want %d", got, want)
		}
		if got, want := len(hexStr), want.hexLen; got != want {
			return fmt.Errorf("hex string length is %d; want %d", got, want)
		}
		if len(decStr) < 1 {
			return fmt.Errorf("decimal string is empty; want at least one digit")
		}

		return nil
	}
}

const (
	testAccResourceIDConfig = `
resource "random_id" "foo" {
  byte_length = 4
}`

	testAccResourceIDConfigWithPrefix = `
resource "random_id" "bar" {
  byte_length = 4
  prefix      = "cloud-"
}
`
)
