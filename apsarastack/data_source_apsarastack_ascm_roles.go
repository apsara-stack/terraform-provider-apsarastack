package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"regexp"
	"strings"
)

func dataSourceApsaraStackAscmRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackAscmRolesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"role_type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_level": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"role_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram_role": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"role_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackAscmRolesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	id := d.Get("id").(int)
	request := requests.NewCommonRequest()
	request.Product = "ascm"
	request.Version = "2019-05-10"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.ApiName = "ListRoles"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeyId":     client.AccessKey,
		"AccessKeySecret": client.SecretKey,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Product":         "ascm",
		"RegionId":        client.RegionId,
		"Action":          "ListRoles",
		"Version":         "2019-05-10",
		"pageSize":        "100000",
	}
	response := AscmRoles{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw ListRoles : %s", raw)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_ascm_roles", request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return WrapError(err)
		}
		if response.AsapiErrorCode == "200" || len(response.Data) < 1 {
			break
		}
	}

	var r *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok && nameRegex.(string) != "" {
		r = regexp.MustCompile(nameRegex.(string))
	}
	var ids []string
	var s []map[string]interface{}

	for _, rg := range response.Data {
		if r != nil && !r.MatchString(rg.RoleName) {
			continue
		}
		if response.Data[0].ID == id {
			mapping := map[string]interface{}{
				"id":          rg.ID,
				"name":        rg.RoleName,
				"description": rg.Description,
				"user_count":  rg.UserCount,
				"role_level":  rg.RoleLevel,
				"role_type":   rg.RoleType,
				"role_range":  rg.RoleRange,
				"ram_role":    rg.RAMRole,
			}
			//ids = append(ids, string(rune(rg.ID)))
			s = append(s, mapping)
			break
		} else {
			mapping := map[string]interface{}{
				"id":          rg.ID,
				"name":        rg.RoleName,
				"description": rg.Description,
				"user_count":  rg.UserCount,
				"role_level":  rg.RoleLevel,
				"role_type":   rg.RoleType,
				"role_range":  rg.RoleRange,
				"ram_role":    rg.RAMRole,
			}
			ids = append(ids, string(rune(rg.ID)))
			s = append(s, mapping)
		}
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("roles", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
