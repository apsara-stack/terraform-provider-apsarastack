package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"regexp"

	//"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

const (
	KubernetesClusterNetworkTypeFlannel = "flannel"
	KubernetesClusterNetworkTypeTerway  = "terway"

	KubernetesClusterLoggingTypeSLS = "SLS"
)

var (
	KubernetesClusterNodeCIDRMasksByDefault = 24
)

func resourceApsaraStackCSKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackCSKubernetesCreate,
		Read:   resourceApsaraStackCSKubernetesRead,
		Update: resourceApsaraStackCSKubernetesUpdate,
		Delete: resourceApsaraStackCSKubernetesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(1, 63),
				ConflictsWith: []string{"name_prefix"},
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "Terraform-Creation",
				ValidateFunc:  validation.StringLenBetween(0, 37),
				ConflictsWith: []string{"name"},
				Deprecated:    "Field 'name_prefix' has been deprecated from provider version 1.75.0.",
			},
			// master configurations
			//"master_vswitch_ids": {
			//	Type:     schema.TypeList,
			//	Required: true,
			//	Elem: &schema.Schema{
			//		Type:         schema.TypeString,
			//		ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
			//	},
			//	MinItems:         3,
			//	MaxItems:         5,
			//	DiffSuppressFunc: csForceUpdateSuppressFunc,
			//},
			//"master_instance_types": {
			//	Type:     schema.TypeList,
			//	Required: true,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeString,
			//	},
			//	MinItems:         3,
			//	MaxItems:         5,
			//	DiffSuppressFunc: csForceUpdateSuppressFunc,
			//},
			"master_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          40,
				ValidateFunc:     validation.IntBetween(40, 500),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"master_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudEfficiency,
				ValidateFunc: validation.StringInSlice([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD)}, false),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"master_instance_charge_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Default:          PostPaid,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"master_period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          Month,
				ValidateFunc:     validation.StringInSlice([]string{"Week", "Month"}, false),
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"master_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				// must be a valid period, expected [1-9], 12, 24, 36, 48 or 60,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36, 48, 60})),
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"master_auto_renew": {
				Type:             schema.TypeBool,
				Default:          false,
				Optional:         true,
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			"master_auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: csKubernetesMasterPostPaidDiffSuppressFunc,
			},
			// worker configurations
			//"worker_vswitch_ids": {
			//	Type:     schema.TypeList,
			//	Required: true,
			//	Elem: &schema.Schema{
			//		Type:         schema.TypeString,
			//		ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
			//	},
			//	MinItems:         1,
			//	DiffSuppressFunc: csForceUpdateSuppressFunc,
			//},
			//"worker_instance_types": {
			//	Type:     schema.TypeList,
			//	Required: true,
			//	Elem: &schema.Schema{
			//		Type: schema.TypeString,
			//	},
			//	MinItems: 1,
			//	MaxItems: 10,
			//},
			"worker_number": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"worker_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          40,
				ValidateFunc:     validation.IntBetween(20, 32768),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DiskCloudEfficiency,
				ValidateFunc: validation.StringInSlice([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD)}, false),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_data_disk_size": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          40,
				ValidateFunc:     validation.IntBetween(20, 32768),
				DiffSuppressFunc: workerDataDiskSizeSuppressFunc,
			},
			"worker_data_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(DiskCloudEfficiency), string(DiskCloudSSD)}, false),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_data_disks": {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"category": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"all", "cloud", "ephemeral_ssd", "cloud_essd", "cloud_efficiency", "cloud_ssd", "local_disk"}, false),
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"device": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"encrypted": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_snapshot_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"worker_instance_charge_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
				Default:          PostPaid,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"worker_period_unit": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          Month,
				ValidateFunc:     validation.StringInSlice([]string{"Week", "Month"}, false),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"worker_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36, 48, 60})),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"worker_auto_renew": {
				Type:             schema.TypeBool,
				Default:          false,
				Optional:         true,
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"worker_auto_renew_period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: csKubernetesWorkerPostPaidDiffSuppressFunc,
			},
			"exclude_autoscaler_nodes": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			// global configurations
			// Terway network
			"pod_vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MaxItems:         10,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			// Flannel network
			"pod_cidr": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"service_cidr": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"node_cidr_mask": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          KubernetesClusterNodeCIDRMasksByDefault,
				ValidateFunc:     validation.IntBetween(24, 28),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"new_nat_gateway": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"password": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ConflictsWith:    []string{"key_name", "kms_encrypted_password"},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"key_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"password", "kms_encrypted_password"},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"kms_encrypted_password": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"password", "key_name"},
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"user_ca": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"enable_ssh": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"image_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: imageIdSuppressFunc,
			},
			"install_cloud_monitor": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// cpu policy options of kubelet
			"cpu_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "static"}, false),
			},
			"proxy_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"iptables", "ipvs"}, false),
			},
			"addons": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"slb_internet_enabled": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
			},
			// computed parameters
			"kube_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_ca_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connections": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_server_internet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_server_intranet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"slb_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Field 'slb_id' has been deprecated from provider version 1.9.2. New field 'slb_internet' replaces it.",
			},
			"slb_internet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slb_intranet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_enterprise_security_group": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"security_group_id"},
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"worker_nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			// remove parameters below
			// mix vswitch_ids between master and worker is not a good guidance to create cluster
			"worker_instance_type": {
				Type:     schema.TypeString,
				Required: true,
				//Removed:  "Field 'worker_instance_type' has been removed from provider version 1.75.0. New field 'worker_instance_types' replaces it.",
			},
			"vswitch_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`^vsw-[a-z0-9]*$`), "should start with 'vsw-'."),
				},
				MinItems:         3,
				MaxItems:         5,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
				Removed:          "Field 'vswitch_ids' has been removed from provider version 1.75.0. New field 'master_vswitch_ids' and 'worker_vswitch_ids' replace it.",
			},
			// single instance type would cause extra troubles
			"master_instance_type": {
				Type:     schema.TypeString,
				Required: true,
				//Removed:  "Field 'master_instance_type' has been removed from provider version 1.75.0. New field 'master_instance_types' replaces it.",
			},
			// force update is a high risk operation
			"force_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Removed:  "Field 'force_update' has been removed from provider version 1.75.0.",
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// single az would be never supported.
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				//Removed:  "Field 'vswitch_id' has been removed from provider version 1.75.0. New field 'master_vswitch_ids' and 'worker_vswitch_ids' replaces it.",
			},
			// worker_numbers in array is a hell of management
			"worker_numbers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:    schema.TypeInt,
					Default: 3,
				},
				MinItems:         1,
				MaxItems:         3,
				DiffSuppressFunc: csForceUpdateSuppressFunc,
				Removed:          "Field 'worker_numbers' has been removed from provider version 1.75.0. New field 'worker_number' replaces it.",
			},
			"nodes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Removed:  "Field 'nodes' has been removed from provider version 1.9.4. New field 'master_nodes' replaces it.",
			},
			// too hard to use this config
			"log_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{KubernetesClusterLoggingTypeSLS}, false),
							Required:     true,
						},
						"project": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				DiffSuppressFunc: csForceUpdateSuppressFunc,
				Removed:          "Field 'log_config' has been removed from provider version 1.75.0. New field 'addons' replaces it.",
			},
			"cluster_network_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{KubernetesClusterNetworkTypeFlannel, KubernetesClusterNetworkTypeTerway}, false),
				DiffSuppressFunc: csForceUpdateSuppressFunc,
				Removed:          "Field 'cluster_network_type' has been removed from provider version 1.75.0. New field 'addons' replaces it.",
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_name_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^customized,[a-z0-9]([-a-z0-9\.])*,([5-9]|[1][0-2]),([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`), "Each node name consists of a prefix, an IP substring, and a suffix. For example, if the node IP address is 192.168.0.55, the prefix is aliyun.com, IP substring length is 5, and the suffix is test, the node name will be aliyun.com00055test."),
			},
			"worker_ram_role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_account_issuer": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"api_audiences": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},
		},
	}
}

//func resourceApsaraStackCSKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*connectivity.ApsaraStackClient)
//
//	//csService := CsService{client}
//	Name:= d.Get("name").(string)
//	vpcId:= d.Get("vpc_id").(string)
//	vswitchID:= d.Get("vswitch_id").(string)
//	request:= requests.NewCommonRequest()
//	request.Method = "POST"                // Set request method
//	request.Product = "CS"            // Specify product
//	request.Version = "2015-12-15"            // Specify product version
//	request.Scheme = "http"
//	request.ServiceCode = "cs"
//	request.ApiName = "CreateCluster"
//	request.Headers = map[string]string{"RegionId": client.RegionId}
//	//request.QueryParams = map[string]string{"AccessKeyId": client.AccessKey, "AccessKeySecret": client.SecretKey, "Product": "CS", "RegionId": client.RegionId, "Action": "CreateCluster", "Version": cs.CSAPIVersion,  "Department": client.Department, "ResourceGroup": client.ResourceGroup}
//	//"X-acs-body": "{\"Product\":\"Cs\",\"OsType\":\"Linux\",\"Platform\":\"CentOS\",\"KubernetesVersion\":\"1.14.8-aliyun.1\",\"cluster_type\":\"Kubernetes\",\"RegionId\":\"cn-neimeng-env30-d01\",\"timeout_mins\":60,\"disable_rollback\":true,\"kubernetes_version\":\"1.14.8-aliyun.1\",\"container_cidr\":\"10.16.0.1/16\",\"service_cidr\":\"10.17.0.1/16\",\"name\":\"apsara_Suraj1\",\"vpcid\":\"vpc-0rvua7x0f1ho7pzqi3aap\",\"vswitchid\":\"vsw-0rvhzq0ja4yqxf8q9i6tg\",\"master_instance_type\":\"ecs.n4.4xlarge\",\"master_system_disk_category\":\"cloud_ssd\",\"worker_instance_type\":\"ecs.n4.4xlarge\",\"worker_system_disk_category\":\"cloud_efficiency\",\"worker_data_disk_category\":\"\",\"login_Password\":\"test@123\",\"master_system_disk_size\":200,\"worker_data_disk_size\":0,\"worker_system_disk_size\":40,\"num_of_nodes\":3,\"master_count\":3,\"worker_data_disk\":true,\"snat_entry\":true,\"endpoint_public_access\":true,\"ssh_flags\":true,\"deletion_protection\":true}"
//	request.QueryParams = map[string]string{
//		"RegionId":         client.RegionId,
//		"AccessKeySecret":  client.SecretKey,
//		"Product":          "CS",
//		"Department":       client.Department,
//		"ResourceGroup":    client.ResourceGroup,
//		"Action":           "CreateCluster",
//		"AccountInfo":      "123456",
//		"Version":          "2015-12-15",
//		"SignatureVersion": "1.0",
//		"ProductName":      "cs",
//		"X-acs-body": fmt.Sprintf("{\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%d,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%d,\"%s\":%d,\"%s\":%d,\"%s\":%d,\"%s\":%d,\"%s\":%t,\"%s\":%t,\"%s\":%t,\"%s\":%t,\"%s\":%t}",
//			"Product", "Cs",
//			"OsType", "Linux",
//			"Platform", "CentOS",
//			"KubernetesVersion", KubernetesVersion,
//			"cluster_type", "Kubernetes",
//			"RegionId", "cn-neimeng-env30-d01",
//			"timeout_mins", 60,
//			"disable_rollback", true,
//			"kubernetes_version", "1.14.8-aliyun.1",
//			"container_cidr", "10.14.0.0/16",
//			"service_cidr", "10.15.0.0/16",
//			"name", Name,
//			"vpcid", vpcId,
//			"vswitchid", vswitchID,
//			"master_instance_type", MasterInstanceType,
//			"master_system_disk_category", "cloud_ssd",
//			"worker_instance_type", WorkerInstanceType,
//			"worker_system_disk_category", wsysdiskcat,
//			"worker_data_disk_category", wdatadiskcat,
//			"login_Password", LoginPassword,
//			"master_system_disk_size", 200,
//			"worker_data_disk_size", wdatadisksize,
//			"worker_system_disk_size", wsysdisksize,
//			"num_of_nodes", NumOfNodes,
//			"master_count", 3,
//			"worker_data_disk", true,
//			"snat_entry", SnatEntry,
//			"endpoint_public_access", end,
//			"ssh_flags", true,
//			"deletion_protection", true,
//		),
//	}
//	request.Method = "POST" // Set request method
//	request.Product = "CS"  // Specify product
//	//request.Domain = endpointsSchema().     // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
//	request.Version = "2015-12-15" // Specify product version
//	request.ServiceCode = "cs"
//	request.Scheme = "http" // Set request scheme. Default: http
//	request.ApiName = "CreateCluster"
//	request.Headers = map[string]string{"RegionId": client.RegionId}
//
//	request.RegionId = client.RegionId
//
//
//
//	Clusterresponse := cs.Cluster{}
//
//
//	for {
//		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
//			return csClient.ProcessCommonRequest(request)
//		})
//		if err != nil {
//			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_cs_kubernetes", request.GetActionName(), ApsaraStackSdkGoERROR)
//		}
//		resp, _ := raw.(*responses.CommonResponse)
//		request.TransToAcsRequest()
//
//		err = json.Unmarshal(resp.GetHttpContentBytes(), &Clusterresponse)
//		if err != nil {
//			//fmt.Printf("Response %s", resp.GetHttpContentString())
//			//return WrapError(err)
//		}
//		log.Printf("Suraj success Response data %s", Clusterresponse[0].ClusterID)
//
//		//resp := responses.BaseResponse{}
//		if Clusterresponse[0].ClusterID == "" || len(Clusterresponse) < 1 {
//			break
//		}
//	}
//	//Clusterresponse[0].Name = d.Get("name").(string)
//	//Clusterresponse[0].ClusterID = d.Get("cluster_id").(string)
//	return resourceApsaraStackCSKubernetesUpdate(d, meta)
//}

func resourceApsaraStackCSKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	csService := CsService{client}
	invoker := NewInvoker()
	var requestInfo *cs.Client
	var raw interface{}
	//request := map[string]string{"bucketName": d.Get("bucket").(string), "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	//request.RegionId = client.RegionId
	//ecsService := EcsService{client}
	//var requestInfo *cn.Client
	vpcService := VpcService{client}
	var vswitchID string
	//if list := expandStringList(d.Get("worker_vswitch_ids").([]interface{})); len(list) > 0 {
	//	vswitchID = list[0]
	//} else {
	//	vswitchID = ""
	//}
	vswitchID = d.Get("vswitch_id").(string)
	var vpcId string
	vpcId = d.Get("vpc_id").(string)
	if vpcId == "" {
		vsw, err := vpcService.DescribeVSwitch(vswitchID)
		if err != nil {
			return err
		}
		vpcId = vsw.VpcId
	}

	Name := d.Get("name").(string)
	OsType := "Linux"
	Platform := "CentOS"
	//wsysdisksize := d.Get("worker_disk_size").(int)
	//wsysdiskcat := d.Get("worker_disk_category").(string)
	//wdatadisksize := d.Get("worker_data_disk_size").(int)
	//wdatadiskcat := d.Get("worker_data_disk_category").(string)
	KubernetesVersion := d.Get("version").(string)
	//NodeCidrMask:= strconv.Itoa(d.Get("node_cidr_mask").(int))
	VpcId := vpcId
	//ImageId:=d.Get("image_id").(string)
	var LoginPassword string
	if password := d.Get("password").(string); password == "" {
		if v := d.Get("kms_encrypted_password").(string); v != "" {
			kmsService := KmsService{client}
			decryptResp, err := kmsService.Decrypt(v, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			password = decryptResp.Plaintext
		}
		LoginPassword = password
	} else {
		LoginPassword = password
	}

	end := d.Get("slb_internet_enabled").(bool)
	SnatEntry := d.Get("new_nat_gateway").(bool)
	//MasterVSwitchIds := expandStringList(d.Get("master_vswitch_ids").([]interface{}))
	MasterInstanceType := d.Get("master_instance_type").(string)
	//MasterSystemDiskCategory := d.Get("master_disk_category").(string)
	//MasterSystemDiskSize := int64(d.Get("master_disk_size").(int))
	//if v, ok := d.GetOk("master_instance_charge_type"); ok {
	// MasterInstanceChargeType:= v.(string)
	// if MasterInstanceChargeType == string(PrePaid) {
	//    MasterAutoRenew:= d.Get("master_auto_renew").(bool)
	//    MasterAutoRenewPeriod:= d.Get("master_auto_renew_period").(int)
	//    MasterPeriod:= d.Get("master_period").(int)
	//    MasterPeriodUnit:= d.Get("master_period_unit").(string)
	// }
	//}
	//WorkerVSwitchIds:= expandStringList(d.Get("worker_vswitch_ids").([]interface{}))
	WorkerInstanceType := d.Get("worker_instance_type").(string)
	NumOfNodes := int64(d.Get("worker_number").(int))
	////   WorkerSystemDiskCategory:= d.Get("worker_disk_category").(string)
	// WorkerSystemDiskSize:=     int64(d.Get("worker_disk_size").(int))
	//if v, ok := d.GetOk("worker_instance_charge_type"); ok {
	// WorkerInstanceChargeType:= v.(string)
	// if WorkerInstanceChargeType == string(PrePaid) {
	//    WorkerAutoRenew := d.Get("worker_auto_renew").(bool)
	//    WorkerAutoRenewPeriod:= d.Get("worker_auto_renew_period").(int)
	//    WorkerPeriod:= d.Get("worker_period").(int)
	//    WorkerPeriodUnit:= d.Get("worker_period_unit").(string)
	// }
	//}
	request := requests.NewCommonRequest()
	//request.QueryParams = map[string]string{
	// "AccessKeySecret":  client.SecretKey,
	// "Product":          "cs",
	// "Department":       client.Department,
	// "ResourceGroup":    client.ResourceGroup,
	// "RegionId":         client.RegionId,
	// "Action":           "DoOpenApi",
	// "Params":           fmt.Sprintf("{\"%s\":%s,\"%s\":%s,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\"}",
	//    "Department", client.Department,
	//    "ResourceGroup", client.ResourceGroup,
	//    "RegionId", client.RegionId, "Name", Name,
	//    "OsType", OsType,
	//    "Platform", Platform ,
	//    "KubernetesVersion", KubernetesVersion,
	//    "VpcId", VpcId,
	//    "LoginPassword", LoginPassword,
	//    "Language", "en",
	//    "","",
	// ),
	//}
	scdir := d.Get("service_cidr").(string)
	pcidr := d.Get("pod_cidr").(string)
	request.QueryParams = map[string]string{
		"RegionId":         client.RegionId,
		"AccessKeySecret":  client.SecretKey,
		"Product":          "CS",
		"Department":       client.Department,
		"ResourceGroup":    client.ResourceGroup,
		"Action":           "CreateCluster",
		"AccountInfo":      "123456",
		"Version":          "2015-12-15",
		"SignatureVersion": "1.0",
		"ProductName":      "cs",
		"X-acs-body": fmt.Sprintf("{\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%d,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%d,\"%s\":%d,\"%s\":%t,\"%s\":%t,\"%s\":%t,\"%s\":%t,\"%s\":%t}",
			"Product", "Cs",
			"OsType", OsType,
			"Platform", Platform,
			"KubernetesVersion", KubernetesVersion,
			"cluster_type", "Kubernetes",
			"RegionId", client.RegionId,
			"timeout_mins", 60,
			"disable_rollback", true,
			"kubernetes_version", "1.14.8-aliyun.1",
			"container_cidr", pcidr,
			"service_cidr", scdir,
			"name", Name,
			"vpcid", VpcId,
			"vswitchid", vswitchID,
			"master_instance_type", MasterInstanceType,
			"worker_instance_type", WorkerInstanceType,
			"login_Password", LoginPassword,
			"num_of_nodes", NumOfNodes,
			"master_count", 3,
			"worker_data_disk", true,
			"snat_entry", SnatEntry,
			"endpoint_public_access", end,
			"ssh_flags", true,
			"deletion_protection", false,
		),
	}
	//request.QueryParams = map[string]string{
	//	"RegionId":         client.RegionId,
	//	"AccessKeySecret":  client.SecretKey,
	//	"Product":          "CS",
	//	"Department":       client.Department,
	//	"ResourceGroup":    client.ResourceGroup,
	//	"Action":           "CreateCluster",
	//	"AccountInfo":      "123456",
	//	"Version":          "2015-12-15",
	//	"SignatureVersion": "1.0",
	//	"ProductName":      "cs",
	//	"X-acs-body": fmt.Sprintf("{\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%d,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%d,\"%s\":%d,\"%s\":%d,\"%s\":%d,\"%s\":%d,\"%s\":%t,\"%s\":%t,\"%s\":%t,\"%s\":%t,\"%s\":%t}",
	//		"Product", "CS",
	//		"OsType", OsType,
	//		"Platform", Platform,
	//		"KubernetesVersion", KubernetesVersion,
	//		"cluster_type", "Kubernetes",
	//		"RegionId", "cn-neimeng-env30-d01",
	//		"timeout_mins", 60,
	//		"disable_rollback", true,
	//		"kubernetes_version", "1.14.8-aliyun.1",
	//		"container_cidr",pcidr,
	//		"service_cidr", scdir,
	//		"name", Name,
	//		"vpcid", VpcId,
	//		"vswitchid", vswitchID,
	//		"master_instance_type", MasterInstanceType,
	//		"master_system_disk_category", "cloud_ssd",
	//		"worker_instance_type", WorkerInstanceType,
	//		"worker_system_disk_category", wsysdiskcat,
	//		"worker_data_disk_category", wdatadiskcat,
	//		"login_Password", LoginPassword,
	//		"master_system_disk_size", 200,
	//		"worker_data_disk_size", wdatadisksize,
	//		"worker_system_disk_size", wsysdisksize,
	//		"num_of_nodes", NumOfNodes,
	//		"master_count", 3,
	//		"worker_data_disk", true,
	//		"snat_entry", SnatEntry,
	//		"endpoint_public_access", end,
	//		"ssh_flags", true,
	//		"deletion_protection", true,
	//	),
	//}
	request.Method = "POST" // Set request method
	request.Product = "CS"  // Specify product
	//request.Domain = endpointsSchema().     // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2015-12-15" // Specify product version
	request.ServiceCode = "cs"
	request.Scheme = "http" // Set request scheme. Default: http
	request.ApiName = "CreateCluster"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	//raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
	//	return ecsClient.ProcessCommonRequest(request)
	//})

	var err error
	err = nil
	if err = invoker.Run(func() error {
		raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		return err
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsaracloud_cs_kubernetes", "CreateKubernetesCluster", raw)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["RegionId"] = common.Region(client.RegionId)
		requestMap["Params"] = request.GetQueryParams()
		addDebug("CreateKubernetesCluster", raw, requestInfo, requestMap)
	}
	//if err != nil {
	///    if ossNotFoundError(err) {
	//    return WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
	// }
	// return WrapErrorf(err, DefaultErrorMsg, Name, "CreateCluster", ApsaraStackSdkGoERROR)
	//}
	//addDebug("CreateBucketInfo", raw, requestInfo, request)
	//bresponse, _ := raw.(*responses.CommonResponse)
	//if bresponse.GetHttpStatus() != 200 {
	// return WrapErrorf(err, DefaultErrorMsg, "apsarastack_cs_kubernetes", "CreateCluster", ApsaraStackSdkGoERROR)
	//}
	//addDebug("CreateBucket", raw, requestInfo, bresponse.GetHttpContentString())
	//err = resource.Retry(3*time.Minute, func() *resource.RetryError {
	// det, err := csService.DescribeCsKubernetes(d.Id())
	// //det, err := ossService.DescribeOssBucket(bucketName)
	// if err != nil {
	//    return resource.NonRetryableError(err)
	// }
	// if det.Name == Name {
	//    return resource.RetryableError(Error("Trying to ensure new OSS bucket %#v has been created successfully.", Name))
	// }
	// return nil
	//})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_kubernetes", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	//stateConf := BuildStateConf([]string{"Pending", "Starting", "Stopped"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, CsService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	//
	//if _, err := stateConf.WaitForState(); err != nil {
	// return WrapErrorf(err, IdMsg, d.Id())
	//}
	//requestInfo.CreateCluster()
	//d.SetId(Name)
	//return resourceApsaraStackCSKubernetesUpdate(d, meta)
	if debugOn() {
		//requestMap := make(map[string]interface{})
		//requestMap["RegionId"] = common.Region(client.RegionId)
		//requestMap["Params"] = request
		addDebug("CreateKubernetesCluster", raw, request)
	}
	clusterresponse := ClusterCommonResponse{}
	cluster, _ := raw.(*responses.CommonResponse)
	ok := json.Unmarshal(cluster.GetHttpContentBytes(), &clusterresponse)
	if ok != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_cs_kubernetes", "ParseKubernetesClusterResponse", raw)
	}
	d.SetId(clusterresponse.ClusterID)
	//stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	//if _, err := stateConf.WaitForState(); err != nil {
	// return WrapErrorf(err, IdMsg, d.Id())
	//}
	//// reset interval to 10s
	//stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second,(d.Id(),ecsService.InstanceStateRefreshFunc([]string{"deleting", "failed"})) )
	stateConf := BuildStateConf([]string{"initial"}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"deleting", "failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	////
	//response:=responses.NewCommonResponse()
	//ok=response.IsSuccess()
	//if ok!=true{
	// return WrapErrorf(err, DefaultErrorMsg, "apsarastack_instance", request.GetActionName(), ApsaraStackSdkGoERROR)
	//}

	return resourceApsaraStackCSKubernetesUpdate(d, meta)
}

type Response struct {
	RequestId string `json:"request_id"`
}
type ClusterCommonResponse struct {
	Response
	ClusterID  string `json:"cluster_id"`
	Token      string `json:"token,omitempty"`
	TaskId     string `json:"task_id,omitempty"`
	InstanceId string `json:"instanceId"`
}

func resourceApsaraStackCSKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {

	//client := meta.(*connectivity.ApsaraStackClient)
	//d.Partial(true)
	//request:= requests.NewCommonRequest()
	//request.Method = "POST"                // Set request method
	//request.Product = "Cs"            // Specify product
	//request.Version = "2015-12-15"            // Specify product version
	//request.Scheme = "http"
	//request.ServiceCode = "cs"
	//request.ApiName = "CreateCluster"
	//request.Headers = map[string]string{"RegionId": client.RegionId}
	//request.QueryParams = map[string]string{"AccessKeyId": client.AccessKey, "AccessKeySecret": client.SecretKey, "Product": "Cs", "RegionId": client.RegionId, "Action": "CreateCluster", "Version": cl.CSAPIVersion,  "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	//

	return resourceApsaraStackCSKubernetesRead(d, meta)

}

func resourceApsaraStackCSKubernetesRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.ApsaraStackClient)
	csService := CsService{client}
	//invoker := NewInvoker()
	object, err := csService.DescribeCsKubernetes(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	//client1,err:= sdk.NewClientWithAccessKey(region,access,secret)
	//client1.Domain=endpoint
	//if err!=nil{
	//	fmt.Print("Error in client")
	//}

	//request:= requests.NewCommonRequest()
	//request.Method = "GET"                // Set request method
	//request.Product = "Cs"            // Specify product
	//request.Version = "2015-12-15"            // Specify product version
	//request.Scheme = "http"
	//request.ServiceCode = "cs"
	//request.ApiName = "DescribeClusters"
	//request.Headers = map[string]string{"RegionId": client.RegionId}
	//request.QueryParams = map[string]string{"AccessKeyId": client.AccessKey, "AccessKeySecret": client.SecretKey, "Product": "CS", "RegionId": client.RegionId, "Action": "DescribeClusters", "Version": cl.CSAPIVersion,  "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	//request.RegionId = client.RegionId
	//
	//
	//Clusterresponse := cl.Cluster{}

	//for {
	//	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
	//		return ecsClient.ProcessCommonRequest(request)
	//	})
	//	if err != nil {
	//		return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_cs_kubernetes", request.GetActionName(), ApsaraStackSdkGoERROR)
	//	}
	//	resp, _ := raw.(*responses.CommonResponse)
	//	request.TransToAcsRequest()
	//
	//	err = json.Unmarshal(resp.GetHttpContentBytes(), &Clusterresponse)
	//	if err != nil {
	//	}
	//	log.Printf("Suraj success Response data %s", Clusterresponse.ClusterID)
	//
	//	//resp := responses.BaseResponse{}
	//	if Clusterresponse.Name!="" || len(Clusterresponse.ClusterID) < 1 {
	//		break
	//	}
	//}
	d.Set("name", object.Name)
	d.Set("id", object.ClusterId)
	//d.Set("availability_zone", object.ZoneID)
	d.Set("state", object.State)
	//d.Set("availability_zone", object[0].ZoneID)

	//var names []string
	//var s []map[string]interface{}
	//for _, kc := range Clusterresponse {
	//	mapping := map[string]interface{}{
	//		"id":              kc.ClusterID,
	//		"name":            kc.Name,
	//		"vpc_id": 			kc.VpcID,
	//		"security_group_id": 			kc.SecurityGroupID,
	//		"availability_zone": 			kc.ZoneID,
	//		"state" : kc.State,
	//		"master_instance_types" : []string{kc.Parameters.MasterInstanceType},
	//
	//		"nat_gateway_id" : kc.Parameters.NatGatewayID,
	//		"vswitch_ids" : []string{kc.Parameters.VSwitchID},
	//
	//		"master_disk_category": 			kc.Parameters.MasterSystemDiskCategory,
	//		"cluster_network_type": 			kc.Parameters.Network,
	//	}
	//
	//	names = append(names, kc.Name)
	//
	//	s = append(s, mapping)
	//}

	return nil
}

func resourceApsaraStackCSKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	csService := CsService{client}
	invoker := NewInvoker()
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":         client.RegionId,
		"AccessKeySecret":  client.SecretKey,
		"Product":          "CS",
		"Department":       client.Department,
		"ResourceGroup":    client.ResourceGroup,
		"Action":           "DeleteCluster",
		"Version":          "2015-12-15",
		"SignatureVersion": "1.0",
		"ProductName":      "cs",
		"ClusterId":        d.Id(),
		"X-acs-body":       fmt.Sprintf("{\"%s\":\"%t\",\"%s\":\"%s\"}", "keep_slb", false, "ClusterId", d.Id()),
	}
	request.Method = "POST" // Set request method
	request.Product = "Cs"  // Specify product
	// request.Domain =       // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2015-12-15" // Specify product version
	request.ServiceCode = "cs"
	request.Scheme = "http" // Set request scheme. Default: http
	request.ApiName = "DeleteCluster"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	var response interface{}
	//var raw interface{}
	err := resource.Retry(30*time.Minute, func() *resource.RetryError {
		if err := invoker.Run(func() error {
			raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
				return csClient.ProcessCommonRequest(request)
			})
			response = raw
			return err
		}); err != nil {
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			requestMap := make(map[string]interface{})
			requestMap["ClusterId"] = d.Id()
			addDebug("DeleteCluster", response, d.Id(), requestMap)
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteCluster", ApsaraStackLogGoSdkERROR)
	}

	stateConf := BuildStateConf([]string{"running", "deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, csService.CsKubernetesInstanceStateRefreshFunc(d.Id(), []string{"delete_failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil

}
