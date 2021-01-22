package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_Enviroment_DataSource(t *testing.T) { //not completed
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_Enviromenttest,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_environment_services.default"),
					/*		resource.TestCheckResourceAttrSet("data.apsarastack_ascm_environment_services.default", "ids.#"),*/
					//resource.TestCheckNoResourceAttr("data.apsarastack_ascm_specific_fields.default", "group_filed"),
				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_Enviromenttest = `

data "apsarastack_ascm_environment_services" "default" {

}
`
