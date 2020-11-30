package apsarastack

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// functions for a given region
func sharedClientForRegion(region string) (interface{}, error) {
	var accessKey, secretKey, proxy, domain, ossEndpoint, essEndpoint, slbEndpoint, crEndpoint, vpcEndpoint, rdsEndpoint, ecsEndpoint, rgsName string
	var insecure bool
	if accessKey = os.Getenv("APSARASTACK_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty APSARASTACK_ACCESS_KEY")
	}

	if secretKey = os.Getenv("APSARASTACK_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty APSARASTACK_SECRET_KEY")
	}
	insecure, _ = strconv.ParseBool(os.Getenv("APSARASTACK_INSECURE"))

	if proxy = os.Getenv("APSARASTACK_PROXY"); proxy == "" {
		return nil, fmt.Errorf("empty APSARASTACK_PROXY")
	}
	if domain = os.Getenv("APSARASTACK_DOMAIN"); domain == "" {
		return nil, fmt.Errorf("empty APSARASTACK_DOMAIN")
	}
	if ossEndpoint = os.Getenv("OSS_ENDPOINT"); ossEndpoint == "" {
		//return nil, fmt.Errorf("empty OSS_ENDPOINT")
	}
	if rdsEndpoint := os.Getenv("RDS_ENDPOINT"); rdsEndpoint == "" {
		//eturn nil, fmt.Errorf("empty RDS_ENDPOINT")
	}
	if essEndpoint = os.Getenv("ESS_ENDPOINT"); essEndpoint == "" {
		//return nil, fmt.Errorf("empty OSS_ENDPOINT")
	}
	if ecsEndpoint = os.Getenv("ECS_ENDPOINT"); ecsEndpoint == "" {
		//return nil, fmt.Errorf("empty OSS_ENDPOINT")
	}
	if vpcEndpoint = os.Getenv("VPC_ENDPOINT"); vpcEndpoint == "" {
		//return nil, fmt.Errorf("empty OSS_ENDPOINT")
	}
	if slbEndpoint = os.Getenv("SLB_ENDPOINT"); slbEndpoint == "" {
		//return nil, fmt.Errorf("empty OSS_ENDPOINT")
	}
	if crEndpoint = os.Getenv("CR_ENDPOINT"); crEndpoint == "" {
		//return nil, fmt.Errorf("empty OSS_ENDPOINT")
	}
	if domain = os.Getenv("APSARASTACK_DOMAIN"); domain == "" {
		//return nil, fmt.Errorf("empty APSARASTACK_DOMAIN")
	}
	if rgsName = os.Getenv("APSARASTACK_RESOURCE_GROUP_SET"); rgsName == "" {
		return nil, fmt.Errorf("empty APSARASTACK_RESOURCE_GROUP_SET")
	}

	conf := connectivity.Config{
		Region:          connectivity.Region(region),
		RegionId:        region,
		AccessKey:       accessKey,
		SecretKey:       secretKey,
		Proxy:           proxy,
		Insecure:        insecure,
		Domain:          domain,
		Protocol:        "HTTP",
		OssEndpoint:     ossEndpoint,
		EssEndpoint:     essEndpoint,
		RdsEndpoint:     rdsEndpoint,
		EcsEndpoint:     ecsEndpoint,
		VpcEndpoint:     vpcEndpoint,
		CrEndpoint:      crEndpoint,
		SlbEndpoint:     slbEndpoint,
		ResourceSetName: rgsName,
	}
	if accountId := os.Getenv("APSARASTACK_ACCOUNT_ID"); accountId != "" {
		conf.AccountId = accountId
	}

	// configures a default client for the region, using the above env vars
	client, err := conf.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}
