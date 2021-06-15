---
subcategory: "Container Service (CS)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_cs_kubernetes"
sidebar_current: "docs-apsarastack-resource-cs-kubernetes"
description: |-
  Provides a Apsarastack resource to manage container kubernetes cluster.
---

# apsarastack\_cs\_kubernetes

This resource will help you to manage a Kubernetes Cluster in Apsarastack Kubernetes Service. 

-> **NOTE:** Each kubernetes cluster contains 3 master nodes and those number cannot be changed at now.

-> **NOTE:** Creating kubernetes cluster need to install several packages and it will cost about 15 minutes. Please be patient.
## Example Usage
```$xslt

// If there is not specifying vpc_id, the module will launch a new vpc
resource "apsarastack_vpc" "vpc" {
   name = "testing_cs"
   cidr_block = "10.0.0.0/8"
}

// According to the vswitch cidr blocks to launch several vswitches
resource "apsarastack_vswitch" "vswitches" {
   vpc_id = "${apsarastack_vpc.default.id}"
   cidr_block        = "10.1.0.0/16"
   name = "apsara_vswitch
   availability_zone = var.availability_zone
}

resource "apsarastack_cs_kubernetes" "k8s" {
   name="apsara_test"
   master_vswitch_ids=[apsarastack_vswitch.vswitches.id,apsarastack_vswitch.vswitches.id,apsarastack_vswitch.vswitches.id]
   worker_vswitch_ids=[apsarastack_vswitch.vswitches.id,apsarastack_vswitch.vswitches.id]
   version="1.14.8-aliyun.1"
   master_count=3
   timeout_mins=60
   master_disk_category="cloud_ssd"
   master_disk_size=45
   worker_disk_category="cloud_ssd"
   worker_disk_size=30
   delete_protection=false
   worker_data_disk=true
   worker_data_disk_category="cloud_ssd"
   worker_data_disk_size=100
   new_nat_gateway=false
   vpc_id=apsarastack_vpc.vpc.id
   slb_internet_enabled=false
   proxy_mode="ipvs"
   master_instance_types = ["ecs.e4.large","ecs.e4.large","ecs.e4.large"]
   worker_instance_types = ["ecs.e4.large","ecs.e4.large"]
   num_of_nodes         = 2
   enable_ssh            = true
   password              = "Test@123"
   pod_cidr              = "172.23.0.0/16"
   service_cidr          = "172.24.0.0/20"
   node_cidr_mask="26"
   dynamic "addons" {
      for_each = var.cluster_addons
      content {
        name                    = lookup(addons.value, "name", var.cluster_addons)
        config                  = lookup(addons.value, "config", var.cluster_addons)
      }
   }
  user_data="ZWNobyBoZWxsbw=="
}
```

## Argument Reference

The following arguments are supported:

#### Global params
* `name` - (Optional) The kubernetes cluster's name. It is unique in one Apsarastack account.
* `version` - (Optional) Desired Kubernetes version. If you do not specify a value, the latest available version at resource creation is used and no upgrades will occur except you set a higher version number. The value must be configured and increased to upgrade the version when desired. Downgrades are not supported by ACK.
* `password` - (Required, Sensitive) The password of ssh login cluster node. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `kms_encrypted_password` - (Required) An KMS encrypts password used to a cs kubernetes. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `enable_ssh` - (Optional) Enable login to the node through SSH. default: false 
* `cpu_policy` - kubelet cpu policy. options: static|none. default: none.
* `proxy_mode` - Proxy mode is option of kube-proxy. options: iptables|ipvs. default: ipvs.
* `user_data` - (Optional) Windows instances support batch and PowerShell scripts. If your script file is larger than 1 KB, we recommend that you upload the script to Object Storage Service (OSS) and pull it through the internal endpoint of your OSS bucket.
* `instances`- (Optional) A list of instances that can be attached as worker nodes in the same Vpc.
* `runtime`-  (Optional) The platform on which the clusters are going to run.
    * `name`- (Optional) Name of the runtime platform
    * `version`- (Optional) Version of the runtime platform
    
#### Network
* `pod_cidr` - (Required) [Flannel Specific] The CIDR block for the pod network when using Flannel. 
* `pod_vswitch_ids` - (Required) [Terway Specific] The vswitches for the pod network when using Terway.Be careful the `pod_vswitch_ids` can not equal to `worker_vswtich_ids` or `master_vswtich_ids` but must be in same availability zones.
* `new_nat_gateway` - (Optional) Whether to create a new nat gateway while creating kubernetes cluster. Default to true. Then openapi in Apsarastack are not all on intranet, So turn this option on is a good choice.
* `service_cidr` - (Optional) The CIDR block for the service network. It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
* `node_cidr_mask` - (Optional) The node cidr block to specific how many pods can run on single node. 24-28 is allowed. 24 means 2^(32-24)-1=255 and the node can run at most 255 pods. default: 24
* `slb_internet_enabled` - (Optional) Whether to create internet load balancer for API Server. Default to true.

If you want to use `Terway` as CNI network plugin, You need to specific the `pod_vswitch_ids` field and addons with `terway-eniip`.    
If you want to use `Flannel` as CNI network plugin, You need to specific the `pod_cidr` field and addons with `flannel`.

#### Master params
* `master_disk_category` - (Optional) The system disk category of master node. Its valid value are `cloud_ssd` and `cloud_efficiency`. Default to `cloud_efficiency`.
* `master_disk_size` - (Optional) The system disk size of master node. Its valid value range [20~500] in GB. Default to 20.
* `master_vswtich_ids` - (Required) The vswitches used by master, you can specific 3 or 5 vswitches because of the amount of masters. Detailed below.
* `master_instance_types` - (Required) The instance type of master node. Specify one type for single AZ Cluster, three types for MultiAZ Cluster.

#### Worker params 
* `num_of_nodes` - (Required) The worker node number of the kubernetes cluster. Default to 3. It is limited up to 50 and if you want to enlarge it, please apply white list or contact with us.
* `worker_vswtich_ids` - (Required) The vswitches used by worker, you can specific 1 or more than 1 vswitches.
* `worker_disk_size` - (Optional) The system disk size of worker node. Its valid value range [20~32768] in GB. 
* `worker_disk_category` - (Optional) The system disk category of worker node. Its valid value are cloud, cloud_ssd, cloud_essd and cloud_efficiency. 
* `worker_data_disk_category` - (Optional) The data disk category of worker node. Its valid value are cloud, cloud_ssd, cloud_essd and cloud_efficiency. 
* `worker_data_disk_size` - (Optional) The data disk size of worker node. Its valid value range [40~500] in GB. 
* `worker_instance_types` - (Required) The instance type of worker node. Specify one type for single AZ Cluster, three types for MultiAZ Cluster.

#### Computed params (No need to configure) 
* `kube_config` - (Optional) The path of kube config, like `~/.kube/config`.
* `client_cert` - (Optional) The path of client certificate, like `~/.kube/client-cert.pem`.
* `client_key` - (Optional) The path of client key, like `~/.kube/client-key.pem`.
* `cluster_ca_cert` - (Optional) The path of cluster ca certificate, like `~/.kube/cluster-ca-cert.pem`
* `availability_zone` - (Optional) The Zone where new kubernetes cluster will be located. If it is not be specified, the `vswitch_ids` should be set, its value will be vswitch's zone.
  
### Timeouts
The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when creating the kubernetes cluster (until it reaches the initial `running` status). 
* `update` - (Defaults to 60 mins) Used when activating the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when terminating the kubernetes cluster. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the container cluster.
* `name` - The name of the container cluster.
* `availability_zone` - The ID of availability zone.
* `vpc_id` - The ID of VPC where the current cluster is located.
* `slb_intranet` - The ID of private load balancer where the current cluster master node is located.
* `security_group_id` - The ID of security group where the current cluster worker node is located.
* `nat_gateway_id` - The ID of nat gateway used to launch kubernetes cluster.
* `master_nodes` - List of cluster master nodes. It contains several attributes to `Block Nodes`.
* `worker_nodes` - List of cluster worker nodes. It contains several attributes to `Block Nodes`.
* `connections` - Map of kubernetes cluster connection information. It contains several attributes to `Block Connections`.
* `version` - The Kubernetes server version for the cluster.
* `worker_ram_role_name` - The RamRole Name attached to worker node.

### Block Nodes
* `id` - ID of the node.
* `name` - Node name.
* `private_ip` - The private IP address of node.


### Block Connections
* `api_server_internet` - API Server Internet endpoint.
* `api_server_intranet` - API Server Intranet endpoint.
* `master_public_ip` - Master node SSH IP address.
* `service_domain` - Service Access Domain.

