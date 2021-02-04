package apsarastack

import (
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccApsaraStackAscm_UserRoleBinding(t *testing.T) {
	var v *User
	resourceId := "apsarastack_ascm_user_role_binding.default"
	ra := resourceAttrInit(resourceId, testAccCheckUserRoleBinding)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		//CheckDestroy:  rac.checkResourceDestroy(),
		CheckDestroy: testAccCheckAscm_UserRoleBinding_Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAscm_UserRoleBinding,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckAscm_UserRoleBinding_Destroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "apsarastack_ascm_user_role_binding" || rs.Type != "apsarastack_ascm_user_role_binding" {
			continue
		}
		ascm, err := ascmService.DescribeAscmUserRoleBinding(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		if ascm.Message != "" {
			return WrapError(Error("resource  still exist"))
		}
	}

	return nil
}

const testAccCheckAscm_UserRoleBinding = `
resource "apsarastack_ascm_organization" "default" {
 name = "Test_binder"
 parent_id = "1"
}

resource "apsarastack_ascm_user" "default" {
 cellphone_number = "8675757834"
 email = "test@gmail.com"
 display_name = "C2C-DELTA"
 organization_id = apsarastack_ascm_organization.default.org_id
 mobile_nation_code = "91"
 login_name = "User_Role_Test"
 login_policy_id = 1
}

resource "apsarastack_ascm_user_role_binding" "default" {
  role_id = 5
  login_name = apsarastack_ascm_user.default.login_name
}
`

var testAccCheckUserRoleBinding = map[string]string{
	"login_name": CHECKSET,
	"role_id":    CHECKSET,
}
