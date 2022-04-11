package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceApsaraStackNetworkAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackNetworkAclCreate,
		Read:   resourceApsaraStackNetworkAclRead,
		Update: resourceApsaraStackNetworkAclUpdate,
		Delete: resourceApsaraStackNetworkAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"egress_acl_entries": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"destination_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network_acl_entry_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"accept", "drop"}, false),
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"icmp", "gre", "tcp", "udp", "all"}, false),
						},
					},
				},
			},
			"ingress_acl_entries": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network_acl_entry_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"policy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"accept", "drop"}, false),
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"icmp", "gre", "tcp", "udp", "all"}, false),
						},
						"source_cidr_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"network_acl_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validation.StringLenBetween(2, 128),
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.122.0. New field 'network_acl_name' instead",
				ConflictsWith: []string{"network_acl_name"},
				ValidateFunc:  validation.StringLenBetween(2, 128),
			},
			"resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceApsaraStackNetworkAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := vpc.CreateCreateNetworkAclRequest()
	action := "CreateNetworkAcl"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}

	//非必须参数
	if val, ok := d.GetOk("description"); ok {
		//断言
		request.Description = val.(string)
	}
	request.NetworkAclName = d.Get("network_acl_name").(string)
	request.VpcId = d.Get("vpc_id").(string)
	request.ClientToken = buildClientToken("CreateNetworkAcl")
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.Domain = client.Domain
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateNetworkAcl(request)
	})
	if err != nil {
		//打印错误
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_network_acl", action, ApsaraStackSdkGoERROR)
	}
	//类型转换（断言）
	response := raw.(*vpc.CreateNetworkAclResponse)
	addDebug(action, raw, request.RpcRequest, request)
	err = json.Unmarshal(response.GetHttpContentBytes(), &response)
	if err != nil {
		return WrapError(err)
	}
	// 设置newState.ID
	d.SetId(fmt.Sprint(response.NetworkAclId))
	//更新函数
	return resourceApsaraStackNetworkAclUpdate(d, meta)
}
func resourceApsaraStackNetworkAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeNetworkAcl(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource apsarastack_network_acl vpcService.DescribeNetworkAcl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])

	egressAclEntry := make([]map[string]interface{}, 0)
	if egressAclEntryList, ok := object["EgressAclEntries"].(map[string]interface{})["EgressAclEntry"].([]interface{}); ok {
		for _, v := range egressAclEntryList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"description":            m1["Description"],
					"destination_cidr_ip":    m1["DestinationCidrIp"],
					"network_acl_entry_name": m1["NetworkAclEntryName"],
					"policy":                 m1["Policy"],
					"port":                   m1["Port"],
					"protocol":               m1["Protocol"],
				}
				egressAclEntry = append(egressAclEntry, temp1)

			}
		}
	}
	if err := d.Set("egress_acl_entries", egressAclEntry); err != nil {
		return WrapError(err)
	}

	ingressAclEntry := make([]map[string]interface{}, 0)
	if ingressAclEntryList, ok := object["IngressAclEntries"].(map[string]interface{})["IngressAclEntry"].([]interface{}); ok {
		for _, v := range ingressAclEntryList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"description":            m1["Description"],
					"network_acl_entry_name": m1["NetworkAclEntryName"],
					"policy":                 m1["Policy"],
					"port":                   m1["Port"],
					"protocol":               m1["Protocol"],
					"source_cidr_ip":         m1["SourceCidrIp"],
				}
				ingressAclEntry = append(ingressAclEntry, temp1)

			}
		}
	}
	if err := d.Set("ingress_acl_entries", ingressAclEntry); err != nil {
		return WrapError(err)
	}
	d.Set("network_acl_name", object["NetworkAclName"])
	d.Set("name", object["NetworkAclName"])

	resourceMap := make([]map[string]interface{}, 0)
	if resourceMapList, ok := object["Resources"].(map[string]interface{})["Resource"].([]interface{}); ok {
		for _, v := range resourceMapList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"resource_id":   m1["ResourceId"],
					"resource_type": m1["ResourceType"],
				}
				resourceMap = append(resourceMap, temp1)

			}
		}
	}
	if err := d.Set("resources", resourceMap); err != nil {
		return WrapError(err)
	}
	d.Set("status", object["Status"])
	d.Set("vpc_id", object["VpcId"])
	return nil
}
func resourceApsaraStackNetworkAclUpdate(d *schema.ResourceData, meta interface{}) error {
	//获取客户端
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}
	//获取请求对象
	request := vpc.CreateModifyNetworkAclAttributesRequest()
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	d.Partial(true)

	update := false
	request.NetworkAclId = d.Id()
	request.RegionId = client.RegionId
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		//断言，没有值会报错
		request.Description = d.Get("description").(string)
	}
	if !d.IsNewResource() && d.HasChange("network_acl_name") {
		update = true
		request.NetworkAclName = d.Get("network_acl_name").(string)
	}
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request.NetworkAclName = d.Get("name").(string)
	}

	if update {
		action := "ModifyNetworkAclAttributes"
		request.ClientToken = buildClientToken("ModifyNetworkAclAttributes")
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyNetworkAclAttributes(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
		}
		//断言
		response := raw.(*vpc.ModifyNetworkAclAttributesResponse)
		addDebug(action, response, request.RpcRequest, request)
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.NetworkAclStateRefreshFunc(d.Id(), []string{"Modifying"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		//d.SetPartial("description")
		//d.SetPartial("name")
		//d.SetPartial("network_acl_name")
	}

	update = false
	/*设置请求参数*/
	updateNetworkAclEntriesRequest := vpc.CreateUpdateNetworkAclEntriesRequest()
	updateNetworkAclEntriesRequest.NetworkAclId = d.Id()
	updateNetworkAclEntriesRequest.RegionId = client.RegionId
	updateNetworkAclEntriesRequest.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	updateNetworkAclEntriesRequest.ClientToken = buildClientToken("UpdateNetworkAclEntries")
	if d.HasChange("egress_acl_entries") {
		updateNetworkAclEntriesRequest.UpdateEgressAclEntries = requests.NewBoolean(true)
		update = true
		var EgressAclEntries []vpc.UpdateNetworkAclEntriesEgressAclEntries
		for _, EgressAclEntriesValue := range d.Get("egress_acl_entries").([]interface{}) {
			var egress vpc.UpdateNetworkAclEntriesEgressAclEntries
			EgressAclEntriesMap := EgressAclEntriesValue.(map[string]interface{})
			egress.Description = EgressAclEntriesMap["description"].(string)
			egress.DestinationCidrIp = EgressAclEntriesMap["destination_cidr_ip"].(string)
			egress.NetworkAclEntryName = EgressAclEntriesMap["network_acl_entry_name"].(string)
			egress.Policy = EgressAclEntriesMap["policy"].(string)
			egress.Port = EgressAclEntriesMap["port"].(string)
			egress.Protocol = EgressAclEntriesMap["protocol"].(string)
			EgressAclEntries = append(EgressAclEntries, egress)
		}
		updateNetworkAclEntriesRequest.EgressAclEntries = &EgressAclEntries
	}
	if d.HasChange("ingress_acl_entries") {
		updateNetworkAclEntriesRequest.UpdateIngressAclEntries = requests.NewBoolean(true)
		update = true
		var IngressAclEntries []vpc.UpdateNetworkAclEntriesIngressAclEntries
		for _, IngressAclEntriesValue := range d.Get("ingress_acl_entries").([]interface{}) {
			IngressAclEntriesMap := IngressAclEntriesValue.(map[string]interface{})
			var ingress vpc.UpdateNetworkAclEntriesIngressAclEntries
			ingress.Description = IngressAclEntriesMap["description"].(string)
			ingress.NetworkAclEntryName = IngressAclEntriesMap["network_acl_entry_name"].(string)
			ingress.Policy = IngressAclEntriesMap["policy"].(string)
			ingress.Port = IngressAclEntriesMap["port"].(string)
			ingress.Protocol = IngressAclEntriesMap["protocol"].(string)
			ingress.SourceCidrIp = IngressAclEntriesMap["source_cidr_ip"].(string)
			IngressAclEntries = append(IngressAclEntries, ingress)
		}
		updateNetworkAclEntriesRequest.IngressAclEntries = &IngressAclEntries
	}
	if update {
		action := "UpdateNetworkAclEntries"
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UpdateNetworkAclEntries(updateNetworkAclEntriesRequest)
		})
		//断言
		response := raw.(*vpc.UpdateNetworkAclEntriesResponse)
		addDebug(action, response, updateNetworkAclEntriesRequest.RpcRequest, updateNetworkAclEntriesRequest)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.NetworkAclStateRefreshFunc(d.Id(), []string{"Modifying"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		//d.SetPartial("egress_acl_entries")
		//d.SetPartial("ingress_acl_entries")
	}
	d.Partial(false)
	if d.HasChange("resources") {
		oldResources, newResources := d.GetChange("resources")
		oldResourcesSet := oldResources.(*schema.Set)
		newResourcesSet := newResources.(*schema.Set)

		removed := oldResourcesSet.Difference(newResourcesSet)
		added := newResourcesSet.Difference(oldResourcesSet)
		if added.Len() > 0 {
			var aclResources []vpc.AssociateNetworkAclResource
			for _, resources := range added.List() {
				resourcesArg := resources.(map[string]interface{})
				resourcesMap := vpc.AssociateNetworkAclResource{
					ResourceId:   resourcesArg["resource_id"].(string),
					ResourceType: resourcesArg["resource_type"].(string),
				}
				aclResources = append(aclResources, resourcesMap)
			}
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			request := vpc.CreateAssociateNetworkAclRequest()
			request.Resource = &aclResources
			request.NetworkAclId = d.Id()
			request.RegionId = client.RegionId
			request.ClientToken = buildClientToken("AssociateNetworkAcl")
			request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
			action := "AssociateNetworkAcl"
			response := vpc.CreateAssociateNetworkAclResponse()
			err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
					return vpcClient.AssociateNetworkAcl(request)
				})
				//断言
				response = raw.(*vpc.AssociateNetworkAclResponse)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request.RpcRequest, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.NetworkAclStateRefreshFunc(d.Id(), []string{"Modifying"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			//d.SetPartial("resources")
		}
		if removed.Len() > 0 {
			var resourcesMaps []vpc.UnassociateNetworkAclResource
			for _, resources := range removed.List() {
				resourcesArg := resources.(map[string]interface{})
				resourcesMap := vpc.UnassociateNetworkAclResource{
					ResourceId:   resourcesArg["resource_id"].(string),
					ResourceType: resourcesArg["resource_type"].(string),
				}
				resourcesMaps = append(resourcesMaps, resourcesMap)
			}

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			request := vpc.CreateUnassociateNetworkAclRequest()
			request.Resource = &resourcesMaps
			request.RegionId = client.RegionId
			request.ClientToken = buildClientToken("UnassociateNetworkAcl")
			request.NetworkAclId = d.Id()
			action := "UnassociateNetworkAcl"
			request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
			response := vpc.CreateUnassociateNetworkAclResponse()
			err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
					return vpcClient.UnassociateNetworkAcl(request)
				})
				//断言
				response = raw.(*vpc.UnassociateNetworkAclResponse)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request.RpcRequest, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.NetworkAclStateRefreshFunc(d.Id(), []string{"Modifying"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
			//d.SetPartial("resources")
		}
	}
	return resourceApsaraStackNetworkAclRead(d, meta)
}
func resourceApsaraStackNetworkAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	vpcService := VpcService{client}
	request := vpc.CreateDeleteNetworkAclRequest()
	var response = vpc.CreateDeleteNetworkAclResponse()
	// Delete binging resources before delete the ACL
	_, err := vpcService.DeleteAclResources(d.Id())
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteNetworkAcl"
	request.NetworkAclId = d.Id()
	request.RegionId = client.RegionId
	request.ClientToken = buildClientToken("DeleteNetworkAcl")
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteNetworkAcl(request)
		})
		response = raw.(*vpc.DeleteNetworkAclResponse)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request.RpcRequest, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
	}
	/*stateConf := BuildStateConf([]string{"Modifying"}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.NetworkAclStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}*/
	return nil
}

func mapToStr(object interface{}) (string, error) {
	//把map 转换为json
	marshal, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	//强制类型转换
	return string(marshal), nil
}
