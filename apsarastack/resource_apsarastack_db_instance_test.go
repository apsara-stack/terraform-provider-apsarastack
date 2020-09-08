package apsarastack

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("apsarastack_db_instance", &resource.Sweeper{
		Name: "apsarastack_db_instance",
		F:    testSweepDBInstances,
	})
}

func testSweepDBInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting ApsaraStack client: %s", err)
	}
	client := rawClient.(*connectivity.ApsaraStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []rds.DBInstance
	req := rds.CreateDescribeDBInstancesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeDBInstances(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving RDS Instances: %s", err)
		}
		resp, _ := raw.(*rds.DescribeDBInstancesResponse)
		if resp == nil || len(resp.Items.DBInstance) < 1 {
			break
		}
		insts = append(insts, resp.Items.DBInstance...)

		if len(resp.Items.DBInstance) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	sweeped := false
	vpcService := VpcService{client}
	for _, v := range insts {
		name := v.DBInstanceDescription
		id := v.DBInstanceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
		if skip {
			if need, err := vpcService.needSweepVpc(v.VpcId, v.VSwitchId); err == nil {
				skip = !need
			}

		}

		if skip {
			log.Printf("[INFO] Skipping RDS Instance: %s (%s)", name, id)
			continue
		}

		log.Printf("[INFO] Deleting RDS Instance: %s (%s)", name, id)
		if len(v.ReadOnlyDBInstanceIds.ReadOnlyDBInstanceId) > 0 {
			request := rds.CreateReleaseReadWriteSplittingConnectionRequest()
			request.DBInstanceId = id
			if _, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ReleaseReadWriteSplittingConnection(request)
			}); err != nil {
				log.Printf("[ERROR] ReleaseReadWriteSplittingConnection error: %#v", err)
			} else {
				time.Sleep(5 * time.Second)
			}
		}
		req := rds.CreateDeleteDBInstanceRequest()
		req.DBInstanceId = id
		_, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DeleteDBInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete RDS Instance (%s (%s)): %s", name, id, err)
		} else {
			sweeped = true
		}
	}
	if sweeped {
		// Waiting 30 seconds to eusure these DB instances have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccApsaraStackDBInstanceMysql(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	var ips []map[string]interface{}

	resourceId := "apsarastack_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigDependence)
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
					"engine":           "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":   "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":    "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage": "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_name":    "${var.name}",
					"vswitch_id":       "${apsarastack_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":           "MySQL",
						"engine_version":   "5.6",
						"instance_type":    CHECKSET,
						"instance_storage": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_restart"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time": "22:00Z-02:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time": "22:00Z-02:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_storage": "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min + data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.step}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_storage": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "tf-testAccDBInstance_instance_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testAccDBInstance_instance_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.apsarastack_db_instance_classes.default.instance_classes.1.instance_class}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(testAccCheckSecurityIpExists("apsarastack_db_instance.default", ips)),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": "${apsarastack_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    CHECKSET,
						"security_group_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       REMOVEKEY,
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":               "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":       "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":        "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":     "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min * 3}",
					"instance_name":        "tf-testAccDBInstanceConfig",
					"instance_charge_type": "Postpaid",
					"security_group_ids":   []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "MySQL",
						"engine_version":       "5.6",
						"instance_type":        CHECKSET,
						"instance_storage":     "15",
						"instance_name":        "tf-testAccDBInstanceConfig",
						"zone_id":              CHECKSET,
						"connection_string":    CHECKSET,
						"port":                 CHECKSET,
						"security_group_ids.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_mode": SafetyMode,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_mode": SafetyMode,
					}),
				),
			},
		},
	})
}

func resourceDBInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "%s"
}
variable "creation" {
		default = "Rds"
}

data "apsarastack_db_instance_engines" "default" {
  instance_charge_type = "PostPaid"
  engine               = "MySQL"
  engine_version       = "5.6"
}

data "apsarastack_db_instance_classes" "default" {
  instance_charge_type = "PostPaid"
  engine               = "MySQL"
  engine_version       = "5.6"
}

resource "apsarastack_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${apsarastack_vpc.default.id}"
}
`, RdsCommonTestCase, name)
}

func TestAccApsaraStackDBInstanceMultiInstance(t *testing.T) {
	var instance *rds.DBInstanceAttribute

	resourceId := "apsarastack_db_instance.default.4"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceConfigDependence)

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
					"count":            "5",
					"engine":           "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":   "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":    "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage": "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_name":    "${var.name}",
					"vswitch_id":       "${apsarastack_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

// Unknown current resource exists
func TestAccApsaraStackDBInstanceSQLServer(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	var ips []map[string]interface{}

	resourceId := "apsarastack_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceSQLServerConfigDependence)
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
					"engine":           "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":   "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":    "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage": "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"instance_name":    "${var.name}",
					"vswitch_id":       "${apsarastack_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":           "SQLServer",
						"engine_version":   "2012",
						"instance_type":    CHECKSET,
						"instance_storage": "20",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "tf-testAccDBInstance_instance_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testAccDBInstance_instance_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.apsarastack_db_instance_classes.default.instance_classes.1.instance_class}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"monitoring_period": "300",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"monitoring_period": "300",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": "${apsarastack_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    CHECKSET,
						"security_group_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":             "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":     "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":      "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":   "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min + data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"instance_name":      "${var.name}",
					"vswitch_id":         "${apsarastack_vswitch.default.id}",
					"security_group_ids": []string{},
					"monitoring_period":  "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "SQLServer",
						"engine_version":       "2012",
						"instance_type":        CHECKSET,
						"instance_storage":     "25",
						"instance_name":        "tf-testAccDBInstanceConfig",
						"zone_id":              CHECKSET,
						"connection_string":    CHECKSET,
						"port":                 CHECKSET,
						"security_group_ids.#": "0",
					}),
				),
			},
		},
	})
}

func resourceDBInstanceSQLServerConfigDependence(name string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "%s"
}
variable "creation" {
		default = "Rds"
}

data "apsarastack_db_instance_engines" "default" {
  instance_charge_type = "PostPaid"
  engine               = "SQLServer"
  engine_version       = "2012"
}

data "apsarastack_db_instance_classes" "default" {
  instance_charge_type = "PostPaid"
  engine               = "SQLServer"
  engine_version       = "2012"
}

resource "apsarastack_security_group" "default" {
	count = 2
	name   = "${var.name}"
	vpc_id = "${apsarastack_vpc.default.id}"
}
`, RdsCommonTestCase, name)
}

func TestAccApsaraStackDBInstancePostgreSQL(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	var ips []map[string]interface{}

	resourceId := "apsarastack_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstancePostgreSQLConfigDependence)
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
					"engine":           "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":   "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":    "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage": "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"zone_id":          "${data.apsarastack_db_instance_classes.default.instance_classes.0.zone_ids.0.id}",
					"instance_name":    "${var.name}",
					"vswitch_id":       "${apsarastack_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":           "PostgreSQL",
						"engine_version":   "9.4",
						"instance_storage": "20",
						"instance_type":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "tf-testAccDBInstance_instance_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testAccDBInstance_instance_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.apsarastack_db_instance_classes.default.instance_classes.1.instance_class}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": "${apsarastack_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    CHECKSET,
						"security_group_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":             "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":     "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":      "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":   "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min + data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"instance_name":      "${var.name}",
					"vswitch_id":         "${apsarastack_vswitch.default.id}",
					"security_group_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "PostgreSQL",
						"engine_version":       "9.4",
						"instance_type":        CHECKSET,
						"instance_storage":     "25",
						"instance_name":        "tf-testAccDBInstanceConfig",
						"zone_id":              CHECKSET,
						"connection_string":    CHECKSET,
						"port":                 CHECKSET,
						"security_group_ids.#": "0",
					}),
				),
			},
		},
	})
}

func resourceDBInstancePostgreSQLConfigDependence(name string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "%s"
}
variable "creation" {
		default = "Rds"
}

data "apsarastack_db_instance_engines" "default" {
  	instance_charge_type = "PostPaid"
  	engine               = "PostgreSQL"
  	engine_version       = "9.4"
	multi_zone           = true
}

data "apsarastack_db_instance_classes" "default" {
  	instance_charge_type = "PostPaid"
  	engine               = "PostgreSQL"
  	engine_version       = "9.4"
  	multi_zone           = true
}

resource "apsarastack_security_group" "default" {
	count = 2
	name   = var.name
	vpc_id = apsarastack_vpc.default.id
}
`, RdsCommonTestCase, name)
}

// Unknown current resource exists
func TestAccApsaraStackDBInstancePPAS(t *testing.T) {
	var instance *rds.DBInstanceAttribute
	var ips []map[string]interface{}

	resourceId := "apsarastack_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceAZConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.RdsPPASNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":           "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":   "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":    "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage": "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"zone_id":          "${data.apsarastack_db_instance_classes.default.instance_classes.0.zone_ids.0.id}",
					"instance_name":    "${var.name}",
					"vswitch_id":       "${apsarastack_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":           "PPAS",
						"engine_version":   "9.3",
						"instance_storage": "250",
						"instance_type":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "tf-testAccDBInstance_instance_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": "tf-testAccDBInstance_instance_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.apsarastack_db_instance_classes.default.instance_classes.1.instance_class}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12", "100.69.7.112"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyValueInMaps(ips, "security ip", "security_ips", "10.168.1.12,100.69.7.112"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_ids": "${apsarastack_security_group.default.*.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id":    CHECKSET,
						"security_group_ids.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":             "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":     "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":      "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage":   "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min + data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.step}",
					"instance_name":      "${var.name}",
					"vswitch_id":         "${apsarastack_vswitch.default.id}",
					"security_group_ids": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine":               "PPAS",
						"engine_version":       "9.3",
						"instance_type":        CHECKSET,
						"instance_storage":     "500",
						"instance_name":        "tf-testAccDBInstanceConfig",
						"zone_id":              CHECKSET,
						"connection_string":    CHECKSET,
						"port":                 CHECKSET,
						"security_group_ids.#": "0",
					}),
				),
			},
		},
	})
}

func resourceDBInstanceAZConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "apsarastack_db_instance_engines" "default" {
  instance_charge_type = "PostPaid"
  engine               = "PPAS"
  engine_version       = "9.3"
  multi_zone           = true
}

data "apsarastack_db_instance_classes" "default" {
  instance_charge_type = "PostPaid"
  engine               = "PPAS"
  engine_version       = "9.3"
  multi_zone           = true
}

resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_db_instance_classes.default.instance_classes.0.zone_ids.0.sub_zone_ids.0}"
  name              = "${var.name}"
}

resource "apsarastack_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${apsarastack_vpc.default.id}"
}
`, name)
}

// Unknown current resource exists
func TestAccApsaraStackDBInstanceMultiAZ(t *testing.T) {
	var instance = &rds.DBInstanceAttribute{}
	resourceId := "apsarastack_db_instance.default"
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDBInstance")
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstance_multiAZ"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceMysqlAZConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.RdsMultiAzNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":           "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":   "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":    "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage": "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"zone_id":          "${data.apsarastack_db_instance_classes.default.instance_classes.0.zone_ids.0.id}",
					"instance_name":    "${var.name}",
					"vswitch_id":       "${apsarastack_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":       REGEXMATCH + ".*" + MULTI_IZ_SYMBOL + ".*",
						"instance_name": "tf-testAccDBInstance_multiAZ",
					}),
				),
			},
		},
	})

}

func resourceDBInstanceMysqlAZConfigDependence(name string) string {
	return fmt.Sprintf(`
%s
variable "name" {
	default = "%s"
}
variable "creation" {
		default = "Rds"
}
data "apsarastack_db_instance_engines" "default" {
  	engine               = "MySQL"
  	engine_version       = "5.6"
	multi_zone           = true
}

data "apsarastack_db_instance_classes" "default" {
  	engine               = "MySQL"
  	engine_version       = "5.6"
	multi_zone           = true
}
resource "apsarastack_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${apsarastack_vpc.default.id}"
}
`, RdsCommonTestCase, name)
}

func TestAccApsaraStackDBInstanceClassic(t *testing.T) {
	var instance *rds.DBInstanceAttribute

	resourceId := "apsarastack_db_instance.default"
	ra := resourceAttrInit(resourceId, instanceBasicMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &instance, func() interface{} {
		return &RdsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeDBInstance")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tf-testAccDBInstanceConfig"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDBInstanceClassicConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.RdsClassicNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"engine":           "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine}",
					"engine_version":   "${data.apsarastack_db_instance_engines.default.instance_engines.0.engine_version}",
					"instance_type":    "${data.apsarastack_db_instance_classes.default.instance_classes.0.instance_class}",
					"instance_storage": "${data.apsarastack_db_instance_classes.default.instance_classes.0.storage_range.min}",
					"zone_id":          `${lookup(data.apsarastack_db_instance_classes.default.instance_classes.0.zone_ids[length(data.apsarastack_db_instance_classes.default.instance_classes.0.zone_ids)-1], "id")}`,
					"instance_name":    "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func resourceDBInstanceClassicConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "apsarastack_db_instance_engines" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
}

data "apsarastack_db_instance_classes" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
}

data "apsarastack_zones" "default" {
  	available_resource_creation= "Rds"
}`, name)
}

func testAccCheckSecurityIpExists(n string, ips []map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DB Instance ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.ApsaraStackClient)
		rdsService := RdsService{client}
		resp, err := rdsService.DescribeDBSecurityIps(rs.Primary.ID)
		log.Printf("[DEBUG] check instance %s security ip %#v", rs.Primary.ID, resp)

		if err != nil {
			return err
		}

		if len(resp) < 1 {
			return fmt.Errorf("DB security ip not found")
		}

		ips = rdsService.flattenDBSecurityIPs(resp)
		return nil
	}
}

func testAccCheckKeyValueInMaps(ps []map[string]interface{}, propName, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, policy := range ps {
			if policy[key].(string) != value {
				return fmt.Errorf("DB %s attribute '%s' expected %#v, got %#v", propName, key, value, policy[key])
			}
		}
		return nil
	}
}

var instanceBasicMap = map[string]string{
	"engine":            "MySQL",
	"engine_version":    "5.6",
	"instance_type":     CHECKSET,
	"instance_storage":  "5",
	"instance_name":     "tf-testAccDBInstanceConfig",
	"zone_id":           CHECKSET,
	"connection_string": CHECKSET,
	"port":              CHECKSET,
}
