package apsarastack

/*func TestAccApsaraStackAscm_UserBasic(t *testing.T) {
	var v ecs.KeyPair
	resourceId := "apsarastack_ascm_user.default"
	ra := resourceAttrInit(resourceId, nil)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
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
		CheckDestroy:  rac.checkResourceDestroy(),
		//	CheckDestroy:  testAccCheckAscm_E_OrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config:testAccAscm_USer_Resource_Basic ,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

func testAccCheckAscm_UserDestroy(s *terraform.State) error { //destroy function
	client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "apsarastack_ascm_user" || rs.Type != "apsarastack_ascm_user" {
			continue
		}
		ascm, err := ascmService.DescribeAscmUser(rs.Primary.ID)
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

const testAccAscm_USer_Resource_Basic = `
resource "apsarastack_ascm_user" "user" {
  cellphone_number = "899999537"
  email = "test@gmail.com"
  display_name = "C2C-DEL3"
  organization_id = "54437"
  mobile_nation_code = "91"
  login_name = "C2C_apsara_C2C"
}`

*/
