package apsarastack

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/maxcompute"

	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"time"
)

func resourceApsaraStackMaxComputeProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackMaxComputeProjectCreate,
		Read:   resourceApsarastackMaxComputeProjectRead,
		Delete: resourceApsarastackMaxComputeProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 27),
			},

			"specification_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"OdpsStandard"}, false),
			},

			"order_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo"}, false),
			},
		},
	}
}

func resourceApsaraStackMaxComputeProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := maxcompute.CreateCreateProjectRequest()

	request.OdpsRegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "maxcompute", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.ProjectName = d.Get("name").(string)
	request.OdpsSpecificationType = d.Get("specification_type").(string)
	request.OrderType = d.Get("order_type").(string)

	raw, err := client.WithMaxComputeClient(func(MaxComputeClient *maxcompute.Client) (interface{}, error) {
		return MaxComputeClient.CreateProject(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_maxcompute_project", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*maxcompute.CreateProjectResponse)

	if response.Code != "200" {
		return WrapError(Error("%v", response))
	}

	d.SetId(request.ProjectName)

	return resourceApsarastackMaxComputeProjectRead(d, meta)
}

func resourceApsarastackMaxComputeProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	maxcomputeService := MaxComputeService{client}
	response, err := maxcomputeService.DescribeMaxComputeProject(d.Id())
	if err != nil {
		return WrapError(err)
	}

	var dat map[string]interface{}

	if err := json.Unmarshal([]byte(response.Data), &dat); err != nil {
		return WrapError(Error("%v", response))
	}
	d.Set("order_type", dat["orderType"].(string))
	d.Set("name", dat["projectName"].(string))

	return nil
}

func resourceApsarastackMaxComputeProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	maxcomputeService := MaxComputeService{client}

	request := maxcompute.CreateDeleteProjectRequest()

	request.RegionIdName = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "maxcompute", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.ProjectName = d.Get("name").(string)

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithMaxComputeClient(func(MaxComputeClient *maxcompute.Client) (interface{}, error) {
			return MaxComputeClient.DeleteProject(request)
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}

		response := raw.(*maxcompute.DeleteProjectResponse)
		if response.Code == "500" {
			return resource.RetryableError(nil)
		}

		if response.Code != "200" {
			return resource.NonRetryableError(err)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		if isProjectNotExistError(response.Data) {
			return nil
		}

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, request.ProjectName, "DeleteProject", ApsarastackMaxComputeSdkGo)
	}
	return WrapError(maxcomputeService.WaitForMaxComputeProject(request.ProjectName, Deleted, DefaultTimeout))

}
