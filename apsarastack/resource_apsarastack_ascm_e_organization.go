package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strings"
	"time"
)

func resourceApsaraStackAscmOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmOrganizationCreate,
		Read:   resourceApsaraStackAscmOrganizationRead,
		Update: resourceApsaraStackAscmOrganizationUpdate,
		Delete: resourceApsaraStackAscmOrganizationDelete,
		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"person_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApsaraStackAscmOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}

	var requestInfo *ecs.Client
	name := d.Get("name").(string)
	check, err := ascmService.DescribeAscmOrganization(name)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_organization", "ORG alreadyExist", ApsaraStackSdkGoERROR)
	}
	parentid := d.Get("parent_id").(string)

	if len(check.Data) == 0 {
		request := requests.NewCommonRequest()
		request.Method = "POST"
		request.Domain = client.Domain
		request.RegionId = client.RegionId
		request.Product = "Ascm"
		request.Version = "2019-05-10"
		request.Scheme = "http"
		request.ApiName = "CreateOrganization"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"Product":         "ascm",
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"RegionId":        client.RegionId,
			"Action":          "CreateOrganization",
			"ParentId":        parentid,
			"name":            name,
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_organization", "CreateOrganization", raw)
		}
		addDebug("CreateOrganization", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_organization", "CreateOrganization", ApsaraStackSdkGoERROR)
		}
		addDebug("CreateOrganization", raw, requestInfo, bresponse.GetHttpContentString())
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		check, err = ascmService.DescribeAscmOrganization(name)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(check.Data) != 0 {
			return nil
		}
		return resource.RetryableError(Error("New Organization has been created successfully."))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_organization", "Failed to create Organization", ApsaraStackSdkGoERROR)
	}

	d.SetId(check.Data[0].Name + SLASH_SEPARATED + fmt.Sprint(check.Data[0].ID))

	return resourceApsaraStackAscmOrganizationUpdate(d, meta)

}

func resourceApsaraStackAscmOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceApsaraStackAscmOrganizationRead(d, meta)

}

func resourceApsaraStackAscmOrganizationRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmOrganization(d.Id())
	did := strings.Split(d.Id(), SLASH_SEPARATED)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if len(object.Data) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("org_id", did[1])
	d.Set("name", did[0])
	d.Set("parent_id", object.Data[0].ParentID)

	return nil

}
func resourceApsaraStackAscmOrganizationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmOrganization(d.Id())
	did := strings.Split(d.Id(), SLASH_SEPARATED)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsOrganizationExist", ApsaraStackSdkGoERROR)
	}

	addDebug("IsOrganizationExist", check, requestInfo, map[string]string{"id": did[1]})
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		if len(check.Data) != 0 {
			request := requests.NewCommonRequest()
			request.QueryParams = map[string]string{
				"RegionId":        client.RegionId,
				"AccessKeySecret": client.SecretKey,
				"Department":      client.Department,
				"ResourceGroup":   client.ResourceGroup,
				"Product":         "ascm",
				"Action":          "RemoveOrganization",
				"Version":         "2019-05-10",
				"ProductName":     "ascm",
				"id":              did[1],
			}

			request.Method = "POST"
			request.Product = "ascm"
			request.Version = "2019-05-10"
			request.ServiceCode = "ascm"
			request.Domain = client.Domain
			request.Scheme = "http"
			request.ApiName = "RemoveOrganization"
			request.Headers = map[string]string{"RegionId": client.RegionId}
			request.RegionId = client.RegionId

			_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
				return csClient.ProcessCommonRequest(request)
			})
			if err != nil {
				return resource.RetryableError(err)
			}
			check, err = ascmService.DescribeAscmOrganization(d.Id())

			if err != nil {
				return resource.NonRetryableError(err)
			}
			if did[0] != "" {
				return resource.RetryableError(Error("Trying to delete Organization %#v successfully.", did[0]))
			}
		}
		return nil
	})
	return nil
}
