package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"log"
	"strings"

	//	"github.com/aliyun/alibaba-cloud-sdk-go/services/cs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func resourceApsaraStackCSKubernetesNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackCSKubernetesNodePoolCreate,
		Read:   resourceApsaraStackCSKubernetesNodePoolRead,
		Update: resourceApsaraStackCSKubernetesNodePoolUpdate,
		Delete: resourceApsaraStackCSKubernetesNodePoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"password": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instances": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"format_disk": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"keep_instance_name": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceApsaraStackCSKubernetesNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("check meta %v", meta)
	csService := CsService{meta.(*connectivity.ApsaraStackClient)}
	client := meta.(*connectivity.ApsaraStackClient)
	var raw interface{}
	invoker := NewInvoker()
	//nodepoolid := d.Get("nodepool_id").(string)
	var inst string
	var clusterId string
	var insts []string
	var formatDisk, retainIname bool
	clusterId = d.Get("cluster_id").(string)
	if v, ok := d.GetOk("instances"); ok {
		formatDisk = d.Get("format_disk").(bool)
		retainIname = d.Get("keep_instance_name").(bool)
		insts = expandStringList(v.(*schema.Set).List())
		fmt.Print("checking instances attached: ", insts)
		for i, k := range insts {
			if i != 0 {
				inst = fmt.Sprintf("%s\",\"%s", inst, k)

			} else {
				inst = k
			}
		}
	}
	password := d.Get("password").(string)
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":         client.RegionId,
		"AccessKeySecret":  client.SecretKey,
		"Product":          "CS",
		"Department":       client.Department,
		"ResourceGroup":    client.ResourceGroup,
		"Action":           "AttachInstances",
		"AccountInfo":      "123456",
		"Version":          "2015-12-15",
		"SignatureVersion": "1.0",
		"ProductName":      "cs",
		"ClusterId":        clusterId,
		"X-acs-body": fmt.Sprintf("{\"%s\":[\"%s\"],\"%s\":%t,\"%s\":%t,\"%s\":\"%s\"}",

			"instances", inst,
			"format_disk", formatDisk,
			"keep_instance_name", retainIname,
			"password", password,
		),
	}
	request.Method = "POST"        // Set request method
	request.Product = "CS"         // Specify product
	request.Version = "2015-12-15" // Specify product version
	request.ServiceCode = "cs"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	} // Set request scheme. Default: http
	request.ApiName = "AttachInstances"
	request.Headers = map[string]string{"RegionId": client.RegionId}

	var err error
	if err = invoker.Run(func() error {
		raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		addDebug("AttachInstances", raw)

		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_cs_kubernetes_nodepool", "AttachInstances", raw)
	}

	if debugOn() {
		resizeRequestMap := make(map[string]interface{})
		resizeRequestMap["ClusterId"] = clusterId
		resizeRequestMap["Args"] = request.GetQueryParams()
		addDebug("AttachInstances", raw, resizeRequestMap)
	}

	stateConf := BuildStateConf([]string{"scaling"}, []string{"running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, csService.CsKubernetesInstanceStateRefreshFunc(clusterId, []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, clusterId)
	}

	d.SetId(clusterId)

	return resourceApsaraStackCSKubernetesNodePoolUpdate(d, meta)
}

func resourceApsaraStackCSKubernetesNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {

	return resourceApsaraStackCSKubernetesNodePoolRead(d, meta)

}

func resourceApsaraStackCSKubernetesNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceApsaraStackCSKubernetesNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
