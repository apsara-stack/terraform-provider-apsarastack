package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackCms_Projectmeta_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackcms_project,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_cms_project_meta.default"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_project_meta.default", "resources.description"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_project_meta.default", "resources.labels"),
					resource.TestCheckNoResourceAttr("data.apsarastack_cms_project_meta.default", "resources.namespace"),
				),
			},
		},
	})
}

const dataSourceApsaraStackcms_project = `
data "apsarastack_cms_project_meta" "default" {
}
`
