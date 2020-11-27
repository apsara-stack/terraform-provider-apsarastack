package connectivity

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	cdn_new "github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/location"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity/ascm"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/denverdino/aliyungo/cdn"

	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"sync"

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
	Region            Region
	RegionId          string
	AccessKey         string
	SecretKey         string
	Department        string
	ResourceGroup     string
	config            *Config
	accountId         string
	ascmconn          *ascm.Client
	ecsconn           *ecs.Client
	accountIdMutex    sync.RWMutex
	vpcconn           *vpc.Client
	slbconn           *slb.Client
	csconn            *cs.Client
	polarDBconn       *polardb.Client
	cdnconn           *cdn.CdnClient
	cdnconn_new       *cdn_new.Client
	kmsconn           *kms.Client
	bssopenapiconn    *bssopenapi.Client
	rdsconn           *rds.Client
	ramconn           *ram.Client
	essconn           *ess.Client
	gpdbconn          *gpdb.Client
	elasticsearchconn *elasticsearch.Client
	hbaseconn         *hbase.Client
	adbconn           *adb.Client
	ossconn           *oss.Client
	rkvconn           *r_kvstore.Client
	fcconn            *fc.Client
	ddsconn           *dds.Client
	onsconn           *ons.Client
	logconn           *sls.Client
	logpopconn        *slsPop.Client
	dnsconn           *alidns.Client
	creeconn          *cr_ee.Client
	crconn            *cr.Client
}

const (
	ApiVersion20140526 = ApiVersion("2014-05-26")
	ApiVersion20160815 = ApiVersion("2016-08-15")
	ApiVersion20140515 = ApiVersion("2014-05-15")
	ApiVersion20190510 = ApiVersion("2019-05-10")
)

const DefaultClientRetryCountSmall = 5

const Terraform = "HashiCorp-Terraform"

const Provider = "Terraform-Provider"

const Module = "Terraform-Module"

type ApiVersion string

// The main version number that is being run at the moment.
var ProviderVersion = "1.94.0"
var TerraformVersion = strings.TrimSuffix(schema.Provider{}.TerraformVersion, "-dev")
var goSdkMutex = sync.RWMutex{} // The Go SDK is not thread-safe

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
		config:        c,
		Region:        c.Region,
		RegionId:      c.RegionId,
		AccessKey:     c.AccessKey,
		SecretKey:     c.SecretKey,
		Department:    c.Department,
		ResourceGroup: c.ResourceGroup,
	}, nil
}
func (client *ApsaraStackClient) WithAscmClient(do func(*ascm.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ASCM client if necessary
	if client.ascmconn == nil {
		endpoint := client.config.AscmEndpoint
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ASCMCode), endpoint)
		}
		//ascmconn, err := ascm.NewClientWithOptions(client.config.RegionId, client.getSdkConfig().WithTimeout(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
		ascmconn, err := sdk.NewClientWithAccessKey(client.RegionId, client.AccessKey, client.SecretKey)

		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ASCM client AccessKey: %#v", err)
		}
		ascmconn.Domain = endpoint
		ascmconn.AppendUserAgent(Terraform, TerraformVersion)
		ascmconn.AppendUserAgent(Provider, ProviderVersion)
		ascmconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		ascmconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			ascmconn.SetHttpsProxy(client.config.Proxy)
			ascmconn.SetHttpProxy(client.config.Proxy)
		}

	}
	return do(client.ascmconn)
}

func (client *ApsaraStackClient) WithEcsClient(do func(*ecs.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ECS client if necessary
	if client.ecsconn == nil {
		endpoint := client.config.EcsEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the ecs client: endpoint or domain is not provided for ecs service")
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ECSCode), endpoint)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
		}
		ecsconn, err := ecs.NewClientWithOptions(client.config.RegionId, client.getSdkConfig().WithTimeout(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ECS client: %#v", err)
		}

		ecsconn.Domain = endpoint
		ecsconn.AppendUserAgent(Terraform, TerraformVersion)
		ecsconn.AppendUserAgent(Provider, ProviderVersion)
		ecsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		ecsconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			ecsconn.SetHttpsProxy(client.config.Proxy)
			ecsconn.SetHttpProxy(client.config.Proxy)
		}
		client.ecsconn = ecsconn
	}

	return do(client.ecsconn)
}

func (client *ApsaraStackClient) WithPolarDBClient(do func(*polardb.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the PolarDB client if necessary
	if client.polarDBconn == nil {
		endpoint := client.config.PolarDBEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the polardb client: endpoint or domain is not provided for polardb service")
		}
		polarDBconn, err := polardb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the PolarDB client: %#v", err)

		}
		polarDBconn.Domain = endpoint
		polarDBconn.AppendUserAgent(Terraform, TerraformVersion)
		polarDBconn.AppendUserAgent(Provider, ProviderVersion)
		polarDBconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		polarDBconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			polarDBconn.SetHttpProxy(client.config.Proxy)
			polarDBconn.SetHTTPSInsecure(client.config.Insecure)
		}

		client.polarDBconn = polarDBconn
	}

	return do(client.polarDBconn)
}
func (client *ApsaraStackClient) WithElasticsearchClient(do func(*elasticsearch.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the Elasticsearch client if necessary
	if client.elasticsearchconn == nil {
		endpoint := client.config.ElasticsearchEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the ElasticSearch client: endpoint or domain is not provided for ElasticSearch service")
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ELASTICSEARCHCode), endpoint)
		}
		elasticsearchconn, err := elasticsearch.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Elasticsearch client: %#v", err)
		}

		elasticsearchconn.AppendUserAgent(Terraform, TerraformVersion)
		elasticsearchconn.AppendUserAgent(Provider, ProviderVersion)
		elasticsearchconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		elasticsearchconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			elasticsearchconn.SetHttpProxy(client.config.Proxy)
		}
		client.elasticsearchconn = elasticsearchconn
	}

	return do(client.elasticsearchconn)
}
func (client *ApsaraStackClient) WithEssClient(do func(*ess.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ESS client if necessary
	if client.essconn == nil {
		endpoint := client.config.EssEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the ess client: endpoint or domain is not provided for ess service")
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ESSCode), endpoint)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
		}
		essconn, err := ess.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ESS client: %#v", err)
		}
		essconn.Domain = endpoint
		essconn.AppendUserAgent(Terraform, TerraformVersion)
		essconn.AppendUserAgent(Provider, ProviderVersion)
		essconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		essconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			essconn.SetHttpsProxy(client.config.Proxy)
			essconn.SetHttpProxy(client.config.Proxy)
		}
		client.essconn = essconn
	}

	return do(client.essconn)
}

func (client *ApsaraStackClient) WithRkvClient(do func(*r_kvstore.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the RKV client if necessary
	if client.rkvconn == nil {
		endpoint := client.config.KVStoreEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the kvstore client: endpoint or domain is not provided for logpop service")
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, fmt.Sprintf("R-%s", string(KVSTORECode)), endpoint)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
		}
		rkvconn, err := r_kvstore.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the RKV client: %#v", err)
		}
		rkvconn.Domain = endpoint
		rkvconn.AppendUserAgent(Terraform, TerraformVersion)
		rkvconn.AppendUserAgent(Provider, ProviderVersion)
		rkvconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		rkvconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			rkvconn.SetHttpProxy(client.config.Proxy)
		}
		client.rkvconn = rkvconn
	}

	return do(client.rkvconn)
}

func (client *ApsaraStackClient) WithGpdbClient(do func(*gpdb.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the GPDB client if necessary
	if client.gpdbconn == nil {
		endpoint := client.config.GpdbEndpoint
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(GPDBCode), endpoint)
		}
		gpdbconn, err := gpdb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the GPDB client: %#v", err)
		}

		gpdbconn.Domain = endpoint
		gpdbconn.AppendUserAgent(Terraform, TerraformVersion)
		gpdbconn.AppendUserAgent(Provider, ProviderVersion)

		gpdbconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		gpdbconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			gpdbconn.SetHttpProxy(client.config.Proxy)
		}
		client.gpdbconn = gpdbconn
	}

	return do(client.gpdbconn)
}
func (client *ApsaraStackClient) WithAdbClient(do func(*adb.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the adb client if necessary
	if client.adbconn == nil {
		endpoint := client.config.AdbEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the  client: endpoint or domain is not provided for  service")
		}
		adbconn, err := adb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the adb client: %#v", err)

		}
		adbconn.Domain = endpoint
		adbconn.AppendUserAgent(Terraform, TerraformVersion)
		adbconn.AppendUserAgent(Provider, ProviderVersion)
		adbconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		adbconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			adbconn.SetHttpProxy(client.config.Proxy)
		}
		client.adbconn = adbconn
	}

	return do(client.adbconn)
}
func (client *ApsaraStackClient) WithHbaseClient(do func(*hbase.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the HBase client if necessary
	if client.hbaseconn == nil {
		endpoint := client.config.HBaseEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the  client: endpoint or domain is not provided for  service")
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(HBASECode), endpoint)
		}
		hbaseconn, err := hbase.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the hbase client: %#v", err)
		}

		hbaseconn.AppendUserAgent(Terraform, TerraformVersion)
		hbaseconn.AppendUserAgent(Provider, ProviderVersion)
		hbaseconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		hbaseconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			hbaseconn.SetHttpProxy(client.config.Proxy)
		}
		client.hbaseconn = hbaseconn
	}

	return do(client.hbaseconn)
}
func (client *ApsaraStackClient) WithFcClient(do func(*fc.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the FC client if necessary
	if client.fcconn == nil {
		endpoint := client.config.FcEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the  client: endpoint or domain is not provided for  service")
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
		}
		accountId, err := client.AccountId()
		if err != nil {
			return nil, err
		}

		config := client.getSdkConfig()
		clientOptions := []fc.ClientOption{fc.WithSecurityToken(client.config.SecurityToken), fc.WithTransport(config.HttpTransport),
			fc.WithTimeout(30), fc.WithRetryCount(DefaultClientRetryCountSmall)}
		fcconn, err := fc.NewClient(fmt.Sprintf("https://%s.%s", accountId, endpoint), string(ApiVersion20160815), client.config.AccessKey, client.config.SecretKey, clientOptions...)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the FC client: %#v", err)
		}

		fcconn.Config.UserAgent = client.getUserAgent()
		fcconn.Config.SecurityToken = client.config.SecurityToken
		client.fcconn = fcconn
	}

	return do(client.fcconn)
}
func (client *ApsaraStackClient) WithVpcClient(do func(*vpc.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the VPC client if necessary
	if client.vpcconn == nil {
		endpoint := client.config.VpcEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the vpc client: endpoint or domain is not provided for vpc service")
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(VPCCode), endpoint)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
		}
		vpcconn, err := vpc.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the VPC client: %#v", err)
		}
		vpcconn.Domain = endpoint
		vpcconn.AppendUserAgent(Terraform, TerraformVersion)
		vpcconn.AppendUserAgent(Provider, ProviderVersion)
		vpcconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		vpcconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			vpcconn.SetHttpsProxy(client.config.Proxy)
			vpcconn.SetHttpProxy(client.config.Proxy)
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
			return nil, fmt.Errorf("unable to initialize the slb client: endpoint or domain is not provided for slb service")
		}

		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(SLBCode), endpoint)
		}
		slbconn, err := slb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the SLB client: %#v", err)
		}
		slbconn.Domain = endpoint
		slbconn.AppendUserAgent(Terraform, TerraformVersion)
		slbconn.AppendUserAgent(Provider, ProviderVersion)
		slbconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		slbconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			slbconn.SetHttpsProxy(client.config.Proxy)
			slbconn.SetHttpProxy(client.config.Proxy)
		}
		client.slbconn = slbconn
	}

	return do(client.slbconn)
}
func (client *ApsaraStackClient) WithDdsClient(do func(*dds.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the DDS client if necessary
	if client.ddsconn == nil {
		endpoint := client.config.DdsEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the  client: endpoint or domain is not provided for  service")
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(DDSCode), endpoint)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
		}
		ddsconn, err := dds.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DDS client: %#v", err)
		}
		ddsconn.Domain = endpoint

		ddsconn.AppendUserAgent(Terraform, TerraformVersion)
		ddsconn.AppendUserAgent(Provider, ProviderVersion)
		ddsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		ddsconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			ddsconn.SetHttpProxy(client.config.Proxy)
		}
		client.ddsconn = ddsconn
	}

	return do(client.ddsconn)
}

func (client *ApsaraStackClient) WithOssNewClient(do func(*ecs.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ECS client if necessary
	if client.ecsconn == nil {
		endpoint := client.config.OssEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the oss client: endpoint or domain is not provided for ecs service")
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ECSCode), endpoint)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
		}
		ecsconn, err := ecs.NewClientWithOptions(client.config.RegionId, client.getSdkConfig().WithTimeout(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ECS client: %#v", err)
		}

		ecsconn.Domain = endpoint
		ecsconn.AppendUserAgent(Terraform, TerraformVersion)
		ecsconn.AppendUserAgent(Provider, ProviderVersion)
		ecsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		ecsconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			ecsconn.SetHttpsProxy(client.config.Proxy)
			ecsconn.SetHttpProxy(client.config.Proxy)
		}
		client.ecsconn = ecsconn
	}

	return do(client.ecsconn)
}

func (client *ApsaraStackClient) describeEndpointForService(serviceCode string) (*location.Endpoint, error) {
	args := location.CreateDescribeEndpointsRequest()
	args.ServiceCode = serviceCode
	args.Id = client.config.RegionId
	args.Domain = client.config.LocationEndpoint

	if args.Domain == "" {
		args.Domain = "location-readonly.aliyuncs.com"
	}

	locationClient, err := location.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("Unable to initialize the location client: %#v", err)

	}
	locationClient.AppendUserAgent(Terraform, TerraformVersion)
	locationClient.AppendUserAgent(Provider, ProviderVersion)
	locationClient.AppendUserAgent(Module, client.config.ConfigurationSource)
	locationClient.SetHTTPSInsecure(client.config.Insecure)
	if client.config.Proxy != "" {
		locationClient.SetHttpsProxy(client.config.Proxy)
	}
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
	var endpoint string
	if strings.ToUpper(product) == "SLB" {
		endpoint = client.config.SlbEndpoint
	}
	if strings.ToUpper(product) == "ECS" {
		endpoint = client.config.EcsEndpoint
	}
	if strings.ToUpper(product) == "ASCM" {
		endpoint = client.config.AscmEndpoint
	}

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

	if strings.ToUpper(product) == "SLB" {
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup, "Version": string(apiVersion)}
	}
	if strings.ToUpper(product) == "ECS" {
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ecs", "Department": client.Department, "ResourceGroup": client.ResourceGroup, "Version": string(apiVersion)}
	}
	if strings.ToUpper(product) == "ASCM" {
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "ascm", "Department": client.Department, "ResourceGroup": client.ResourceGroup, "Version": string(apiVersion)}
	}

	request.AppendUserAgent(Terraform, TerraformVersion)
	request.AppendUserAgent(Provider, ProviderVersion)
	request.AppendUserAgent(Module, client.config.ConfigurationSource)
	request.SetHTTPSInsecure(client.config.Insecure)
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
		WithScheme("http")
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
func (client *ApsaraStackClient) AccountId() (string, error) {
	client.accountIdMutex.Lock()
	defer client.accountIdMutex.Unlock()

	if client.accountId == "" {
		log.Printf("[DEBUG] account_id not provided, attempting to retrieve it automatically...")
		identity, err := client.GetCallerIdentity()
		if err != nil {
			return "", err
		}
		if identity == "" {
			return "", fmt.Errorf("caller identity doesn't contain any AccountId")
		}
		client.accountId = identity
	}
	return client.accountId, nil
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
			return nil, fmt.Errorf("unable to initialize the kms client: endpoint or domain is not provided for KMS service")
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(KMSCode), endpoint)
		}
		kmsconn, err := kms.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the kms client: %#v", err)
		}
		kmsconn.AppendUserAgent(Terraform, TerraformVersion)
		kmsconn.Domain = endpoint
		kmsconn.AppendUserAgent(Provider, ProviderVersion)
		kmsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		kmsconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			kmsconn.SetHttpProxy(client.config.Proxy)
		}
		client.kmsconn = kmsconn
	}
	return do(client.kmsconn)
}
func (client *ApsaraStackClient) GetCallerIdentity() (string, error) {

	endpoint := client.config.AscmEndpoint
	if endpoint == "" {
		return "", fmt.Errorf("unable to initialize the ascm client: endpoint or domain is not provided for ascm service")
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, string(ASCMCode), endpoint)
	}
	ascmClient, err := sdk.NewClientWithAccessKey(client.config.RegionId, client.config.AccessKey, client.config.SecretKey)
	if err != nil {
		return "", fmt.Errorf("unable to initialize the ascm client: %#v", err)
	}

	ascmClient.AppendUserAgent(Terraform, TerraformVersion)
	ascmClient.AppendUserAgent(Provider, ProviderVersion)
	ascmClient.AppendUserAgent(Module, client.config.ConfigurationSource)
	ascmClient.SetHTTPSInsecure(client.config.Insecure)
	ascmClient.Domain = endpoint
	if client.config.Proxy != "" {
		ascmClient.SetHttpProxy(client.config.Proxy)
	}
	if client.config.Department == "" || client.config.ResourceGroup == "" {
		return "", fmt.Errorf("unable to initialize the ascm client: department or resource_group is not provided")
	}
	request := requests.NewCommonRequest()
	request.Method = "GET"         // Set request method
	request.Product = "ascm"       // Specify product
	request.Domain = endpoint      // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2019-05-10" // Specify product version
	request.Scheme = "http"        // Set request scheme. Default: http
	request.ApiName = "GetUserInfo"
	request.QueryParams = map[string]string{
		"AccessKeySecret":  client.config.SecretKey,
		"Product":          "ascm",
		"Department":       client.config.Department,
		"ResourceGroup":    client.config.ResourceGroup,
		"RegionId":         client.RegionId,
		"Action":           "GetAllNavigationInfo",
		"Version":          "2019-05-10",
		"SignatureVersion": "1.0",
	}
	resp := responses.BaseResponse{}
	request.TransToAcsRequest()
	err = ascmClient.DoAction(request, &resp)
	if err != nil {
		return "", err
	}
	response := &AccountId{}
	err = json.Unmarshal(resp.GetHttpContentBytes(), response)
	ownerId := response.Data.PrimaryKey

	if ownerId == "" {
		return "", fmt.Errorf("ownerId not found")
	}
	return ownerId, err
}

type AccountId struct {
	Data struct {
		PrimaryKey string `json:"primaryKey"`
	} `json:"data"`
}

func (client *ApsaraStackClient) WithBssopenapiClient(do func(*bssopenapi.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the bssopenapi client if necessary
	if client.bssopenapiconn == nil {
		endpoint := client.config.BssOpenApiEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the bss client: endpoint or domain is not provided for bss service")
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(BSSOPENAPICode), endpoint)
		}

		bssopenapiconn, err := bssopenapi.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the BSSOPENAPI client: %#v", err)
		}
		bssopenapiconn.AppendUserAgent(Terraform, TerraformVersion)
		bssopenapiconn.AppendUserAgent(Provider, ProviderVersion)
		bssopenapiconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		bssopenapiconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			bssopenapiconn.SetHttpsProxy(client.config.Proxy)
		}
		client.bssopenapiconn = bssopenapiconn
	}

	return do(client.bssopenapiconn)
}
func (client *ApsaraStackClient) WithOssClient(do func(*oss.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the OSS client if necessary
	if client.ossconn == nil {
		schma := "http"
		endpoint := client.config.OssEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the oss client: endpoint or domain is not provided for OSS service")
		}
		if endpoint == "" {
			endpointItem, _ := client.describeEndpointForService(strings.ToLower(string(OSSCode)))
			if endpointItem != nil {
				if len(endpointItem.Protocols.Protocols) > 0 {
					// HTTP or HTTPS
					schma = strings.ToLower(endpointItem.Protocols.Protocols[0])
					for _, p := range endpointItem.Protocols.Protocols {
						if strings.ToLower(p) == "http" {
							schma = strings.ToLower(p)
							break
						}
					}
				}
				endpoint = endpointItem.Endpoint
			}
		}
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("%s://%s", schma, endpoint)
		}

		clientOptions := []oss.ClientOption{oss.UserAgent(client.getUserAgent()),
			oss.SecurityToken(client.config.SecurityToken)}
		if client.config.Proxy != "" {
			clientOptions = append(clientOptions, oss.Proxy(client.config.Proxy))
		}

		clientOptions = append(clientOptions, oss.UseCname(false))

		ossconn, err := oss.New(endpoint, client.config.AccessKey, client.config.SecretKey, clientOptions...)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the OSS client: %#v", err)
		}

		client.ossconn = ossconn
	}

	return do(client.ossconn)
}

func (client *ApsaraStackClient) WithRamClient(do func(*ram.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the RAM client if necessary
	if client.ramconn == nil {
		endpoint := client.config.RamEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the ram client: endpoint or domain is not provided for ram operation")
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
		ramconn.AppendUserAgent(Terraform, TerraformVersion)
		ramconn.AppendUserAgent(Provider, ProviderVersion)
		ramconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		ramconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			ramconn.SetHttpsProxy(client.config.Proxy)
		}
		client.ramconn = ramconn
	}

	return do(client.ramconn)
}

func (client *ApsaraStackClient) WithRdsClient(do func(*rds.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the RDS client if necessary
	if client.rdsconn == nil {
		endpoint := client.config.RdsEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the rds client: endpoint or domain is not provided for RDS service")
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(RDSCode), endpoint)
		}
		rdsconn, err := rds.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the RDS client: %#v", err)
		}
		rdsconn.Domain = endpoint
		rdsconn.AppendUserAgent(Terraform, TerraformVersion)
		rdsconn.AppendUserAgent(Provider, ProviderVersion)
		rdsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		rdsconn.SetHTTPSInsecure(client.config.Insecure)

		if client.config.Proxy != "" {
			rdsconn.SetHttpProxy(client.config.Proxy)
		}

		client.rdsconn = rdsconn
	}

	return do(client.rdsconn)
}

func (client *ApsaraStackClient) WithCdnClient_new(do func(*cdn_new.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the CDN client if necessary
	if client.cdnconn_new == nil {
		endpoint := client.config.CdnEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the CDN client: endpoint or domain is not provided for CDN service")
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CDNCode), endpoint)
		}
		cdnconn, err := cdn_new.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CDN client: %#v", err)
		}

		cdnconn.AppendUserAgent(Terraform, TerraformVersion)
		cdnconn.AppendUserAgent(Provider, ProviderVersion)
		cdnconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		cdnconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			cdnconn.SetHttpsProxy(client.config.Proxy)
		}
		client.cdnconn_new = cdnconn
	}

	return do(client.cdnconn_new)
}
func (client *ApsaraStackClient) getUserAgent() string {
	return fmt.Sprintf("%s/%s %s/%s %s/%s", Terraform, TerraformVersion, Provider, ProviderVersion, Module, client.config.ConfigurationSource)
}
func (client *ApsaraStackClient) WithCsClient(do func(*cs.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CS client if necessary
	if client.csconn == nil {
		csconn := cs.NewClientForAussumeRole(client.config.AccessKey, client.config.SecretKey, client.config.SecurityToken)
		csconn.SetUserAgent(client.getUserAgent())
		endpoint := client.config.CsEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the cs client: endpoint or domain is not provided for cs service")
		}
		if endpoint != "" {
			if !strings.HasPrefix(endpoint, "http") {
				endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "://"))
			}
			csconn.SetEndpoint(endpoint)
		}
		if client.config.Proxy != "" {
			os.Setenv("http_proxy", client.config.Proxy)
		}
		client.csconn = csconn
	}

	return do(client.csconn)
}

func (client *ApsaraStackClient) getHttpProxyUrl() *url.URL {
	for _, v := range []string{"HTTPS_PROXY", "https_proxy", "HTTP_PROXY", "http_proxy"} {
		value := strings.Trim(os.Getenv(v), " ")
		if value != "" {
			if !regexp.MustCompile(`^http(s)?://`).MatchString(value) {
				value = fmt.Sprintf("https://%s", value)
			}
			proxyUrl, err := url.Parse(value)
			if err == nil {
				return proxyUrl
			}
			break
		}
	}
	return nil
}

func (client *ApsaraStackClient) WithOssBucketByName(bucketName string, do func(*oss.Bucket) (interface{}, error)) (interface{}, error) {
	return client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		bucket, err := client.ossconn.Bucket(bucketName)

		if err != nil {
			return nil, fmt.Errorf("unable to get the bucket %s: %#v", bucketName, err)
		}
		return do(bucket)
	})
}

func (client *ApsaraStackClient) WithOnsClient(do func(*ons.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ons client if necessary
	if client.onsconn == nil {
		endpoint := client.config.OnsEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the ons client: endpoint or domain is not provided for ons service")
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ONSCode), endpoint)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
		}
		onsconn, err := ons.NewClientWithAccessKey(client.RegionId, client.AccessKey, client.SecretKey)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ONS client: %#v", err)
		}

		onsconn.AppendUserAgent(Terraform, TerraformVersion)
		onsconn.AppendUserAgent(Provider, ProviderVersion)
		onsconn.Domain = endpoint

		onsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		onsconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			onsconn.SetHttpProxy(client.config.Proxy)
		}
		client.onsconn = onsconn
	}

	return do(client.onsconn)
}

func (client *ApsaraStackClient) WithLogClient(do func(*sls.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the LOG client if necessary
	if client.logconn == nil {
		endpoint := client.config.LogEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the log client: endpoint or domain is not provided for log service")
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
		}
		if client.config.Proxy != "" {
			os.Setenv("http_proxy", client.config.Proxy)
		}
		client.logconn = &sls.Client{
			AccessKeyID:     client.config.AccessKey,
			AccessKeySecret: client.config.SecretKey,
			Endpoint:        endpoint,
			SecurityToken:   client.config.SecurityToken,
			UserAgent:       client.getUserAgent(),
		}
	}

	return do(client.logconn)
}
func (client *ApsaraStackClient) WithLogPopClient(do func(*slsPop.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the HBase client if necessary
	if client.logpopconn == nil {
		endpoint := client.config.LogEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the lopgpop client: endpoint or domain is not provided for logpop service")
		}
		if endpoint != "" {
			endpoint = fmt.Sprintf("%s."+endpoint, client.config.RegionId)
		}
		logpopconn, err := slsPop.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))

		if err != nil {
			return nil, fmt.Errorf("unable to initialize the sls client: %#v", err)
		}

		logpopconn.AppendUserAgent(Terraform, TerraformVersion)
		logpopconn.AppendUserAgent(Provider, ProviderVersion)
		logpopconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.logpopconn = logpopconn
	}

	return do(client.logpopconn)
}

func (client *ApsaraStackClient) WithCrEEClient(do func(*cr_ee.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the CR EE client if necessary
	if client.creeconn == nil {
		endpoint := client.config.CrEndpoint
		if endpoint == "" {
			return nil, fmt.Errorf("unable to initialize the CRee client: endpoint or domain is not provided for CR service")
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CRCode), endpoint)
		}
		creeconn, err := cr_ee.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CR EE client: %#v", err)
		}
		creeconn.AppendUserAgent(Terraform, TerraformVersion)
		creeconn.AppendUserAgent(Provider, ProviderVersion)
		creeconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		if client.config.Proxy != "" {
			creeconn.SetHttpProxy(client.config.Proxy)
		}
		client.creeconn = creeconn
	}

	return do(client.creeconn)
}

func (client *ApsaraStackClient) WithCrClient(do func(*cr.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the CR client if necessary
	if client.crconn == nil {
		endpoint := client.config.CrEndpoint

		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CRCode), endpoint)
		}

		if strings.HasPrefix(endpoint, "http") {
			endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
		}
		crconn, err := cr.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CR client: %#v", err)
		}
		crconn.Domain = endpoint
		if client.config.Proxy != "" {
			crconn.SetHttpProxy(client.config.Proxy)
		}
		crconn.AppendUserAgent(Terraform, TerraformVersion)
		crconn.AppendUserAgent(Provider, ProviderVersion)
		crconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.crconn = crconn
	}

	return do(client.crconn)
}
func (client *ApsaraStackClient) WithDnsClient(do func(*alidns.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the DNS client if necessary
	if client.dnsconn == nil {
		endpoint := client.config.DnsEndpoint
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(DNSCode), endpoint)
		}

		dnsconn, err := alidns.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DNS client: %#v", err)
		}
		dnsconn.AppendUserAgent(Terraform, TerraformVersion)
		dnsconn.AppendUserAgent(Provider, ProviderVersion)
		dnsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		dnsconn.Domain = endpoint
		if client.config.Proxy != "" {
			dnsconn.SetHttpProxy(client.config.Proxy)
		}
		client.dnsconn = dnsconn
	}

	return do(client.dnsconn)
}
