package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_Roles_DataSource(t *testing.T) { // not completed
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_Roles_Organization,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_roles.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_roles.default", "roles.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_roles.default", "roles.name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_roles.default", "roles.role_level"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_roles.default", "roles.role_type"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_Roles_Organization = `


data "apsarastack_ascm_roles" "default" {
  name_regex = "datahub_full_access"
}
`
