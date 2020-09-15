package apsarastack

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
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

			"cors_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_methods": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_origins": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"expose_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"max_age_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
				MaxItems: 10,
			},

			"website": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index_document": {
							Type:     schema.TypeString,
							Required: true,
						},

						"error_document": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				MaxItems: 1,
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
			},

			"referer_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_empty": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"referers": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				MaxItems: 1,
			},

			"lifecycle_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(0, 255),
						},
						"prefix": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"expiration": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      expirationHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validateOssBucketDateTimestamp,
									},
									"days": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"transitions": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      transitionsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"created_before_date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validateOssBucketDateTimestamp,
									},
									"days": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"storage_class": {
										Type:     schema.TypeString,
										Default:  oss.StorageStandard,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(oss.StorageStandard),
											string(oss.StorageIA),
											string(oss.StorageArchive),
										}, false),
									},
								},
							},
						},
					},
				},
				MaxItems: 1000,
			},

			"policy": {
				Type:     schema.TypeString,
				Optional: true,
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
				ValidateFunc: validation.StringInSlice([]string{
					string(oss.StorageStandard),
					string(oss.StorageIA),
					string(oss.StorageArchive),
				}, false),
			},
			"server_side_encryption_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sse_algorithm": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								ServerSideEncryptionAes256,
								ServerSideEncryptionKMS,
							}, false),
						},
					},
				},
				MaxItems: 1,
			},

			"tags": tagsSchema(),

			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"versioning": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Enabled",
								"Suspended",
							}, false),
						},
					},
				},
				MaxItems: 1,
			},
		},
	}
}

func resourceApsaraStackOssBucketCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := map[string]string{"bucketName": d.Get("bucket").(string)}
	var requestInfo *oss.Client
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.IsBucketExist(request["bucketName"])
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "IsBucketExist", ApsaraStackOssGoSdk)
	}
	addDebug("IsBucketExist", raw, requestInfo, request)
	isExist, _ := raw.(bool)
	if isExist {
		return WrapError(Error("[ERROR] The specified bucket name: %#v is not available. The bucket namespace is shared by all users of the OSS system. Please select a different name and try again.", request["bucketName"]))
	}
	type Request struct {
		BucketName string
		Option     oss.Option
	}

	req := Request{
		d.Get("bucket").(string),
		oss.StorageClass(oss.StorageClassType(d.Get("storage_class").(string))),
	}
	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return nil, ossClient.CreateBucket(req.BucketName) //req.Option)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "CreateBucket", ApsaraStackOssGoSdk)
	}
	addDebug("CreateBucket", raw, requestInfo, req)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return ossClient.IsBucketExist(request["bucketName"])
		})

		if err != nil {
			return resource.NonRetryableError(err)
		}
		isExist, _ := raw.(bool)
		if !isExist {
			return resource.RetryableError(Error("Trying to ensure new OSS bucket %#v has been created successfully.", request["bucketName"]))
		}
		addDebug("IsBucketExist", raw, requestInfo, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_oss_bucket", "IsBucketExist", ApsaraStackOssGoSdk)
	}

	// Assign the bucket name as the resource ID
	d.SetId(request["bucketName"])

	return resourceApsaraStackOssBucketUpdate(d, meta)
}

func resourceApsaraStackOssBucketRead(d *schema.ResourceData, meta interface{}) error {
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

	d.Set("bucket", d.Id())

	d.Set("acl", object.BucketInfo.ACL)
	d.Set("creation_date", object.BucketInfo.CreationDate.Format("2006-01-02"))
	d.Set("extranet_endpoint", object.BucketInfo.ExtranetEndpoint)
	d.Set("intranet_endpoint", object.BucketInfo.IntranetEndpoint)
	d.Set("location", object.BucketInfo.Location)
	d.Set("owner", object.BucketInfo.Owner.ID)
	d.Set("storage_class", object.BucketInfo.StorageClass)

	if &object.BucketInfo.SseRule != nil {
		if len(object.BucketInfo.SseRule.SSEAlgorithm) > 0 && object.BucketInfo.SseRule.SSEAlgorithm != "None" {
			rule := make(map[string]interface{})
			rule["sse_algorithm"] = object.BucketInfo.SseRule.SSEAlgorithm
			data := make([]map[string]interface{}, 0)
			data = append(data, rule)
			d.Set("server_side_encryption_rule", data)
		}
	}

	if object.BucketInfo.Versioning != "" {
		data := map[string]interface{}{
			"status": object.BucketInfo.Versioning,
		}
		versioning := make([]map[string]interface{}, 0)
		versioning = append(versioning, data)
		d.Set("versioning", versioning)
	}
	request := map[string]string{"bucketName": d.Id()}
	var requestInfo *oss.Client

	// Read the CORS
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.GetBucketCORS(request["bucketName"])
	})
	if err != nil && !IsExpectedErrors(err, []string{"NoSuchCORSConfiguration"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketCORS", ApsaraStackOssGoSdk)
	}
	addDebug("GetBucketCORS", raw, requestInfo, request)
	cors, _ := raw.(oss.GetBucketCORSResult)
	rules := make([]map[string]interface{}, 0, len(cors.CORSRules))
	for _, r := range cors.CORSRules {
		rule := make(map[string]interface{})
		rule["allowed_headers"] = r.AllowedHeader
		rule["allowed_methods"] = r.AllowedMethod
		rule["allowed_origins"] = r.AllowedOrigin
		rule["expose_headers"] = r.ExposeHeader
		rule["max_age_seconds"] = r.MaxAgeSeconds

		rules = append(rules, rule)
	}
	if err := d.Set("cors_rule", rules); err != nil {
		return WrapError(err)
	}

	// Read the website configuration
	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketWebsite(d.Id())
	})
	if err != nil && !IsExpectedErrors(err, []string{"NoSuchWebsiteConfiguration"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketWebsite", ApsaraStackOssGoSdk)
	}
	addDebug("GetBucketWebsite", raw, requestInfo, request)
	ws, _ := raw.(oss.GetBucketWebsiteResult)
	websites := make([]map[string]interface{}, 0)
	if err == nil && &ws != nil {
		w := make(map[string]interface{})

		if v := &ws.IndexDocument; v != nil {
			w["index_document"] = v.Suffix
		}

		if v := &ws.ErrorDocument; v != nil {
			w["error_document"] = v.Key
		}
		websites = append(websites, w)
	}
	if err := d.Set("website", websites); err != nil {
		return WrapError(err)
	}

	// Read the logging configuration
	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketLogging(d.Id())
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketLogging", ApsaraStackOssGoSdk)
	}
	addDebug("GetBucketLogging", raw, requestInfo, request)
	logging, _ := raw.(oss.GetBucketLoggingResult)

	if &logging != nil {
		enable := logging.LoggingEnabled
		if &enable != nil {
			lgs := make([]map[string]interface{}, 0)
			tb := logging.LoggingEnabled.TargetBucket
			tp := logging.LoggingEnabled.TargetPrefix
			if tb != "" || tp != "" {
				lgs = append(lgs, map[string]interface{}{
					"target_bucket": tb,
					"target_prefix": tp,
				})
			}
			if err := d.Set("logging", lgs); err != nil {
				return WrapError(err)
			}
		}
	}

	// Read the bucket referer
	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketReferer(d.Id())
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketReferer", ApsaraStackOssGoSdk)
	}
	addDebug("GetBucketReferer", raw, requestInfo, request)
	referers := make([]map[string]interface{}, 0)
	referer, _ := raw.(oss.GetBucketRefererResult)
	if len(referer.RefererList) > 0 {
		referers = append(referers, map[string]interface{}{
			"allow_empty": referer.AllowEmptyReferer,
			"referers":    referer.RefererList,
		})
		if err := d.Set("referer_config", referers); err != nil {
			return WrapError(err)
		}
	}

	// Read the lifecycle rule configuration
	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketLifecycle(d.Id())
	})
	if err != nil && !ossNotFoundError(err) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketLifecycle", ApsaraStackOssGoSdk)
	}
	addDebug("GetBucketLifecycle", raw, requestInfo, request)
	lrules := make([]map[string]interface{}, 0)
	lifecycle, _ := raw.(oss.GetBucketLifecycleResult)
	for _, lifecycleRule := range lifecycle.Rules {
		rule := make(map[string]interface{})
		rule["id"] = lifecycleRule.ID
		rule["prefix"] = lifecycleRule.Prefix
		if LifecycleRuleStatus(lifecycleRule.Status) == ExpirationStatusEnabled {
			rule["enabled"] = true
		} else {
			rule["enabled"] = false
		}
		// expiration
		if lifecycleRule.Expiration != nil {
			e := make(map[string]interface{})
			if lifecycleRule.Expiration.Date != "" {
				t, err := time.Parse("2006-01-02T15:04:05.000Z", lifecycleRule.Expiration.Date)
				if err != nil {
					return WrapError(err)
				}
				e["date"] = t.Format("2006-01-02")
			}
			e["days"] = int(lifecycleRule.Expiration.Days)
			rule["expiration"] = schema.NewSet(expirationHash, []interface{}{e})
		}
		// transitions
		if len(lifecycleRule.Transitions) != 0 {
			var eSli []interface{}
			for _, transition := range lifecycleRule.Transitions {
				e := make(map[string]interface{})
				if transition.CreatedBeforeDate != "" {
					t, err := time.Parse("2006-01-02T15:04:05.000Z", transition.CreatedBeforeDate)
					if err != nil {
						return WrapError(err)
					}
					e["created_before_date"] = t.Format("2006-01-02")
				}
				e["days"] = transition.Days
				e["storage_class"] = string(transition.StorageClass)
				eSli = append(eSli, e)
			}
			rule["transitions"] = schema.NewSet(transitionsHash, eSli)
		}

		lrules = append(lrules, rule)
	}

	if err := d.Set("lifecycle_rule", lrules); err != nil {
		return WrapError(err)
	}

	// Read Policy
	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		params := map[string]interface{}{}
		params["policy"] = nil
		return ossClient.Conn.Do("GET", d.Id(), "", params, nil, nil, 0, nil)
	})

	if err != nil && !ossNotFoundError(err) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetPolicyByConn", ApsaraStackOssGoSdk)
	}
	addDebug("GetPolicyByConn", raw, requestInfo, request)
	policy := ""
	if err == nil {
		rawResp := raw.(*oss.Response)
		defer rawResp.Body.Close()
		rawData, err := ioutil.ReadAll(rawResp.Body)
		if err != nil {
			return WrapError(err)
		}
		policy = string(rawData)
	}

	if err := d.Set("policy", policy); err != nil {
		return WrapError(err)
	}

	// Read tags
	raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketTagging(d.Id())
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "GetBucketTagging", ApsaraStackOssGoSdk)
	}
	addDebug("GetBucketTagging", raw, requestInfo, request)
	tagging, _ := raw.(oss.GetBucketTaggingResult)
	tagsMap := make(map[string]string)
	if len(tagging.Tags) > 0 {
		for _, t := range tagging.Tags {
			tagsMap[t.Key] = t.Value
		}
	}
	if err := d.Set("tags", tagsMap); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceApsaraStackOssBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	d.Partial(true)

	if d.HasChange("acl") {
		request := map[string]string{"bucketName": d.Id(), "bucketACL": d.Get("acl").(string)}
		var requestInfo *oss.Client
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			requestInfo = ossClient
			return nil, ossClient.SetBucketACL(d.Id(), oss.ACLType(d.Get("acl").(string)))
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketACL", ApsaraStackOssGoSdk)
		}
		addDebug("SetBucketACL", raw, requestInfo, request)
		d.SetPartial("acl")
	}

	if d.HasChange("cors_rule") {
		if err := resourceApsaraStackOssBucketCorsUpdate(client, d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("cors_rule")
	}

	if d.HasChange("website") {
		if err := resourceApsaraStackOssBucketWebsiteUpdate(client, d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("website")
	}

	if d.HasChange("logging") {
		if err := resourceApsaraStackOssBucketLoggingUpdate(client, d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("logging")
	}

	if d.HasChange("referer_config") {
		if err := resourceApsaraStackOssBucketRefererUpdate(client, d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("referer_config")
	}

	if d.HasChange("lifecycle_rule") {
		if err := resourceApsaraStackOssBucketLifecycleRuleUpdate(client, d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("lifecycle_rule")
	}

	if d.HasChange("policy") {
		if err := resourceApsaraStackOssBucketPolicyUpdate(client, d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("policy")
	}

	if d.HasChange("server_side_encryption_rule") {
		if err := resourceApsaraStackOssBucketEncryptionUpdate(client, d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("server_side_encryption_rule")
	}

	if d.HasChange("tags") {
		if err := resourceApsaraStackOssBucketTaggingUpdate(client, d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	if d.HasChange("versioning") {
		if err := resourceApsaraStackOssBucketVersioningUpdate(client, d); err != nil {
			return WrapError(err)
		}
		d.SetPartial("versioning")
	}

	d.Partial(false)
	return resourceApsaraStackOssBucketRead(d, meta)
}

func resourceApsaraStackOssBucketCorsUpdate(client *connectivity.ApsaraStackClient, d *schema.ResourceData) error {
	cors := d.Get("cors_rule").([]interface{})
	var requestInfo *oss.Client
	if cors == nil || len(cors) == 0 {
		err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
				requestInfo = ossClient
				return nil, ossClient.DeleteBucketCORS(d.Id())
			})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			addDebug("DeleteBucketCORS", raw, requestInfo, map[string]string{"bucketName": d.Id()})
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteBucketCORS", ApsaraStackOssGoSdk)
		}
		return nil
	}
	// Put CORS
	rules := make([]oss.CORSRule, 0, len(cors))
	for _, c := range cors {
		corsMap := c.(map[string]interface{})
		rule := oss.CORSRule{}
		for k, v := range corsMap {
			log.Printf("[DEBUG] OSS bucket: %s, put CORS: %#v, %#v", d.Id(), k, v)
			if k == "max_age_seconds" {
				rule.MaxAgeSeconds = v.(int)
			} else {
				rMap := make([]string, len(v.([]interface{})))
				for i, vv := range v.([]interface{}) {
					rMap[i] = vv.(string)
				}
				switch k {
				case "allowed_headers":
					rule.AllowedHeader = rMap
				case "allowed_methods":
					rule.AllowedMethod = rMap
				case "allowed_origins":
					rule.AllowedOrigin = rMap
				case "expose_headers":
					rule.ExposeHeader = rMap
				}
			}
		}
		rules = append(rules, rule)
	}

	log.Printf("[DEBUG] Oss bucket: %s, put CORS: %#v", d.Id(), cors)
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return nil, ossClient.SetBucketCORS(d.Id(), rules)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketCORS", ApsaraStackOssGoSdk)
	}
	addDebug("SetBucketCORS", raw, requestInfo, map[string]interface{}{
		"bucketName": d.Id(),
		"corsRules":  rules,
	})
	return nil
}
func resourceApsaraStackOssBucketWebsiteUpdate(client *connectivity.ApsaraStackClient, d *schema.ResourceData) error {
	ws := d.Get("website").([]interface{})
	var requestInfo *oss.Client
	if ws == nil || len(ws) == 0 {
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			requestInfo = ossClient
			return nil, ossClient.DeleteBucketWebsite(d.Id())
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteBucketWebsite", ApsaraStackOssGoSdk)
		}
		addDebug("DeleteBucketWebsite", raw, requestInfo, map[string]string{"bucketName": d.Id()})
		return nil
	}

	var index_document, error_document string
	w := ws[0].(map[string]interface{})

	if v, ok := w["index_document"]; ok {
		index_document = v.(string)
	}
	if v, ok := w["error_document"]; ok {
		error_document = v.(string)
	}
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return nil, ossClient.SetBucketWebsite(d.Id(), index_document, error_document)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketWebsite", ApsaraStackOssGoSdk)
	}
	addDebug("SetBucketWebsite", raw, requestInfo, map[string]interface{}{
		"bucketName":    d.Id(),
		"indexDocument": index_document,
		"errorDocument": error_document,
	})
	return nil
}

func resourceApsaraStackOssBucketLoggingUpdate(client *connectivity.ApsaraStackClient, d *schema.ResourceData) error {
	logging := d.Get("logging").([]interface{})
	var requestInfo *oss.Client
	if logging == nil || len(logging) == 0 {
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			requestInfo = ossClient
			return nil, ossClient.DeleteBucketLogging(d.Id())
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteBucketLogging", ApsaraStackOssGoSdk)
		}
		addDebug("DeleteBucketLogging", raw, requestInfo, map[string]string{"bucketName": d.Id()})
		return nil
	}

	c := logging[0].(map[string]interface{})
	var target_bucket, target_prefix string
	if v, ok := c["target_bucket"]; ok {
		target_bucket = v.(string)
	}
	if v, ok := c["target_prefix"]; ok {
		target_prefix = v.(string)
	}
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return nil, ossClient.SetBucketLogging(d.Id(), target_bucket, target_prefix, target_bucket != "" || target_prefix != "")
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketLogging", ApsaraStackOssGoSdk)
	}
	addDebug("SetBucketLogging", raw, requestInfo, map[string]interface{}{
		"bucketName":   d.Id(),
		"targetBucket": target_bucket,
		"targetPrefix": target_prefix,
		"isEnable":     target_bucket != "",
	})
	return nil
}

func resourceApsaraStackOssBucketRefererUpdate(client *connectivity.ApsaraStackClient, d *schema.ResourceData) error {
	config := d.Get("referer_config").([]interface{})
	var requestInfo *oss.Client
	if config == nil || len(config) < 1 {
		log.Printf("[DEBUG] OSS set bucket referer as nil")
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			requestInfo = ossClient
			return nil, ossClient.SetBucketReferer(d.Id(), nil, true)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketReferer", ApsaraStackOssGoSdk)
		}
		addDebug("SetBucketReferer", raw, requestInfo, map[string]interface{}{
			"allowEmptyReferer": true,
			"bucketName":        d.Id(),
		})
		return nil
	}

	c := config[0].(map[string]interface{})

	var allow bool
	var referers []string
	if v, ok := c["allow_empty"]; ok {
		allow = v.(bool)
	}
	if v, ok := c["referers"]; ok {
		for _, referer := range v.([]interface{}) {
			referers = append(referers, referer.(string))
		}
	}
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return nil, ossClient.SetBucketReferer(d.Id(), referers, allow)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketReferer", ApsaraStackOssGoSdk)
	}
	addDebug("SetBucketReferer", raw, requestInfo, map[string]interface{}{
		"bucketName":        d.Id(),
		"referers":          referers,
		"allowEmptyReferer": allow,
	})
	return nil
}

func resourceApsaraStackOssBucketLifecycleRuleUpdate(client *connectivity.ApsaraStackClient, d *schema.ResourceData) error {
	bucket := d.Id()
	lifecycleRules := d.Get("lifecycle_rule").([]interface{})
	var requestInfo *oss.Client
	if lifecycleRules == nil || len(lifecycleRules) == 0 {
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			requestInfo = ossClient
			return nil, ossClient.DeleteBucketLifecycle(bucket)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteBucketLifecycle", ApsaraStackOssGoSdk)

		}
		addDebug("DeleteBucketLifecycle", raw, requestInfo, map[string]interface{}{
			"bucketName": bucket,
		})
		return nil
	}

	rules := make([]oss.LifecycleRule, 0, len(lifecycleRules))

	for i, lifecycleRule := range lifecycleRules {
		r := lifecycleRule.(map[string]interface{})

		rule := oss.LifecycleRule{
			Prefix: r["prefix"].(string),
		}

		// ID
		if val, ok := r["id"].(string); ok && val != "" {
			rule.ID = val
		}

		// Enabled
		if val, ok := r["enabled"].(bool); ok && val {
			rule.Status = string(ExpirationStatusEnabled)
		} else {
			rule.Status = string(ExpirationStatusDisabled)
		}

		// Expiration
		expiration := d.Get(fmt.Sprintf("lifecycle_rule.%d.expiration", i)).(*schema.Set).List()
		if len(expiration) > 0 {
			e := expiration[0].(map[string]interface{})
			i := oss.LifecycleExpiration{}
			valDate, _ := e["date"].(string)
			valDays, _ := e["days"].(int)

			if (valDate != "" && valDays > 0) || (valDate == "" && valDays <= 0) {
				return WrapError(Error("'date' conflicts with 'days'. One and only one of them can be specified in one expiration configuration."))
			}

			if valDate != "" {
				i.Date = fmt.Sprintf("%sT00:00:00.000Z", valDate)
			}
			if valDays > 0 {
				i.Days = valDays
			}
			rule.Expiration = &i
		}

		//Transitions
		transitions := d.Get(fmt.Sprintf("lifecycle_rule.%d.transitions", i)).(*schema.Set).List()
		if len(transitions) > 0 {
			for _, transition := range transitions {
				i := oss.LifecycleTransition{}

				valCreatedBeforeDate := transition.(map[string]interface{})["created_before_date"].(string)
				valDays := transition.(map[string]interface{})["days"].(int)
				valStorageClass := transition.(map[string]interface{})["storage_class"].(string)

				if (valCreatedBeforeDate != "" && valDays > 0) || (valCreatedBeforeDate == "" && valDays <= 0) || (valStorageClass == "") {
					return WrapError(Error("'CreatedBeforeDate' conflicts with 'Days'. One and only one of them can be specified in one transition configuration. 'storage_class' must be set."))
				}

				if valCreatedBeforeDate != "" {
					i.CreatedBeforeDate = fmt.Sprintf("%sT00:00:00.000Z", valCreatedBeforeDate)
				}
				if valDays > 0 {
					i.Days = valDays
				}

				if valStorageClass != "" {
					i.StorageClass = oss.StorageClassType(valStorageClass)
				}
				rule.Transitions = append(rule.Transitions, i)
			}
		}

		rules = append(rules, rule)
	}

	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return nil, ossClient.SetBucketLifecycle(bucket, rules)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketLifecycle", ApsaraStackOssGoSdk)
	}
	addDebug("SetBucketLifecycle", raw, requestInfo, map[string]interface{}{
		"bucketName": bucket,
		"rules":      rules,
	})
	return nil
}

func resourceApsaraStackOssBucketPolicyUpdate(client *connectivity.ApsaraStackClient, d *schema.ResourceData) error {
	bucket := d.Id()
	policy := d.Get("policy").(string)
	var requestInfo *oss.Client
	if len(policy) == 0 {
		params := map[string]interface{}{}
		params["policy"] = nil
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			requestInfo = ossClient
			return ossClient.Conn.Do("DELETE", bucket, "", params, nil, nil, 0, nil)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeletePolicyByConn", ApsaraStackOssGoSdk)
		}
		addDebug("DeletePolicyByConn", raw, requestInfo, params)
		return nil
	}
	params := map[string]interface{}{}
	params["policy"] = nil
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		buffer := new(bytes.Buffer)
		buffer.Write([]byte(policy))
		return ossClient.Conn.Do("PUT", bucket, "", params, nil, buffer, 0, nil)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "PutPolicyByConn", ApsaraStackOssGoSdk)
	}
	addDebug("PutPolicyByConn", raw, requestInfo, params)
	return nil
}

func resourceApsaraStackOssBucketEncryptionUpdate(client *connectivity.ApsaraStackClient, d *schema.ResourceData) error {
	encryption_rule := d.Get("server_side_encryption_rule").([]interface{})
	var requestInfo *oss.Client
	if encryption_rule == nil || len(encryption_rule) == 0 {
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			requestInfo = ossClient
			return nil, ossClient.DeleteBucketEncryption(d.Id())
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteBucketEncryption", ApsaraStackOssGoSdk)
		}
		addDebug("DeleteBucketEncryption", raw, requestInfo, map[string]string{"bucketName": d.Id()})
		return nil
	}

	var sseRule oss.ServerEncryptionRule
	c := encryption_rule[0].(map[string]interface{})
	if v, ok := c["sse_algorithm"]; ok {
		sseRule.SSEDefault.SSEAlgorithm = v.(string)
	}

	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return nil, ossClient.SetBucketEncryption(d.Id(), sseRule)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketEncryption", ApsaraStackOssGoSdk)
	}
	addDebug("SetBucketEncryption", raw, requestInfo, map[string]interface{}{
		"bucketName":     d.Id(),
		"encryptionRule": sseRule,
	})
	return nil
}

func resourceApsaraStackOssBucketTaggingUpdate(client *connectivity.ApsaraStackClient, d *schema.ResourceData) error {
	tagsMap := d.Get("tags").(map[string]interface{})
	var requestInfo *oss.Client
	if tagsMap == nil || len(tagsMap) == 0 {
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			requestInfo = ossClient
			return nil, ossClient.DeleteBucketTagging(d.Id())
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteBucketTagging", ApsaraStackOssGoSdk)
		}
		addDebug("DeleteBucketTagging", raw, requestInfo, map[string]string{"bucketName": d.Id()})
		return nil
	}

	// Put tagging
	var bTagging oss.Tagging
	for k, v := range tagsMap {
		bTagging.Tags = append(bTagging.Tags, oss.Tag{
			Key:   k,
			Value: v.(string),
		})
	}
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return nil, ossClient.SetBucketTagging(d.Id(), bTagging)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketTagging", ApsaraStackOssGoSdk)
	}
	addDebug("SetBucketTagging", raw, requestInfo, map[string]interface{}{
		"bucketName": d.Id(),
		"tagging":    bTagging,
	})
	return nil
}

func resourceApsaraStackOssBucketVersioningUpdate(client *connectivity.ApsaraStackClient, d *schema.ResourceData) error {
	versioning := d.Get("versioning").([]interface{})
	if len(versioning) == 1 {
		var status string
		c := versioning[0].(map[string]interface{})
		if v, ok := c["status"]; ok {
			status = v.(string)
		}

		versioningCfg := oss.VersioningConfig{}
		versioningCfg.Status = status
		var requestInfo *oss.Client
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			requestInfo = ossClient
			return nil, ossClient.SetBucketVersioning(d.Id(), versioningCfg)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "SetBucketVersioning", ApsaraStackOssGoSdk)
		}
		addDebug("SetBucketVersioning", raw, requestInfo, map[string]interface{}{
			"bucketName":       d.Id(),
			"versioningConfig": versioningCfg,
		})
	}

	return nil
}

func resourceApsaraStackOssBucketDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ossService := OssService{client}
	var requestInfo *oss.Client
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.IsBucketExist(d.Id())
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsBucketExist", ApsaraStackOssGoSdk)
	}
	addDebug("IsBucketExist", raw, requestInfo, map[string]string{"bucketName": d.Id()})

	exist, _ := raw.(bool)
	if !exist {
		return nil
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return nil, ossClient.DeleteBucket(d.Id())
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"BucketNotEmpty"}) {
				if d.Get("force_destroy").(bool) {
					raw, er := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
						bucket, _ := ossClient.Bucket(d.Get("bucket").(string))
						lor, err := bucket.ListObjectVersions()
						if err != nil {
							return nil, WrapErrorf(err, DefaultErrorMsg, d.Id(), "ListObjectVersions", ApsaraStackOssGoSdk)
						}
						addDebug("ListObjectVersions", lor, requestInfo)
						objectsToDelete := make([]oss.DeleteObject, 0)
						for _, object := range lor.ObjectDeleteMarkers {
							objectsToDelete = append(objectsToDelete, oss.DeleteObject{
								Key:       object.Key,
								VersionId: object.VersionId,
							})
						}

						for _, object := range lor.ObjectVersions {
							objectsToDelete = append(objectsToDelete, oss.DeleteObject{
								Key:       object.Key,
								VersionId: object.VersionId,
							})
						}
						return bucket.DeleteObjectVersions(objectsToDelete)
					})
					if er != nil {
						return resource.NonRetryableError(er)
					}
					addDebug("DeleteObjectVersions", raw, requestInfo, map[string]string{"bucketName": d.Id()})
					return resource.RetryableError(err)
				}
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteBucket", raw, requestInfo, map[string]string{"bucketName": d.Id()})
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteBucket", ApsaraStackOssGoSdk)
	}
	return WrapError(ossService.WaitForOssBucket(d.Id(), Deleted, DefaultTimeoutMedium))
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
