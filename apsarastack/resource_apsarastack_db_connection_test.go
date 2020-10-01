package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccApsaraStackDBConnectionConfigUpdate(t *testing.T) {
	var v *rds.DBInstanceNetInfo
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccDBconnection%s", rand)
	var basicMap = map[string]string{
		"instance_id":       CHECKSET,
		"connection_string": REGEXMATCH + fmt.Sprintf("^tf-testacc%s.mysql.([a-z-A-Z-0-9]+.){0,1}rds.aliyuncs.com", rand),
		"port":              "3306",
		"ip_address":        CHECKSET,
	}
	resourceId := "apsarastack_db_connection.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeDBConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBConnectionConfigDependence)
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
					"instance_id":       "${apsarastack_db_instance.instance.id}",
					"connection_prefix": fmt.Sprintf("tf-testacc%s", rand),
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
					"port": "3333",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port": "3333",
					}),
				),
			},
		},
	})
}

func resourceDBConnectionConfigDependence(name string) string {
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

	resource "apsarastack_db_instance" "instance" {
		engine = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}"
		engine_version = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}"
		instance_type = "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}"
		instance_storage = "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min}"
		vswitch_id = "${apsarastack_vswitch.default.id}"
		instance_name = "${var.name}"
	}
	`, RdsCommonTestCase, name)
}
