package apsarastack

import (
	"fmt"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceApsaraStackKmsSecrets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackKmsSecretsRead,
		Schema: map[string]*schema.Schema{
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
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"fetch_tags": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secrets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"planned_delete_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_name": {
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

func dataSourceApsaraStackKmsSecretsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := kms.CreateListSecretsRequest()
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "kms"}
	request.QueryParams["Department"] = client.Department
	request.QueryParams["ResourceGroup"] = client.ResourceGroup

	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []kms.Secret
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[vv.(string)] = vv.(string)
		}
	}
	if v, ok := d.GetOk("fetch_tags"); ok {
		request.FetchTags = fmt.Sprintf("%v", v.(bool))
	}

	tagsMap := make(map[string]interface{})
	if v, ok := d.GetOk("tags"); ok && len(v.(map[string]interface{})) > 0 {
		tagsMap = v.(map[string]interface{})
	}
	for {
		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.ListSecrets(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_kms_secrets", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*kms.ListSecretsResponse)

		for _, item := range response.SecretList.Secret {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.SecretName) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.SecretName]; !ok {
					continue
				}
			}
			if len(tagsMap) > 0 {
				if len(item.Tags.Tag) != len(tagsMap) {
					continue
				}
				match := true
				for _, tag := range item.Tags.Tag {
					if v, ok := tagsMap[tag.TagKey]; !ok || v.(string) != tag.TagValue {
						match = false
						break
					}
				}
				if !match {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(response.SecretList.Secret) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, len(objects))
	names := make([]string, len(objects))
	s := make([]map[string]interface{}, len(objects))

	for i, object := range objects {
		mapping := map[string]interface{}{
			"planned_delete_time": object.PlannedDeleteTime,
			"id":                  object.SecretName,
			"secret_name":         object.SecretName,
		}
		tags := make(map[string]string)
		for _, t := range object.Tags.Tag {
			tags[t.TagKey] = t.TagValue
		}
		mapping["tags"] = tags
		ids[i] = object.SecretName
		names[i] = object.SecretName
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("secrets", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
