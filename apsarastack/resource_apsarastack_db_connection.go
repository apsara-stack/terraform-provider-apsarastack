package apsarastack

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const dbConnectionSuffixRegex = "\\.mysql\\.([a-zA-Z0-9\\-]+\\.){0,1}rds\\.aliyuncs\\.com"
const dbConnectionIdWithSuffixRegex = "^([a-zA-Z0-9\\-_]+:[a-zA-Z0-9\\-_]+)" + dbConnectionSuffixRegex + "$"

var dbConnectionIdWithSuffixRegexp = regexp.MustCompile(dbConnectionIdWithSuffixRegex)

func resourceApsaraStackDBConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackDBConnectionCreate,
		Read:   resourceApsaraStackDBConnectionRead,
		Update: resourceApsaraStackDBConnectionUpdate,
		Delete: resourceApsaraStackDBConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"connection_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 31),
			},
			"port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDBConnectionPort,
				Default:      "3306",
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApsaraStackDBConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	rdsService := RdsService{client}
	instanceId := d.Get("instance_id").(string)
	prefix := d.Get("connection_prefix").(string)
	if prefix == "" {
		prefix = fmt.Sprintf("%stf", instanceId)
	}

	request := rds.CreateAllocateInstancePublicConnectionRequest()
	request.RegionId = client.RegionId

	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.DBInstanceId = instanceId
	request.ConnectionStringPrefix = prefix
	request.Port = d.Get("port").(string)
	var raw interface{}
	var err error
	err = resource.Retry(8*time.Minute, func() *resource.RetryError {
		raw, err = client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.AllocateInstancePublicConnection(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_db_connection", request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", instanceId, COLON_SEPARATED, request.ConnectionStringPrefix))

	if err := rdsService.WaitForDBConnection(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	// wait instance running after allocating
	if err := rdsService.WaitForDBInstance(instanceId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceApsaraStackDBConnectionRead(d, meta)
}

func resourceApsaraStackDBConnectionRead(d *schema.ResourceData, meta interface{}) error {
	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	client := meta.(*connectivity.ApsaraStackClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeDBConnection(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("instance_id", parts[0])
	d.Set("connection_prefix", parts[1])
	d.Set("port", object.Port)
	d.Set("connection_string", object.ConnectionString)
	d.Set("ip_address", object.IPAddress)

	return nil
}

func resourceApsaraStackDBConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	rdsService := RdsService{client}

	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("port") {
		request := rds.CreateModifyDBInstanceConnectionStringRequest()
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.DBInstanceId = parts[0]
		object, err := rdsService.DescribeDBConnection(d.Id())
		if err != nil {
			return WrapError(err)
		}
		request.CurrentConnectionString = object.ConnectionString
		request.ConnectionStringPrefix = parts[1]
		request.Port = d.Get("port").(string)
		if err := resource.Retry(8*time.Minute, func() *resource.RetryError {
			raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
				return rdsClient.ModifyDBInstanceConnectionString(request)
			})
			if err != nil {
				if IsExpectedErrors(err, OperationDeniedDBStatus) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
		}

		// wait instance running after modifying
		if err := rdsService.WaitForDBInstance(request.DBInstanceId, Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
	}
	return resourceApsaraStackDBConnectionRead(d, meta)
}

func resourceApsaraStackDBConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	rdsService := RdsService{client}

	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	split := strings.Split(d.Id(), COLON_SEPARATED)
	request := rds.CreateReleaseInstancePublicConnectionRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "rds", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	request.DBInstanceId = split[0]

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		object, err := rdsService.DescribeDBConnection(d.Id())
		if err != nil {
			return resource.NonRetryableError(WrapError(err))
		}
		request.CurrentConnectionString = object.ConnectionString
		var raw interface{}
		raw, err = client.WithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.ReleaseInstancePublicConnection(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidCurrentConnectionString.NotFound", "AtLeastOneNetTypeExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), ApsaraStackSdkGoERROR)
	}

	return rdsService.WaitForDBConnection(d.Id(), Deleted, DefaultTimeoutMedium)
}
