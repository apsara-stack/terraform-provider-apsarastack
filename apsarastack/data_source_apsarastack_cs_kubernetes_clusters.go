package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"regexp"
)

func dataSourceApsaraStackCSKubernetesClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackCSKubernetesClustersRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slb_internet_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"master_instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"worker_instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"worker_numbers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pod_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_cidr_mask": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"project": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"worker_disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_data_disk_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_data_disk_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_period_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"master_auto_renew": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"master_auto_renew_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"worker_instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_period_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"worker_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"worker_auto_renew": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"worker_auto_renew_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"master_nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"worker_nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"connections": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_server_internet": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"api_server_intranet": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"master_public_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service_domain": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackCSKubernetesClustersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.Method = "GET"
	request.Product = "Cs"
	request.Version = "2015-12-15"

	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ServiceCode = "cs"
	request.ApiName = "DescribeClusters"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeyId": client.AccessKey, "AccessKeySecret": client.SecretKey, "Product": "Cs", "RegionId": client.RegionId, "Action": "DescribeClusters", "Version": cs.CSAPIVersion, "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.RegionId = client.RegionId
	Clusterresponse := []Cluster{}

	for {
		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_cs_kubernetes_clusters", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		resp, _ := raw.(*responses.CommonResponse)
		request.TransToAcsRequest()

		err = json.Unmarshal(resp.GetHttpContentBytes(), &Clusterresponse)
		if err != nil {
			return WrapError(err)
		}
		if Clusterresponse[0].Name != "" || len(Clusterresponse) < 1 {
			break
		}

	}
	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, kc := range Clusterresponse {
		if r != nil && !r.MatchString(kc.Name) {
			continue
		}
		mapping := map[string]interface{}{
			"id":                          kc.ClusterID,
			"name":                        kc.Name,
			"vpc_id":                      kc.VpcID,
			"security_group_id":           kc.SecurityGroupID,
			"availability_zone":           kc.ZoneID,
			"state":                       kc.State,
			"master_instance_types":       []string{kc.Parameters.MasterInstanceType},
			"nat_gateway_id":              kc.Parameters.NatGatewayID,
			"vswitch_ids":                 []string{kc.Parameters.VSwitchID},
			"master_disk_category":        kc.Parameters.MasterSystemDiskCategory,
			"cluster_network_type":        kc.Parameters.Network,
			"pod_cidr":                    kc.SubnetCidr,
			"worker_data_disk_size":       kc.Parameters.WorkerDataDiskSize,
			"worker_disk_category":        kc.Parameters.WorkerDataDiskCategory,
			"worker_instance_types":       []string{kc.Parameters.WorkerInstanceType},
			"worker_instance_charge_type": kc.Parameters.WorkerInstanceChargeType,
			"node_cidr_mask":              kc.Parameters.NodeCIDRMask,
		}

		ids = append(ids, string(kc.ClusterID))
		names = append(names, kc.Name)

		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
