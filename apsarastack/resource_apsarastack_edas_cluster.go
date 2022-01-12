package apsarastack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"time"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackEdasCluster() *schema.Resource {
	return &schema.Resource{
		Create: rresourceApsaraStackEdasClusterCreate,
		Read:   resourceApsaraStackEdasClusterRead,
		Delete: resourceApsaraStackEdasClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3}),
			},
			"network_mode": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2}),
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func rresourceApsaraStackEdasClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	edasService := EdasService{client}

	request := edas.CreateInsertClusterRequest()
	request.RegionId = client.RegionId
	request.ClusterName = d.Get("cluster_name").(string)
	request.ClusterType = requests.NewInteger(d.Get("cluster_type").(int))
	request.NetworkMode = requests.NewInteger(d.Get("network_mode").(int))
	request.Headers["x-ascm-product-name"] = "Edas"
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"
	log.Printf("rresourceApsaraStackEdasClusterCreate  client.Department = %s", client.Department)
	log.Printf("rresourceApsaraStackEdasClusterCreate  client.RegionId = %s", client.RegionId)
	if v, ok := d.GetOk("vpc_id"); !ok {
		if d.Get("network_mode") == 2 {
			return WrapError(Error("vpcId is required for vpc network mode"))
		}
	} else {
		request.VpcId = v.(string)
	}

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.InsertCluster(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_edas_cluster", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.InsertClusterResponse)
	if response.Code != 200 {
		return WrapError(Error("create cluster failed for " + response.Message))
	}
	d.SetId(response.Cluster.ClusterId)

	return resourceApsaraStackEdasClusterRead(d, meta)
}

func resourceApsaraStackEdasClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	edasService := EdasService{client}

	clusterId := d.Id()
	regionId := client.RegionId

	request := edas.CreateGetClusterRequest()
	request.RegionId = regionId
	request.ClusterId = clusterId
	request.Headers["x-ascm-product-name"] = "Edas"
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"

	raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
		return edasClient.GetCluster(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_edas_cluster", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*edas.GetClusterResponse)
	if response.Code != 200 {
		return WrapError(Error("create cluster failed for " + response.Message))
	}

	d.Set("cluster_name", response.Cluster.ClusterName)
	d.Set("cluster_type", response.Cluster.ClusterType)
	d.Set("network_mode", response.Cluster.NetworkMode)
	d.Set("region_id", response.Cluster.RegionId)
	d.Set("vpc_id", response.Cluster.VpcId)

	return nil
}

func resourceApsaraStackEdasClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	edasService := EdasService{client}

	clusterId := d.Id()
	regionId := client.RegionId

	request := edas.CreateDeleteClusterRequest()
	request.RegionId = regionId
	request.ClusterId = clusterId
	request.Headers["x-ascm-product-name"] = "Edas"
	request.Headers["x-acs-organizationid"] = client.Department
	request.Headers["x-acs-content-type"] = "application/x-www-form-urlencoded"

	wait := incrementalWait(1*time.Second, 2*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := edasService.client.WithEdasClient(func(edasClient *edas.Client) (interface{}, error) {
			return edasClient.DeleteCluster(request)
		})
		response, _ := raw.(*edas.DeleteClusterResponse)
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if response.Code != 200 {
			if strings.Contains(response.Message, "there are still instances in it") {
				return resource.RetryableError(Error("delete cluster failed for " + response.Message))
			}
			return resource.NonRetryableError(Error("delete cluster failed for " + response.Message))
		}

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	return nil
}
