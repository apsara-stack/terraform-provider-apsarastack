package apsarastack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_ACCESS_KEY", os.Getenv("APSARASTACK_ACCESS_KEY")),
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_SECRET_KEY", os.Getenv("APSARASTACK_SECRET_KEY")),
				Description: descriptions["secret_key"],
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_REGION", os.Getenv("APSARASTACK_REGION")),
				Description: descriptions["region"],
			},
			"security_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_SECURITY_TOKEN", os.Getenv("SECURITY_TOKEN")),
				Description: descriptions["security_token"],
			},
			"ecs_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_ECS_ROLE_NAME", os.Getenv("APSARASTACK_ECS_ROLE_NAME")),
				Description: descriptions["ecs_role_name"],
			},
			"skip_region_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: descriptions["skip_region_validation"],
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["profile"],
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_PROFILE", ""),
			},
			"endpoints": endpointsSchema(),
			"shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["shared_credentials_file"],
				DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_SHARED_CREDENTIALS_FILE", ""),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				DefaultFunc: schema.EnvDefaultFunc("AS_INSECURE", nil),
				Description: descriptions["insecure"],
			},
			"assume_role": assumeRoleSchema(),
			"fc": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'fc' has been deprecated from provider version 1.28.0. New field 'fc' which in nested endpoints instead.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "HTTPS",
				Description:  descriptions["protocol"],
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
			},
			"configuration_source": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  descriptions["configuration_source"],
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"proxy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["proxy"],
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["domain"],
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"apsarastack_instances":                      dataSourceApsaraStackInstances(),
			"apsarastack_disks":                          dataSourceApsaraStackDisks(),
			"apsarastack_key_pairs":                      dataSourceApsaraStackKeyPairs(),
			"apsarastack_network_interfaces":             dataSourceApsaraStackNetworkInterfaces(),
			"apsarastack_instance_type_families":         dataSourceApsaraStackInstanceTypeFamilies(),
			"apsarastack_instance_types":                 dataSourceApsaraStackInstanceTypes(),
			"apsarastack_security_groups":                dataSourceApsaraStackSecurityGroups(),
			"apsarastack_security_group_rules":           dataSourceApsaraStackSecurityGroupRules(),
			"apsarastack_snapshots":                      dataSourceApsaraStackSnapshots(),
			"apsarastack_images":                         dataSourceApsaraStackImages(),
			"apsarastack_vswitches":                      dataSourceApsaraStackVSwitches(),
			"apsarastack_vpcs":                           dataSourceApsaraStackVpcs(),
			"apsarastack_eips":                           dataSourceApsaraStackEips(),
			"apsarastack_slb_listeners":                  dataSourceApsaraStackSlbListeners(),
			"apsarastack_slb_server_groups":              dataSourceApsaraStackSlbServerGroups(),
			"apsarastack_slb_acls":                       dataSourceApsaraStackSlbAcls(),
			"apsarastack_slb_domain_extensions":          dataSourceApsaraStackSlbDomainExtensions(),
			"apsarastack_slb_rules":                      dataSourceApsaraStackSlbRules(),
			"apsarastack_route_tables":                   dataSourceApsaraStackRouteTables(),
			"apsarastack_slb_master_slave_server_groups": dataSourceApsaraStackSlbMasterSlaveServerGroups(),
			"apsarastack_slbs":                           dataSourceApsaraStackSlbs(),
			"apsarastack_slb_zones":                      dataSourceApsaraStackSlbZones(),
			"apsarastack_common_bandwidth_packages":      dataSourceApsaraStackCommonBandwidthPackages(),
			"apsarastack_forward_entries":                dataSourceApsaraStackForwardEntries(),
			"apsarastack_nat_gateways":                   dataSourceApsaraStackNatGateways(),
			"apsarastack_snat_entries":                   dataSourceApsaraStackSnatEntries(),
			"apsarastack_db_instances":                   dataSourceApsaraStackDBInstances(),
			"apsarastack_db_zones":                       dataSourceApsaraStackDBZones(),
			"apsarastack_slb_server_certificates":        dataSourceApsaraStackSlbServerCertificates(),
			"apsarastack_slb_backend_servers":            dataSourceApsaraStackSlbBackendServers(),
			"apsarastack_zones":                          dataSourceApsaraStackZones(),
			"apsarastack_oss_buckets":                    dataSourceApsaraStackOssBuckets(),
			"apsarastack_oss_bucket_objects":             dataSourceApsaraStackOssBucketObjects(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"apsarastack_network_interface":                   resourceApsaraStackNetworkInterface(),
			"apsarastack_network_interface_attachment":        resourceNetworkInterfaceAttachment(),
			"apsarastack_disk":                                resourceApsaraStackDisk(),
			"apsarastack_disk_attachment":                     resourceApsaraStackDiskAttachment(),
			"apsarastack_key_pair":                            resourceApsaraStackKeyPair(),
			"apsarastack_key_pair_attachment":                 resourceApsaraStackKeyPairAttachment(),
			"apsarastack_instance":                            resourceApsaraStackInstance(),
			"apsarastack_ram_role_attachment":                 resourceApsaraStackRamRoleAttachment(),
			"apsarastack_security_group":                      resourceApsaraStackSecurityGroup(),
			"apsarastack_security_group_rule":                 resourceApsaraStackSecurityGroupRule(),
			"apsarastack_launch_template":                     resourceApsaraStackLaunchTemplate(),
			"apsarastack_reserved_instance":                   resourceApsaraStackReservedInstance(),
			"apsarastack_image":                               resourceApsaraStackImage(),
			"apsarastack_image_export":                        resourceApsaraStackImageExport(),
			"apsarastack_image_copy":                          resourceApsaraStackImageCopy(),
			"apsarastack_image_import":                        resourceApsaraStackImageImport(),
			"apsarastack_image_share_permission":              resourceApsaraStackImageSharePermission(),
			"apsarastack_snapshot":                            resourceApsaraStackSnapshot(),
			"apsarastack_snapshot_policy":                     resourceApsaraStackSnapshotPolicy(),
			"apsarastack_vswitch":                             resourceApsaraStackSwitch(),
			"apsarastack_vpc":                                 resourceApsaraStackVpc(),
			"apsarastack_eip":                                 resourceApsaraStackEip(),
			"apsarastack_eip_association":                     resourceApsaraStackEipAssociation(),
			"apsarastack_slb_listener":                        resourceApsaraStackSlbListener(),
			"apsarastack_slb_server_group":                    resourceApsaraStackSlbServerGroup(),
			"apsarastack_slb_acl":                             resourceApsaraStackSlbAcl(),
			"apsarastack_slb_domain_extensions":               resourceApsaraStackSlbDomainExtension(),
			"apsarastack_slb_rule":                            resourceApsaraStackSlbRule(),
			"apsarastack_route_table":                         resourceApsaraStackRouteTable(),
			"apsarastack_route_table_attachment":              resourceApsaraStackRouteTableAttachment(),
			"apsarastack_slb_master_slave_server_group":       resourceApsaraStackSlbMasterSlaveServerGroup(),
			"apsarastack_slb":                                 resourceApsaraStackSlb(),
			"apsarastack_common_bandwidth_package":            resourceApsaraStackCommonBandwidthPackage(),
			"apsarastack_common_bandwidth_package_attachment": resourceApsaraStackCommonBandwidthPackageAttachment(),
			"apsarastack_forward_entry":                       resourceApsaraStackForwardEntry(),
			"apsarastack_nat_gateway":                         resourceApsaraStackNatGateway(),
			"apsarastack_snat_entry":                          resourceApsaraStackSnatEntry(),
			"apsarastack_db_instance":                         resourceApsaraStackDBInstance(),
			"apsarastack_slb_server_certificate":              resourceApsaraStackSlbServerCertificate(),
			"apsarastack_slb_backend_server":                  resourceApsaraStackSlbBackendServer(),
			"apsarastack_oss_bucket":                          resourceApsaraStackOssBucket(),
			"apsarastack_oss_bucket_object":                   resourceApsaraStackOssBucketObject(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var providerConfig map[string]interface{}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var getProviderConfig = func(str string, key string) string {
		if str == "" {
			value, err := getConfigFromProfile(d, key)
			if err == nil && value != nil {
				str = value.(string)
			}
		}
		return str
	}

	accessKey := getProviderConfig(d.Get("access_key").(string), "access_key_id")
	secretKey := getProviderConfig(d.Get("secret_key").(string), "access_key_secret")
	region := getProviderConfig(d.Get("region").(string), "region_id")
	if region == "" {
		region = DEFAULT_REGION
	}

	ecsRoleName := getProviderConfig(d.Get("ecs_role_name").(string), "ram_role_name")

	config := &connectivity.Config{
		AccessKey:            strings.TrimSpace(accessKey),
		SecretKey:            strings.TrimSpace(secretKey),
		EcsRoleName:          strings.TrimSpace(ecsRoleName),
		Region:               connectivity.Region(strings.TrimSpace(region)),
		RegionId:             strings.TrimSpace(region),
		SkipRegionValidation: d.Get("skip_region_validation").(bool),
		ConfigurationSource:  d.Get("configuration_source").(string),
		Protocol:             d.Get("protocol").(string),
		Insecure:             d.Get("insecure").(bool),
		Proxy:                d.Get("proxy").(string),
	}
	token := getProviderConfig(d.Get("security_token").(string), "sts_token")
	config.SecurityToken = strings.TrimSpace(token)

	config.RamRoleArn = getProviderConfig("", "ram_role_arn")
	config.RamRoleSessionName = getProviderConfig("", "ram_session_name")
	expiredSeconds, err := getConfigFromProfile(d, "expired_seconds")
	if err == nil && expiredSeconds != nil {
		config.RamRoleSessionExpiration = (int)(expiredSeconds.(float64))
	}

	assumeRoleList := d.Get("assume_role").(*schema.Set).List()
	if len(assumeRoleList) == 1 {
		assumeRole := assumeRoleList[0].(map[string]interface{})
		if assumeRole["role_arn"].(string) != "" {
			config.RamRoleArn = assumeRole["role_arn"].(string)
		}
		if assumeRole["session_name"].(string) != "" {
			config.RamRoleSessionName = assumeRole["session_name"].(string)
		}
		if config.RamRoleSessionName == "" {
			config.RamRoleSessionName = "terraform"
		}
		config.RamRolePolicy = assumeRole["policy"].(string)
		if assumeRole["session_expiration"].(int) == 0 {
			if v := os.Getenv("APSARASTACK_ASSUME_ROLE_SESSION_EXPIRATION"); v != "" {
				if expiredSeconds, err := strconv.Atoi(v); err == nil {
					config.RamRoleSessionExpiration = expiredSeconds
				}
			}
			if config.RamRoleSessionExpiration == 0 {
				config.RamRoleSessionExpiration = 3600
			}
		} else {
			config.RamRoleSessionExpiration = assumeRole["session_expiration"].(int)
		}

		log.Printf("[INFO] assume_role configuration set: (RamRoleArn: %q, RamRoleSessionName: %q, RamRolePolicy: %q, RamRoleSessionExpiration: %d)",
			config.RamRoleArn, config.RamRoleSessionName, config.RamRolePolicy, config.RamRoleSessionExpiration)
	}

	if err := config.MakeConfigByEcsRoleName(); err != nil {
		return nil, err
	}
	domain := d.Get("domain").(string)
	if domain != "" {
		config.EcsEndpoint = "ecs." + domain
		config.VpcEndpoint = "vpc." + domain
		config.StsEndpoint = "sts." + domain

	} else {

		endpointsSet := d.Get("endpoints").(*schema.Set)

		for _, endpointsSetI := range endpointsSet.List() {
			endpoints := endpointsSetI.(map[string]interface{})
			config.EcsEndpoint = strings.TrimSpace(endpoints["ecs"].(string))
			config.VpcEndpoint = strings.TrimSpace(endpoints["vpc"].(string))
			config.StsEndpoint = strings.TrimSpace(endpoints["sts"].(string))

		}
	}

	if config.RamRoleArn != "" {
		config.AccessKey, config.SecretKey, config.SecurityToken, err = getAssumeRoleAK(config.AccessKey, config.SecretKey, config.SecurityToken, region, config.RamRoleArn, config.RamRoleSessionName, config.RamRolePolicy, config.RamRoleSessionExpiration, config.StsEndpoint)
		if err != nil {
			return nil, err
		}
	}

	if ots_instance_name, ok := d.GetOk("ots_instance_name"); ok && ots_instance_name.(string) != "" {
		config.OtsInstanceName = strings.TrimSpace(ots_instance_name.(string))
	}

	if account, ok := d.GetOk("account_id"); ok && account.(string) != "" {
		config.AccountId = strings.TrimSpace(account.(string))
	}

	if config.ConfigurationSource == "" {
		sourceName := fmt.Sprintf("Default/%s:%s", config.AccessKey, strings.Trim(uuid.New().String(), "-"))
		if len(sourceName) > 64 {
			sourceName = sourceName[:64]
		}
		config.ConfigurationSource = sourceName
	}
	client, err := config.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "The access key for API operations. You can retrieve this from the 'Security Management' section of the ApsaraStack console.",

		"secret_key": "The secret key for API operations. You can retrieve this from the 'Security Management' section of the ApsaraStackconsole.",

		"security_token": "security token. A security token is only required if you are using Security Token Service.",

		"insecure": "Use this to Trust self-signed certificates. It's typically used to allow insecure connections",

		"proxy": "Use this to set proxy connection",

		"domain": "Use this to override the default domain. It's typically used to connect to custom domain.",
	}
}
func endpointsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cbn": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cbn_endpoint"],
				},

				"ecs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ecs_endpoint"],
				},
				"rds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["rds_endpoint"],
				},
				"slb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["slb_endpoint"],
				},
				"vpc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["vpc_endpoint"],
				},
				"cen": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cen_endpoint"],
				},
				"ess": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ess_endpoint"],
				},
				"oss": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["oss_endpoint"],
				},
				"ons": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ons_endpoint"],
				},
				"alikafka": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["alikafka_endpoint"],
				},
				"dns": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dns_endpoint"],
				},
				"ram": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ram_endpoint"],
				},
				"cs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cs_endpoint"],
				},
				"cr": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cr_endpoint"],
				},
				"cdn": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cdn_endpoint"],
				},

				"kms": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["kms_endpoint"],
				},

				"ots": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ots_endpoint"],
				},

				"cms": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cms_endpoint"],
				},

				"pvtz": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["pvtz_endpoint"],
				},

				"sts": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["sts_endpoint"],
				},
				// log service is sls service
				"log": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["log_endpoint"],
				},
				"drds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["drds_endpoint"],
				},
				"dds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dds_endpoint"],
				},
				"polardb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["polardb_endpoint"],
				},
				"gpdb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["gpdb_endpoint"],
				},
				"kvstore": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["kvstore_endpoint"],
				},
				"fc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["fc_endpoint"],
				},
				"apigateway": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["apigateway_endpoint"],
				},
				"datahub": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["datahub_endpoint"],
				},
				"mns": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["mns_endpoint"],
				},
				"location": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["location_endpoint"],
				},
				"elasticsearch": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["elasticsearch_endpoint"],
				},
				"nas": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["nas_endpoint"],
				},
				"actiontrail": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["actiontrail_endpoint"],
				},
				"cas": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cas_endpoint"],
				},
				"bssopenapi": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["bssopenapi_endpoint"],
				},
				"ddoscoo": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ddoscoo_endpoint"],
				},
				"ddosbgp": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ddosbgp_endpoint"],
				},
				"emr": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["emr_endpoint"],
				},
				"market": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["market_endpoint"],
				},
				"adb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["adb_endpoint"],
				},
				"maxcompute": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["maxcompute_endpoint"],
				},
			},
		},
		Set: endpointsToHash,
	}
}
func endpointsToHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["ecs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["rds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["slb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["vpc"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cen"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ess"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["oss"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ons"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["alikafka"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dns"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ram"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cdn"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["kms"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ots"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cms"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["pvtz"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["sts"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["log"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["drds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["gpdb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["kvstore"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["polardb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["fc"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["apigateway"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["datahub"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["mns"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["location"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["elasticsearch"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["nas"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["actiontrail"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cas"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["bssopenapi"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ddoscoo"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ddosbgp"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["emr"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["market"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["adb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cbn"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["maxcompute"].(string)))

	return hashcode.String(buf.String())
}

func getConfigFromProfile(d *schema.ResourceData, ProfileKey string) (interface{}, error) {

	if providerConfig == nil {
		if v, ok := d.GetOk("profile"); !ok && v.(string) == "" {
			return nil, nil
		}
		current := d.Get("profile").(string)
		// Set CredsFilename, expanding home directory
		profilePath, err := homedir.Expand(d.Get("shared_credentials_file").(string))
		if err != nil {
			return nil, WrapError(err)
		}
		if profilePath == "" {
			profilePath = fmt.Sprintf("%s/.apsarastack/config.json", os.Getenv("HOME"))
			if runtime.GOOS == "windows" {
				profilePath = fmt.Sprintf("%s/.apsarastack/config.json", os.Getenv("USERPROFILE"))
			}
		}
		providerConfig = make(map[string]interface{})
		_, err = os.Stat(profilePath)
		if !os.IsNotExist(err) {
			data, err := ioutil.ReadFile(profilePath)
			if err != nil {
				return nil, WrapError(err)
			}
			config := map[string]interface{}{}
			err = json.Unmarshal(data, &config)
			if err != nil {
				return nil, WrapError(err)
			}
			for _, v := range config["profiles"].([]interface{}) {
				if current == v.(map[string]interface{})["name"] {
					providerConfig = v.(map[string]interface{})
				}
			}
		}
	}

	mode := ""
	if v, ok := providerConfig["mode"]; ok {
		mode = v.(string)
	} else {
		return v, nil
	}
	switch ProfileKey {
	case "access_key_id", "access_key_secret":
		if mode == "EcsRamRole" {
			return "", nil
		}
	case "ram_role_name":
		if mode != "EcsRamRole" {
			return "", nil
		}
	case "sts_token":
		if mode != "StsToken" {
			return "", nil
		}
	case "ram_role_arn", "ram_session_name":
		if mode != "RamRoleArn" {
			return "", nil
		}
	case "expired_seconds":
		if mode != "RamRoleArn" {
			return float64(0), nil
		}
	}

	return providerConfig[ProfileKey], nil
}
func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_arn": {
					Type:        schema.TypeString,
					Required:    true,
					Description: descriptions["assume_role_role_arn"],
					DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_ASSUME_ROLE_ARN", ""),
				},
				"session_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_session_name"],
					DefaultFunc: schema.EnvDefaultFunc("APSARASTACK_ASSUME_ROLE_SESSION_NAME", ""),
				},
				"policy": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_policy"],
				},
				"session_expiration": {
					Type:         schema.TypeInt,
					Optional:     true,
					Description:  descriptions["assume_role_session_expiration"],
					ValidateFunc: intBetween(900, 3600),
				},
			},
		},
	}
}

func getAssumeRoleAK(accessKey, secretKey, stsToken, region, roleArn, sessionName, policy string, sessionExpiration int, stsEndpoint string) (string, string, string, error) {
	request := sts.CreateAssumeRoleRequest()
	request.RoleArn = roleArn
	request.RoleSessionName = sessionName
	request.DurationSeconds = requests.NewInteger(sessionExpiration)
	request.Policy = policy
	request.Scheme = "https"
	request.Domain = stsEndpoint

	var client *sts.Client
	var err error
	if stsToken == "" {
		client, err = sts.NewClientWithAccessKey(region, accessKey, secretKey)
	} else {
		client, err = sts.NewClientWithStsToken(region, accessKey, secretKey, stsToken)
	}

	if err != nil {
		return "", "", "", err
	}

	response, err := client.AssumeRole(request)
	if err != nil {
		return "", "", "", err
	}

	return response.Credentials.AccessKeyId, response.Credentials.AccessKeySecret, response.Credentials.SecurityToken, nil
}
