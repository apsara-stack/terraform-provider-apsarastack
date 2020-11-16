package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
)

// OssService *connectivity.ApsaraStackClient
type OssService struct {
	client *connectivity.ApsaraStackClient
}

func (s *OssService) DescribeOssBucket(id string) (response oss.GetBucketInfoResult, err error) {
	//request := map[string]string{"bucketName": id, "Department": s.client.Department, "ResourceGroup": s.client.ResourceGroup}
	var requestInfo *oss.Client

	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{

		"AccessKeySecret":  s.client.SecretKey,
		"Product":          "OneRouter",
		"Department":       s.client.Department,
		"ResourceGroup":    s.client.ResourceGroup,
		"RegionId":         s.client.RegionId,
		"Action":           "DoOpenApi",
		"AccountInfo":      "123456",
		"Version":          "2018-12-12",
		"SignatureVersion": "1.0",
		"OpenApiAction":    "GetService",
		"ProductName":      "oss",
	}
	request.Method = "POST"        // Set request method
	request.Product = "OneRouter"  // Specify product
	request.Version = "2018-12-12" // Specify product version
	request.ServiceCode = "OneRouter"
	request.Scheme = "http" // Set request scheme. Default: http
	request.ApiName = "DoOpenApi"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}

	var bucketList = &BucketList{}
	raw, err := s.client.WithOssNewClient(func(ossClient *ecs.Client) (interface{}, error) {

		return ossClient.ProcessCommonRequest(request)
	})

	if err != nil {
		if ossNotFoundError(err) {
			return response, WrapErrorf(err, NotFoundMsg, ApsaraStackOssGoSdk)
		}
		return response, WrapErrorf(err, DefaultErrorMsg, id, "GetBucketInfo", ApsaraStackOssGoSdk)
	}
	addDebug("GetBucketInfo", raw, requestInfo, request)
	bresponse, _ := raw.(*responses.CommonResponse)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), bucketList)
	if err != nil {
		return response, WrapError(err)
	}
	if bucketList.Code != "200" || len(bucketList.Data.ListAllMyBucketsResult.Buckets.Bucket) < 1 {
		return response, WrapError(err)
	}

	var found = false
	for _, j := range bucketList.Data.ListAllMyBucketsResult.Buckets.Bucket {
		if j.Name == id {
			response.BucketInfo.Name = j.Name
			response.BucketInfo.StorageClass = j.StorageClass
			response.BucketInfo.ExtranetEndpoint = j.ExtranetEndpoint
			response.BucketInfo.IntranetEndpoint = j.IntranetEndpoint
			response.BucketInfo.Owner.ID = fmt.Sprint(j.ResourceGroupName)
			//response.BucketInfo.CreationDate=fmt.Sprint(j.CreationDate.
			response.BucketInfo.Location = j.Location
			found = true
			break
		}
	}
	if !found {
		response.BucketInfo.Name = ""
	}
	return
}

type BucketList struct {
	Data struct {
		ListAllMyBucketsResult struct {
			Buckets struct {
				Bucket []struct {
					Comment           string `json:"Comment"`
					CreationDate      string `json:"CreationDate"`
					Department        int64  `json:"Department"`
					DepartmentName    string `json:"DepartmentName"`
					ExtranetEndpoint  string `json:"ExtranetEndpoint"`
					IntranetEndpoint  string `json:"IntranetEndpoint"`
					Location          string `json:"Location"`
					Name              string `json:"Name"`
					ResourceGroup     int64  `json:"ResourceGroup"`
					ResourceGroupName string `json:"ResourceGroupName"`
					StorageClass      string `json:"StorageClass"`
				} `json:"Bucket"`
			} `json:"Buckets"`
			Owner struct{} `json:"Owner"`
		} `json:"ListAllMyBucketsResult"`
	} `json:"Data"`
	Code         string `json:"code"`
	Cost         int64  `json:"cost"`
	Message      string `json:"message"`
	PureListData bool   `json:"pureListData"`
	Redirect     bool   `json:"redirect"`
	Success      bool   `json:"success"`
}

func (s *OssService) WaitForOssBucket(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeOssBucket(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}

		if object.BucketInfo.Name != "" && status != Deleted {
			return nil
		}
		if object.BucketInfo.Name == "" && status == Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.BucketInfo.Name, status, ProviderERROR)
		}
	}
}

func (s *OssService) WaitForOssBucketObject(bucket *oss.Bucket, id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		exist, err := bucket.IsObjectExist(id)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, id, "IsObjectExist", ApsaraStackOssGoSdk)
		}
		addDebug("IsObjectExist", exist)

		if !exist {
			return nil
		}

		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.FormatBool(exist), status, ProviderERROR)
		}
	}
}
