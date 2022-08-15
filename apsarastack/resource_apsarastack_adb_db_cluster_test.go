package apsarastack

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("ApsaraStack_adb_db_instance", &resource.Sweeper{
		Name: "ApsaraStack_adb_db_instance",
		F:    testSweepAdbDbInstances,
	})
}

func testSweepAdbDbInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting ApsaraStack client: %s", err)
	}
	client := rawClient.(*connectivity.ApsaraStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	action := "DescribeDBClusters"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var response map[string]interface{}
	conn, err := client.NewAdsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			log.Println(WrapErrorf(err, DataDefaultErrorMsg, "ApsaraStack_adb_db_clusters", action, ApsaraStackSdkGoERROR))
			break
		}

		resp, err := jsonpath.Get("$.Items.DBCluster", response)
		if err != nil {
			log.Println(WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items.DBCluster", response))
			break
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["DBClusterDescription"])
			id := fmt.Sprint(item["DBClusterId"])
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping ADB Instance: %s (%s)", name, id)
				continue
			}
			log.Printf("[INFO] Deleting adb Instance: %s (%s)", name, id)
			action := "DeleteDBCluster"
			conn, err := client.NewAdsClient()
			if err != nil {
				log.Println(WrapError(err))
				break
			}
			request := map[string]interface{}{
				"DBClusterId": id,
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				log.Printf("[ERROR] Deleting ADB cluster failed with error: %#v", err)
				return nil
			})
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

// 316 基础版 ClusterType:analyticdb
func TestAccApsaraStackAdbDbCluster_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "apsarastack_adb_db_cluster.default"
	ra := resourceAttrInit(resourceId, ApsaraStackAdbDbClusterMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeAdbDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sadbCluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackAdbDbClusterBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.AdbReserverUnSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_category": "Basic",
					"db_cluster_version":  "3.0",
					"db_node_class":       "C8",
					"description":         "${var.name}",
					"db_node_count":       "2",
					"db_node_storage":     "200",
					"mode":                "reserver",
					"vswitch_id":          "${apsarastack_vswitch.default.id}",
					"cluster_type":        "analyticdb",
					"cpu_type":            "intel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_category": "Basic",
						"db_node_class":       "C8",
						"description":         name,
						"db_node_count":       "2",
						"db_node_storage":     "200",
						"mode":                "reserver",
						"vswitch_id":          CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_type", "cpu_type"},
			},
			/*{
				Config: testAccConfig(map[string]interface{}{
					"db_node_class": "C20",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_class": "C20",
					}),
				),
			},*/
			/*{
				Config: testAccConfig(map[string]interface{}{
					"db_node_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_count": "2",
					}),
				),
			},*/
			/*{
				Config: testAccConfig(map[string]interface{}{
					"db_node_storage": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_storage": "200",
					}),
				),
			},*/
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			/*{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time": "23:00Z-00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time": "23:00Z-00:00Z",
					}),
				),
			},*/
			/*{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.apsarastack_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},*/
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "1",
					}),
				),
			},
			/*{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},*/
			{
				Config: testAccConfig(map[string]interface{}{
					"db_node_class":   "C8",
					"db_node_count":   "1",
					"db_node_storage": "200",
					"description":     name,
					//"maintain_time":   "01:00Z-02:00Z",
					"security_ips": []string{"10.168.1.13"},
					/*"tags": map[string]string{
						"Created": "TF-update",
						"For":     "test-update",
					},*/
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_node_class":   "C8",
						"db_node_count":   "1",
						"db_node_storage": "200",
						"description":     name,
						//"maintain_time":   "01:00Z-02:00Z",
						"security_ips.#": "1",
						//"tags.%":          "2",
						//"tags.Created":    "TF-update",
						//"tags.For":        "test-update",
					}),
				),
			},
		},
	})
}

// 316 集群版 ClusterType:AnalyticdbOnPanguHybrid
func TestAccApsaraStackAdbDbCluster_flexible(t *testing.T) {
	var v map[string]interface{}
	resourceId := "apsarastack_adb_db_cluster.default"
	ra := resourceAttrInit(resourceId, ApsaraStackAdbDbClusterMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AdbService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeAdbDbCluster")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sadbCluster%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackAdbDbClusterBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_category": "cluster",
					"storage_resource":    "8Core45GB",
					"storage_type":        "SSD",
					"description":         "${var.name}",
					"mode":                "flexible",
					"compute_resource":    "8Core40GB",
					"db_node_count":       "2",
					"vswitch_id":          "${apsarastack_vswitch.default.id}",
					"cluster_type":        "AnalyticdbOnPanguHybrid",
					"cpu_type":            "intel",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_cluster_category": "cluster",
						"storage_resource":    "8Core45GB",
						"storage_type":        "SSD",
						"description":         name,
						"mode":                "flexible",
						"compute_resource":    "8Core40GB",
						"db_node_count":       "2",
						"vswitch_id":          CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cluster_type", "cpu_type", "storage_resource"},
			},
			// API does not support to updating the compute_resource
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"compute_resource": "16Core64GB",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"compute_resource": "16Core64GB",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"elastic_io_resource": "1",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"elastic_io_resource": "1",
			//		}),
			//	),
			//},

			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			/*{
				Config: testAccConfig(map[string]interface{}{
					"maintain_time": "23:00Z-00:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_time": "23:00Z-00:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					//"resource_group_id": "${data.apsarastack_resource_manager_resource_groups.default.ids.0}",
					"resource_group_id": "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},*/
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips": []string{"10.168.1.12"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips.#": "1",
					}),
				),
			},
			/*{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},*/
			{
				Config: testAccConfig(map[string]interface{}{
					"compute_resource": "8Core40GB",
					//"elastic_io_resource": "1",
					"description": name,
					//"maintain_time": "01:00Z-02:00Z",
					"security_ips": []string{"10.168.1.13"},
					//"tags": map[string]string{
					//	"Created": "TF-update",
					//	"For":     "test-update",
					//},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"compute_resource": "8Core40GB",
						//"elastic_io_resource": "1",
						"description": name,
						//"maintain_time":  "01:00Z-02:00Z",
						"security_ips.#": "0",
						//"tags.%":         "2",
						//"tags.Created":   "TF-update",
						//"tags.For":       "test-update",
					}),
				),
			},
		},
	})
}

var ApsaraStackAdbDbClusterMap0 = map[string]string{
	"auto_renew_period": NOSET,
	"compute_resource":  "",
	//"connection_string":  CHECKSET,
	"db_cluster_version": "3.0",
	"db_node_storage":    "0",
	//"elastic_io_resource": "0",
	"maintain_time":  CHECKSET,
	"modify_type":    NOSET,
	"payment_type":   "PayAsYouGo",
	"pay_type":       "PostPaid",
	"period":         NOSET,
	"renewal_status": NOSET,
	//"resource_group_id": CHECKSET,
	"security_ips.#": "1",
	"status":         "Running",
	//"tags.%":            "0",
	"zone_id": CHECKSET,
}

func ApsaraStackAdbDbClusterBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "apsarastack_ascm_resource_groups" "default" {
  name_regex = ""
}
%s
`, name, AdbCommonTestCase)
}

var ApsaraStackAdbDbClusterMap1 = map[string]string{
	"auto_renew_period": NOSET,
	"compute_resource":  "8Core40GB",
	//"connection_string":  CHECKSET,
	"db_cluster_version": "3.0",
	"db_node_class":      "B7",
	"db_node_count":      "1",
	"db_node_storage":    "500",
	//"elastic_io_resource": "0",
	"maintain_time":  CHECKSET,
	"modify_type":    NOSET,
	"payment_type":   "PayAsYouGo",
	"pay_type":       "PostPaid",
	"period":         NOSET,
	"renewal_status": NOSET,
	//"resource_group_id": CHECKSET,
	"security_ips.#": "1",
	"status":         "Running",
	//"tags.%":         "0",
	"zone_id": CHECKSET,
}

func ApsaraStackAdbDbClusterBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s
`, name, AdbCommonTestCase)
}
