package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackSecurityGroupsDataSourceBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApsaraStackSecurityGroupsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_security_groups.default"),

				),
			},
		},
	})
}

const testAccCheckApsaraStackSecurityGroupsDataSourceConfig = `

variable "name" {
  default = "tf-securityGroupdatasource"
}
resource "apsarastack_security_group" "group" {
  name        = var.name
  description = "foo"
  vpc_id      = "vpc-rt7ruq390e0yywjo7wgpr"
}
data "apsarastack_security_groups" "default" {
  ids = [apsarastack_security_group.group.id]
}
`
