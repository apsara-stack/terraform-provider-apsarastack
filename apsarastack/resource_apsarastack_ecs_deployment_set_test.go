package apsarastack

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"apsarastack_ecs_deployment_set",
		&resource.Sweeper{
			Name: "apsarastack_ecs_deployment_set",
			F:    testSweepEcsDeploymentSet,
		})
}

func testSweepEcsDeploymentSet(region string) error {

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting ApsaraStack client: %s", err)
	}
	client := rawClient.(*connectivity.ApsaraStackClient)
	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}
	action := "DescribeDeploymentSets"
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}

	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}
		resp, err := jsonpath.Get("$.DeploymentSets.DeploymentSet", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.DeploymentSets.DeploymentSet", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			if _, ok := item["DeploymentSetName"]; !ok {
				continue
			}
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["DeploymentSetName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Ecs DeploymentSet: %s", item["DeploymentSetName"].(string))
				continue
			}
			action := "DeleteDeploymentSet"
			request := map[string]interface{}{
				"DeploymentSetId": item["DeploymentSetId"],
				"RegionId":        client.RegionId,
			}
			request["ClientToken"] = buildClientToken("DeleteDeploymentSet")
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ecs DeploymentSet (%s): %s", item["DeploymentSetId"].(string), err)
			}
			log.Printf("[INFO] Delete Ecs DeploymentSet success: %s ", item["DeploymentSetId"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	return nil
}

//type EcsDescribeDeploymentSetsResult struct {
//	DeploymentSets struct {
//		DeploymentSet []struct {
//			Granularity              string `json:"Granularity"`
//			DeploymentStrategy       string `json:"DeploymentStrategy"`
//			DeploymentSetDescription string `json:"DeploymentSetDescription"`
//			DeploymentSetName        string `json:"DeploymentSetName"`
//			Domain                   string `json:"Domain"`
//		} `json:"DeploymentSet"`
//	} `json:"DeploymentSets"`
//}
func TestAccApsaraStackECSDeploymentSet_basic0(t *testing.T) {
	var v *EcsDescribeDeploymentSetsResult
	resourceId := "apsarastack_ecs_deployment_set.default"
	ra := resourceAttrInit(resourceId, ApsaraStackECSDeploymentSetMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeEcsDeploymentSet")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdeploymentset%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackECSDeploymentSetBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"strategy":            "Availability",
					"domain":              "default",
					"granularity":         "host",
					"deployment_set_name": name,
					"description":         name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"strategy":    "Availability",
						"domain":      "default",
						"granularity": "host",
						//"deployment_set_name": name,
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deployment_set_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						//"deployment_set_name": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":         name,
					"deployment_set_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
						//"deployment_set_name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"on_unable_to_redeploy_failed_instance"},
			},
		},
	})
}

var ApsaraStackECSDeploymentSetMap0 = map[string]string{}

func ApsaraStackECSDeploymentSetBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
