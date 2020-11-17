//

package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

//func main() {
//	plugin.Serve(&plugin.ServeOpts{
//		ProviderFunc: apsarastack.Provider,
//	})
//}
//func main(){
//	GetClusterDetails()
//}
//
//func GetClusterDetails(){
//	access := "ckhCs1KpWEQtvYZD"
//	secret := "2lY9uNh155EvHJrmPuqYNzCPEksnx1"
//	region := "cn-neimeng-env30-d01"
//	endpoint := "server.asapi.cn-neimeng-env30-d01.intra.env30.shuguang.com/asapi/v3"
//	proxy := "http://100.67.76.9:53001"
//	department:="54437"
//	resource_group:="571"
//	client, err := sdk.NewClientWithAccessKey(region, access, secret)
//	client.Domain = endpoint
//	if err != nil {
//		fmt.Print("Error in client")
//	}
//	request := requests.NewCommonRequest()
//	request.QueryParams = map[string]string{
//		"RegionId": region,
//		"AccessKeySecret": secret,
//		"Product": "CS",
//		"Department": department,
//		"ResourceGroup": resource_group,
//		"Action": "DescribeClusters",
//		"AccountInfo": "123456",
//		"Version": "2015-12-15",
//		"SignatureVersion": "1.0",
//		"ProductName": "cs",
//		//"name": "afgh",
//		//"vpcid":"vpc-0rvggkz6dnbas2wplcc9k",
//		//"vswitchid": "vsw-0rvjt1f6qvx073g7lxjlc",
//		/*"X-acs-body": fmt.Sprintf("{\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%d,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%d,\"%s\":%d,\"%s\":%d,\"%s\":%d,\"%s\":%d,\"%s\":%t,\"%s\":%t,\"%s\":%t,\"%s\":%t,\"%s\":%t}",
//			"Product","Cs",
//			"cluster_type", "Kubernetes",
//			"RegionId","cn-neimeng-env30-d01",
//			"timeout_mins",60,
//			"disable_rollback", true,
//			"kubernetes_version", "1.14.8-aliyun.1",
//			"container_cidr", "10.14.0.0/16",
//			"service_cidr", "10.15.0.0/16",
//			"name", "k8s-success",
//			"vpcid","vpc-0rvggkz6dnbas2wplcc9k",
//			"vswitchid", "vsw-0rvjt1f6qvx073g7lxjlc",
//			"master_instance_type", "ecs.n4.2xlarge",
//			"master_system_disk_category", "cloud_ssd",
//			"worker_instance_type", "ecs.n4.2xlarge",
//			"worker_system_disk_category", "cloud_ssd",
//			"worker_data_disk_category", "cloud_efficiency",
//			"login_Password", "P@ssw0rd",
//			"master_system_disk_size", 200,
//			"worker_data_disk_size", 200,
//			"worker_system_disk_size", 200,
//			"num_of_nodes", 3,
//			"master_count", 3,
//			"worker_data_disk", true,
//			"snat_entry", false,
//			"endpoint_public_access",true,
//			"ssh_flags", true,
//			"deletion_protection", true),
//		//}*/
//		//
//	}
//	request.Method = "GET" // Set request method
//	request.Product = "Cs" // Specify product
//	request.Domain = endpoint // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
//	request.Version = "2015-12-15" // Specify product version
//	request.ServiceCode="cs"
//	request.Scheme = "http" // Set request scheme. Default: http
//	request.ApiName = "DescribeClusters"
//	request.Headers = map[string]string{"RegionId": region}
//
//
//
//	//fmt.Print(request)
//	client.SetHttpProxy(proxy)
//
//
//
//	resp := responses.BaseResponse{}
//	request.TransToAcsRequest()
//
//	//err=client.Init()
//
//	err = client.DoAction(request, &resp)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("Response: %s", resp)
//
//
//
//}
func main() {
	access := "ckhCs1KpWEQtvYZD"
	secret := "2lY9uNh155EvHJrmPuqYNzCPEksnx1"
	region := "cn-neimeng-env30-d01"
	endpoint := "server.asapi.cn-neimeng-env30-d01.intra.env30.shuguang.com/asapi/v3"
	proxy := "http://100.67.76.9:53001"
	department := "54437"
	resource_group := "571"
	client, err := sdk.NewClientWithAccessKey(region, access, secret)
	client.Domain = endpoint
	if err != nil {
		fmt.Print("Error in client")
	}
	request := requests.NewCommonRequest()
	request.QueryParams = map[string]string{
		"RegionId":         region,
		"AccessKeySecret":  secret,
		"Product":          "CS",
		"Department":       department,
		"ResourceGroup":    resource_group,
		"Action":           "CreateCluster",
		"AccountInfo":      "123456",
		"Version":          "2015-12-15",
		"SignatureVersion": "1.0",
		"ProductName":      "cs",
		//"name": "afgh",
		"vpcid":     "vpc-0rvuzw2wep2wfofle24qb",
		"vswitchid": "vsw-0rvw06mtyu2yerpmc3a6e",
		"X-acs-body": fmt.Sprintf("{\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%d,\"%s\":%t,\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":\"%s\",\"%s\":%d,\"%s\":%d,\"%s\":%d,\"%s\":%d,\"%s\":%d,\"%s\":%t,\"%s\":%t,\"%s\":%t,\"%s\":%t,\"%s\":%t}",
			"Product", "Cs",
			"cluster_type", "Kubernetes",
			"RegionId", "cn-neimeng-env30-d01",
			"timeout_mins", 60,
			"disable_rollback", true,
			"kubernetes_version", "1.14.8-aliyun.1",
			"container_cidr", "172.20.0.0/16",
			"service_cidr", "172.20.0.0/20",
			"name", "k8ss-success",
			"vpcid", "vpc-0rvggkz6dnbas2wplcc9k",
			"vswitchid", "vsw-0rvjt1f6qvx073g7lxjlc",
			"master_instance_type", "ecs.n4.2xlarge",

			"master_system_disk_category", "cloud_ssd",

			"worker_instance_type", "ecs.n4.2xlarge",

			"worker_system_disk_category", "cloud_ssd",

			"worker_data_disk_category", "cloud_efficiency",

			//

			"login_Password", "P@ssw0rd",
			"master_system_disk_size", 200,
			"worker_data_disk_size", 200,
			"worker_system_disk_size", 200,
			"num_of_nodes", 3,
			"master_count", 3,
			"worker_data_disk", true,
			"snat_entry", false,
			"endpoint_public_access", true,
			"ssh_flags", true,
			"deletion_protection", true),
		//}
		//
	}
	request.Method = "POST"        // Set request method
	request.Product = "Cs"         // Specify product
	request.Domain = endpoint      // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2015-12-15" // Specify product version
	request.ServiceCode = "cs"
	request.Scheme = "http" // Set request scheme. Default: http
	request.ApiName = "CreateCluster"
	request.Headers = map[string]string{"RegionId": region}

	//fmt.Print(request)
	client.SetHttpProxy(proxy)

	resp := responses.BaseResponse{}
	request.TransToAcsRequest()

	err = client.DoAction(request, &resp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %s", resp)

}

//func main() {
//	//access := "n5OOaQDHYOWdhNcT"
//	//secret := "r38I511wXWopAE1f7FnhbUY4QBEsXv"
//	access := "ckhCs1KpWEQtvYZD"
//	secret := "2lY9uNh155EvHJrmPuqYNzCPEksnx1"
//	region := "cn-neimeng-env30-d01"
//	//region := "cn-qingdao-env66-d01"
//	//endpoint:= "server.asapi.cn-qingdao-env66-d01.intra.env66.shuguang.com/asapi/v3"
//	endpoint:= "server.asapi.cn-neimeng-env30-d01.intra.env30.shuguang.com/asapi/v3"
//	//endpoint:= "https://asc.inter.env30.shuguang.com/module/workbench"
//	proxy := "http://100.67.76.9:53001"
//	department:="54437"
//	//department:="303707"
//	resource_group:="571"
//	//resource_group:="1740"
//	client,err:= sdk.NewClientWithAccessKey(region,access,secret)
//	client.Domain=endpoint
//	if err!=nil{
//		fmt.Print("Error in client")
//	}
//	request:= requests.NewCommonRequest()
//	request.Method = "GET"                // Set request method
//	request.Product = "CS"            // Specify product
//	request.Domain = endpoint          // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
//	request.Version = "2015-12-15"            // Specify product version
//	request.Scheme = "http"                // Set request scheme. Default: http
//	request.ApiName = "DescribeClusters"
//	request.Headers = map[string]string{"RegionId": region}
//	request.QueryParams = map[string]string{
//		"AccessKeySecret": secret,
//		"AccessKeyId": access,
//		"Product": "CS",
//		"Department": department,
//		"ResourceGroup": resource_group,
//		"RegionId": region,
//		"Action": "DescribeClusters",
//		"Version":"2015-12-15",
//		"vpc_id":"vpc-0rvuzw2wep2wfofle24qb",
//		//"ParentId": "17",
//		"Name": "wangk8s1",
//		//"Id":"54438",
//	}
//	//fmt.Print(request)
//	client.SetHttpProxy(proxy)
//	resp := responses.BaseResponse{}
//	request.TransToAcsRequest()
//	err= client.DoAction(request,&resp)
//	fmt.Print(request)
//	if err!=nil{
//		fmt.Printf("Response %s", resp.GetHttpContentString())
//		panic(err)
//	}
//	fmt.Printf("Response: %s",resp)
//}
