package apsarastack

import (
	"fmt"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackImageExport() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackImageExportCreate,
		Read:   resourceApsaraStackImageExportRead,
		Delete: resourceApsaraStackImageExportDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"oss_bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"oss_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceApsaraStackImageExportCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client: client}

	request := ecs.CreateExportImageRequest()
	request.RegionId = client.RegionId
	request.ImageId = d.Get("image_id").(string)
	request.OSSBucket = d.Get("oss_bucket").(string)
	request.OSSPrefix = d.Get("oss_prefix").(string)
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ExportImage(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_image_export", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response := raw.(*ecs.ExportImageResponse)
	taskId := response.TaskId
	d.SetId(request.ImageId)
	stateConf := BuildStateConf([]string{"ExportImage", "Waiting", "Processing"}, []string{"Finished"}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, ecsService.TaskStateRefreshFunc(taskId, []string{"CreateFailed", "UnAvailable"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceApsaraStackImageExportRead(d, meta)

}

func resourceApsaraStackImageExportRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client: client}

	object, err := ecsService.DescribeImageById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("image_id", object.ImageId)
	return WrapError(err)
}

func resourceApsaraStackImageExportDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.Bucket(d.Get("oss_bucket").(string))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Get("oss_bucket").(string), "OSS Bucket", ApsaraStackOssGoSdk)
	}
	addDebug("OSS Bucket", raw, requestInfo, map[string]string{"bucketName": d.Get("oss_bucket").(string)})
	bucket, _ := raw.(*oss.Bucket)
	var objectName string
	if d.Get("oss_prefix").(string) != "" {
		objectName = fmt.Sprintf(d.Get("oss_prefix").(string) + "_" + d.Id() + "_system.raw.tar.gz")
	} else {
		objectName = fmt.Sprintf(d.Id() + "_system.raw.tar.gz")
	}
	err = bucket.DeleteObject(objectName)
	if err != nil {
		if IsExpectedErrors(err, []string{"No Content", "Not Found"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteObject", ApsaraStackOssGoSdk)
	}

	return WrapError(ossService.WaitForOssBucketObject(bucket, d.Id(), Deleted, DefaultTimeoutMedium))
}
