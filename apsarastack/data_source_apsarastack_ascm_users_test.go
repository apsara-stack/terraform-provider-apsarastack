package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccApsaraStackAscm_Users_DataSource(t *testing.T) { // not completed
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceApsaraStackAscm_User_Organization,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckApsaraStackDataSourceID("data.apsarastack_ascm_users.default"),
					resource.TestCheckResourceAttr("data.apsarastack_ascm_users.default", "users.#", "0"),
				/*	resource.TestCheckResourceAttrSet("data.apsarastack_ascm_users.default", "users.id"),
					    resource.TestCheckResourceAttrSet("data.apsarastack_ascm_users.default", "users.name"),
					    resource.TestCheckResourceAttrSet("data.apsarastack_ascm_users.default", "users.organization_id"),
						resource.TestCheckResourceAttrSet("data.apsarastack_ascm_users.default", "users.cell_phone_number"),
						resource.TestCheckResourceAttrSet("data.apsarastack_ascm_users.default", "users.display_name"),
						resource.TestCheckResourceAttrSet("data.apsarastack_ascm_users.default", "users.email"),
						resource.TestCheckResourceAttrSet("data.apsarastack_ascm_users.default", "users.mobile_nation_code"),
						resource.TestCheckResourceAttrSet("data.apsarastack_ascm_users.default", "users.login_policy_id"),*/

				),
			},
		},
	})
}

const dataSourceApsaraStackAscm_User_Organization = `
resource "apsarastack_ascm_user" "default" {
  cellphone_number = "899999567"
  email = "testing@mail.com"
  display_name = "C2C-DataSource"
  organization_id = "54437"
  mobile_nation_code = "91"
  login_name = "C2C_apsara_C2"
}
data "apsarastack_ascm_users" "default" {
  name_regex = apsarastack_ascm_user.default.login_name
}

`
