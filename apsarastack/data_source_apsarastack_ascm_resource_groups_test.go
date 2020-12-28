package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_Resource_Groups_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_Resource_Group_Organization,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_resource_groups.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_resource_groups.default", "groups.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_resource_groups.default", "groups.name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_resource_groups.default", "groups.organization_id"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_Resource_Group_Organization = `
resource "apsarastack_ascm_resource_group" "default" {
  organization_id = "54438"
  name = "apsarastack-Datasource-resourceGroup"
}
data "apsarastack_ascm_resource_groups" "default" {
  name_regex = apsarastack_ascm_resource_group.default.name
}
`
