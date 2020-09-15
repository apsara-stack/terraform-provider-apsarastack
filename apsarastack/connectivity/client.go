package connectivity

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	cdn_new "github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/location"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/fc-go-sdk"
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
	config            *Config
	accountId         string
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

func (client *ApsaraStackClient) WithPolarDBClient(do func(*polardb.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the PolarDB client if necessary
	if client.polarDBconn == nil {
		endpoint := client.config.PolarDBEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, POLARDBCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.polardb.aliyuncs.com", client.config.RegionId)
			}
		}

		polarDBconn, err := polardb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the PolarDB client: %#v", err)

		}

		polarDBconn.AppendUserAgent(Terraform, terraformVersion)
		polarDBconn.AppendUserAgent(Provider, providerVersion)
		polarDBconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		polarDBconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			polarDBconn.SetHttpsProxy(client.config.Proxy)
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
			endpoint = loadEndpoint(client.config.RegionId, ELASTICSEARCHCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ELASTICSEARCHCode), endpoint)
		}
		elasticsearchconn, err := elasticsearch.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Elasticsearch client: %#v", err)
		}

		elasticsearchconn.AppendUserAgent(Terraform, terraformVersion)
		elasticsearchconn.AppendUserAgent(Provider, providerVersion)
		elasticsearchconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		elasticsearchconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			elasticsearchconn.SetHttpsProxy(client.config.Proxy)
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
			endpoint = loadEndpoint(client.config.RegionId, ESSCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ESSCode), endpoint)
		}
		essconn, err := ess.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ESS client: %#v", err)
		}

		essconn.AppendUserAgent(Terraform, terraformVersion)
		essconn.AppendUserAgent(Provider, providerVersion)
		essconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		essconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			essconn.SetHttpsProxy(client.config.Proxy)
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
			endpoint = loadEndpoint(client.config.RegionId, KVSTORECode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, fmt.Sprintf("R-%s", string(KVSTORECode)), endpoint)
		}
		rkvconn, err := r_kvstore.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the RKV client: %#v", err)
		}

		rkvconn.AppendUserAgent(Terraform, terraformVersion)
		rkvconn.AppendUserAgent(Provider, providerVersion)
		rkvconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		rkvconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			rkvconn.SetHttpsProxy(client.config.Proxy)
		}
		client.rkvconn = rkvconn
	}

	return do(client.rkvconn)
}

func (client *ApsaraStackClient) WithGpdbClient(do func(*gpdb.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the GPDB client if necessary
	if client.gpdbconn == nil {
		endpoint := client.config.GpdbEnpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, GPDBCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(GPDBCode), endpoint)
		}
		gpdbconn, err := gpdb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the GPDB client: %#v", err)
		}

		gpdbconn.AppendUserAgent(Terraform, terraformVersion)
		gpdbconn.AppendUserAgent(Provider, providerVersion)
		gpdbconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		gpdbconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			gpdbconn.SetHttpsProxy(client.config.Proxy)
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
			endpoint = loadEndpoint(client.config.RegionId, ADBCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.adb.aliyuncs.com", client.config.RegionId)
			}
		}

		adbconn, err := adb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the adb client: %#v", err)

		}

		adbconn.AppendUserAgent(Terraform, terraformVersion)
		adbconn.AppendUserAgent(Provider, providerVersion)
		adbconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		adbconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			adbconn.SetHttpsProxy(client.config.Proxy)
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
			endpoint = loadEndpoint(client.config.RegionId, HBASECode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(HBASECode), endpoint)
		}
		hbaseconn, err := hbase.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the hbase client: %#v", err)
		}

		hbaseconn.AppendUserAgent(Terraform, terraformVersion)
		hbaseconn.AppendUserAgent(Provider, providerVersion)
		hbaseconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		hbaseconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			hbaseconn.SetHttpsProxy(client.config.Proxy)
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
			endpoint = loadEndpoint(client.config.RegionId, FCCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.fc.aliyuncs.com", client.config.RegionId)
			}
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
		slbconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			slbconn.SetHttpsProxy(client.config.Proxy)
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
			endpoint = loadEndpoint(client.config.RegionId, DDSCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(DDSCode), endpoint)
		}
		ddsconn, err := dds.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DDS client: %#v", err)
		}

		ddsconn.AppendUserAgent(Terraform, terraformVersion)
		ddsconn.AppendUserAgent(Provider, providerVersion)
		ddsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		ddsconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			ddsconn.SetHttpsProxy(client.config.Proxy)
		}
		client.ddsconn = ddsconn
	}

	return do(client.ddsconn)
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
func (client *ApsaraStackClient) AccountId() (string, error) {
	client.accountIdMutex.Lock()
	defer client.accountIdMutex.Unlock()

	if client.accountId == "" {
		log.Printf("[DEBUG] account_id not provided, attempting to retrieve it automatically...")
		identity, err := client.GetCallerIdentity()
		if err != nil {
			return "", err
		}
		if identity.AccountId == "" {
			return "", fmt.Errorf("caller identity doesn't contain any AccountId")
		}
		client.accountId = identity.AccountId
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
func (client *ApsaraStackClient) GetCallerIdentity() (*sts.GetCallerIdentityResponse, error) {
	args := sts.CreateGetCallerIdentityRequest()

	endpoint := client.config.StsEndpoint
	if endpoint == "" {
		endpoint = loadEndpoint(client.config.RegionId, STSCode)
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, string(STSCode), endpoint)
	}
	stsClient, err := sts.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the STS client: %#v", err)
	}

	stsClient.AppendUserAgent(Terraform, terraformVersion)
	stsClient.AppendUserAgent(Provider, providerVersion)
	stsClient.AppendUserAgent(Module, client.config.ConfigurationSource)
	stsClient.SetHTTPSInsecure(client.config.Insecure)
	if client.config.Proxy != "" {
		stsClient.SetHttpsProxy(client.config.Proxy)
	}

	identity, err := stsClient.GetCallerIdentity(args)
	if err != nil {
		return nil, err
	}
	if identity == nil {
		return nil, fmt.Errorf("caller identity not found")
	}
	return identity, err
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
func (client *ApsaraStackClient) WithOssClient(do func(*oss.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the OSS client if necessary
	if client.ossconn == nil {
		schma := "https"
		endpoint := client.config.OssEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, OSSCode)
		}
		if endpoint == "" {
			endpointItem, _ := client.describeEndpointForService(strings.ToLower(string(OSSCode)))
			if endpointItem != nil {
				if len(endpointItem.Protocols.Protocols) > 0 {
					// HTTP or HTTPS
					schma = strings.ToLower(endpointItem.Protocols.Protocols[0])
					for _, p := range endpointItem.Protocols.Protocols {
						if strings.ToLower(p) == "https" {
							schma = strings.ToLower(p)
							break
						}
					}
				}
				endpoint = endpointItem.Endpoint
			} else {
				endpoint = fmt.Sprintf("oss-%s.aliyuncs.com", client.RegionId)
			}
		}
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("%s://%s", schma, endpoint)
		}

		clientOptions := []oss.ClientOption{oss.UserAgent(client.getUserAgent()),
			oss.SecurityToken(client.config.SecurityToken)}
		proxy, err := client.getHttpProxy()
		if proxy != nil {
			skip, err := client.skipProxy(endpoint)
			if err != nil {
				return nil, err
			}
			if !skip {
				clientOptions = append(clientOptions, oss.Proxy(proxy.String()))
			}
		}

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

func (client *ApsaraStackClient) WithRdsClient(do func(*rds.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the RDS client if necessary
	if client.rdsconn == nil {
		endpoint := client.config.RdsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, RDSCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(RDSCode), endpoint)
		}
		rdsconn, err := rds.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the RDS client: %#v", err)
		}

		rdsconn.AppendUserAgent(Terraform, terraformVersion)
		rdsconn.AppendUserAgent(Provider, providerVersion)
		rdsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		rdsconn.SetHTTPSInsecure(client.config.Insecure)
		if client.config.Proxy != "" {
			rdsconn.SetHttpsProxy(client.config.Proxy)
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
func (client *ApsaraStackClient) getUserAgent() string {
	return fmt.Sprintf("%s/%s %s/%s %s/%s", Terraform, terraformVersion, Provider, providerVersion, Module, client.config.ConfigurationSource)
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
			endpoint = loadEndpoint(client.config.RegionId, CONTAINCode)
		}
		if endpoint != "" {
			if !strings.HasPrefix(endpoint, "http") {
				endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "://"))
			}
			csconn.SetEndpoint(endpoint)
		}
		client.csconn = csconn
	}

	return do(client.csconn)
}
