package apsarastack

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// functions for a given region
func sharedClientForRegion(region string) (interface{}, error) {
	var accessKey, secretKey string
	if accessKey = os.Getenv("APSARASTACK_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty APSARASTACK_ACCESS_KEY")
	}

	if secretKey = os.Getenv("APSARASTACK_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty APSARASTACK_SECRET_KEY")
	}

	conf := connectivity.Config{
		Region:    connectivity.Region(region),
		RegionId:  region,
		AccessKey: accessKey,
		SecretKey: secretKey,
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
