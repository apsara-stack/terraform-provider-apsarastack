package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("apsarastack_network_acl", &resource.Sweeper{
		Name: "apsarastack_network_acl",
		F:    testSweepNetworkAcl,
	})
}

func testSweepNetworkAcl(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting ApsaraStack client: %s", err)
	}
	client := rawClient.(*connectivity.ApsaraStackClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "DescribeNetworkAcls"
	request1 := vpc.CreateDescribeNetworkAclsRequest()
	response := vpc.CreateDescribeNetworkAclsResponse()
	params := make(map[string]string)
	request1.QueryParams = params
	params["RegionId"] = client.RegionId
	params["PageSize"] = strconv.Itoa(PageSizeLarge)
	params["PageNumber"] = "1"
	params["Action"] = action
	networkAclIds := make([]string, 0)
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNetworkAcls(request1)
		})
		response = raw.(*vpc.DescribeNetworkAclsResponse)
		if err != nil {
			log.Printf("Error retrieving network acl: %s", err)
			return nil
		}
		var networkAcl = response.NetworkAcls.NetworkAcl
		for _, v := range networkAcl {
			name := v.NetworkAclName
			id := v.NetworkAclId
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Network Acl: %s (%s)", name, id)
				continue
			}
			networkAclIds = append(networkAclIds, id)
		}
		if len(networkAcl) < PageSizeLarge {
			break
		}
		num, err := strconv.Atoi(params["PageNumber"])
		if err != nil {
			return WrapError(err)
		}
		params["PageNumber"] = strconv.Itoa(num + 1)
	}

	vpcService := VpcService{client}
	for _, id := range networkAclIds {
		//	Delete attach resources
		object, err := vpcService.DescribeNetworkAcl(id)
		if err != nil {
			log.Println("DescribeNetworkAcl failed", err)
		}
		deleteResources, _ := object["Resources"].(map[string]interface{})["Resource"].([]interface{})
		if len(deleteResources) > 0 {
			params = make(map[string]string)
			params["NetworkAclId"] = id
			resourcesMaps := make([]map[string]interface{}, 0)
			for _, resources := range deleteResources {
				resourcesArg := resources.(map[string]interface{})
				resourcesMap := map[string]interface{}{
					"ResourceId":   resourcesArg["ResourceId"],
					"ResourceType": resourcesArg["ResourceType"],
				}
				resourcesMaps = append(resourcesMaps, resourcesMap)
			}
			str, err := mapToStr(resourcesMaps)
			if err != nil {
				return WrapErrorf(err, "map转换json异常", resourcesMaps)
			}
			params["Resource"] = str
			params["RegionId"] = client.RegionId
			action := "UnassociateNetworkAcl"
			params["Action"] = action
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			request1 := vpc.CreateUnassociateNetworkAclRequest()
			request1.QueryParams = params
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
					return vpcClient.UnassociateNetworkAcl(request1)
				})
				log.Println("UnassociateNetworkAcl response:", raw)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				log.Println("UnassociateNetworkAcl failed", err)
			}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, 5*time.Minute, 5*time.Second, vpcService.NetworkAclStateRefreshFunc(id, []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				log.Println("UnassociateNetworkAcl failed", err)
			}
		}

		log.Printf("[INFO] Deleting Network Acl: (%s)", id)
		request1 := vpc.CreateDeleteNetworkAclRequest()
		request1.QueryParams = params
		params = make(map[string]string)
		params["NetworkAclId"] = id
		action := "DeleteNetworkAcl"
		params["Action"] = action
		params["RegionId"] = client.RegionId
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DeleteNetworkAcl(request1)
			})
			log.Println("DeleteNetworkAcl responce", raw)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Network Acl (%s): %s", id, err)
		}
	}
	return nil
}

func TestAccApsaraStackVpcNetworkAcl_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "apsarastack_network_acl.default"
	ra := resourceAttrInit(resourceId, ApsaraStackNetworkAclMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeNetworkAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snetworkacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackNetworkAclBasicDependence0)
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
					"vpc_id":           "${apsarastack_vpc.default.id}",
					"network_acl_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":           CHECKSET,
						"network_acl_name": name,
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
					"description": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"egress_acl_entries": []map[string]interface{}{
						{
							"description":            "engress test",
							"destination_cidr_ip":    "10.0.0.0/24",
							"network_acl_entry_name": "tf-testacc78924",
							"policy":                 "accept",
							"port":                   "20/80",
							"protocol":               "tcp",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"egress_acl_entries.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ingress_acl_entries": []map[string]interface{}{
						{
							"description":            "ingress test",
							"network_acl_entry_name": "tf-testacc78999",
							"policy":                 "accept",
							"port":                   "20/80",
							"protocol":               "tcp",
							"source_cidr_ip":         "10.0.0.0/24",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ingress_acl_entries.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_acl_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"network_acl_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resources": []map[string]interface{}{
						{
							"resource_id":   "${apsarastack_vswitch.default0.id}",
							"resource_type": "VSwitch",
						},
						{
							"resource_id":   "${apsarastack_vswitch.default1.id}",
							"resource_type": "VSwitch",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resources.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resources": []map[string]interface{}{
						{
							"resource_id":   "${apsarastack_vswitch.default0.id}",
							"resource_type": "VSwitch",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resources.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      name,
					"network_acl_name": name,
					"ingress_acl_entries": []map[string]interface{}{
						{
							"description":            "ingress test change",
							"network_acl_entry_name": "tf-testacc78999",
							"policy":                 "accept",
							"port":                   "20/80",
							"protocol":               "tcp",
							"source_cidr_ip":         "10.0.0.0/24",
						},
					},
					"egress_acl_entries": []map[string]interface{}{
						{
							"description":            "engress test change",
							"destination_cidr_ip":    "10.0.0.0/24",
							"network_acl_entry_name": "tf-testacc78924",
							"policy":                 "accept",
							"port":                   "20/80",
							"protocol":               "tcp",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":           name,
						"network_acl_name":      name,
						"ingress_acl_entries.#": "1",
						"egress_acl_entries.#":  "1",
					}),
				),
			},
		},
	})
}

var ApsaraStackNetworkAclMap0 = map[string]string{}

func ApsaraStackNetworkAclBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%[1]s"
		}
variable "name_change" {
			default = "%[1]s_change"
		}
data "apsarastack_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "apsarastack_vpc" "default" {
  cidr_block = "192.168.0.0/16"
  name = var.name
}
resource "apsarastack_vswitch" "default0" {
  vpc_id            = apsarastack_vpc.default.id
  name      = var.name
  cidr_block        = cidrsubnets(apsarastack_vpc.default.cidr_block, 4, 4)[0]
  availability_zone           = data.apsarastack_zones.default.ids.0
}
resource "apsarastack_vswitch" "default1" {
  vpc_id            = apsarastack_vpc.default.id
  name      = var.name_change
  cidr_block        = cidrsubnets(apsarastack_vpc.default.cidr_block, 4, 4)[1]
  availability_zone           = data.apsarastack_zones.default.ids.0
}

`, name)
}

func Test_resourceApsaraStackNetworkAclUpdate(t *testing.T) {
	type args struct {
		d    *schema.ResourceData
		meta interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := resourceApsaraStackNetworkAclUpdate(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceApsaraStackNetworkAclUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
