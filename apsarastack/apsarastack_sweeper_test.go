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
	var accessKey, secretKey, proxy, domain string
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
	if ossEndpoint := os.Getenv("OSS_ENDPOINT"); ossEndpoint == "" {
		//return nil, fmt.Errorf("empty OSS_ENDPOINT")
	}
	if rdsEndpoint := os.Getenv("RDS_ENDPOINT"); rdsEndpoint == "" {
		//return nil, fmt.Errorf("empty OSS_ENDPOINT")
	}
	if domain = os.Getenv("APSARASTACK_DOMAIN"); domain == "" {
		//return nil, fmt.Errorf("empty APSARASTACK_DOMAIN")
	}

	conf := connectivity.Config{
		Region:    connectivity.Region(region),
		RegionId:  region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Proxy:     proxy,
		Insecure:  insecure,
		Domain:    domain,
		Protocol:  "HTTPS",
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
