package apsarastack

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackImageCreate,
		Read:   resourceApsaraStackImageRead,
		Update: resourceApsaraStackImageUpdate,
		Delete: resourceApsaraStackImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{

			"instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"disk_device_mapping", "snapshot_id"},
			},
			"snapshot_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"instance_id", "disk_device_mapping"},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"image_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"disk_device_mapping": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"instance_id", "snapshot_id"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(5, 2000),
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tags": tagsSchema(),
		},
	}
}
func resourceApsaraStackImageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}

	// Make sure the instance status is Running or Stopped
	if v, ok := d.GetOk("instance_id"); ok {
		instance, err := ecsService.DescribeInstance(v.(string))
		if err != nil {
			return WrapError(err)
		}
		status := Status(instance.Status)
		if status != Running && status != Stopped {
			return WrapError(Error("You must make sure that the status of the specified instance is Running or Stopped. "))
		}
	}

	// The snapshot cannot be a snapshot created before July 15, 2013 (inclusive)
	if snapshotId, ok := d.GetOk("snapshot_id"); ok {
		snapshot, err := ecsService.DescribeSnapshot(snapshotId.(string))
		if err != nil {
			return WrapError(err)
		}
		snapshotCreationTime, err := time.Parse("2006-01-02T15:04:05Z", snapshot.CreationTime)
		if err != nil {
			return WrapErrorf(err, IdMsg, snapshotId)
		}
		deadlineTime, _ := time.Parse("2006-01-02T15:04:05Z", "2013-07-16T00:00:00Z")
		if deadlineTime.After(snapshotCreationTime) {
			return WrapError(Error("the specified snapshot cannot be created on or before July 15, 2013."))
		}
	}
	request := ecs.CreateCreateImageRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	if instanceId, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = instanceId.(string)
	}
	if value, ok := d.GetOk("disk_device_mapping"); ok {
		diskDeviceMappings := value.([]interface{})
		if diskDeviceMappings != nil && len(diskDeviceMappings) > 0 {
			mappings := make([]ecs.CreateImageDiskDeviceMapping, 0, len(diskDeviceMappings))
			for _, diskDeviceMapping := range diskDeviceMappings {
				mapping := diskDeviceMapping.(map[string]interface{})
				deviceMapping := ecs.CreateImageDiskDeviceMapping{
					SnapshotId: mapping["snapshot_id"].(string),
					Size:       strconv.Itoa(mapping["size"].(int)),
				}
				mappings = append(mappings, deviceMapping)
			}
			request.DiskDeviceMapping = &mappings
		}
	}

	tags := d.Get("tags").(map[string]interface{})
	if tags != nil && len(tags) > 0 {
		imageTags := make([]ecs.CreateImageTag, 0, len(tags))
		for k, v := range tags {
			imageTag := ecs.CreateImageTag{
				Key:   k,
				Value: v.(string),
			}
			imageTags = append(imageTags, imageTag)
		}
		request.Tag = &imageTags
	}
	if snapshotId, ok := d.GetOk("snapshot_id"); ok {
		request.SnapshotId = snapshotId.(string)
	}

	request.ImageName = d.Get("image_name").(string)
	request.Description = d.Get("description").(string)

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.CreateImage(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectInstanceStatus"}) {
				time.Sleep(time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ecs.CreateImageResponse)
		d.SetId(response.ImageId)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"Creating"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 3*time.Minute, ecsService.ImageStateRefreshFunc(d.Id(), []string{"CreateFailed", "UnAvailable"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceApsaraStackImageRead(d, meta)
}
func resourceApsaraStackImageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}
	err := ecsService.updateImage(d)
	if err != nil {
		return WrapError(err)
	}
	return resourceApsaraStackImageRead(d, meta)
}
func resourceApsaraStackImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	ecsService := EcsService{client}
	object, err := ecsService.DescribeImageById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("image_name", object.ImageName)
	d.Set("description", object.Description)
	d.Set("disk_device_mapping", FlattenImageDiskDeviceMappings(object.DiskDeviceMappings.DiskDeviceMapping))
	tags := object.Tags.Tag
	if len(tags) > 0 {
		err = d.Set("tags", ecsService.tagsToMap(tags))
	}
	return WrapError(err)
}

func resourceApsaraStackImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ecsService := EcsService{client}
	return ecsService.deleteImage(d)
}

func FlattenImageDiskDeviceMappings(list []ecs.DiskDeviceMapping) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		size, _ := strconv.Atoi(i.Size)
		l := map[string]interface{}{
			"size":        size,
			"snapshot_id": i.SnapshotId,
		}
		result = append(result, l)
	}

	return result
}
