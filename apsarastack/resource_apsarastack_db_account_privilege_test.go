package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccApsaraStackDBAccountPrivilege_mysql(t *testing.T) {

	var v *rds.DBInstanceAccount
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sdnsrecordbasic%v.abc", defaultRegionToTest, rand)
	resourceId := "apsarastack_db_account_privilege.default"
	var basicMap = map[string]string{
		"instance_id":  CHECKSET,
		"account_name": "tftestprivilege",
		"privilege":    "ReadOnly",
		"db_names.#":   "2",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBAccountPrivilege")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBAccountPrivilegeConfigDependenceForMySql)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${apsarastack_db_instance.default.id}",
					"account_name": "${apsarastack_db_account.default.name}",
					"privilege":    "ReadOnly",
					"db_names":     "${apsarastack_db_database.default.*.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${apsarastack_db_instance.default.id}",
					"account_name": "${apsarastack_db_account.default.name}",
					"privilege":    "ReadOnly",
					"db_names":     []string{"${apsarastack_db_database.default.0.name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_names.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${apsarastack_db_instance.default.id}",
					"account_name": "${apsarastack_db_account.default.name}",
					"privilege":    "ReadOnly",
					"db_names":     "${apsarastack_db_database.default.*.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_names.#": "2",
					}),
				),
			},
		},
	})

}

func TestAccApsaraStackDBAccountPrivilege_PostgreSql(t *testing.T) {

	var v *rds.DBInstanceAccount
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%sdnsrecordbasic%v.abc", defaultRegionToTest, rand)
	resourceId := "apsarastack_db_account_privilege.default"
	var basicMap = map[string]string{
		"instance_id":  CHECKSET,
		"account_name": "tftestprivilege",
		"privilege":    "DBOwner",
		"db_names.#":   "1",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBAccountPrivilege")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBAccountPrivilegeConfigDependenceForPostgreSql)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${apsarastack_db_instance.default.id}",
					"account_name": "${apsarastack_db_account.default.name}",
					"privilege":    "DBOwner",
					"db_names":     []string{"${apsarastack_db_database.default.0.name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${apsarastack_db_instance.default.id}",
					"account_name": "${apsarastack_db_account.default.name}",
					"privilege":    "DBOwner",
					"db_names":     "${apsarastack_db_database.default.*.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_names.#": "2",
					}),
				),
			},
		},
	})

}

func resourceDBAccountPrivilegeConfigDependenceForMySql(name string) string {
	return fmt.Sprintf(`
%s
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "%s"
	}

	data "apsarastack_db_instance_engines" "default" {
  		instance_charge_type = "PostPaid"
  		engine               = "MySQL"
  		engine_version       = "5.6"
	}

	data "apsarastack_db_instance_classes" "default" {
 	 	engine = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}"
		engine_version = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}"
	}

	resource "apsarastack_db_instance" "default" {
		engine = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}"
		engine_version = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}"
		instance_type = "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}"
		instance_storage = "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min}"
		vswitch_id = "${apsarastack_vswitch.default.id}"
		instance_name = "${var.name}"
	}

	resource "apsarastack_db_database" "default" {
	  count = 2
	  instance_id = "${apsarastack_db_instance.default.id}"
	  name = "tfaccountpri_${count.index}"
	  description = "from terraform"
	}

	resource "apsarastack_db_account" "default" {
	  instance_id = "${apsarastack_db_instance.default.id}"
	  name = "tftestprivilege"
	  password = "Test12345"
	  description = "from terraform"
	}
`, RdsCommonTestCase, name)
}

func resourceDBAccountPrivilegeConfigDependenceForPostgreSql(name string) string {
	return fmt.Sprintf(`
%s
	variable "creation" {
		default = "Rds"
	}

	variable "name" {
		default = "%s"
	}
	
	data "apsarastack_db_instance_classes" "default" {
		instance_charge_type = "PostPaid"
		engine               = "PostgreSQL"
		engine_version       = "10.0"
		storage_type         = "cloud_ssd"
	}

	resource "apsarastack_db_instance" "default" {
		engine = "PostgreSQL"
		engine_version = "10.0"
		instance_type = "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}"
		instance_storage = "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min}"
		vswitch_id = "${apsarastack_vswitch.default.id}"
		instance_name = "${var.name}"
	}

	resource "apsarastack_db_database" "default" {
	  count = 2
	  instance_id = "${apsarastack_db_instance.default.id}"
	  name = "tfaccountpri_${count.index}"
	  description = "from terraform"
      character_set = "UTF8"
	}

	resource "apsarastack_db_account" "default" {
	  instance_id = "${apsarastack_db_instance.default.id}"
	  name = "tftestprivilege"
	  password = "Test12345"
	  description = "from terraform"
	}
`, RdsCommonTestCase, name)
}
