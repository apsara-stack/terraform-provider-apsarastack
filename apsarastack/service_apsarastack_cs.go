package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"log"
	"time"
)

type CsService struct {
	client *connectivity.ApsaraStackClient
}

const UpgradeClusterTimeout = 30 * time.Minute

func (s *CsService) DescribeCsKubernetes(id string) (cl *cs.KubernetesClusterDetail, err error) {
	invoker := NewInvoker()
	cluster := &cs.KubernetesClusterDetail{}
	cluster.ClusterId = ""
	var requestInfo *cs.Client
	var response interface{}
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":         s.client.RegionId,
		"AccessKeySecret":  s.client.SecretKey,
		"Product":          "CS",
		"Department":       s.client.Department,
		"ResourceGroup":    s.client.ResourceGroup,
		"Action":           "DescribeClusters",
		"AccountInfo":      "123456",
		"Version":          "2015-12-15",
		"SignatureVersion": "1.0",
		"ProductName":      "cs",
	}
	request.Method = "POST" // Set request method
	request.Product = "Cs"  // Specify product
	// request.Domain =       // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2015-12-15" // Specify product version
	request.ServiceCode = "cs"
	request.Scheme = "http" // Set request scheme. Default: http
	request.ApiName = "DescribeClusters"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	if err := invoker.Run(func() error {
		raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		response = raw
		return err
	}); err != nil {
		if IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
			return cluster, WrapErrorf(err, NotFoundMsg, DenverdinoApsaraStackgo)
		}
		return cluster, WrapErrorf(err, DefaultErrorMsg, id, "DescribeKubernetesCluster", DenverdinoApsaraStackgo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = id
		addDebug("DescribeKubernetesCluster", response, requestInfo, requestMap)
	}
	Cdetails := []Cluster{}
	clusterdetails, _ := response.(*responses.CommonResponse)
	_ = json.Unmarshal(clusterdetails.GetHttpContentBytes(), &Cdetails)

	if len(Cdetails) < 1 {
		return
	}

	cluster = &cs.KubernetesClusterDetail{}
	for _, k := range Cdetails {
		if k.ClusterID == id {
			cluster.Name = k.Name
			cluster.State = k.State
			cluster.ClusterId = k.ClusterID
			cluster.ClusterType = cs.KubernetesClusterType(k.ClusterType)
			cluster.VpcId = k.VpcID
			cluster.ResourceGroupId = k.ResourceGroupID
			cluster.ContainerCIDR = k.SubnetCidr
			cluster.CurrentVersion = k.CurrentVersion
			cluster.DeletionProtection = k.DeletionProtection
			cluster.RegionId = common.Region(k.RegionID)
			cluster.Size = k.Size
			cluster.IngressLoadbalancerId = k.ExternalLoadbalancerID
			cluster.InitVersion = k.InitVersion
			cluster.MetaData = k.MetaData
			cluster.NetworkMode = k.NetworkMode
			cluster.PrivateZone = k.PrivateZone
			cluster.Profile = k.Profile
			cluster.VSwitchIds = k.VswitchID
			//cluster.Updated=k.Updated
			//cluster.Created= k.Created.
			break
		}
	}
	if cluster.ClusterId != id {
		return cluster, WrapErrorf(Error(GetNotFoundMessage("CsKubernetes", id)), NotFoundMsg, ProviderERROR)
	}
	return cluster, nil
}

func (s *CsService) CsKubernetesInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsKubernetes(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if string(object.State) == failState {
				return object, string(object.State), WrapError(Error(FailedToReachTargetStatus, string(object.State)))
			}
		}
		return object, string(object.State), nil
	}
}

type Cluster struct {
	_                      string `json:"-"`
	Department             int64  `json:"Department"`
	DepartmentName         string `json:"DepartmentName"`
	ResourceGroup          int64  `json:"ResourceGroup"`
	ResourceGroupName      string `json:"ResourceGroupName"`
	ClusterHealthy         string `json:"cluster_healthy"`
	ClusterID              string `json:"cluster_id"`
	ClusterType            string `json:"cluster_type"`
	Created                string `json:"created"`
	CurrentVersion         string `json:"current_version"`
	DataDiskCategory       string `json:"data_disk_category"`
	DataDiskSize           int64  `json:"data_disk_size"`
	DeletionProtection     bool   `json:"deletion_protection"`
	DockerVersion          string `json:"docker_version"`
	EnabledMigration       bool   `json:"enabled_migration"`
	ErrMsg                 string `json:"err_msg"`
	ExternalLoadbalancerID string `json:"external_loadbalancer_id"`
	GwBridge               string `json:"gw_bridge"`
	InitVersion            string `json:"init_version"`
	InstanceType           string `json:"instance_type"`
	MasterURL              string `json:"master_url"`
	MetaData               string `json:"meta_data"`
	Name                   string `json:"name"`
	NeedUpdateAgent        bool   `json:"need_update_agent"`
	NetworkMode            string `json:"network_mode"`
	NodeStatus             string `json:"node_status"`
	Outputs                []struct {
		Description string      `json:"Description"`
		OutputKey   string      `json:"OutputKey"`
		OutputValue interface{} `json:"OutputValue"`
	} `json:"outputs"`
	Parameters struct {
		ALIYUN__AccountID        string `json:"ALIYUN::AccountId"`
		ALIYUN__NoValue          string `json:"ALIYUN::NoValue"`
		ALIYUN__Region           string `json:"ALIYUN::Region"`
		ALIYUN__StackID          string `json:"ALIYUN::StackId"`
		ALIYUN__StackName        string `json:"ALIYUN::StackName"`
		AdjustmentType           string `json:"AdjustmentType"`
		AuditFlags               string `json:"AuditFlags"`
		BetaVersion              string `json:"BetaVersion"`
		Ca                       string `json:"CA"`
		ClientCA                 string `json:"ClientCA"`
		CloudMonitorFlags        string `json:"CloudMonitorFlags"`
		CloudMonitorVersion      string `json:"CloudMonitorVersion"`
		ContainerCIDR            string `json:"ContainerCIDR"`
		DockerVersion            string `json:"DockerVersion"`
		Eip                      string `json:"Eip"`
		EipAddress               string `json:"EipAddress"`
		ElasticSearchHost        string `json:"ElasticSearchHost"`
		ElasticSearchPass        string `json:"ElasticSearchPass"`
		ElasticSearchPort        string `json:"ElasticSearchPort"`
		ElasticSearchUser        string `json:"ElasticSearchUser"`
		EtcdVersion              string `json:"EtcdVersion"`
		ExecuteVersion           string `json:"ExecuteVersion"`
		GPUFlags                 string `json:"GPUFlags"`
		HealthCheckType          string `json:"HealthCheckType"`
		IPVSEnable               string `json:"IPVSEnable"`
		ImageID                  string `json:"ImageId"`
		K8SMasterPolicyDocument  string `json:"K8SMasterPolicyDocument"`
		K8sWorkerPolicyDocument  string `json:"K8sWorkerPolicyDocument"`
		Key                      string `json:"Key"`
		KeyPair                  string `json:"KeyPair"`
		KubernetesVersion        string `json:"KubernetesVersion"`
		LoggingType              string `json:"LoggingType"`
		MasterAutoRenew          string `json:"MasterAutoRenew"`
		MasterAutoRenewPeriod    string `json:"MasterAutoRenewPeriod"`
		MasterDataDisk           string `json:"MasterDataDisk"`
		MasterDataDiskCategory   string `json:"MasterDataDiskCategory"`
		MasterDataDiskDevice     string `json:"MasterDataDiskDevice"`
		MasterDataDiskSize       string `json:"MasterDataDiskSize"`
		MasterImageID            string `json:"MasterImageId"`
		MasterInstanceChargeType string `json:"MasterInstanceChargeType"`
		MasterInstanceType       string `json:"MasterInstanceType"`
		MasterKeyPair            string `json:"MasterKeyPair"`
		MasterLoginPassword      string `json:"MasterLoginPassword"`
		MasterPeriod             string `json:"MasterPeriod"`
		MasterPeriodUnit         string `json:"MasterPeriodUnit"`
		MasterSystemDiskCategory string `json:"MasterSystemDiskCategory"`
		MasterSystemDiskSize     string `json:"MasterSystemDiskSize"`
		NatGateway               string `json:"NatGateway"`
		NatGatewayID             string `json:"NatGatewayId"`
		Network                  string `json:"Network"`
		NodeCIDRMask             string `json:"NodeCIDRMask"`
		NumOfNodes               string `json:"NumOfNodes"`
		Password                 string `json:"Password"`
		ProtectedInstances       string `json:"ProtectedInstances"`
		PublicSLB                string `json:"PublicSLB"`
		RemoveInstanceIds        string `json:"RemoveInstanceIds"`
		SLSProjectName           string `json:"SLSProjectName"`
		SNatEntry                string `json:"SNatEntry"`
		SSHFlags                 string `json:"SSHFlags"`
		ServiceCIDR              string `json:"ServiceCIDR"`
		SnatTableID              string `json:"SnatTableId"`
		UserCA                   string `json:"UserCA"`
		VSwitchID                string `json:"VSwitchId"`
		VpcID                    string `json:"VpcId"`
		WillReplace              string `json:"WillReplace"`
		WorkerAutoRenew          string `json:"WorkerAutoRenew"`
		WorkerAutoRenewPeriod    string `json:"WorkerAutoRenewPeriod"`
		WorkerDataDisk           string `json:"WorkerDataDisk"`
		WorkerDataDiskCategory   string `json:"WorkerDataDiskCategory"`
		WorkerDataDiskDevice     string `json:"WorkerDataDiskDevice"`
		WorkerDataDiskSize       string `json:"WorkerDataDiskSize"`
		WorkerImageID            string `json:"WorkerImageId"`
		WorkerInstanceChargeType string `json:"WorkerInstanceChargeType"`
		WorkerInstanceType       string `json:"WorkerInstanceType"`
		WorkerKeyPair            string `json:"WorkerKeyPair"`
		WorkerLoginPassword      string `json:"WorkerLoginPassword"`
		WorkerPeriod             string `json:"WorkerPeriod"`
		WorkerPeriodUnit         string `json:"WorkerPeriodUnit"`
		WorkerSystemDiskCategory string `json:"WorkerSystemDiskCategory"`
		WorkerSystemDiskSize     string `json:"WorkerSystemDiskSize"`
		ZoneID                   string `json:"ZoneId"`
	} `json:"parameters"`
	Port              int64  `json:"port"`
	PrivateZone       bool   `json:"private_zone"`
	Profile           string `json:"profile"`
	RegionID          string `json:"region_id"`
	ResourceGroupID   string `json:"resource_group_id"`
	SecurityGroupID   string `json:"security_group_id"`
	Size              int64  `json:"size"`
	State             string `json:"state"`
	SubnetCidr        string `json:"subnet_cidr"`
	SwarmMode         bool   `json:"swarm_mode"`
	Updated           string `json:"updated"`
	UpgradeComponents struct {
		Kubernetes struct {
			CanUpgrade     bool   `json:"can_upgrade"`
			Changed        string `json:"changed"`
			ComponentName  string `json:"component_name"`
			Exist          bool   `json:"exist"`
			Force          bool   `json:"force"`
			Message        string `json:"message"`
			NextVersion    string `json:"next_version"`
			Policy         string `json:"policy"`
			ReadyToUpgrade string `json:"ready_to_upgrade"`
			Required       bool   `json:"required"`
			Version        string `json:"version"`
		} `json:"Kubernetes"`
	} `json:"upgrade_components"`
	VpcID       string `json:"vpc_id"`
	VswitchCidr string `json:"vswitch_cidr"`
	VswitchID   string `json:"vswitch_id"`
	ZoneID      string `json:"zone_id"`
}

func (s *CsService) UpgradeCluster(clusterId string, args *cs.UpgradeClusterArgs) error {
	invoker := NewInvoker()
	err := invoker.Run(func() error {
		_, e := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.UpgradeCluster(clusterId, args)
		})
		if e != nil {
			return e
		}
		return nil
	})

	if err != nil {
		return WrapError(err)
	}

	state, upgradeError := s.WaitForUpgradeCluster(clusterId, "Upgrade")
	if state == cs.Task_Status_Success && upgradeError == nil {
		return nil
	}

	// if upgrade failed cancel the task
	err = invoker.Run(func() error {
		_, e := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.CancelUpgradeCluster(clusterId)
		})
		if e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		return WrapError(upgradeError)
	}

	if state, err := s.WaitForUpgradeCluster(clusterId, "CancelUpgrade"); err != nil || state != cs.Task_Status_Success {
		log.Printf("[WARN] %s ACK Cluster cancel upgrade error: %#v", clusterId, err)
	}

	return WrapError(upgradeError)
}

func (s *CsService) WaitForUpgradeCluster(clusterId string, action string) (string, error) {
	err := resource.Retry(UpgradeClusterTimeout, func() *resource.RetryError {
		resp, err := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.QueryUpgradeClusterResult(clusterId)
		})
		if err != nil || resp == nil {
			return resource.RetryableError(err)
		}

		upgradeResult := resp.(*cs.UpgradeClusterResult)
		if upgradeResult.UpgradeStep == cs.UpgradeStep_Success {
			return nil
		}

		if upgradeResult.UpgradeStep == cs.UpgradeStep_Pause && upgradeResult.UpgradeStatus.Failed == "true" {
			msg := ""
			events := upgradeResult.UpgradeStatus.Events
			if len(events) > 0 {
				msg = events[len(events)-1].Message
			}
			return resource.NonRetryableError(fmt.Errorf("faild to %s cluster, error: %s", action, msg))
		}
		return resource.RetryableError(fmt.Errorf("%s cluster state not matched", action))
	})

	if err == nil {
		log.Printf("[INFO] %s ACK Cluster %s successed", action, clusterId)
		return cs.Task_Status_Success, nil
	}

	return cs.Task_Status_Failed, WrapError(err)
}
