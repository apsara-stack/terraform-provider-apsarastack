package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strings"
)

func dataSourceApsaraStackSlbs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackSlbsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"master_availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"slave_availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values
			"slbs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slave_availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackSlbsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	slbService := &SlbService{client}

	request := slb.CreateDescribeLoadBalancersRequest()
	request.RegionId = client.RegionId
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	//request.ResourceGroupId = d.Get("resource_group_id").(string)
	if v, ok := d.GetOk("master_availability_zone"); ok && v.(string) != "" {
		request.MasterZoneId = v.(string)
	}
	if v, ok := d.GetOk("slave_availability_zone"); ok && v.(string) != "" {
		request.SlaveZoneId = v.(string)
	}
	if v, ok := d.GetOk("network_type"); ok && v.(string) != "" {
		request.NetworkType = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		request.VpcId = v.(string)
	}
	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		request.VSwitchId = v.(string)
	}
	if v, ok := d.GetOk("address"); ok && v.(string) != "" {
		request.Address = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		var tags []Tag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, Tag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tags = toSlbTagsString(tags)
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	//var allLoadBalancers []LoadBalancer
	var filteredLoadBalancersTemp []LoadBalancer
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeLoadBalancers(request)
		})
		//if err != nil {
		//	return WrapError(err)
		//}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*slb.DescribeLoadBalancersResponse)
		if len(response.LoadBalancers.LoadBalancer) < 1 {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page

		for _, item := range response.LoadBalancers.LoadBalancer {
			filteredLoadBalancersTemp = append(filteredLoadBalancersTemp,
				LoadBalancer{
					LoadBalancerId:   item.LoadBalancerId,
					RegionId:         item.RegionId,
					MasterZoneId:     item.MasterZoneId,
					SlaveZoneId:      item.SlaveZoneId,
					LoadBalancerName: item.LoadBalancerName,
					NetworkType:      item.NetworkType,
					VpcId:            item.VpcId,
					VSwitchId:        item.VSwitchId,
					Address:          item.Address,
					CreateTime:       item.CreateTime,
				},
			)
		}
	}
	return slbsDescriptionAttributes(d, filteredLoadBalancersTemp, slbService)
}

func slbsDescriptionAttributes(d *schema.ResourceData, loadBalancers []LoadBalancer, slbService *SlbService) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, loadBalancer := range loadBalancers {
		tags, _ := slbService.DescribeTags(loadBalancer.LoadBalancerId, nil, TagResourceInstance)
		mapping := map[string]interface{}{
			"id":                       loadBalancer.LoadBalancerId,
			"region_id":                loadBalancer.RegionId,
			"master_availability_zone": loadBalancer.MasterZoneId,
			"slave_availability_zone":  loadBalancer.SlaveZoneId,
			"name":                     loadBalancer.LoadBalancerName,
			"network_type":             loadBalancer.NetworkType,
			"vpc_id":                   loadBalancer.VpcId,
			"vswitch_id":               loadBalancer.VSwitchId,
			"address":                  loadBalancer.Address,
			"creation_time":            loadBalancer.CreateTime,
			"tags":                     slbService.tagsToMap(tags),
		}

		ids = append(ids, loadBalancer.LoadBalancerId)
		names = append(names, loadBalancer.LoadBalancerName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("slbs", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

// DescribeLoadBalancersResponse is the response struct for api DescribeLoadBalancers
type DescribeLoadBalancersResponse struct {
	*responses.BaseResponse
	RequestId     string        `json:"RequestId" xml:"RequestId"`
	PageNumber    int           `json:"PageNumber" xml:"PageNumber"`
	PageSize      int           `json:"PageSize" xml:"PageSize"`
	TotalCount    int           `json:"TotalCount" xml:"TotalCount"`
	LoadBalancers LoadBalancers `json:"LoadBalancers" xml:"LoadBalancers"`
}
type LoadBalancers struct {
	LoadBalancer []LoadBalancer `json:"LoadBalancer" xml:"LoadBalancer"`
}
type LoadBalancer struct {
	VpcId                        string `json:"VpcId" xml:"VpcId"`
	CreateTimeStamp              int64  `json:"CreateTimeStamp" xml:"CreateTimeStamp"`
	LoadBalancerId               string `json:"LoadBalancerId" xml:"LoadBalancerId"`
	CreateTime                   string `json:"CreateTime" xml:"CreateTime"`
	PayType                      string `json:"PayType" xml:"PayType"`
	AddressType                  string `json:"AddressType" xml:"AddressType"`
	NetworkType                  string `json:"NetworkType" xml:"NetworkType"`
	ServiceManagedMode           string `json:"ServiceManagedMode" xml:"ServiceManagedMode"`
	SpecBpsFlag                  bool   `json:"SpecBpsFlag" xml:"SpecBpsFlag"`
	AddressIPVersion             string `json:"AddressIPVersion" xml:"AddressIPVersion"`
	LoadBalancerName             string `json:"LoadBalancerName" xml:"LoadBalancerName"`
	Bandwidth                    int    `json:"Bandwidth" xml:"Bandwidth"`
	Address                      string `json:"Address" xml:"Address"`
	SlaveZoneId                  string `json:"SlaveZoneId" xml:"SlaveZoneId"`
	MasterZoneId                 string `json:"MasterZoneId" xml:"MasterZoneId"`
	InternetChargeTypeAlias      string `json:"InternetChargeTypeAlias" xml:"InternetChargeTypeAlias"`
	LoadBalancerSpec             string `json:"LoadBalancerSpec" xml:"LoadBalancerSpec"`
	SpecType                     string `json:"SpecType" xml:"SpecType"`
	RegionId                     string `json:"RegionId" xml:"RegionId"`
	ModificationProtectionReason string `json:"ModificationProtectionReason" xml:"ModificationProtectionReason"`
	ModificationProtectionStatus string `json:"ModificationProtectionStatus" xml:"ModificationProtectionStatus"`
	VSwitchId                    string `json:"VSwitchId" xml:"VSwitchId"`
	LoadBalancerStatus           string `json:"LoadBalancerStatus" xml:"LoadBalancerStatus"`
	ResourceGroupId              string `json:"ResourceGroupId" xml:"ResourceGroupId"`
	InternetChargeType           string `json:"InternetChargeType" xml:"InternetChargeType"`
	BusinessStatus               string `json:"BusinessStatus" xml:"BusinessStatus"`
	DeleteProtection             string `json:"DeleteProtection" xml:"DeleteProtection"`
	RegionIdAlias                string `json:"RegionIdAlias" xml:"RegionIdAlias"`
	Tags                         []Tag  `json:"Tags" xml:"Tags"`
}
