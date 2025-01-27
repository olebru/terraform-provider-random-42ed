package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// These results are current as of Go 1.6. The Go
// "rand" package does not guarantee that the random
// number generator will generate the same results
// forever, but the maintainers endeavor not to change
// it gratuitously.
// These tests allow us to detect such changes and
// document them when they arise, but the docs for this
// resource specifically warn that results are not
// guaranteed consistent across Terraform releases.
func TestAccResourceShuffleDefault(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceShuffleConfigDefault,
				Check: resource.ComposeTestCheckFunc(
					testAccResourceShuffleCheck(
						"random_shuffle.default_length",
						[]string{"a", "c", "b", "e", "d"},
					),
				),
			},
		},
	})
}

func TestAccResourceShuffleShorter(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceShuffleConfigShorter,
				Check: resource.ComposeTestCheckFunc(
					testAccResourceShuffleCheck(
						"random_shuffle.shorter_length",
						[]string{"a", "c", "b"},
					),
				),
			},
		},
	})
}

func TestAccResourceShuffleLonger(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceShuffleConfigLonger,
				Check: resource.ComposeTestCheckFunc(
					testAccResourceShuffleCheck(
						"random_shuffle.longer_length",
						[]string{"a", "c", "b", "e", "d", "a", "e", "d", "c", "b", "a", "b"},
					),
				),
			},
		},
	})
}

func TestAccResourceShuffleEmpty(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceShuffleConfigEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccResourceShuffleCheck(
						"random_shuffle.empty_length",
						[]string{},
					),
				),
			},
		},
	})
}

func TestAccResourceShuffleOne(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceShuffleConfigOne,
				Check: resource.ComposeTestCheckFunc(
					testAccResourceShuffleCheck(
						"random_shuffle.one_length",
						[]string{"a"},
					),
				),
			},
		},
	})
}

func testAccResourceShuffleCheck(id string, wants []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		attrs := rs.Primary.Attributes

		gotLen := attrs["result.#"]
		wantLen := strconv.Itoa(len(wants))
		if gotLen != wantLen {
			return fmt.Errorf("got %s result items; want %s", gotLen, wantLen)
		}

		for i, want := range wants {
			key := fmt.Sprintf("result.%d", i)
			if got := attrs[key]; got != want {
				return fmt.Errorf("index %d is %q; want %q", i, got, want)
			}
		}

		return nil
	}
}

const (
	testAccResourceShuffleConfigDefault = `
resource "random_shuffle" "default_length" {
    input = ["a", "b", "c", "d", "e"]
    seed = "-"
}`

	testAccResourceShuffleConfigShorter = `
resource "random_shuffle" "shorter_length" {
    input = ["a", "b", "c", "d", "e"]
    seed = "-"
    result_count = 3
}
`

	testAccResourceShuffleConfigLonger = `
resource "random_shuffle" "longer_length" {
    input = ["a", "b", "c", "d", "e"]
    seed = "-"
    result_count = 12
}
`

	testAccResourceShuffleConfigEmpty = `
resource "random_shuffle" "empty_length" {
    input = []
    seed = "-"
    result_count = 12
}
`

	testAccResourceShuffleConfigOne = `
resource "random_shuffle" "one_length" {
    input = ["a"]
    seed = "-"
    result_count = 1
}
`
)
