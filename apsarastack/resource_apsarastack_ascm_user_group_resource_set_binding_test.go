package apsarastack

import (
	"fmt"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

func TestAccApsaraStackAscmUserGroupResourceSetBinding(t *testing.T) {
	var v *ListResourceGroup
	resourceId := "apsarastack_ascm_user_group_resource_set_binding.default"
	ra := resourceAttrInit(resourceId, testAccCheckUserGroupResourceSetBinding)
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
		CheckDestroy: testAccCheckAscmUserGroupResourceSetBindingDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAscmUserGroupResourceSetRoleBinding),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckAscmUserGroupResourceSetBindingDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "apsarastack_ascm_user_group_resource_set_binding" || rs.Type != "apsarastack_ascm_user_group_resource_set_binding" {
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

const testAccCheckAscmUserGroupResourceSetRoleBinding = `
resource "apsarastack_ascm_organization" "default" {
 name = "Test_binder"
 parent_id = "1"
}

resource "apsarastack_ascm_resource_group" "default" {
  organization_id = apsarastack_ascm_organization.default.org_id
  name = "apsarastack-terraform-resourceGroup"
}

resource "apsarastack_ascm_user_group_resource_set_binding" "default" {
  resource_set_id = apsarastack_ascm_resource_group.default.rg_id
  user_group_id = "82"
}
`

var testAccCheckUserGroupResourceSetBinding = map[string]string{
	"user_group_id":   CHECKSET,
	"resource_set_id": CHECKSET,
}
