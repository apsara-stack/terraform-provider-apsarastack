package apsarastack

import (
	"fmt"
	"testing"
	"time"

	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("apsarastack_edas_application", &resource.Sweeper{
		Name: "apsarastack_edas_application",
		F:    testSweepEdasApplication,
	})
}

func testSweepEdasApplication(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting ApsaraStack client.")
	}
	client := rawClient.(*connectivity.ApsaraStackClient)
	edasService := EdasService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}

	applicationListRq := edas.CreateListApplicationRequest()
	applicationListRq.RegionId = region

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.ListApplication(applicationListRq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve edas application in service list: %s", err)
	}

	listApplicationResponse, _ := raw.(*edas.ListApplicationResponse)
	if listApplicationResponse.Code != 200 {
		log.Printf("[ERROR] Failed to retrieve edas application in service list: %v", listApplicationResponse)
		return nil
	}

	for _, v := range listApplicationResponse.ApplicationList.Application {
		name := v.Name
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}

		if skip {
			log.Printf("[INFO] Skipping edas application: %s", name)
			continue
		}
		log.Printf("[INFO] delete edas application: %s", name)
		// stop it before delete
		stopAppRequest := edas.CreateStopApplicationRequest()
		stopAppRequest.RegionId = region
		stopAppRequest.AppId = v.AppId

		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.StopApplication(stopAppRequest)
		})
		if err != nil {
			return err
		}
		addDebug(stopAppRequest.GetActionName(), raw, stopAppRequest.RoaRequest, stopAppRequest)
		stopAppResponse, _ := raw.(*edas.StopApplicationResponse)
		changeOrderId := stopAppResponse.ChangeOrderId

		if len(changeOrderId) > 0 {
			stateConf := BuildStateConf([]string{"0", "1"}, []string{"2"}, 5*time.Minute, 5*time.Second, edasService.EdasChangeOrderStatusRefreshFunc(changeOrderId, []string{"3", "6", "10"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return err
			}
		}

		deleteApplicationRequest := edas.CreateDeleteApplicationRequest()
		deleteApplicationRequest.RegionId = region
		deleteApplicationRequest.AppId = v.AppId

		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
				return edasClient.DeleteApplication(deleteApplicationRequest)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(deleteApplicationRequest.GetActionName(), raw, deleteApplicationRequest.RoaRequest, deleteApplicationRequest)
			rsp := raw.(*edas.DeleteApplicationResponse)
			if rsp.Code == 601 && strings.Contains(rsp.Message, "Operation cannot be processed because there are running instances.") {
				err = Error("Operation cannot be processed because there are running instances.")
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] DeleteApplication failed: %v", err)
		}
	}

	return nil
}

func TestAccApsaraStackEdasApplication_basic(t *testing.T) {
	var v *edas.Applcation
	resourceId := "apsarastack_edas_application.default"
	ra := resourceAttrInit(resourceId, edasApplicationBasicMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasapplicationbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasApplicationConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEdasApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"application_name": "${var.name}",
					"package_type":     "WAR",
					"cluster_id":       "${apsarastack_edas_cluster.default.id}",
					//"build_pack_id":    "-1",
					//"region_id":        "cn-neimeng-env30-d01",
					"component_id": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name": name,
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"application_name": fmt.Sprintf("tf-testacc-edasapplicationchange%v", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"application_name": fmt.Sprintf("tf-testacc-edasapplicationchange%v", rand)}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"package_type": "WAR",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"package_type": "WAR"}),
				),
			},
		},
	})
}

func TestAccApsaraStackEdasApplication_multi(t *testing.T) {
	var v *edas.Applcation
	resourceId := "apsarastack_edas_application.default.1"
	ra := resourceAttrInit(resourceId, edasApplicationBasicMap)
	serviceFunc := func() interface{} {
		return &EdasService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(100, 999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-edasapplicationbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceEdasApplicationConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EdasSupportedRegions)
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEdasApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":            "2",
					"application_name": "${var.name}-${count.index}",
					"package_type":     "JAR",
					"cluster_id":       "${apsarastack_edas_cluster.default.id}",
					"build_pack_id":    "15",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var edasApplicationBasicMap = map[string]string{
	"application_name": CHECKSET,
	"package_type":     CHECKSET,
	"cluster_id":       CHECKSET,
}

func testAccCheckEdasApplicationDestroy(s *terraform.State) error {
	return nil
}

func resourceEdasApplicationConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		  default = "%v"
		}

		resource "apsarastack_vpc" "default" {
		  cidr_block = "172.16.0.0/12"
		  name       = "${var.name}"
		}

		resource "apsarastack_edas_cluster" "default" {
		  cluster_name = "${var.name}"
		  cluster_type = 2
		  network_mode = 2
		  vpc_id       = "${apsarastack_vpc.default.id}"
          region_id    = "cn-neimeng-env30-d01"
		}
		`, name)
}
