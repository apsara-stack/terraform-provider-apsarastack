package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackDRDSInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackDRDSInstanceCreate,
		Read:   resourceApsaraStackDRDSInstanceRead,
		Update: resourceApsaraStackDRDSInstanceUpdate,
		Delete: resourceApsaraStackDRDSInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 129),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"specification": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{string(PostPaid), string(PrePaid)}, false),
				ForceNew:     true,
				Default:      PostPaid,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_series": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"drds.sn2.4c16g", "drds.sn2.8c32g", "drds.sn2.16c64g", "drds.sn1.32c64g"}, false),
				ForceNew:     true,
			},
		},
	}
}

func resourceApsaraStackDRDSInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	drdsService := DrdsService{client}

	request := drds.CreateCreateDrdsInstanceRequest()
	request.RegionId = client.RegionId
	request.Description = d.Get("description").(string)
	//request.Type = "1"
	request.Type = "PRIVATE"
	request.ZoneId = d.Get("zone_id").(string)
	request.Specification = d.Get("specification").(string)
	request.PayType = d.Get("instance_charge_type").(string)
	request.VswitchId = d.Get("vswitch_id").(string)
	request.InstanceSeries = d.Get("instance_series").(string)
	request.Quantity = "1"

	if request.VswitchId != "" {

		vpcService := VpcService{client}
		vsw, err := vpcService.DescribeVSwitch(request.VswitchId)
		if err != nil {
			return WrapError(err)
		}

		request.VpcId = vsw.VpcId
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	if request.PayType == string(PostPaid) {
		request.PayType = "drdsPost"
	}
	if request.PayType == string(PrePaid) {
		request.PayType = "drdsPre"
	}
	request.ResourceGroupId = client.Department

	request.Headers["x-ascm-product-name"] = "Drds"
	request.Headers["x-acs-organizationId"] = client.Department
	raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.CreateDrdsInstance(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_drds_instance", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(CreateDrdsInstanceResponse2)
	idList := response.Data.DrdsInstanceIdList.DrdsInstanceId
	if len(idList) != 1 {
		return WrapError(Error("failed to get DRDS instance id and response. DrdsInstanceIdList is %#v", idList))
	}
	d.SetId(idList[0])

	// wait instance status change from DO_CREATE to RUN
	stateConf := BuildStateConf([]string{"DO_CREATE"}, []string{"RUN"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, drdsService.DrdsInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceApsaraStackDRDSInstanceUpdate(d, meta)

}

type CreateDrdsInstanceResponse2 struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Success   bool   `json:"Success" xml:"Success"`
	Data      Data2  `json:"Data" xml:"Data"`
}
type Data2 struct {
	OrderId            int64              `json:"OrderId" xml:"OrderId"`
	NewestVersion      string             `json:"NewestVersion" xml:"NewestVersion"`
	CreateTime         string             `json:"CreateTime" xml:"CreateTime"`
	Mode               string             `json:"Mode" xml:"Mode"`
	InstRole           string             `json:"InstRole" xml:"InstRole"`
	ShardTbKey         string             `json:"ShardTbKey" xml:"ShardTbKey"`
	Expired            string             `json:"Expired" xml:"Expired"`
	IsActive           bool               `json:"IsActive" xml:"IsActive"`
	Schema             string             `json:"Schema" xml:"Schema"`
	DbInstType         string             `json:"DbInstType" xml:"DbInstType"`
	SourceTableName    string             `json:"SourceTableName" xml:"SourceTableName"`
	ShardDbKey         string             `json:"ShardDbKey" xml:"ShardDbKey"`
	TableName          string             `json:"TableName" xml:"TableName"`
	DbName             string             `json:"DbName" xml:"DbName"`
	Stage              string             `json:"Stage" xml:"Stage"`
	Progress           string             `json:"Progress" xml:"Progress"`
	InstanceVersion    string             `json:"InstanceVersion" xml:"InstanceVersion"`
	RandomCode         string             `json:"RandomCode" xml:"RandomCode"`
	TargetTableName    string             `json:"TargetTableName" xml:"TargetTableName"`
	Msg                string             `json:"Msg" xml:"Msg"`
	Status             string             `json:"Status" xml:"Status"`
	DrdsInstanceIdList DrdsInstanceIdList `json:"DrdsInstanceIdList" xml:"DrdsInstanceIdList"`
	FullRevise         FullRevise         `json:"FullRevise" xml:"FullRevise"`
	Increment          Increment          `json:"Increment" xml:"Increment"`
	//FullCheck          FullCheck               `json:"FullCheck" xml:"FullCheck"`
	//Full               Full                    `json:"Full" xml:"Full"`
	//Review             Review                  `json:"Review" xml:"Review"`
	//List               ListInDescribeHotDbList `json:"List" xml:"List"`
}
type DrdsInstanceIdList struct {
	DrdsInstanceId []string `json:"DrdsInstanceId" xml:"DrdsInstanceId"`
}
type FullRevise struct {
	Expired   int    `json:"Expired" xml:"Expired"`
	Progress  int    `json:"Progress" xml:"Progress"`
	Total     int    `json:"Total" xml:"Total"`
	Tps       int    `json:"Tps" xml:"Tps"`
	StartTime string `json:"StartTime" xml:"StartTime"`
}
type Increment struct {
	Delay     int    `json:"Delay" xml:"Delay"`
	Tps       int    `json:"Tps" xml:"Tps"`
	StartTime string `json:"StartTime" xml:"StartTime"`
}

func resourceApsaraStackDRDSInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	drdsService := DrdsService{client}

	configItem := make(map[string]string)
	if d.HasChange("description") {
		request := drds.CreateModifyDrdsInstanceDescriptionRequest()
		request.DrdsInstanceId = d.Id()
		request.Description = d.Get("description").(string)
		configItem["description"] = request.Description
		client := meta.(*connectivity.ApsaraStackClient)
		request.RegionId = client.RegionId
		request.Headers["x-ascm-product-name"] = "Drds"
		request.Headers["x-acs-organizationId"] = client.Department
		raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
			return drdsClient.ModifyDrdsInstanceDescription(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	//wait for update effected and instance status returning to run
	if err := drdsService.WaitDrdsInstanceConfigEffect(
		d.Id(), configItem, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return WrapError(err)
	}
	stateConf := BuildStateConf([]string{}, []string{"RUN"}, d.Timeout(schema.TimeoutUpdate), 3*time.Second, drdsService.DrdsInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceApsaraStackDRDSInstanceRead(d, meta)
}

func resourceApsaraStackDRDSInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	drdsService := DrdsService{client}

	object, err := drdsService.DescribeDrdsInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	data := object.Data
	//other attribute not set,because these attribute from `data` can't  get
	d.Set("zone_id", data.ZoneId)
	d.Set("description", data.Description)

	return nil
}

func resourceApsaraStackDRDSInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	drdsService := DrdsService{client}
	request := drds.CreateRemoveDrdsInstanceRequest()
	request.RegionId = client.RegionId
	request.DrdsInstanceId = d.Id()

	request.Headers["x-ascm-product-name"] = "Drds"
	request.Headers["x-acs-organizationId"] = client.Department
	raw, err := client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.RemoveDrdsInstance(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDrdsInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*drds.RemoveDrdsInstanceResponse)

	if !response.Success {
		return WrapError(Error("failed to delete instance timeout "+"and got an error: %#v", err))
	}

	//0 -> RUN, 1->DO_CREATE, 2->EXCEPTION, 3->EXPIRE, 4->DO_RELEASE, 5->RELEASE, 6->UPGRADE, 7->DOWNGRADE, 10->VersionUpgrade, 11->VersionRollback, 14->RESTART
	stateConf := BuildStateConf([]string{"RUN", "DO_CREATE", "EXCEPTION", "EXPIRE", "DO_RELEASE", "RELEASE", "UPGRADE", "DOWNGRADE", "VersionUpgrade", "VersionRollback", "RESTART"}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, drdsService.DrdsInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
