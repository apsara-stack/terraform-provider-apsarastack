package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackSlbsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackSlbsDataSource,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_slbs.default"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

const testAccCheckApsaraStackSlbsDataSource = `
variable "name" {
	default = "tf-SlbDataSourceSlbsx"
}

resource "apsarastack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "vsw-rt7pgwfikxd2g3ujwtppt"
 tags = {
           Created = "TF"
           For = "Test"
         }
}
resource "apsarastack_slb" "defaultt" {
  name = "${var.name}"
  vswitch_id = "vsw-rt7pgwfikxd2g3ujwtppt"
 tags = {
           Created = "TF"
           For = "Test"
         }
}
data "apsarastack_slbs" "default" {
 ids = ["${apsarastack_slb.default.id}","${apsarastack_slb.defaultt.id}"]
}
`
