package apsarastack

import (
	"fmt"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccApsaraStackAscm_UserGroupRoleBinding(t *testing.T) {
	var v *UserGroup
	resourceId := "apsarastack_ascm_user_group_role_binding.default"
	ra := resourceAttrInit(resourceId, testAccCheckUserGroupRoleBinding)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	//rand := acctest.RandInt()
	//name := fmt.Sprintf("tf-ascmusergroup%v", rand)
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
		CheckDestroy: testAccCheckAscm_UserGroupRoleBinding_Destroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAscm_UserGroupRoleBinding),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckAscm_UserGroupRoleBinding_Destroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "apsarastack_ascm_user_group_role_binding" || rs.Type != "apsarastack_ascm_user_group_role_binding" {
			continue
		}
		ascm, err := ascmService.DescribeAscmUserGroup(rs.Primary.ID)
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

const testAccCheckAscm_UserGroupRoleBinding = `

resource "apsarastack_ascm_user_group_role_binding" "default" {
  role_ids = [5]
  user_group_id = "82"
}
`

var testAccCheckUserGroupRoleBinding = map[string]string{
	"user_group_id": CHECKSET,
}
