package connectivity

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	cdn_new "github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/location"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/denverdino/aliyungo/cdn"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ApsaraStackClient struct {
	Region    Region
	RegionId  string
	AccessKey string
	SecretKey string
	config    *Config
	ecsconn   *ecs.Client
	vpcconn   *vpc.Client
	slbconn   *slb.Client

	cdnconn        *cdn.CdnClient
	cdnconn_new    *cdn_new.Client
	kmsconn        *kms.Client
	bssopenapiconn *bssopenapi.Client
	ramconn        *ram.Client
}

const (
	ApiVersion20140526 = ApiVersion("2014-05-26")
	ApiVersion20160815 = ApiVersion("2016-08-15")
	ApiVersion20140515 = ApiVersion("2014-05-15")
)

const DefaultClientRetryCountSmall = 5

const Terraform = "HashiCorp-Terraform"

const Provider = "Terraform-Provider"

const Module = "Terraform-Module"

type ApiVersion string

// The main version number that is being run at the moment.
var providerVersion = "1.94.0"
var terraformVersion = strings.TrimSuffix(schema.Provider{}.TerraformVersion, "-dev")

// Client for ApsaraStackClient
func (c *Config) Client() (*ApsaraStackClient, error) {
	// Get the auth and region. This can fail if keys/regions were not
	// specified and we're attempting to use the environment.
	if !c.SkipRegionValidation {
		err := c.loadAndValidate()
		if err != nil {
			return nil, err
		}
	}

	return &ApsaraStackClient{
		config:    c,
		Region:    c.Region,
		RegionId:  c.RegionId,
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
	}, nil
}

func (client *ApsaraStackClient) WithEcsClient(do func(*ecs.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ECS client if necessary
	if client.ecsconn == nil {
		endpoint := client.config.EcsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, ECSCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ECSCode), endpoint)
		}
		ecsconn, err := ecs.NewClientWithOptions(client.config.RegionId, client.getSdkConfig().WithTimeout(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ECS client: %#v", err)
		}

		//if _, err := ecsconn.DescribeRegions(ecs.CreateDescribeRegionsRequest()); err != nil {
		//	return nil, err
		//}
		ecsconn.AppendUserAgent(Terraform, terraformVersion)
		ecsconn.AppendUserAgent(Provider, providerVersion)
		ecsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		ecsconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			ecsconn.SetHttpsProxy(client.config.Proxy)
		}
		client.ecsconn = ecsconn
	}

	return do(client.ecsconn)
}

func (client *ApsaraStackClient) WithVpcClient(do func(*vpc.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the VPC client if necessary
	if client.vpcconn == nil {
		endpoint := client.config.VpcEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, VPCCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(VPCCode), endpoint)
		}
		vpcconn, err := vpc.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the VPC client: %#v", err)
		}

		vpcconn.AppendUserAgent(Terraform, terraformVersion)
		vpcconn.AppendUserAgent(Provider, providerVersion)
		vpcconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		vpcconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			vpcconn.SetHttpsProxy(client.config.Proxy)
		}
		client.vpcconn = vpcconn
	}

	return do(client.vpcconn)
}
func (client *ApsaraStackClient) WithSlbClient(do func(*slb.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the SLB client if necessary
	if client.slbconn == nil {
		endpoint := client.config.SlbEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, SLBCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(SLBCode), endpoint)
		}
		slbconn, err := slb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the SLB client: %#v", err)
		}

		slbconn.AppendUserAgent(Terraform, terraformVersion)
		slbconn.AppendUserAgent(Provider, providerVersion)
		slbconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.slbconn = slbconn
	}

	return do(client.slbconn)
}

func (client *ApsaraStackClient) describeEndpointForService(serviceCode string) (*location.Endpoint, error) {
	args := location.CreateDescribeEndpointsRequest()
	args.ServiceCode = serviceCode
	args.Id = client.config.RegionId
	args.Domain = client.config.LocationEndpoint
	if args.Domain == "" {
		args.Domain = loadEndpoint(client.RegionId, LOCATIONCode)
	}
	if args.Domain == "" {
		args.Domain = "location-readonly.aliyuncs.com"
	}

	locationClient, err := location.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("Unable to initialize the location client: %#v", err)

	}
	locationClient.AppendUserAgent(Terraform, terraformVersion)
	locationClient.AppendUserAgent(Provider, providerVersion)
	locationClient.AppendUserAgent(Module, client.config.ConfigurationSource)
	endpointsResponse, err := locationClient.DescribeEndpoints(args)
	if err != nil {
		return nil, fmt.Errorf("Describe %s endpoint using region: %#v got an error: %#v.", serviceCode, client.RegionId, err)
	}
	if endpointsResponse != nil && len(endpointsResponse.Endpoints.Endpoint) > 0 {
		for _, e := range endpointsResponse.Endpoints.Endpoint {
			if e.Type == "openAPI" {
				return &e, nil
			}
		}
	}
	return nil, fmt.Errorf("There is no any available endpoint for %s in region %s.", serviceCode, client.RegionId)
}

func (client *ApsaraStackClient) NewCommonRequest(product, serviceCode, schema string, apiVersion ApiVersion) (*requests.CommonRequest, error) {
	request := requests.NewCommonRequest()
	endpoint := loadEndpoint(client.RegionId, ServiceCode(strings.ToUpper(product)))
	if endpoint == "" {
		endpointItem, err := client.describeEndpointForService(serviceCode)
		if err != nil {
			return nil, fmt.Errorf("describeEndpointForService got an error: %#v.", err)
		}
		if endpointItem != nil {
			endpoint = endpointItem.Endpoint
		}
	}
	// Use product code to find product domain
	if endpoint != "" {
		request.Domain = endpoint
	} else {
		// When getting endpoint failed by location, using custom endpoint instead
		request.Domain = fmt.Sprintf("%s.%s.aliyuncs.com", strings.ToLower(serviceCode), client.RegionId)
	}
	request.Version = string(apiVersion)
	request.RegionId = client.RegionId
	request.Product = product
	request.Scheme = schema
	request.AppendUserAgent(Terraform, terraformVersion)
	request.AppendUserAgent(Provider, providerVersion)
	request.AppendUserAgent(Module, client.config.ConfigurationSource)
	return request, nil
}

func (client *ApsaraStackClient) getSdkConfig() *sdk.Config {
	return sdk.NewConfig().
		WithMaxRetryTime(DefaultClientRetryCountSmall).
		WithTimeout(time.Duration(30) * time.Second).
		WithEnableAsync(true).
		WithGoRoutinePoolSize(100).
		WithMaxTaskQueueSize(10000).
		WithDebug(false).
		WithHttpTransport(client.getTransport()).
		WithScheme(client.config.Protocol)
}

func (client *ApsaraStackClient) getTransport() *http.Transport {
	handshakeTimeout, err := strconv.Atoi(os.Getenv("TLSHandshakeTimeout"))
	if err != nil {
		handshakeTimeout = 120
	}
	transport := &http.Transport{}
	transport.TLSHandshakeTimeout = time.Duration(handshakeTimeout) * time.Second

	return transport
}

func (client *ApsaraStackClient) getHttpProxy() (proxy *url.URL, err error) {
	if client.config.Protocol == "HTTPS" {
		if rawurl := os.Getenv("HTTPS_PROXY"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		} else if rawurl := os.Getenv("https_proxy"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		}
	} else {
		if rawurl := os.Getenv("HTTP_PROXY"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		} else if rawurl := os.Getenv("http_proxy"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		}
	}
	return proxy, err
}

func (client *ApsaraStackClient) skipProxy(endpoint string) (bool, error) {
	var urls []string
	if rawurl := os.Getenv("NO_PROXY"); rawurl != "" {
		urls = strings.Split(rawurl, ",")
	} else if rawurl := os.Getenv("no_proxy"); rawurl != "" {
		urls = strings.Split(rawurl, ",")
	}
	for _, value := range urls {
		if strings.HasPrefix(value, "*") {
			value = fmt.Sprintf(".%s", value)
		}
		noProxyReg, err := regexp.Compile(value)
		if err != nil {
			return false, err
		}
		if noProxyReg.MatchString(endpoint) {
			return true, nil
		}
	}
	return false, nil
}
func (client *ApsaraStackClient) WithKmsClient(do func(*kms.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the KMS client if necessary
	if client.kmsconn == nil {

		endpoint := client.config.KmsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, KMSCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(KMSCode), endpoint)
		}
		kmsconn, err := kms.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the kms client: %#v", err)
		}
		kmsconn.AppendUserAgent(Terraform, terraformVersion)
		kmsconn.AppendUserAgent(Provider, providerVersion)
		kmsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		kmsconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			kmsconn.SetHttpsProxy(client.config.Proxy)
		}
		client.kmsconn = kmsconn
	}
	return do(client.kmsconn)
}
func (client *ApsaraStackClient) WithBssopenapiClient(do func(*bssopenapi.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the bssopenapi client if necessary
	if client.bssopenapiconn == nil {
		endpoint := client.config.BssOpenApiEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, BSSOPENAPICode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(BSSOPENAPICode), endpoint)
		}

		bssopenapiconn, err := bssopenapi.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the BSSOPENAPI client: %#v", err)
		}
		bssopenapiconn.AppendUserAgent(Terraform, terraformVersion)
		bssopenapiconn.AppendUserAgent(Provider, providerVersion)
		bssopenapiconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		bssopenapiconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			bssopenapiconn.SetHttpsProxy(client.config.Proxy)
		}
		client.bssopenapiconn = bssopenapiconn
	}

	return do(client.bssopenapiconn)
}
func (client *ApsaraStackClient) WithRamClient(do func(*ram.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the RAM client if necessary
	if client.ramconn == nil {
		endpoint := client.config.RamEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, RAMCode)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "http://"))
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(RAMCode), endpoint)
		}

		ramconn, err := ram.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the RAM client: %#v", err)
		}
		ramconn.AppendUserAgent(Terraform, terraformVersion)
		ramconn.AppendUserAgent(Provider, providerVersion)
		ramconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		ramconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			ramconn.SetHttpsProxy(client.config.Proxy)
		}
		client.ramconn = ramconn
	}

	return do(client.ramconn)
}

func (client *ApsaraStackClient) WithCdnClient_new(do func(*cdn_new.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the CDN client if necessary
	if client.cdnconn_new == nil {
		endpoint := client.config.CdnEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CDNCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CDNCode), endpoint)
		}
		cdnconn, err := cdn_new.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CDN client: %#v", err)
		}

		cdnconn.AppendUserAgent(Terraform, terraformVersion)
		cdnconn.AppendUserAgent(Provider, providerVersion)
		cdnconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		cdnconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			cdnconn.SetHttpsProxy(client.config.Proxy)
		}
		client.cdnconn_new = cdnconn
	}

	return do(client.cdnconn_new)
}
