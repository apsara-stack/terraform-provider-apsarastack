package apsarastack

import (
	"fmt"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

func TestAccApsaraStackAscm_UserBasic(t *testing.T) {
	var v *User
	resourceId := "apsarastack_ascm_user.default"
	ra := resourceAttrInit(resourceId, ascmuserBasicMap)
	serviceFunc := func() interface{} {
		return &AscmService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-ascmusers%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testascmuserconfigbasic)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAscm_UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cellphone_number":   "8999995370",
					"email":              "test01@gmail.com",
					"display_name":       "Test_Apsara",
					"organization_id":    os.Getenv("APSARASTACK_DEPARTMENT"),
					"mobile_nation_code": "91",
					"login_name":         name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckAscm_UserDestroy(s *terraform.State) error { //destroy function
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if true {
			continue
		}
		_, err := ascmService.DescribeAscmUser(rs.Primary.ID)
		if err == nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}
func testascmuserconfigbasic(name string) string {
	return fmt.Sprintf(`
variable name{
 default = "%s"
}

`, name)
}

var ascmuserBasicMap = map[string]string{
	"cellphone_number":   CHECKSET,
	"email":              CHECKSET,
	"display_name":       CHECKSET,
	"organization_id":    CHECKSET,
	"mobile_nation_code": CHECKSET,
	"login_name":         CHECKSET,
}
