package ucloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/ucloud/ucloud-sdk-go/services/unet"
)

func TestAccUCloudEIP_basic(t *testing.T) {
	var eip unet.UnetEIPSet

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ucloud_eip.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEIPDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEIPConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists("ucloud_eip.foo", &eip),
					testAccCheckEIPAttributes(&eip),
					resource.TestCheckResourceAttr("ucloud_eip.foo", "bandwidth", "1"),
					resource.TestCheckResourceAttr("ucloud_eip.foo", "name", "tf-acc-eip"),
					resource.TestCheckResourceAttr("ucloud_eip.foo", "charge_mode", "bandwidth"),
				),
			},

			resource.TestStep{
				Config: testAccEIPConfigTwo,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists("ucloud_eip.foo", &eip),
					testAccCheckEIPAttributes(&eip),
					resource.TestCheckResourceAttr("ucloud_eip.foo", "bandwidth", "2"),
					resource.TestCheckResourceAttr("ucloud_eip.foo", "name", "tf-acc-eip-two"),
					resource.TestCheckResourceAttr("ucloud_eip.foo", "charge_mode", "traffic"),
				),
			},
		},
	})

}

func testAccCheckEIPExists(n string, eip *unet.UnetEIPSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("eip id is empty")
		}

		client := testAccProvider.Meta().(*UCloudClient)
		ptr, err := client.describeEIPById(rs.Primary.ID)

		log.Printf("[INFO] eip id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		*eip = *ptr
		return nil
	}
}

func testAccCheckEIPAttributes(eip *unet.UnetEIPSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if eip.EIPId == "" {
			return fmt.Errorf("eip id is empty")
		}
		return nil
	}
}

func testAccCheckEIPDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ucloud_eip" {
			continue
		}

		client := testAccProvider.Meta().(*UCloudClient)
		d, err := client.describeEIPById(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if isNotFoundError(err) {
				continue
			}
			return err
		}

		if d.EIPId != "" {
			return fmt.Errorf("EIP still exist")
		}
	}

	return nil
}

const testAccEIPConfig = `
resource "ucloud_eip" "foo" {
	name          = "tf-acc-eip"
	bandwidth     = 1
	internet_type = "bgp"
	charge_mode   = "bandwidth"
}
`
const testAccEIPConfigTwo = `
resource "ucloud_eip" "foo" {
	name          = "tf-acc-eip-two"
	bandwidth     = 2
	internet_type = "bgp"
	charge_mode   = "traffic"
}
`
