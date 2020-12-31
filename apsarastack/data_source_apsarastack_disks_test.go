package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackDisksDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackDisksDataSourceConfigWithCommon,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_disks.default"),
					resource.TestCheckResourceAttr("data.apsarastack_disks.default", "disks.#", "1"),
					resource.TestCheckResourceAttrSet("data.apsarastack_disks.default", "disks.0.name"),
					resource.TestCheckResourceAttrSet("data.apsarastack_disks.default", "disks.0.description"),
					resource.TestCheckResourceAttrSet("data.apsarastack_disks.default", "disks.0.status"),
					resource.TestCheckResourceAttrSet("data.apsarastack_disks.default", "disks.0.size"),
					resource.TestCheckResourceAttrSet("data.apsarastack_disks.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckApsaraStackDisksDataSourceConfigWithCommon = `

data "apsarastack_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "apsarastack_disk" "default" {
  availability_zone = data.apsarastack_zones.default.zones[0].id
  name              = "tf-testdisk"
  description       = "ECS-Disk"
  category          = "cloud_efficiency"
  size              = "30"

  tags = {
    Name = "TerraformTest"
  }
}


data "apsarastack_disks" "default" {
	ids = [apsarastack_disk.default.id]
}
`
