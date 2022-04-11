package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"regexp"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceApsaraStackNetworkAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackNetworkAclsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"network_acl_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Available", "Modifying"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"acls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"egress_acl_entries": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination_cidr_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_acl_entry_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ingress_acl_entries": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_acl_entry_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_cidr_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_acl_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_acl_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
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
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackNetworkAclsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	action := "DescribeNetworkAcls"
	request1 := vpc.CreateDescribeNetworkAclsRequest()
	var response = vpc.CreateDescribeNetworkAclsResponse()
	params := make(map[string]string)
	request1.QueryParams = params
	params["Product"] = "Vpc"
	params["OrganizationId"] = client.Department
	if v, ok := d.GetOk("network_acl_name"); ok {
		params["NetworkAclName"] = v.(string)
	}
	params["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_id"); ok {
		params["ResourceId"] = v.(string)
	}
	if v, ok := d.GetOk("resource_type"); ok {
		params["ResourceType"] = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		params["VpcId"] = v.(string)
	}
	params["PageSize"] = strconv.Itoa(PageSizeLarge)
	params["PageNumber"] = "1"
	var objects []map[string]interface{}
	var networkAclNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		networkAclNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeNetworkAcls(request1)
			})
			response = raw.(*vpc.DescribeNetworkAclsResponse)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request1)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_network_acls", action, ApsaraStackSdkGoERROR)
		}
		var networkAcl = response.NetworkAcls.NetworkAcl
		var result = make(map[string]interface{})
		for _, v := range networkAcl {
			if networkAclNameRegex != nil {
				if !networkAclNameRegex.MatchString(v.NetworkAclName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[v.NetworkAclId]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != v.Status {
				continue
			}
			//把对象转换为map
			b, err := json.Marshal(v)
			if err != nil {
				return WrapError(err)
			}
			err = json.Unmarshal(b, &result)
			if err != nil {
				return WrapError(err)
			}
			objects = append(objects, result)
		}
		if len(networkAcl) < PageSizeLarge {
			break
		}
		var pageNumber, _ = strconv.Atoi(params["PageNumber"])
		params["PageNumber"] = strconv.Itoa(pageNumber + 1)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"description":      object["Description"],
			"id":               fmt.Sprint(object["NetworkAclId"]),
			"network_acl_id":   fmt.Sprint(object["NetworkAclId"]),
			"network_acl_name": object["NetworkAclName"],
			"status":           object["Status"],
			"vpc_id":           object["VpcId"],
		}

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
		mapping["egress_acl_entries"] = egressAclEntry

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
		mapping["ingress_acl_entries"] = ingressAclEntry

		resourceMap := make([]map[string]interface{}, 0)
		if resourceMapList, ok := object["Resources"].(map[string]interface{})["Resource"].([]interface{}); ok {
			for _, v := range resourceMapList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"resource_id":   m1["ResourceId"],
						"resource_type": m1["ResourceType"],
						"status":        m1["Status"],
					}
					resourceMap = append(resourceMap, temp1)
				}
			}
		}
		mapping["resources"] = resourceMap
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["NetworkAclName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("acls", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
