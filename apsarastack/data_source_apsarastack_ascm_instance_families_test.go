package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_Instance_families_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_Instance_families,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_ecs_instance_families.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ecs_instance_families.default", "families.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ecs_instance_families.default", "families.status"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_ecs_instance_families.default", "families.resource_type"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_Instance_families = `

data "apsarastack_ascm_ecs_instance_families" "default" {
status = "Available"
}
`
