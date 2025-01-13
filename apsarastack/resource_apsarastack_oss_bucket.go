package apsarastack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApsaraStackOssBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackOssBucketCreate,
		Read:   resourceApsaraStackOssBucketRead,
		Update: resourceApsaraStackOssBucketUpdate,
		Delete: resourceApsaraStackOssBucketDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 63),
				Default:      resource.PrefixedUniqueId("tf-oss-bucket-"),
			},
			"acl": {
				Type:         schema.TypeString,
				Default:      oss.ACLPrivate,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "public-read", "public-read-write"}, false),
			},
			"logging": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				MaxItems: 1,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if k == "logging.#" && old == "1" && new == "0" {
						loggings := d.Get("logging").([]interface{})
						logging := loggings[0].(map[string]interface{})
						if logging["target_bucket"] == "" && logging["target_prefix"] == "" {
							return true
						}
					}
					return false
				},
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extranet_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"intranet_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_class": {
				Type:     schema.TypeString,
				Default:  oss.StorageStandard,
				Optional: true,
				ForceNew: true,
			},
			"vpclist": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"bucket_sync": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"storage_capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},
			"sse_algorithm": {
				Type:         schema.TypeString,
				Default:      "",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"", "AES256", "SM4", "KMS"}, false),
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApsaraStackOssBucketCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ossService := OssService{client}
	request := map[string]string{"bucketName": d.Get("bucket").(string), "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	var requestInfo *oss.Client
	bucketName := d.Get("bucket").(string)
	det, err := ossService.DescribeOssBucket(bucketName)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "IsBucketExist", ApsaraStackOssGoSdk)
	}
	acl := d.Get("acl").(string)
	storageClass := d.Get("storage_class")
	if storageClass == "" {
		storageClass = "Standard"
	}
	if acl == "" {
		acl = "private"
	}

	// If not present, Create Bucket
	if det.BucketInfo.Name == "" {
		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"Product":          "OneRouter",
			"Department":       client.Department,
			"ResourceGroup":    client.ResourceGroup,
			"RegionId":         client.RegionId,
			"Action":           "DoOpenApi",
			"AccountInfo":      "123456",
			"Version":          "2018-12-12",
			"SignatureVersion": "1.0",
			"OpenApiAction":    "PutBucket",
			"ProductName":      "oss",
		}
		queryParams := map[string]interface{}{
			"Department":          client.Department,
			"ResourceGroup":       client.ResourceGroup,
			"RegionId":            client.RegionId,
			"asVersion":           "enterprise",
			"asArchitechture":     "x86",
			"haAlibabacloudStack": "true",
			"Language":            "en",
			"BucketName":          bucketName,
			"StorageClass":        storageClass,
			"x-oss-acl":           acl,
		}
		if querybytes, err := json.Marshal(queryParams); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "json Marshal", "CreateBucket", ApsaraStackOssGoSdk)
		} else {
			request.QueryParams["Params"] = string(querybytes)
		}

		request.Method = "POST"        // Set request method
		request.Product = "OneRouter"  // Specify product
		request.Version = "2018-12-12" // Specify product version
		request.ServiceCode = "OneRouter"
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		} // Set request scheme. Default: http
		request.ApiName = "DoOpenApi"
		request.Headers = map[string]string{"RegionId": client.RegionId}

		raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
			return ossClient.ProcessCommonRequest(request)
		})
		log.Printf("Response of Create Bucket: %s", raw)
		log.Printf("Bresponse ossbucket before error")
		if err != nil {
			if ossNotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
			}
			return WrapErrorf(err, DefaultErrorMsg, bucketName, "CreateBucketInfo", ApsaraStackOssGoSdk)
		}
		log.Printf("Bresponse ossbucket after error")
		addDebug("CreateBucketInfo", raw, requestInfo, request)
		log.Printf("Bresponse ossbucket check")
		bresponse, _ := raw.(*responses.CommonResponse)
		log.Printf("Bresponse ossbucket %s", bresponse)
		//headers := bresponse.GetHttpHeaders()
		//if headers["X-Acs-Response-Success"][0] == "false" {
		//	if len(headers["X-Acs-Response-Errorhint"]) > 0 {
		//		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss", "API Action", headers["X-Acs-Response-Errorhint"][0])
		//	} else {
		//		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss", "API Action", bresponse.GetHttpContentString())
		//	}
		//}

		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "CreateBucket", ApsaraStackOssGoSdk)
		}
		//logging:= make(map[string]interface{})
		log.Printf("Enter for logging")

		//addDebug("CreateBucket", raw, requestInfo, bresponse.GetHttpContentString())

	}

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		det, err := ossService.DescribeOssBucket(bucketName)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if det.BucketInfo.Name == "" {
			return resource.RetryableError(Error("Trying to ensure new OSS bucket %#v has been created successfully.", request["bucketName"]))
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "Bucket Not Found", ApsaraStackOssGoSdk)
	}

	// Assign the bucket name as the resource ID
	d.SetId(bucketName)
	//if v := d.Get("logging"); v != nil {
	//	log.Printf("Enter for logging condition passed")
	//
	//	err = resourceApsaraStackOssBucketLoggingCreate(client, d)
	//	if err != nil {
	//		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "Logging Failed", ApsaraStackOssGoSdk)
	//	}
	//}
	//newlist := d.Get("vpclist").([]interface{})
	//bvclient := meta.(*connectivity.ApsaraStackClient)
	//bvserver := BucketVpcService{bvclient}
	//vpclist, binderr := bvserver.BucketVpcList(bucketName)
	//if binderr != nil {
	//	return WrapError(binderr)
	//}
	//oldlist := vpclist.VpcList
	//vpc_err := checkVpcListChange(oldlist, newlist, d, meta)
	//if vpc_err != nil {
	//	return WrapError(vpc_err)
	//}
	return resourceApsaraStackOssBucketUpdate(d, meta)
}

func resourceApsaraStackOssBucketRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	ossService := OssService{client}
	object, err := ossService.DescribeOssBucket(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	logging, err := resourceApsaraStackOssBucketLoggingDescribe(client, d)
	log.Printf("read describe logging %v", logging)
	d.Set("bucket", d.Id())
	if object.BucketInfo.Name == "" {
		log.Print("read: BucketInfo fail!!!!!!")
	}
	d.Set("creation_date", object.BucketInfo.CreationDate.Format("2006-01-02"))
	d.Set("extranet_endpoint", object.BucketInfo.ExtranetEndpoint)
	d.Set("intranet_endpoint", object.BucketInfo.IntranetEndpoint)
	d.Set("location", object.BucketInfo.Location)
	d.Set("owner", object.BucketInfo.Owner.ID)
	d.Set("storage_class", object.BucketInfo.StorageClass)

	var list []map[string]interface{}
	desclog := logging.Data.BucketLoggingStatus.LoggingEnabled
	list = append(list, map[string]interface{}{"target_bucket": desclog.TargetBucket, "target_prefix": desclog.TargetPrefix})

	if err = d.Set("logging", list); err != nil {
		return WrapError(err)
	}
	bvclient := meta.(*connectivity.ApsaraStackClient)
	bvserver := BucketVpcService{bvclient}
	vpclist, binderr := bvserver.BucketVpcList(d.Get("bucket").(string))
	if binderr != nil {
		return WrapError(binderr)
	}
	var vlist []string
	if len(vpclist.VpcList) > 0 {
		for _, v := range vpclist.VpcList {
			vpc := v.(map[string]interface{})
			vlist = append(vlist, vpc["vpcId"].(string))
		}
	}
	d.Set("vpclist", vlist)

	bucketName := d.Get("bucket").(string)

	// 获取同城容灾信息
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"Product":          "OneRouter",
		"Department":       client.Department,
		"ResourceGroup":    client.ResourceGroup,
		"RegionId":         client.RegionId,
		"Action":           "DoOpenApi",
		"AccountInfo":      "123456",
		"Version":          "2018-12-12",
		"SignatureVersion": "1.0",
		"ProductName":      "oss",
		"OpenApiAction":    "GetBucketSync",
		"Params":           fmt.Sprintf("{\"BucketName\":\"%s\"}", bucketName),
	}

	raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
		return ossClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if ossNotFoundError(err) {
			return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
		}
		return WrapErrorf(err, DefaultErrorMsg, bucketName, "CreateBucketInfo", ApsaraStackOssGoSdk)
	}
	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "GetBucketSync", ApsaraStackOssGoSdk)
	}
	bucketSync := BucketSyncResponse{}
	json.Unmarshal([]byte(bresponse.GetHttpContentString()), &bucketSync)
	d.Set("bucket_sync", true)
	for _, rule := range bucketSync.Data.ReplicationConfiguration.Rule {
		if rule.Status == "closing" {
			// 容灾关系是成对出现的
			d.Set("bucket_sync", false)
			break
		}
	}

	// 获取acl信息

	request = requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"Product":          "OneRouter",
		"Department":       client.Department,
		"ResourceGroup":    client.ResourceGroup,
		"RegionId":         client.RegionId,
		"Action":           "DoOpenApi",
		"AccountInfo":      "123456",
		"Version":          "2018-12-12",
		"SignatureVersion": "1.0",
		"ProductName":      "oss",
		"OpenApiAction":    "GetBucketAcl",
		"Params":           fmt.Sprintf("{\"BucketName\":\"%s\", \"acl\":\"acl\"}", bucketName),
	}

	raw, err = client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
		return ossClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if ossNotFoundError(err) {
			return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
		}
		return WrapErrorf(err, DefaultErrorMsg, bucketName, "CreateBucketInfo", ApsaraStackOssGoSdk)
	}
	bresponse, _ = raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "GetBucketSync", ApsaraStackOssGoSdk)
	}
	bucketAcl := BucketAclResponse{}
	json.Unmarshal([]byte(bresponse.GetHttpContentString()), &bucketAcl)
	d.Set("acl", bucketAcl.Data.AccessControlPolicy.AccessControlList.Grant)

	// 获取容量限制信息
	request = requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"Product":          "OneRouter",
		"Department":       client.Department,
		"ResourceGroup":    client.ResourceGroup,
		"RegionId":         client.RegionId,
		"Action":           "DoOpenApi",
		"AccountInfo":      "123456",
		"Version":          "2018-12-12",
		"SignatureVersion": "1.0",
		"ProductName":      "oss",
		"OpenApiAction":    "GetBucketStorageCapacity",
		"Params":           fmt.Sprintf("{\"BucketName\":\"%s\"}", bucketName),
	}

	raw, err = client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
		return ossClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if ossNotFoundError(err) {
			return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
		}
		return WrapErrorf(err, DefaultErrorMsg, bucketName, "CreateBucketInfo", ApsaraStackOssGoSdk)
	}
	bresponse, _ = raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "GetBucketSync", ApsaraStackOssGoSdk)
	}
	storageCapacity := BucketStorageCapacityResponse{}
	json.Unmarshal([]byte(bresponse.GetHttpContentString()), &storageCapacity)
	if v, err := strconv.Atoi(storageCapacity.Data.BucketUserQos.StorageCapacity); err == nil {
		d.Set("storage_capacity", v)
	} else {
		return WrapErrorf(err, "Get storage capacity failed")
	}

	// 获取加密信息
	request = requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"Product":          "OneRouter",
		"Department":       client.Department,
		"ResourceGroup":    client.ResourceGroup,
		"RegionId":         client.RegionId,
		"Action":           "DoOpenApi",
		"AccountInfo":      "123456",
		"Version":          "2018-12-12",
		"SignatureVersion": "1.0",
		"OpenApiAction":    "GetBucketEncryption",
		"ProductName":      "oss",
		"Params":           fmt.Sprintf("{\"BucketName\":\"%s\"}", bucketName),
	}

	raw, err = client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
		return ossClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if ossNotFoundError(err) {
			return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
		}
		return WrapErrorf(err, DefaultErrorMsg, bucketName, "CreateBucketInfo", ApsaraStackOssGoSdk)
	}
	bresponse, _ = raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "GetBucketSync", ApsaraStackOssGoSdk)
	}
	storageEncryption := BucketEncryptionResponse{}
	json.Unmarshal([]byte(bresponse.GetHttpContentString()), &storageEncryption)
	if storageEncryption.Code == "NoSuchServerSideEncryptionRule" {
		d.Set("sse_algorithm", "")
	} else {
		d.Set("sse_algorithm", storageEncryption.Data.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
		if storageEncryption.Data.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.SSEAlgorithm == "KMS" {
			d.Set("kms_key_id", storageEncryption.Data.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.KMSMasterKeyID)
		}
	}

	return nil
}

func resourceApsaraStackOssBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	bucketName := d.Get("bucket").(string)

	if (d.IsNewResource() && !d.Get("bucket_sync").(bool)) || (!d.IsNewResource() && d.HasChange("bucket_sync")) {
		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"Product":          "OneRouter",
			"Department":       client.Department,
			"ResourceGroup":    client.ResourceGroup,
			"RegionId":         client.RegionId,
			"Action":           "DoOpenApi",
			"AccountInfo":      "123456",
			"Version":          "2018-12-12",
			"SignatureVersion": "1.0",
			"ProductName":      "oss",
		}
		if v := d.Get("bucket_sync").(bool); v {
			request.QueryParams["OpenApiAction"] = "PutBucketSync"
		} else {
			request.QueryParams["OpenApiAction"] = "DeleteBucketSync"
		}
		request.QueryParams["Params"] = fmt.Sprintf("{\"BucketName\":\"%s\"}", bucketName)

		raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
			return ossClient.ProcessCommonRequest(request)
		})
		if err != nil {
			if ossNotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
			}
			return WrapErrorf(err, DefaultErrorMsg, bucketName, "CreateBucketInfo", ApsaraStackOssGoSdk)
		}
		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "GetBucketSync", ApsaraStackOssGoSdk)
		}
	}

	if (d.IsNewResource() && d.Get("storage_capacity").(int) != -1) || (!d.IsNewResource() && d.HasChange("storage_capacity")) {
		storageCapacity := d.Get("storage_capacity").(int)
		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"Product":          "OneRouter",
			"Department":       client.Department,
			"ResourceGroup":    client.ResourceGroup,
			"RegionId":         client.RegionId,
			"Action":           "DoOpenApi",
			"AccountInfo":      "123456",
			"Version":          "2018-12-12",
			"SignatureVersion": "1.0",
			"OpenApiAction":    "SetBucketStorageCapacity",
			"ProductName":      "oss",
			"Params":           fmt.Sprintf("{\"BucketName\":\"%s\", \"StorageCapacity\":%d}", bucketName, storageCapacity),
			"Content":          fmt.Sprintf("<BucketUserQos><StorageCapacity>%d</StorageCapacity></BucketUserQos>", storageCapacity),
		}

		raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
			return ossClient.ProcessCommonRequest(request)
		})
		if err != nil {
			if ossNotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
			}
			return WrapErrorf(err, DefaultErrorMsg, bucketName, "CreateBucketInfo", ApsaraStackOssGoSdk)
		}
		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "GetBucketSync", ApsaraStackOssGoSdk)
		}
	}

	if d.HasChange("acl") {
		acl := d.Get("acl").(string)
		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"Product":          "OneRouter",
			"Department":       client.Department,
			"ResourceGroup":    client.ResourceGroup,
			"RegionId":         client.RegionId,
			"Action":           "DoOpenApi",
			"AccountInfo":      "123456",
			"Version":          "2018-12-12",
			"SignatureVersion": "1.0",
			"OpenApiAction":    "PutBucketACL",
			"ProductName":      "oss",
			"Params":           fmt.Sprintf("{\"BucketName\":\"%s\", \"x-oss-acl\":\"%s\"}", bucketName, acl),
		}

		raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
			return ossClient.ProcessCommonRequest(request)
		})
		if err != nil {
			if ossNotFoundError(err) {
				return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
			}
			return WrapErrorf(err, DefaultErrorMsg, bucketName, "CreateBucketInfo", ApsaraStackOssGoSdk)
		}
		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "GetBucketSync", ApsaraStackOssGoSdk)
		}
	}

	if d.HasChanges("sse_algorithm", "kms_key_id") {
		if d.Get("sse_algorithm").(string) == "" {
			request := requests.NewCommonRequest()
			if client.Config.Insecure {
				request.SetHTTPSInsecure(client.Config.Insecure)
			}
			request.QueryParams = map[string]string{
				"Product":          "OneRouter",
				"Department":       client.Department,
				"ResourceGroup":    client.ResourceGroup,
				"RegionId":         client.RegionId,
				"Action":           "DoOpenApi",
				"AccountInfo":      "123456",
				"Version":          "2018-12-12",
				"SignatureVersion": "1.0",
				"OpenApiAction":    "DeleteBucketEncryption",
				"ProductName":      "oss",
				"Params":           fmt.Sprintf("{\"BucketName\":\"%s\"}", bucketName),
			}

			raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
				return ossClient.ProcessCommonRequest(request)
			})
			if err != nil {
				if ossNotFoundError(err) {
					return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
				}
				return WrapErrorf(err, DefaultErrorMsg, bucketName, "CreateBucketInfo", ApsaraStackOssGoSdk)
			}
			bresponse, _ := raw.(*responses.CommonResponse)
			if bresponse.GetHttpStatus() != 200 {
				return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "GetBucketSync", ApsaraStackOssGoSdk)
			}
		} else {
			sse_algorithm := d.Get("sse_algorithm").(string)
			kms_key_id := d.Get("kms_key_id").(string)
			request := requests.NewCommonRequest()
			if client.Config.Insecure {
				request.SetHTTPSInsecure(client.Config.Insecure)
			}
			request.QueryParams = map[string]string{
				"Product":          "OneRouter",
				"Department":       client.Department,
				"ResourceGroup":    client.ResourceGroup,
				"RegionId":         client.RegionId,
				"Action":           "DoOpenApi",
				"AccountInfo":      "123456",
				"Version":          "2018-12-12",
				"SignatureVersion": "1.0",
				"OpenApiAction":    "PutBucketEncryption",
				"ProductName":      "oss",
				"Params":           fmt.Sprintf("{\"BucketName\":\"%s\"}", bucketName),
			}
			if sse_algorithm == "KMS" {
				request.QueryParams["Content"] = fmt.Sprintf("<ServerSideEncryptionRule><ApplyServerSideEncryptionByDefault><SSEAlgorithm>KMS</SSEAlgorithm><KMSMasterKeyID>%s</KMSMasterKeyID></ApplyServerSideEncryptionByDefault></ServerSideEncryptionRule>", kms_key_id)
			} else {
				request.QueryParams["Content"] = fmt.Sprintf("<ServerSideEncryptionRule><ApplyServerSideEncryptionByDefault><SSEAlgorithm>%s</SSEAlgorithm></ApplyServerSideEncryptionByDefault></ServerSideEncryptionRule>", sse_algorithm)
			}

			raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {
				return ossClient.ProcessCommonRequest(request)
			})
			if err != nil {
				if ossNotFoundError(err) {
					return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
				}
				return WrapErrorf(err, DefaultErrorMsg, bucketName, "CreateBucketInfo", ApsaraStackOssGoSdk)
			}
			bresponse, _ := raw.(*responses.CommonResponse)
			if bresponse.GetHttpStatus() != 200 {
				return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "GetBucketSync", ApsaraStackOssGoSdk)
			}
		}
	}
	d.Partial(true)
	if d.HasChange("logging") {
		//if err := resourceApsaraStackOssBucketLoggingUpdate(client, d); err != nil {
		//	return WrapError(err)
		//}
		d.SetPartial("logging")
		log.Print("changes in logging")
		err := resourceApsaraStackOssBucketLoggingCreate(client, d)
		if err != nil {
			return WrapError(err)
		}
	}
	if d.HasChange("vpclist") {
		o, n := d.GetChange("vpclist")
		oldlist := o.([]interface{})
		newlist := n.([]interface{})
		vpc_err := checkVpcListChange(oldlist, newlist, d, meta)
		if vpc_err != nil {
			return WrapError(vpc_err)
		}
	}
	d.Partial(false)
	return resourceApsaraStackOssBucketRead(d, meta)
}

func resourceApsaraStackOssBucketDelete(d *schema.ResourceData, meta interface{}) error {
	//bvclient := meta.(*connectivity.ApsaraStackClient)
	//bvserver := BucketVpcService{bvclient}
	//vpclist, binderr := bvserver.BucketVpcList(d.Id())
	//if binderr != nil {
	//	return WrapError(binderr)
	//}
	//var vlist []string
	//if len(vpclist.VpcList) > 0 {
	//	for _, v := range vpclist.VpcList {
	//		vpc := v.(map[string]interface{})
	//		client2 := meta.(*connectivity.ApsaraStackClient)
	//		bvserver := BucketVpcService{client2}
	//		binderr := bvserver.UnBindBucket(vpc["vpcId"].(string), d.Id())
	//		if binderr != nil {
	//			return WrapError(binderr)
	//		}
	//	}
	//}
	//d.Set("vpclist", vlist)
	client := meta.(*connectivity.ApsaraStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	det, err := ossService.DescribeOssBucket(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsBucketExist", ApsaraStackOssGoSdk)
	}
	addDebug("IsBucketExist", det.BucketInfo, requestInfo, map[string]string{"bucketName": d.Id()})
	if det.BucketInfo.Name == "" {
		return nil
	}

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{

			"Product":          "OneRouter",
			"Department":       client.Department,
			"ResourceGroup":    client.ResourceGroup,
			"RegionId":         client.RegionId,
			"Action":           "DoOpenApi",
			"AccountInfo":      "123456",
			"Version":          "2018-12-12",
			"SignatureVersion": "1.0",
			"OpenApiAction":    "DeleteBucket",
			"ProductName":      "oss",
			"Params":           fmt.Sprintf("{\"%s\":%s,\"%s\":%s,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\"}", "Department", client.Department, "ResourceGroup", client.ResourceGroup, "RegionId", client.RegionId, "asVersion", "enterprise", "asArchitechture", "x86", "haApsaraStack", "true", "Language", "en", "BucketName", d.Id(), "StorageClass", "Standard"), //,"x-one-console-endpoint","http://oss-cn-neimeng-env30-d01-a.intra.env30.shuguang.com"),

		}
		request.Method = "POST"        // Set request method
		request.Product = "OneRouter"  // Specify product
		request.Version = "2018-12-12" // Specify product version
		request.ServiceCode = "OneRouter"
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		} // Set request scheme. Default: http
		request.ApiName = "DoOpenApi"
		request.Headers = map[string]string{"RegionId": client.RegionId}

		_, err := client.WithOssNewClient(func(ossClient *ecs.Client) (interface{}, error) {

			return ossClient.ProcessCommonRequest(request)
		})

		if err != nil {
			if ossNotFoundError(err) {
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		}
		det, err := ossService.DescribeOssBucket(d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if det.BucketInfo.Name != "" {
			return resource.RetryableError(Error("Trying to delete OSS bucket %#v successfully.", d.Id()))
		}
		return nil
	})
	return WrapError(ossService.WaitForOssBucket(d.Id(), Deleted, DefaultTimeoutMedium))
}

func checkVpcListChange(oldlist []interface{}, newlist []interface{}, d *schema.ResourceData, meta interface{}) error {
	vpclist := []string{}
	for _, ovpcid := range oldlist {
		isdelete := true
		for _, nvpcid := range newlist {
			if ovpcid == nvpcid {
				isdelete = false
			}
		}
		if isdelete {
			client2 := meta.(*connectivity.ApsaraStackClient)
			bvserver := BucketVpcService{client2}
			binderr := bvserver.UnBindBucket(ovpcid.(string), d.Id())
			if binderr != nil {
				return WrapError(binderr)
			}
		}
	}
	for _, nvpcid := range newlist {
		iscreate := true
		vpclist = append(vpclist, nvpcid.(string))
		for _, ovpcid := range oldlist {
			if ovpcid == nvpcid {
				iscreate = false
			}
		}
		if iscreate {
			client := meta.(*connectivity.ApsaraStackClient)
			vpcServer := VpcService{client}
			vpcdata, err := vpcServer.DescribeVpc(nvpcid.(string))
			if err != nil {
				return WrapError(err)
			}
			client2 := meta.(*connectivity.ApsaraStackClient)
			bvserver := BucketVpcService{client2}
			binderr := bvserver.BindBucket(vpcdata.VpcId, vpcdata.VpcName, vpcdata.CidrBlock, d.Id())
			if binderr != nil {
				return WrapError(binderr)
			}
		}
	}
	return nil
}

func expirationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	return hashcode.String(buf.String())
}

func transitionsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["created_before_date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["storage_class"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	return hashcode.String(buf.String())
}

func resourceApsaraStackOssBucketLoggingCreate(client *connectivity.ApsaraStackClient, d *schema.ResourceData) error {
	describelogging, err := resourceApsaraStackOssBucketLoggingDescribe(client, d)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketLogging", ApsaraStackOssGoSdk)
	}
	var check Logcheck
	if describelogging.Data.BucketLoggingStatus.LoggingEnabled != check {
		log.Printf("logging is not null %v", d.Get("logging"))
		if _, v := d.GetOk("logging"); v == false {
			log.Print("logging is being disabled")
			logrequest := requests.NewCommonRequest()
			if client.Config.Insecure {
				logrequest.SetHTTPSInsecure(client.Config.Insecure)
			}
			logrequest.QueryParams = map[string]string{

				"Product":          "OneRouter",
				"Department":       client.Department,
				"ResourceGroup":    client.ResourceGroup,
				"RegionId":         client.RegionId,
				"Action":           "DoOpenApi",
				"AccountInfo":      "123456",
				"Version":          "2018-12-12",
				"SignatureVersion": "1.0",
				"OpenApiAction":    "PutBucketLogging",
				"ProductName":      "oss",
				"Content":          fmt.Sprint("<BucketLoggingStatus></BucketLoggingStatus>"),
				//"Content": oss-accesslog/",
				"Params": fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", d.Id()),
				//"Params": "{\"BucketName\":\"source-sample-bucket\"}",

				//"Params":           fmt.Sprintf("{\"%s\":%s,\"%s\":%s,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\"}", "Department", client.Department, "ResourceGroup", client.ResourceGroup, "RegionId", client.RegionId, "asVersion", "enterprise", "asArchitechture", "x86", "haApsaraStack", "true", "Language", "en", "BucketName", bucketName, "StorageClass", storageClass, "x-oss-acl", acl, "SSEAlgorithm", sse_algo), //,"x-one-console-endpoint","http://oss-cn-neimeng-env30-d01-a.intra.env30.shuguang.com"),

			}
			logrequest.Method = "POST"        // Set request method
			logrequest.Product = "OneRouter"  // Specify product
			logrequest.Version = "2018-12-12" // Specify product version
			logrequest.ServiceCode = "OneRouter"
			if strings.ToLower(client.Config.Protocol) == "https" {
				logrequest.Scheme = "https"
			} else {
				logrequest.Scheme = "http"
			} // Set request scheme. Default: http
			logrequest.ApiName = "DoOpenApi"
			logrequest.Headers = map[string]string{"RegionId": client.RegionId}

			raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {

				return ossClient.ProcessCommonRequest(logrequest)
			})
			log.Printf("Response of Logging Bucket: %s", raw)
			if err != nil {
				if ossNotFoundError(err) {
					return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
				}
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "CreateBucketInfo", ApsaraStackOssGoSdk)
			}
			log.Printf("deleting logs oss done")
			//}

		} else {
			logging := make(map[string]interface{})
			log.Print("logging to be updated")
			if v := d.Get("logging"); v != nil {
				log.Print("logging is being enabled")
				all, ok := v.([]interface{})
				if ok {
					log.Printf("printall %v", all)
					for _, a := range all {
						logging, _ = a.(map[string]interface{})
						log.Printf("check target_bucket %v", logging["target_bucket"])
						log.Printf("check target_prefix %v", logging["target_prefix"])
					}
					bucket := fmt.Sprint(logging["target_bucket"])
					log.Printf("checking bucket %v", bucket)
					//b, _ :=json.Marshal(logging)
					//log.Printf("checking b %v",b)
					//bucket:= bytes.NewBuffer(b).String()
					//log.Printf("Checking buckets %v",bucket)
					ossService := OssService{client}
					_, err := ossService.DescribeOssBucket(bucket)

					if err != nil {
						return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "DescribeBucket")
					}
					logrequest := requests.NewCommonRequest()
					if client.Config.Insecure {
						logrequest.SetHTTPSInsecure(client.Config.Insecure)
					}
					logrequest.QueryParams = map[string]string{

						"Product":          "OneRouter",
						"Department":       client.Department,
						"ResourceGroup":    client.ResourceGroup,
						"RegionId":         client.RegionId,
						"Action":           "DoOpenApi",
						"AccountInfo":      "123456",
						"Version":          "2018-12-12",
						"SignatureVersion": "1.0",
						"OpenApiAction":    "PutBucketLogging",
						"ProductName":      "oss",
						"Content":          fmt.Sprint("<BucketLoggingStatus><LoggingEnabled><TargetBucket>", logging["target_bucket"], "</TargetBucket><TargetPrefix>", logging["target_prefix"], "</TargetPrefix></LoggingEnabled></BucketLoggingStatus>"),
						//"Content": oss-accesslog/",
						"Params": fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", d.Id()),
						//"Params": "{\"BucketName\":\"source-sample-bucket\"}",

						//"Params":           fmt.Sprintf("{\"%s\":%s,\"%s\":%s,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\"}", "Department", client.Department, "ResourceGroup", client.ResourceGroup, "RegionId", client.RegionId, "asVersion", "enterprise", "asArchitechture", "x86", "haApsaraStack", "true", "Language", "en", "BucketName", bucketName, "StorageClass", storageClass, "x-oss-acl", acl, "SSEAlgorithm", sse_algo), //,"x-one-console-endpoint","http://oss-cn-neimeng-env30-d01-a.intra.env30.shuguang.com"),

					}
					logrequest.Method = "POST"        // Set request method
					logrequest.Product = "OneRouter"  // Specify product
					logrequest.Version = "2018-12-12" // Specify product version
					logrequest.ServiceCode = "OneRouter"
					if strings.ToLower(client.Config.Protocol) == "https" {
						logrequest.Scheme = "https"
					} else {
						logrequest.Scheme = "http"
					} // Set request scheme. Default: http
					logrequest.ApiName = "DoOpenApi"
					logrequest.Headers = map[string]string{"RegionId": client.RegionId}

					raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {

						return ossClient.ProcessCommonRequest(logrequest)
					})
					log.Printf("Response of Logging Bucket: %s", raw)
					if err != nil {
						if ossNotFoundError(err) {
							return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
						}
						return WrapErrorf(err, DefaultErrorMsg, d.Id(), "CreateBucketInfo", ApsaraStackOssGoSdk)
					}
					log.Printf("logging oss done")
				}

			}
		}
	} else {
		logging := make(map[string]interface{})
		log.Print("logging is  null")
		if v := d.Get("logging"); v != nil {
			log.Print("logging is being enabled")
			all, ok := v.([]interface{})
			if ok {
				log.Printf("printall %v", all)
				for _, a := range all {
					logging, _ = a.(map[string]interface{})
					log.Printf("check target_bucket %v", logging["target_bucket"])
					log.Printf("check target_prefix %v", logging["target_prefix"])
				}
				bucket := fmt.Sprint(logging["target_bucket"])
				log.Printf("checking bucket %v", bucket)
				//b, _ :=json.Marshal(logging)
				//log.Printf("checking b %v",b)
				//bucket:= bytes.NewBuffer(b).String()
				//log.Printf("Checking buckets %v",bucket)
				ossService := OssService{client}
				_, err := ossService.DescribeOssBucket(bucket)

				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "DescribeBucket")
				}
				logrequest := requests.NewCommonRequest()
				if client.Config.Insecure {
					logrequest.SetHTTPSInsecure(client.Config.Insecure)
				}
				logrequest.QueryParams = map[string]string{

					"Product":          "OneRouter",
					"Department":       client.Department,
					"ResourceGroup":    client.ResourceGroup,
					"RegionId":         client.RegionId,
					"Action":           "DoOpenApi",
					"AccountInfo":      "123456",
					"Version":          "2018-12-12",
					"SignatureVersion": "1.0",
					"OpenApiAction":    "PutBucketLogging",
					"ProductName":      "oss",
					"Content":          fmt.Sprint("<BucketLoggingStatus><LoggingEnabled><TargetBucket>", logging["target_bucket"], "</TargetBucket><TargetPrefix>", logging["target_prefix"], "</TargetPrefix></LoggingEnabled></BucketLoggingStatus>"),
					"Params":           fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", d.Id()),
				}
				logrequest.Method = "POST"        // Set request method
				logrequest.Product = "OneRouter"  // Specify product
				logrequest.Version = "2018-12-12" // Specify product version
				logrequest.ServiceCode = "OneRouter"
				if strings.ToLower(client.Config.Protocol) == "https" {
					logrequest.Scheme = "https"
				} else {
					logrequest.Scheme = "http"
				} // Set request scheme. Default: http
				logrequest.ApiName = "DoOpenApi"
				logrequest.Headers = map[string]string{"RegionId": client.RegionId}

				raw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {

					return ossClient.ProcessCommonRequest(logrequest)
				})
				log.Printf("Response of Logging Bucket: %s", raw)
				if err != nil {
					if ossNotFoundError(err) {
						return WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
					}
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), "CreateBucketInfo", ApsaraStackOssGoSdk)
				}
				log.Printf("logging oss done")
			}

		}
	}

	return nil
}
func resourceApsaraStackOssBucketLoggingDescribe(client *connectivity.ApsaraStackClient, d *schema.ResourceData) (*Logging, error) {

	logdescribe := requests.NewCommonRequest()
	if client.Config.Insecure {
		logdescribe.SetHTTPSInsecure(client.Config.Insecure)
	}
	describelogging := Logging{}
	logdescribe.QueryParams = map[string]string{

		"Product":           "OneRouter",
		"Department":        client.Department,
		"ResourceGroup":     client.ResourceGroup,
		"RegionId":          client.RegionId,
		"Action":            "DoOpenApi",
		"AccountInfo":       "123456",
		"Forwardedregionid": client.RegionId,
		"Version":           "2018-12-12",
		"SignatureVersion":  "1.0",
		"OpenApiAction":     "GetBucketLogging",
		"ProductName":       "oss",
		"Params":            fmt.Sprintf("{\"%s\":\"%s\"}", "BucketName", d.Id()),
	}
	logdescribe.Method = "POST"        // Set request method
	logdescribe.Product = "OneRouter"  // Specify product
	logdescribe.Version = "2018-12-12" // Specify product version
	logdescribe.ServiceCode = "OneRouter"
	if strings.ToLower(client.Config.Protocol) == "https" {
		logdescribe.Scheme = "https"
	} else {
		logdescribe.Scheme = "http"
	} // Set request scheme. Default: http
	logdescribe.ApiName = "DoOpenApi"
	logdescribe.Headers = map[string]string{"RegionId": client.RegionId}

	lograw, err := client.WithEcsClient(func(ossClient *ecs.Client) (interface{}, error) {

		return ossClient.ProcessCommonRequest(logdescribe)
	})
	log.Printf("Response of Logging Bucket: %s", lograw)
	if err != nil {
		return &describelogging, WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketLogging", ApsaraStackOssGoSdk)
	}

	osslog, _ := lograw.(*responses.CommonResponse)
	_ = json.Unmarshal(osslog.GetHttpContentBytes(), &describelogging)
	log.Printf("describerawlogging %v", osslog)
	log.Printf("describelogging %v", describelogging)

	return &describelogging, nil
}

type Logcheck struct {
	TargetPrefix string `json:"TargetPrefix"`
	TargetBucket string `json:"TargetBucket"`
}

type Logging struct {
	Data struct {
		BucketLoggingStatus struct {
			LoggingEnabled struct {
				TargetPrefix string `json:"TargetPrefix"`
				TargetBucket string `json:"TargetBucket"`
			} `json:"LoggingEnabled"`
		} `json:"BucketLoggingStatus"`
	} `json:"Data"`
	API string `json:"api"`
}
