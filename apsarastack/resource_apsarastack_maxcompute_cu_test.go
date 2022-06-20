package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackMaxcomputeCu(t *testing.T) {
	resourceId := "apsarastack_maxcompute_cu.default"
	ra := resourceAttrInit(resourceId, nil)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf_testAccApsaraStack%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// Currently does not support creating projects with sub-accounts
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMaxcomputeCu, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cu_name":      name,
						"cu_num":       "1",
						"cluster_name": CHECKSET,
					}),
				),
			},
		},
	})
}

const testAccMaxcomputeCu = `
data "apsarastack_maxcompute_clusters" "default"{
	name_regex = "HYBRIDODPSCLUSTER-.*"
}

resource "apsarastack_maxcompute_cu" "default"{
  cu_name      = "%s"
  cu_num       = "1"
  cluster_name = data.apsarastack_maxcompute_clusters.default.clusters.0.cluster
}
`
