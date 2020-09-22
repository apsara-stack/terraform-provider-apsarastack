package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"log"
	"os"
	"testing"

	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("apsarastack_instance", &resource.Sweeper{
		Name: "apsarastack_instance",
		F:    testSweepInstances,
		// When implemented, these should be removed firstly
		// Now, the resource apsarastack_havip_attachment has been published.
		//Dependencies: []string{
		//	"apsarastack_havip_attachment",
		//},
	})
}

func testSweepInstances(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Apsarastack client: %s", err)
	}
	client := rawClient.(*connectivity.ApsaraStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var insts []ecs.Instance
	req := ecs.CreateDescribeInstancesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Instances: %s", err)
		}
		resp, _ := raw.(*ecs.DescribeInstancesResponse)
		if resp == nil || len(resp.Instances.Instance) < 1 {
			break
		}
		insts = append(insts, resp.Instances.Instance...)

		if len(resp.Instances.Instance) < PageSizeLarge {
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
		name := v.InstanceName
		id := v.InstanceId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		// If a slb name is set by other service, it should be fetched by vswitch name and deleted.
		if skip {
			if need, err := vpcService.needSweepVpc(v.VpcAttributes.VpcId, v.VpcAttributes.VSwitchId); err == nil {
				skip = !need
			}

		}
		if skip {
			log.Printf("[INFO] Skipping Instance: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Instance: %s (%s)", name, id)
		if v.DeletionProtection {
			request := ecs.CreateModifyInstanceAttributeRequest()
			request.InstanceId = id
			request.DeletionProtection = requests.NewBoolean(false)
			_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceAttribute(request)
			})
			if err != nil {
				log.Printf("[ERROR] %#v", WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ApsaraStackSdkGoERROR))
				continue
			}
		}
		if v.InstanceChargeType == string(PrePaid) {
			request := ecs.CreateModifyInstanceChargeTypeRequest()
			request.InstanceIds = convertListToJsonString(append(make([]interface{}, 0, 1), id))
			request.InstanceChargeType = string(PostPaid)
			request.IncludeDataDisks = requests.NewBoolean(true)
			_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceChargeType(request)
			})
			if err != nil {
				log.Printf("[ERROR] %#v", WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ApsaraStackSdkGoERROR))
				continue
			}
			time.Sleep(3 * time.Second)
		}

		req := ecs.CreateDeleteInstanceRequest()
		req.InstanceId = id
		req.Force = requests.NewBoolean(true)
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteInstance(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Instance (%s (%s)): %s", name, id, err)
		} else {
			sweeped = true
		}
	}
	if sweeped {
		// Waiting 20 seconds to eusure these instances have been deleted.
		time.Sleep(20 * time.Second)
	}
	return nil
}

func TestAccApsaraStackInstanceBasic(t *testing.T) {
	var v ecs.Instance

	resourceId := "apsarastack_instance.default"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsInstanceConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceInstanceBasicConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EcsClassicSupportedRegions)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":        "${data.apsarastack_images.default.images.0.id}",
					"security_groups": []string{"${apsarastack_security_group.default.0.id}"},
					"instance_type":   "${data.apsarastack_instance_types.default.instance_types.0.id}",

					"availability_zone":    "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}",
					"system_disk_category": "cloud_efficiency",
					"instance_name":        "${var.name}",
					"key_name":             "${apsarastack_key_pair.default.key_name}",

					"security_enhancement_strategy": "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"key_name":      name,
						"role_name":     NOSET,
						"vswitch_id":    REMOVEKEY,
						"user_data":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_enhancement_strategy", "dry_run"},
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"auto_release_time": time.Now().Add(10 * time.Hour).Format("2006-01-02T15:04:05Z"),
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"auto_release_time": CHECKSET,
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"auto_release_time": "",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"auto_release_time": "",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id": "${data.apsarastack_images.default.images.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_groups": []string{"${apsarastack_security_group.default.0.id}", "${apsarastack_security_group.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_groups.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": fmt.Sprintf("tf-testAccEcsInstanceConfigBasic%d_change", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": fmt.Sprintf("tf-testAccEcsInstanceConfigBasic%d_change", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": fmt.Sprintf("tf-testAccEcsInstanceConfigBasic%d_description", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": fmt.Sprintf("tf-testAccEcsInstanceConfigBasic%d_description", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_max_bandwidth_out": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_out": "50",
						"public_ip":                  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_max_bandwidth_in": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_in": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_name": "hostNameExample",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name": "hostNameExample",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "Password123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "Password123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "50",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"foo": "foo",
						"Bar": "Bar",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":   "2",
						"tags.foo": "foo",
						"tags.Bar": "Bar",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "60",
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"deletion_protection": "true",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"deletion_protection": "true",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":              "${data.apsarastack_instance_types.default.instance_types.0.id}",
					"security_groups":            []string{"${apsarastack_security_group.default.0.id}"},
					"instance_name":              fmt.Sprintf("tf-testAccEcsInstanceConfigBasic%d", rand),
					"description":                fmt.Sprintf("tf-testAccEcsInstanceConfigBasic%d", rand),
					"internet_max_bandwidth_out": REMOVEKEY,
					"internet_max_bandwidth_in":  REMOVEKEY,
					"host_name":                  REMOVEKEY,
					"password":                   REMOVEKEY,

					"system_disk_size": "70",
					"tags":             REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"tags.%":   "0",
						"tags.Bar": REMOVEKEY,
						"tags.foo": REMOVEKEY,

						"instance_name": fmt.Sprintf("tf-testAccEcsInstanceConfigBasic%d", rand),

						"image_id":          CHECKSET,
						"instance_type":     CHECKSET,
						"security_groups.#": "1",

						"availability_zone":             CHECKSET,
						"system_disk_category":          "cloud_efficiency",
						"security_enhancement_strategy": "Active",

						"description": fmt.Sprintf("tf-testAccEcsInstanceConfigBasic%d", rand),
						"host_name":   CHECKSET,
						"password":    "",

						"system_disk_size": "70",

						"private_ip": CHECKSET,
						"status":     "Running",

						"internet_max_bandwidth_in":  "50",
						"internet_max_bandwidth_out": "0",
					}),
				),
			},
		},
	})
}

func TestAccApsaraStackInstanceVpc(t *testing.T) {
	var v ecs.Instance

	resourceId := "apsarastack_instance.default"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAcc%sEcsInstanceConfigVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceInstanceVpcConfigDependence)

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
					"image_id":        "${data.apsarastack_images.default.images.0.id}",
					"security_groups": []string{"${apsarastack_security_group.default.0.id}"},
					"instance_type":   "${data.apsarastack_instance_types.default.instance_types.0.id}",

					"availability_zone":    "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}",
					"system_disk_category": "cloud_efficiency",
					"instance_name":        "${var.name}",
					"key_name":             "${apsarastack_key_pair.default.key_name}",
					//"spot_strategy":                 "NoSpot",
					//"spot_price_limit":              "0",
					"security_enhancement_strategy": "Active",
					"user_data":                     "I_am_user_data",

					"vswitch_id": "${apsarastack_vswitch.default.id}",
					"role_name":  "${apsarastack_ram_role.default.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"key_name":      name,
						"role_name":     name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_enhancement_strategy", "dry_run"},
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"auto_release_time": time.Now().Add(10 * time.Hour).Format("2006-01-02T15:04:05Z"),
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"auto_release_time": CHECKSET,
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"auto_release_time": "",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"auto_release_time": "",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id": "${data.apsarastack_images.default.images.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_data": "I_am_user_data_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_data": "I_am_user_data_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_groups": []string{"${apsarastack_security_group.default.0.id}", "${apsarastack_security_group.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_groups.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_max_bandwidth_out": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_out": "50",
						"private_ip":                 CHECKSET,
						//"public_ip":                  CHECKSET,
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"internet_charge_type": "PayByBandwidth",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"internet_charge_type": "PayByBandwidth",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_max_bandwidth_in": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_in": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_name": "hostNameExample",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name": "hostNameExample",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "Password123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "Password123",
					}),
				),
			},
			// only burstable instances support this attribute.
			/*{
				Config: testAccConfig(map[string]interface{}{
					"credit_specification": "Unlimited",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"credit_specification": "Unlimited",
					}),
				),
			},*/
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_ip": "172.16.0.10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ip": "172.16.0.10",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"foo": "foo",
						"Bar": "Bar",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":   "2",
						"tags.foo": "foo",
						"tags.Bar": "Bar",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":              "${data.apsarastack_instance_types.default.instance_types.0.id}",
					"security_groups":            []string{"${apsarastack_security_group.default.0.id}"},
					"instance_name":              name,
					"description":                name,
					"internet_max_bandwidth_out": REMOVEKEY,
					"internet_max_bandwidth_in":  REMOVEKEY,
					"host_name":                  REMOVEKEY,
					"password":                   REMOVEKEY,

					"system_disk_size": "70",
					"private_ip":       REMOVEKEY,
					"tags":             REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"tags.%":   "0",
						"tags.Bar": REMOVEKEY,
						"tags.foo": REMOVEKEY,

						"instance_name": name,

						"image_id":          CHECKSET,
						"instance_type":     CHECKSET,
						"security_groups.#": "1",

						"availability_zone":    CHECKSET,
						"system_disk_category": "cloud_efficiency",

						"security_enhancement_strategy": "Active",
						"vswitch_id":                    CHECKSET,

						"description":      name,
						"host_name":        CHECKSET,
						"password":         "",
						"system_disk_size": "70",

						"private_ip": CHECKSET,
						"status":     "Running",

						"internet_max_bandwidth_in":  "50",
						"internet_max_bandwidth_out": "0",
					}),
				),
			},
		},
	})
}

func TestAccApsaraStackInstancePrepaid(t *testing.T) {
	var v ecs.Instance

	resourceId := "apsarastack_instance.default"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAcc%sEcsInstanceConfigPrePaid%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceInstancePrePaidConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":        "${data.apsarastack_images.default.images.0.id}",
					"security_groups": []string{"${apsarastack_security_group.default.0.id}"},
					"instance_type":   "${data.apsarastack_instance_types.default.instance_types.0.id}",

					"availability_zone":    "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}",
					"system_disk_category": "cloud_efficiency",
					"instance_name":        "${var.name}",
					"key_name":             "${apsarastack_key_pair.default.key_name}",

					"security_enhancement_strategy": "Active",
					"user_data":                     "I_am_user_data",

					"vswitch_id": "${apsarastack_vswitch.default.id}",
					"role_name":  "${apsarastack_ram_role.default.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"key_name":      name,
						"role_name":     name,

						"force_delete": "false",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"security_enhancement_strategy", "data_disks", "dry_run", "force_delete",
					"include_data_disks"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id": "${data.apsarastack_images.default.images.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_groups": []string{"${apsarastack_security_group.default.0.id}", "${apsarastack_security_group.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_groups.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"force_delete": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"force_delete": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": name + "_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name + "_change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_description",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_max_bandwidth_out": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_out": "50",
						"private_ip":                 CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internet_max_bandwidth_in": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internet_max_bandwidth_in": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_name": "hostNameExample",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_name": "hostNameExample",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "Password123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "Password123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_disk_size": "50",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_disk_size": "50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"private_ip": "172.16.0.10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"private_ip": "172.16.0.10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"foo": "foo",
						"bar": "bar",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":   "2",
						"tags.foo": "foo",
						"tags.bar": "bar",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"security_groups":            []string{"${apsarastack_security_group.default.0.id}"},
					"instance_name":              name,
					"description":                name,
					"internet_max_bandwidth_out": REMOVEKEY,
					"internet_max_bandwidth_in":  REMOVEKEY,
					"host_name":                  REMOVEKEY,
					"password":                   REMOVEKEY,

					"system_disk_size": "70",
					"private_ip":       REMOVEKEY,
					"tags":             REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"tags.%":   "0",
						"tags.bar": REMOVEKEY,
						"tags.foo": REMOVEKEY,

						"instance_name": name,

						"image_id":          CHECKSET,
						"instance_type":     CHECKSET,
						"security_groups.#": "1",

						"availability_zone":             CHECKSET,
						"system_disk_category":          "cloud_efficiency",
						"security_enhancement_strategy": "Active",
						"vswitch_id":                    CHECKSET,

						"description":      name,
						"host_name":        CHECKSET,
						"password":         "",
						"system_disk_size": "70",

						"private_ip": CHECKSET,
						"status":     "Running",

						"internet_max_bandwidth_in":  "50",
						"internet_max_bandwidth_out": "0",
					}),
				),
			},
		},
	})
}

func TestAccApsaraStackInstanceDataDisks(t *testing.T) {
	var v ecs.Instance

	resourceId := "apsarastack_instance.default"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAcc%sEcsInstanceDataDisks%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceInstancePrePaidConfigDependence)

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
					"image_id":        "${data.apsarastack_images.default.images.0.id}",
					"security_groups": []string{"${apsarastack_security_group.default.0.id}"},
					"instance_type":   "${data.apsarastack_instance_types.default.instance_types.0.id}",

					"availability_zone":    "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}",
					"system_disk_category": "cloud_efficiency",
					"instance_name":        "${var.name}",
					"key_name":             "${apsarastack_key_pair.default.key_name}",
					//"spot_strategy":                 "NoSpot",
					//"spot_price_limit":              "0",
					"security_enhancement_strategy": "Active",
					"user_data":                     "I_am_user_data",

					"instance_charge_type": "PrePaid",
					"vswitch_id":           "${apsarastack_vswitch.default.id}",
					"role_name":            "${apsarastack_ram_role.default.name}",
					"data_disks": []map[string]string{
						{
							"name":        "disk1",
							"size":        "20",
							"category":    "cloud_efficiency",
							"description": "disk1",
						},
						{
							"name":        "disk2",
							"size":        "20",
							"category":    "cloud_efficiency",
							"description": "disk2",
						},
					},
					"force_delete": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"key_name":      name,
						"role_name":     name,
						"user_data":     "I_am_user_data",

						"data_disks.#":             "2",
						"data_disks.0.name":        "disk1",
						"data_disks.0.size":        "20",
						"data_disks.0.category":    "cloud_efficiency",
						"data_disks.0.description": "disk1",
						"data_disks.1.name":        "disk2",
						"data_disks.1.size":        "20",
						"data_disks.1.category":    "cloud_efficiency",
						"data_disks.1.description": "disk2",

						"force_delete": "true",
						//"instance_charge_type": "PrePaid",
						//"period":               "1",
						//"period_unit":          "Month",
						//"renewal_status":       "Normal",
						//"auto_renew_period":    "0",
						//"include_data_disks":   "true",
						//"dry_run":              "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_enhancement_strategy", "data_disks", "force_delete"},
			},
		},
	})
}

func TestAccApsaraStackInstanceTypeUpdate(t *testing.T) {
	var v ecs.Instance

	resourceId := "apsarastack_instance.default"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsInstanceConfigInstanceType%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceInstanceTypeConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":             "${data.apsarastack_images.default.images.0.id}",
					"system_disk_category": "cloud_efficiency",
					"system_disk_size":     "40",
					"instance_type":        "${data.apsarastack_instance_types.new1.instance_types.0.id}",
					"instance_name":        "${var.name}",
					"security_groups":      []string{"${apsarastack_security_group.default.id}"},
					"vswitch_id":           "${apsarastack_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":                 REGEXMATCH + "^ecs.t5-[a-z0-9]{1,}.nano",
						"user_data":                     REMOVEKEY,
						"security_enhancement_strategy": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.apsarastack_instance_types.new2.instance_types.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": REGEXMATCH + "^ecs.t5-[a-z0-9]{1,}.small",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type": "${data.apsarastack_instance_types.new4.instance_types.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type": REGEXMATCH + "^ecs.t5-[a-z0-9]{1,}.large",
					}),
				),
			},
		},
	})
}

func TestAccApsaraStackInstanceSpotInstanceLimit(t *testing.T) {
	var v ecs.Instance

	resourceId := "apsarastack_instance.default"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccEcsInstanceConfigSpot%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccCheckSpotInstanceDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.EcsSpotNoSupportedRegions)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vswitch_id":           "${apsarastack_vswitch.default.id}",
					"image_id":             "${data.apsarastack_images.default.images.0.id}",
					"availability_zone":    "${data.apsarastack_instance_types.special.instance_types.0.availability_zones.0}",
					"instance_type":        "${data.apsarastack_instance_types.special.instance_types.0.id}",
					"system_disk_category": "cloud_efficiency",
					//"internet_charge_type":       "PayByTraffic",
					"internet_max_bandwidth_out": "5",
					"security_groups":            []string{"${apsarastack_security_group.default.id}"},
					"instance_name":              "${var.name}",
					//"spot_strategy":              "SpotWithPriceLimit",
					//"spot_price_limit":           "1.002",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						//"spot_strategy":                 "SpotWithPriceLimit",
						//"spot_price_limit":              "1.002",
						"internet_max_bandwidth_out": "5",
						"private_ip":                 CHECKSET,
						//"public_ip":                     CHECKSET,
						"user_data":                     REMOVEKEY,
						"security_enhancement_strategy": REMOVEKEY,
					}),
				),
			},
		},
	})
}

func TestAccApsaraStackInstanceMulti(t *testing.T) {
	var v ecs.Instance

	resourceId := "apsarastack_instance.default.9"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAcc%sEcsInstanceConfigMulti%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceInstanceVpcConfigDependence)

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
					"count":           "10",
					"image_id":        "${data.apsarastack_images.default.images.0.id}",
					"security_groups": []string{"${apsarastack_security_group.default.0.id}"},
					"instance_type":   "${data.apsarastack_instance_types.default.instance_types.0.id}",

					"availability_zone":    "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}",
					"system_disk_category": "cloud_efficiency",
					"instance_name":        "${var.name}",
					"key_name":             "${apsarastack_key_pair.default.key_name}",
					//"spot_strategy":                 "NoSpot",
					//"spot_price_limit":              "0",
					"security_enhancement_strategy": "Active",
					"user_data":                     "I_am_user_data",

					"vswitch_id": "${apsarastack_vswitch.default.id}",
					"role_name":  "${apsarastack_ram_role.default.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"key_name":      name,
						"role_name":     name,
					}),
				),
			},
		},
	})
}

func resourceInstanceVpcConfigDependence(name string) string {
	return fmt.Sprintf(`
data "apsarastack_instance_types" "default" {
 	cpu_core_count    = 1
	memory_size       = 2
}

data "apsarastack_images" "default" {
  name_regex  = "^ubuntu*"
  owners      = "system"
}
resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  name              = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  count = "2"
  name   = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_security_group_rule" "default" {
	count = 2
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${element(apsarastack_security_group.default.*.id,count.index)}"
  	cidr_ip = "172.16.0.0/24"
}

variable "name" {
	default = "%s"
}

resource "apsarastack_ram_role" "default" {
		  name = "${var.name}"
		  services = ["ecs.aliyuncs.com"]
		  force = "true"
}

resource "apsarastack_key_pair" "default" {
	key_name = "${var.name}"
}

`, name)
}

func resourceInstancePrePaidConfigDependence(name string) string {
	return fmt.Sprintf(`
data "apsarastack_instance_types" "default" {
  cpu_core_count    = 2
  memory_size       = 4
  instance_charge_type = "PrePaid"
}
data "apsarastack_images" "default" {
  name_regex  = "^ubuntu*"
  owners      = "system"
}
resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}"
  name              = "${var.name}"
}
resource "apsarastack_security_group" "default" {
  count = "2"
  name   = "${var.name}"
  vpc_id = "${apsarastack_vpc.default.id}"
}
resource "apsarastack_security_group_rule" "default" {
	count = 2
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${element(apsarastack_security_group.default.*.id,count.index)}"
  	cidr_ip = "172.16.0.0/24"
}

variable "name" {
	default = "%s"
}

resource "apsarastack_ram_role" "default" {
		  name = "${var.name}"
		  services = ["ecs.aliyuncs.com"]
		  force = "true"
}

resource "apsarastack_key_pair" "default" {
	key_name = "${var.name}"
}

`, name)
}

func resourceInstanceBasicConfigDependence(name string) string {
	return fmt.Sprintf(`

data "apsarastack_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
}

//variable "resource_group_id" {
//		default = "%s"
//	}

data "apsarastack_images" "default" {
  name_regex  = "^ubuntu*"
  owners      = "system"
}

resource "apsarastack_security_group" "default" {
  count = "2"
  name   = "${var.name}"
}
resource "apsarastack_security_group_rule" "default" {
	count = 2
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${element(apsarastack_security_group.default.*.id,count.index)}"
  	cidr_ip = "172.16.0.0/24"
}

variable "name" {
	default = "%s"
}

resource "apsarastack_key_pair" "default" {
	key_name = "${var.name}"
}

`, os.Getenv("APSARASTaCK_RESOURCE_GROUP_ID"), name)
}

func resourceInstanceTypeConfigDependence(name string) string {
	return fmt.Sprintf(`
    data "apsarastack_zones" "default" {
	  available_disk_category     = "cloud_efficiency"
	  available_resource_creation = "VSwitch"
	}
	data "apsarastack_images" "default" {
	  name_regex  = "^ubuntu_18.*64"
	  most_recent = true
	  owners      = "system"
	}
	resource "apsarastack_vpc" "default" {
	  name       = "${var.name}"
	  cidr_block = "172.16.0.0/16"
	}
	resource "apsarastack_vswitch" "default" {
	  vpc_id            = "${apsarastack_vpc.default.id}"
	  cidr_block        = "172.16.0.0/24"
	  availability_zone = "${reverse(data.apsarastack_zones.default.zones).1.id}"
	  name              = "${var.name}"
	}
	resource "apsarastack_security_group" "default" {
	  name   = "${var.name}"
	  vpc_id = "${apsarastack_vpc.default.id}"
	}
	resource "apsarastack_security_group_rule" "default" {
	  	type = "ingress"
	  	ip_protocol = "tcp"
	  	nic_type = "intranet"
	  	policy = "accept"
	  	port_range = "22/22"
	  	priority = 1
	  	security_group_id = "${apsarastack_security_group.default.id}"
	  	cidr_ip = "172.16.0.0/24"
	}

	variable "name" {
		default = "%s"
	}

	data "apsarastack_instance_types" "new1" {
		availability_zone = "${apsarastack_vswitch.default.availability_zone}"
		cpu_core_count = 1
		memory_size = 0.5
		instance_type_family = "ecs.t5"
	}

	data "apsarastack_instance_types" "new2" {
		availability_zone = "${apsarastack_vswitch.default.availability_zone}"
		cpu_core_count = 1
		memory_size = 1
		instance_type_family = "ecs.t5"
	}

	data "apsarastack_instance_types" "new3" {
		availability_zone = "${apsarastack_vswitch.default.availability_zone}"
		cpu_core_count = 1
		memory_size = 2
		instance_type_family = "ecs.t5"
	}

	data "apsarastack_instance_types" "new4" {
		availability_zone = "${apsarastack_vswitch.default.availability_zone}"
		cpu_core_count = 2
		memory_size = 4
		instance_type_family = "ecs.t5"
	}



`, name)
}

func testAccCheckSpotInstanceDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "%s"
	}
	data "apsarastack_instance_types" "special" {
	  	cpu_core_count    = 2
	  	memory_size       = 4
	  	spot_strategy = "SpotWithPriceLimit"
	}
	
	`, EcsInstanceCommonNoZonesTestCase, name)
}

var testAccInstanceCheckMap = map[string]string{
	"image_id":          CHECKSET,
	"instance_type":     CHECKSET,
	"security_groups.#": "1",

	"availability_zone":    CHECKSET,
	"system_disk_category": "cloud_efficiency",
	//"credit_specification":          "",
	//"spot_strategy":                 "NoSpot",
	//"spot_price_limit":              "0",
	"security_enhancement_strategy": "Active",
	"vswitch_id":                    CHECKSET,
	"user_data":                     "I_am_user_data",

	"description": "",
	"host_name":   CHECKSET,
	"password":    "",
	//"is_outdated":      NOSET,
	"system_disk_size": "40",

	"data_disks.#": NOSET,
	//"volume_tags.%": "0",
	"tags.%": NOSET,

	"private_ip": CHECKSET,
	//"public_ip":  "",
	"status": "Running",

	//"internet_charge_type":       "PayByTraffic",
	"internet_max_bandwidth_in":  "-1",
	"internet_max_bandwidth_out": "0",

	//"instance_charge_type": "PostPaid",
	// the attributes of below are suppressed  when the value of instance_charge_type is `PostPaid`
	//"period":             NOSET,
	//"period_unit":        NOSET,
	//"renewal_status":     NOSET,
	//"auto_renew_period":  NOSET,
	"force_delete": NOSET,
	//"include_data_disks": NOSET,
	//"dry_run":            "false",
}
