package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackLogProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackLogProjectCreate,
		Read:   resourceApsaraStackLogProjectRead,
		Update: resourceApsaraStackLogProjectUpdate,
		Delete: resourceApsaraStackLogProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApsaraStackLogProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	logService := LogService{client}
	name := d.Get("name").(string)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "SLS"
	request.Domain = client.Domain
	request.Version = "2020-03-31"
	request.Scheme = "http"
	request.ApiName = "CreateProject"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "SLS",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "CreateProject",
		"Version":         "2020-03-31",
		"projectName":     name,
		"Description":     d.Get("description").(string),
	}
	raw, err := client.WithEcsClient(func(alidnsClient *ecs.Client) (interface{}, error) {
		return alidnsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_log_project", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug("LogProject", raw)

	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		object, err := logService.DescribeLogProject(name)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if object.ProjectName != "" {
			return nil
		}
		return resource.RetryableError(Error("Failed to describe log project"))
	})
	d.SetId(name)
	return resourceApsaraStackLogProjectRead(d, meta)
}

func resourceApsaraStackLogProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	logService := LogService{client}
	object, err := logService.DescribeLogProject(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object.ProjectName)
	d.Set("description", object.Description)

	return nil
}

func resourceApsaraStackLogProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *sls.Client

	name := d.Id()
	if d.HasChange("description") {
		request := requests.NewCommonRequest()
		request.Method = "POST"
		request.Product = "SLS"
		request.Domain = client.Domain
		request.Version = "2020-03-31"
		request.Scheme = "http"
		request.ApiName = "UpdateProject"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{
			"AccessKeySecret": client.SecretKey,
			"AccessKeyId":     client.AccessKey,
			"Product":         "SLS",
			"RegionId":        client.RegionId,
			"Department":      client.Department,
			"ResourceGroup":   client.ResourceGroup,
			"Action":          "UpdateProject",
			"Version":         "2020-03-31",
			"organizationId":  client.Department,
			"resourceGroupId": client.ResourceGroup,
			"ProjectName":     name,
			"description":     d.Get("description").(string),
		}

		raw, err := client.WithEcsClient(func(slsClient *ecs.Client) (interface{}, error) {
			return slsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateProject", ApsaraStackLogGoSdkERROR)
		}
		addDebug("UpdateProject", raw, requestInfo, request)
	}

	return resourceApsaraStackLogProjectRead(d, meta)
}

func resourceApsaraStackLogProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *sls.Client
	name := d.Get("name").(string)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "SLS"
	request.Domain = client.Domain
	request.Version = "2020-03-31"
	request.Scheme = "http"
	request.ApiName = "DeleteProject"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "SLS",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"RegionId":        client.RegionId,
		"organizationId":  client.Department,
		"resourceGroupId": client.ResourceGroup,
		"Action":          "DeleteProject",
		"Version":         "2020-03-31",
		"ProjectName":     name,
	}
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(slsClient *ecs.Client) (interface{}, error) {
			return slsClient.ProcessCommonRequest(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout, "RequestTimeout"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteProject", raw, requestInfo, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ProjectNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteProject", ApsaraStackLogGoSdkERROR)
	}

	return nil

}
