package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var DBReadWriteMap = map[string]string{
	"port":              "3306",
	"distribution_type": "Standard",
	"weight":            NOSET,
	"max_delay_time":    "30",
	"instance_id":       CHECKSET,
	"connection_string": CHECKSET,
}

func TestAccApsaraStackDBReadWriteSplittingConnection_update(t *testing.T) {
	var connection = &rds.DBInstanceNetInfo{}
	var primary = &rds.DBInstanceAttribute{}
	var readonly = &rds.DBInstanceAttribute{}

	resourceId := "apsarastack_db_read_write_splitting_connection.default"
	ra := resourceAttrInit(resourceId, DBReadWriteMap)

	rc_connection := resourceCheckInitWithDescribeMethod(resourceId, &connection, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDBReadWriteSplittingConnection")
	rc_primary := resourceCheckInitWithDescribeMethod("apsarastack_db_instance.default", &primary, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDBInstance")
	rc_readonly := resourceCheckInitWithDescribeMethod("apsarastack_db_readonly_instance.default", &readonly, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDBReadonlyInstance")
	rand := acctest.RandIntRange(10000, 999999)

	rac := resourceAttrCheckInit(rc_connection, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	prefix := fmt.Sprintf("t-con-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, prefix, resourceDBReadWriteSplittingConfigDependence)
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
					"instance_id":       "${apsarastack_db_readonly_instance.default.master_db_instance_id}",
					"connection_prefix": "${var.prefix}",
					"distribution_type": "Standard",
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
					"max_delay_time":    "300",
					"distribution_type": "Custom",
					"weight": `${map(
						"${apsarastack_db_instance.default.id}", "0",
						"${apsarastack_db_readonly_instance.default.id}", "500"
					)}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					rc_primary.checkResourceExists(),
					rc_readonly.checkResourceExists(),
					testAccCheck(map[string]string{
						"max_delay_time":    "300",
						"weight.%":          "2",
						"distribution_type": "Custom",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       "${apsarastack_db_readonly_instance.default.master_db_instance_id}",
					"connection_prefix": "${var.prefix}",
					"distribution_type": "Standard",
					"max_delay_time":    "30",
					"weight":            REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"port":              "3306",
						"distribution_type": "Standard",
						"weight.%":          REMOVEKEY,
						"max_delay_time":    "30",
						"instance_id":       CHECKSET,
						"connection_string": CHECKSET,
					}),
				),
			},
		},
	})
}

func resourceDBReadWriteSplittingConfigDependence(prefix string) string {
	return fmt.Sprintf(`
	%s
	variable "creation" {
		default = "Rds"
	}
	variable "multi_az" {
		default = "false"
	}
	variable "name" {
		default = "tf-testAccDBInstance_vpc"
	}

	variable "prefix" {
		default = "%s"
	}

	data "apsarastack_db_instance_engines" "default" {
  		instance_charge_type = "PostPaid"
		engine = "MySQL"
		engine_version = "5.6"
	}

	data "apsarastack_db_instance_classes" "default" {
  		instance_charge_type = "PostPaid"
  		engine               = "MySQL"
  		engine_version       = "5.6"
	}

	resource "apsarastack_db_instance" "default" {
		engine = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}"
		engine_version = "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}"
		instance_type = "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}"
		instance_storage = "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min}"
		instance_charge_type = "Postpaid"
		instance_name = "${var.name}"
		vswitch_id = "${apsarastack_vswitch.default.id}"
		security_ips = ["10.168.1.12", "100.69.7.112"]
	}

	resource "apsarastack_db_readonly_instance" "default" {
		master_db_instance_id = "${apsarastack_db_instance.default.id}"
		zone_id = "${apsarastack_db_instance.default.zone_id}"
		engine_version = "${apsarastack_db_instance.default.engine_version}"
		instance_type = "${apsarastack_db_instance.default.instance_type}"
		instance_storage = "${apsarastack_db_instance.default.instance_storage}"
		instance_name = "${var.name}_ro"
		vswitch_id = "${apsarastack_vswitch.default.id}"
	}
`, RdsCommonTestCase, prefix)
}
