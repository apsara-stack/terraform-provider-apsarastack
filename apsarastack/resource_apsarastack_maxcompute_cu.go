package apsarastack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceApsaraStackMaxcomputeCu() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackMaxcomputeCuCreate,
		Read:   resourceApsaraStackMaxcomputeCuRead,
		Delete: resourceApsaraStackMaxcomputeCuDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"max_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"cu_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 27),
			},
			"cu_num": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntAtLeast(1),
				Required:     true,
				ForceNew:     true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceApsaraStackMaxcomputeCuCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var response map[string]interface{}
	action := "CreateUpdateOdpsCu"
	product := "ascm"
	request := make(map[string]interface{})
	conn, err := client.NewAscmClient()
	if err != nil {
		return WrapError(err)
	}
	request["CuName"] = d.Get("cu_name")
	request["CuNum"] = d.Get("cu_num")
	request["ClusterName"] = d.Get("cluster_name")
	//request["Department"] = client.Department
	request["OrganizationId"] = client.Department
	request["ResourceGroupId"] = client.ResourceGroup
	request["RegionId"] = client.RegionId
	request["RegionName"] = client.RegionId
	request["Share"] = "0"
	request["Product"] = product
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_maxcompute_project", action, ApsaraStackSdkGoERROR)
	}
	if fmt.Sprintf(`%v`, response["code"]) != "200" {
		return WrapError(Error("CreateUpdateOdpsCu failed for " + response["asapiErrorMessage"].(string)))
	}

	d.Set("cu_name", request["CuName"])

	return resourceApsaraStackMaxcomputeCuRead(d, meta)
}

func resourceApsaraStackMaxcomputeCuRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	maxcomputeService := MaxcomputeService{client}
	object, err := maxcomputeService.DescribeMaxcomputeCu(d.Get("cu_name").(string))
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource apsarastack_maxcompute_project maxcomputeService.DescribeMaxcomputeCu Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	var data map[string]interface{}
	datas := object["data"].([]interface{})
	if datas == nil || len(datas) < 1 {
		d.SetId(d.Get("max_id").(string))
		d.Set("cluster_name", d.Get("cluster_name").(string))
	}
	s := d.Get("cu_name").(string)
	for _, element := range datas {
		data = element.(map[string]interface{})
		if data["quota_name"].(string) != s {
			continue
		}
		d.SetId(data["id"].(string))
		max_cu, err := data["max_cu"].(json.Number).Float64()
		if err != nil {
			return WrapError(Error("illegal max_cu value"))
		}
		d.Set("cu_num", int64(max_cu))
		d.Set("cluster_name", data["cluster"].(string))
		break
	}
	return nil
}
func resourceApsaraStackMaxcomputeCuDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	action := "DeleteOdpsCu"
	var response map[string]interface{}
	conn, err := client.NewOdpsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"CuId":        d.Id(),
		"CuName":      d.Get("cu_name"),
		"ClusterName": d.Get("cluster_name"),
		"Product":     "ascm",
		"RegionId":    client.RegionId,
		"RegionName":  client.RegionId,
	}

	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-10"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"500"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
	}
	if IsExpectedErrorCodes(fmt.Sprintf("%v", response["code"]), []string{"102", "403"}) {
		return nil
	}
	if fmt.Sprintf(`%v`, response["code"]) != "200" {
		return WrapError(Error("DeleteOdpsCu failed for " + response["Message"].(string)))
	}
	return nil
}
