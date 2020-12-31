package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_quotaDataSource(t *testing.T) { //not completed
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_Quota,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_quota.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_quota.default", "groups.id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_quota.default", "groups.quota_type"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_quota.default", "groups.quota_type_id"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_quota.default", "groups.region"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_quota.default", "groups.cluster_name"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_quota.default", "groups.used_vip_public"),
					resource.TestCheckNoResourceAttr("data.apsarastack_ascm_quota.default", "groups.allocate_vip_internal"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_Quota = `

data "apsarastack_ascm_quota" "default" {
  quota_type = "organization"
  quota_type_id = 54437
  product_name = "SLB"
}
`
